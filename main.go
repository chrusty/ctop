package main

import (
	types "./types"
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"sync"
	"time"
)

var cassandraHost string

// var messageChannel = make(chan types.LogMessage, 10)
var metricsChannel = make(chan types.CFMetric, 100)
var messageChannel = make(chan types.LogMessage, 100)
var stats = make(map[string]types.CFStats)
var statsMutex sync.Mutex
var dataDisplayed = "Metrics"
var dataSortedBy = "Reads"
var termWidth = 80
var termHeight = 25
var refreshTime = 1 * time.Second

const (
	defaultForeGroundColour = termbox.ColorWhite
	defaultBackGroundColour = termbox.ColorBlack
	messageForeGroundColour = termbox.ColorMagenta
)

func init() {
	// Default to localhost (MX4J needs to be configured to listen to this address in cassandra-env.sh though):
	// hostName, _ := os.Hostname()
	hostName := "127.0.0.1"
	flag.StringVar(&cassandraHost, "cassandraHost", hostName, "The address of the Cassandra host to run against")
}

// Do all the things:
func main() {

	// Set the vars from the command-line args:
	flag.Parse()

	err := check_connection(cassandraHost)
	if err != nil {
		fmt.Printf("Can't connect to stats-provider! (%s)\n", cassandraHost)
		os.Exit(2)
	}

	// Initialise "termbox" (console interface):
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termWidth, termHeight = termbox.Size()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	draw_border(termWidth, termHeight)
	termbox.Flush()

	// Run the metrics-collector:
	go MetricsCollector(cassandraHost)
	go handle_metrics()
	go refresh_screen()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		// Key pressed:
		case termbox.EventKey:

			// Handle keypresses:
			if ev.Ch == 113 {
				// "q" (quit):
				printf_tb(2, 1, messageForeGroundColour, termbox.ColorBlack, "Goodbye!: %s", ev.Ch)
				break loop
			} else if ev.Ch == 0 { // "Space"
				show_stats()
			} else if ev.Ch == 109 { // "M"
				dataDisplayed = "Metrics"
				show_stats()
			} else if ev.Ch == 108 { // "L"
				dataDisplayed = "Logs"
			} else if ev.Ch == 49 { // "1"
				dataSortedBy = "Reads"
			} else if ev.Ch == 50 { // "2"
				dataSortedBy = "Writes"
			} else if ev.Ch == 51 { // "3"
				dataSortedBy = "Space"
			} else if ev.Ch == 52 { // "4"
				dataSortedBy = "ReadLatency"
			} else if ev.Ch == 53 { // "5"
				dataSortedBy = "WriteLatency"
			} else {
				// Anything else:
				handle_keypress(&ev)
			}

			draw_border(termWidth, termHeight)
			termbox.Flush()

		// Window is re-sized:
		case termbox.EventResize:
			// Remember the new sizes:
			termWidth = ev.Width
			termHeight = ev.Height

			// Redraw the screen:
			draw_border(termWidth, termHeight)
			termbox.Flush()

		// Error:
		case termbox.EventError:
			panic(ev.Err)

		default:
		}
	}
}
