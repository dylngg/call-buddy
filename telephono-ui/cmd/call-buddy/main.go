package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

var response_body_string string = ""

const (
	TTL_LINE_VIEW = "title_view"
	// CMD_LINE_VIEW The command line view string
	CMD_LINE_VIEW = "command"
	// MTD_BODY_VIEW The method body view string
	MTD_BODY_VIEW = "method_body"
	// RQT_HEAD_VIEW The request header view string
	RQT_HEAD_VIEW = "request_head"
	// RQT_BODY_VIEW The request body view string
	RQT_BODY_VIEW = "request_body"
	// RSP_BODY_VIEW The response body view string
	RSP_BODY_VIEW = "response_body"
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
			response_body_string += fmt.Sprintf("%s: %s\n", key, strings.Trim(strings.Join(value, " "), "[]"))
		}
		response_body_string += "\n"
	}
	response_body_string += string(body)
}

//Setting the manager
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	realMaxX, realMaxY := maxX-1, maxY-1
	verticalSplitX := 27         // Defines the vertical split down to the command line
	horizontalSplitY := maxY - 4 // Defines the horizontal command line split

	// Call-Buddy Title
	titleYStart := 0
	titleYEnd := titleYStart + 2
	if v, err := g.SetView(TTL_LINE_VIEW, 0, titleYStart, verticalSplitX, titleYEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "\u001b[32mTerminal "+"\u001b[29mCall "+"\u001b[29mBuddy")
	}

	// Response Body (e.g. html)
	if v, err := g.SetView(RSP_BODY_VIEW, verticalSplitX+1, titleYStart, realMaxX, horizontalSplitY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Response Body")
	}

	// Method Body (e.g. GET, PUT, HEAD...)
	methodBodyYStart := titleYEnd + 1
	methodBodyYEnd := methodBodyYStart + 10
	if v, err := g.SetView(MTD_BODY_VIEW, 0, methodBodyYStart, verticalSplitX, methodBodyYEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "> https://google.com")
		fmt.Fprintln(v)
		fmt.Fprintln(v, "[ ]"+"GET")
		fmt.Fprintln(v, "[ ]"+"POST")
		fmt.Fprintln(v, "[ ]"+"HEAD")
		fmt.Fprintln(v, "[ ]"+"PUT")
		fmt.Fprintln(v, "[ ]"+"DELETE")
		fmt.Fprintln(v, "[ ]"+"OPTIONS")

	}

	// Request Headers (e.g. Content-type: text/json)
	requestHeadersYStart := methodBodyYEnd + 1
	requestHeadersYEnd := requestHeadersYStart + 6
	if v, err := g.SetView(RQT_HEAD_VIEW, 0, requestHeadersYStart, verticalSplitX, requestHeadersYEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Request Headers")
	}

	// Request Body (e.g. json: {})
	requestBodyYStart := requestHeadersYEnd + 1
	if v, err := g.SetView(RQT_BODY_VIEW, 0, requestBodyYStart, verticalSplitX, horizontalSplitY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Request Body")

	}

	// Command Line (e.g. :get http://httpbin.org/get)
	if v, err := g.SetView(CMD_LINE_VIEW, 0, horizontalSplitY+1, realMaxX, realMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, ":")
	}
	return nil
}

//This is the function to QUIT out of the TUI
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//This function will update the response body (currently) by pressing a variable
func update(g *gocui.Gui, v *gocui.View) error {

	response_view, _ := g.View(RSP_BODY_VIEW)
	response_view.Clear()

	fmt.Fprint(response_view, response_body_string)

	return nil
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

	//Setting up a new TUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Setting a manager, sets the view (defined as another function above)
	g.SetManagerFunc(layout)

	//Setting keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, update); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
