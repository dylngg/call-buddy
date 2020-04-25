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

func bootstrapCallBuddy(localCallBuddyPath, remoteHomeDir string, remoteSftpClient *sftp.Client) string {
	remoteCallBuddyDir := remoteHomeDir + "/.call-buddy"
	remoteCallBuddyPath := remoteCallBuddyDir + "/" + filepath.Base(localCallBuddyPath)
	files, err := remoteSftpClient.ReadDir(remoteHomeDir)
	if err != nil {
		log.Fatal("Failed to read remote homedir: '", remoteHomeDir, "' ", err)
	}
	foundDir := false
	for _, f := range files {
		if f.Name() == filepath.Base(remoteCallBuddyDir) {
			foundDir = true
			break
		}
	}
	if !foundDir {
		err = remoteSftpClient.Mkdir(remoteCallBuddyDir)
		if err != nil {
			log.Fatal("Failed to create remote call buddy dir! ", err)
		}
	} else {
		// We don't know if the existing binary is old, so we'll remove it and
		// replace it with a new one
		// FIXME DG: What if ~/.call-buddy exists, but ~/.call-buddy/call-buddy does not?
		err = remoteSftpClient.Remove(remoteCallBuddyPath)
		if err != nil {
			log.Fatal("Failed to remove existing call buddy path: ", err)
		}
	}
	remoteExe, err := remoteSftpClient.Create(remoteCallBuddyPath)
	if err != nil {
		log.Fatal("Failed to create: ", err)
	}
	defer remoteExe.Close()
	remoteExe.Seek(0, os.SEEK_SET) // Make sure we're pointing to the beginning

	buf, err := ioutil.ReadFile(localCallBuddyPath)
	if err != nil {
		log.Fatal("Failed to read local call buddy path: ", err)
	}
	remoteExe.Write(buf)
	remoteExe.Chmod(0755) // Make the file executable
	return remoteCallBuddyPath
}

func cleanupBootstrapCallBuddy(localCallBuddyPath, remoteHomeDir string, remoteSftpClient *sftp.Client) {
	remoteCallBuddyDir := remoteHomeDir + "/.call-buddy"
	remoteCallBuddyPath := remoteCallBuddyDir + "/" + filepath.Base(localCallBuddyPath)
	remoteSftpClient.Remove(remoteCallBuddyPath)
	remoteSftpClient.RemoveDirectory(remoteCallBuddyDir)
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

	syncingSession, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create syncing session: ", err)
	}
	defer syncingSession.Close()

	syncingStdin, err := syncingSession.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	syncingStdout, err := syncingSession.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	syncingSession.RequestSubsystem("sftp")
	syncingClient, err := sftp.NewClientPipe(syncingStdout, syncingStdin)
	if err != nil {
		// FIXME: Proceed without syncing if we cannot sync
		log.Fatal(err)
	}
	defer syncingClient.Close()

	cwd, _ := os.Getwd()
	localCallBuddyPath := filepath.Clean(cwd + "/../telephono-ui/build/" + *arch + "/call-buddy")
	fmt.Printf("Syncing call-buddy from %s to %s@%s\n", localCallBuddyPath, username, hostname)
	remoteCallBuddyPath := bootstrapCallBuddy(localCallBuddyPath, remoteHomeDir, syncingClient)

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
	fmt.Printf("Removing call-buddy from %s@%s\n", username, hostname)
	cleanupBootstrapCallBuddy(localCallBuddyPath, remoteHomeDir, syncingClient)
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
