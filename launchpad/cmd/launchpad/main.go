package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
	//"../telephono"
)

type remoteSyncClient struct {
	sshSession *ssh.Session
	sftpClient *sftp.Client
}

func parseArgs() (target, arch *string, args []string) {
	arch = flag.String("a", "", "Architecture to use: e.g. linux-amd64, darwin-amd64, linux-arm")
	target = flag.String("t", "", "Target host: user@host")
	flag.Parse()

	if *target != "" && !strings.Contains(*target, "@") {
		user, err := user.Current()
		if err != nil {
			log.Fatal("Failed to get username")
		}
		username := user.Username
		*target = username + "@" + *target
	}
	return target, arch, flag.Args()
}

func getSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func getClient(username, hostname string) *ssh.Client {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			getSSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	return client
}

func remoteMetadata(client *ssh.Client) (string, string) {
	// We need the architecture to copy up the right call-buddy binary
	archSession, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create meta session: ", err)
	}
	defer archSession.Close()
	var out bytes.Buffer
	archSession.Stdout = &out

	archSession.Run("uname -s -m | tr ' ' '-' | tr '[:upper:]' '[:lower:]'")
	// arch should look like "linux-x86_64"
	arch := strings.TrimSpace(out.String())

	// We can't assume ~/ is /home/<username>/ so we have to get the abs path here
	out.Reset()
	homedirSession, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create meta session: ", err)
	}
	defer homedirSession.Close()
	homedirSession.Stdout = &out
	homedirSession.Run("eval echo ~$USER")
	homedir := strings.TrimSpace(out.String())
	return arch, homedir
}

func cleanupBootstrapCallBuddy(remoteSftpClient *sftp.Client, remoteCallBuddyPath string) {
	// FIXME: Handle errors
	fmt.Println("Removing call-buddy from remote client")
	remoteSftpClient.Remove(remoteCallBuddyPath)
	remoteSftpClient.RemoveDirectory(filepath.Dir(remoteCallBuddyPath))
}

func bootstrapCallBuddy(remoteSftpClient *sftp.Client, localCallBuddyPath, remoteCallBuddyPath string) error {
	remoteCallBuddyDir := filepath.Dir(remoteCallBuddyPath)

	// Check whether the call-buddy dir exists
	_, err := remoteSftpClient.Stat(remoteCallBuddyDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = remoteSftpClient.Mkdir(remoteCallBuddyDir)
			if err != nil {
				log.Printf("Failed to create remote call buddy dir '%s': %s\n", remoteCallBuddyDir, err)
				return err
			}
		} else {
			log.Printf("Failed to read remote call-buddy dir '%s': %s\n", remoteCallBuddyDir, err)
			return err
		}
	} else {
		// We don't know if the existing binary is old, so we'll remove it and
		// replace it with a new one
		// FIXME DG: What if ~/.call-buddy exists, but ~/.call-buddy/call-buddy does not?
		err = remoteSftpClient.Remove(remoteCallBuddyPath)
		if err != nil {
			log.Printf("Failed to remove existing remote call-buddy path '%s': %s\n", remoteCallBuddyPath, err)
			return err
		}
	}
	remoteExe, err := remoteSftpClient.Create(remoteCallBuddyPath)
	if err != nil {
		log.Println("Failed to create remote call-buddy exe: ", err)
		return err
	}
	defer remoteExe.Close()
	remoteExe.Seek(0, os.SEEK_SET) // Make sure we're pointing to the beginning

	buf, err := ioutil.ReadFile(localCallBuddyPath)
	if err != nil {
		log.Println("Failed to read local call-buddy path: ", err)
		return err
	}
	remoteExe.Write(buf)
	remoteExe.Chmod(0755) // Make the file executable
	return nil
}

func closeRemoteSyncing(remoteClient *remoteSyncClient, remoteCallBuddyPath string) {
	// TODO: Cleanup syncing goroutine
	cleanupBootstrapCallBuddy(remoteClient.sftpClient, remoteCallBuddyPath)
	remoteClient.sshSession.Close()
	remoteClient.sftpClient.Close()
}

func spawnRemoteSyncing(client *ssh.Client, remoteCallBuddyPath, arch *string) (*remoteSyncClient, error) {
	// FIXME: Log things out

	syncingSession, err := client.NewSession()
	if err != nil {
		log.Println("Failed to create syncing session: ", err)
		return nil, err
	}
	syncingStdin, err := syncingSession.StdinPipe()
	if err != nil {
		syncingSession.Close()
		return nil, err
	}
	syncingStdout, err := syncingSession.StdoutPipe()
	if err != nil {
		syncingSession.Close()
		return nil, err
	}

	syncingSession.RequestSubsystem("sftp")
	syncingClient, err := sftp.NewClientPipe(syncingStdout, syncingStdin)
	if err != nil {
		syncingSession.Close()
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		syncingSession.Close()
		syncingClient.Close()
		return nil, err
	}
	remoteClient := remoteSyncClient{syncingSession, syncingClient}

	localCallBuddyPath := filepath.Clean(cwd + "/../telephono-ui/build/" + *arch + "/call-buddy")
	fmt.Printf("Syncing call-buddy from %s to remote client at %s\n", localCallBuddyPath, *remoteCallBuddyPath)
	err = bootstrapCallBuddy(remoteClient.sftpClient, localCallBuddyPath, *remoteCallBuddyPath)
	if err != nil {
		// FIXME: Use cleanupRemoteSyncing here, but beware of non bootstrap
		remoteClient.sshSession.Close()
		remoteClient.sftpClient.Close()
		return nil, err
	}
	return &remoteClient, nil
}

func remoteRun(target, arch *string, args []string) {
	s := strings.Split(*target, "@")
	username, hostname := s[0], s[1]
	client := getClient(username, hostname)

	detectedArch, remoteHomeDir := remoteMetadata(client)
	if *arch != "" {
		if *arch != detectedArch {
			log.Printf("Given architecture does not match detected architecture: %s vs %s\n (btw -a can be ommited)", *arch, detectedArch)
		}
	} else {
		if detectedArch == "" {
			log.Fatalf("Failed to get architecture for %s! Try using the -a flag.", hostname)
		}
		arch = &detectedArch
	}

	remoteCallBuddyPath := remoteHomeDir + "/.call-buddy/call-buddy"
	syncingClient, err := spawnRemoteSyncing(client, &remoteCallBuddyPath, &detectedArch)
	if err != nil {
		// TODO: Inspect the error and don't always fail
		log.Fatal("Failed to spawn off remote syncing")
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	// FIXME DG: Make adjustments to current terminal propogate to remote pty
	// (also handle signals)
	width, height, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Failed to get terminal size! ", err)
	}
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	session.Run(remoteCallBuddyPath + " " + strings.Join(args, " "))
	// FIXME: When we do ctl-c to quit call-buddy we don't get here!
	closeRemoteSyncing(syncingClient, remoteCallBuddyPath)
}

func localRun(args []string) {
	localCallBuddyPath := "../telephono-ui/call-buddy"
	// We don't have a full argv here since we are missing arg 0: the executable name
	argsWithArg0 := append([]string{filepath.Base(localCallBuddyPath)}, args...)
	exe := &exec.Cmd{
		Path:   localCallBuddyPath,
		Args:   argsWithArg0,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	if err := exe.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
}

func main() {
	target, arch, args := parseArgs()
	if *target != "" {
		remoteRun(target, arch, args)
	} else {
		localRun(args)
	}
}
