package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"telephono"
)

func die(msg string) {
	os.Stderr.WriteString(msg)
	os.Exit(1)
}

func stdinEnvExpanded() string {
	expander := telephono.Expander{}
	expander.AddContributor(telephono.EnvironmentContributor{})
	stdinBuf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	stdinStr := string(stdinBuf)
	stdinStr, err = expander.Expand(stdinStr)
	if err != nil {
		log.Fatal(err)
	}
	return stdinStr
}

func printResponse(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Header) > 1 {
		for key, value := range resp.Header {
			fmt.Printf("%s: %s\n", key, strings.Trim(strings.Join(value, " "), "[]"))
		}
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString(string(body))
}

func main() {
	argLen := len(os.Args[1:])
	if argLen < 2 {
		die("Usage: ./call-buddy <call-type> <url> [content-type]\n")
	}
	callType := strings.ToLower(os.Args[1])
	url := os.Args[2]
	contentType := "text/plain"
	if argLen > 2 {
		contentType = os.Args[3]
	}

	switch callType {
	case "get":
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		printResponse(resp)

	case "post":
		stdin := stdinEnvExpanded()
		resp, err := http.Post(url, contentType, strings.NewReader(stdin))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		printResponse(resp)

	case "head":
		resp, err := http.Head(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		printResponse(resp)

	case "delete":
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Connection", "close")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		printResponse(resp)

	case "put":
		stdin := stdinEnvExpanded()
		req, err := http.NewRequest("PUT", url, strings.NewReader(stdin))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Connection", "close")
		req.Header.Add("Content-type", contentType)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		printResponse(resp)

	default:
		die("Invalid <call-type> given.\n")
	}
}
