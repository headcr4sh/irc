package main

import (
	"flag"
	"fmt"
	"github.com/headcr4sh/irc"
	"github.com/jroimartin/gocui"
	"os"
)

var showHelp = false
var debug = false

var gui *gocui.Gui

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show this help message.")
	flag.BoolVar(&debug, "debug", false, "Enable debug log output.")
	flag.Usage = func() {
		fmt.Println("irc is an IRC command-line client written in Go.")
		fmt.Println("Usage: irc [OPTIONS] <URI>")
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()
	var nArg = flag.NArg()
	var args = flag.Args()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if nArg != 1 {
		fmt.Printf("Exactly one non-flag argument (the URI) is required. %d have been given instead.\n", nArg)
		os.Exit(2)
	}

	_, err := irc.NewURL(args[0])
	if err != nil {
		fmt.Printf("Invalid URL: %s\n", args[0])
		os.Exit(3)
	}
	//hostname := url.Hostname()
	//port := url.Port()
	//conn := irc.NewClientConnection(hostname, port)
	//if conn.Open() != nil {
	//	fmt.Printf("ClientConnection to %s failed.\n", url)
	//}

	if gui, err = gocui.NewGui(gocui.Output256); err != nil {
		fmt.Printf("Unable to open GUI: %v\n", err)
		os.Exit(4)
	}
	defer gui.Close()

	gui.SetManagerFunc(layout)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		fmt.Printf("Unable to attach key listener: %v\n", err)
		os.Exit(5)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("ERROR! %v\n", err)
		os.Exit(6)
	}
}

var channelInfo, status, input *gocui.View

func layout(g *gocui.Gui) error {
	var err error
	minX, minY := -1, -1
	maxX, maxY := g.Size()

	if channelInfo, err = g.SetView("channel_info", minX, minY, maxX, minY+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	channelInfo.Frame = false
	channelInfo.BgColor = gocui.ColorBlue
	channelInfo.FgColor = gocui.ColorWhite
	channelInfo.Editable = false
	fmt.Fprintln(channelInfo, "#some_chatroom | here goes the message of the day...")

	if status, err = g.SetView("status", minX, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	status.Frame = false
	status.BgColor = gocui.ColorBlue
	status.FgColor = gocui.ColorWhite
	status.Editable = false
	fmt.Fprintln(status, "[sw] ~username | irc.example.com")

	if input, err = g.SetView("input", -1, maxY-3, maxX, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	input.Editable = true
	input.Frame = false
	input.BgColor = gocui.ColorWhite
	input.FgColor = gocui.ColorBlack

	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	g.SetCurrentView("input")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
