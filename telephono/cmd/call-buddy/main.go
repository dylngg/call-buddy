package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func die(msg string) {
	os.Stderr.WriteString(msg)
	os.Exit(1)
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

//Setting the manager
func layout(g *gocui.Gui) error {
	//maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 0, 0, 20, 15); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "\u001b[31mTerminal "+"\u001b[32mCall "+"\u001b[33mBuddy")
	}

	return nil
}

//This is the function to QUIT out of the TUI
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//func update(g *gocui.Gui, v *gocui.View) error {
//	v, err := g.View("viewname")
//	if err != nil {
//		// handle error
//	}
//	v.Clear()
//	fmt.Fprintln(v, "THIS IS UPDATED - ALSO DEREK FUCKS")
//	return nil
//}

func main() {

	//Setting up a new TUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Setting a manager, sets the view (defined as another function above)
	g.SetManagerFunc(layout)

	//Setting keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlZ, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	//if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, update); err != nil {
	//	log.Panicln(err)
	//}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

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
		//printResponse(resp)

	case "post":
		resp, err := http.Post(url, contentType, os.Stdin)
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
		req, err := http.NewRequest("PUT", url, os.Stdin)
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
