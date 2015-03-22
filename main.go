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
var cassandraPort string

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

var hostName, _ = os.Hostname()
var portNumber = "8081"

const (
	defaultForeGroundColour = termbox.ColorWhite
	defaultBackGroundColour = termbox.ColorBlack
	messageForeGroundColour = termbox.ColorMagenta
)

func init() {
	// Default to localhost (MX4J needs to be configured to listen to this address in cassandra-env.sh though):
	flag.StringVar(&cassandraHost, "host", hostName, "IP address of the Cassandra host to run against")
	flag.StringVar(&cassandraPort, "port", portNumber, "TCP port of the Cassandra host")
}

// Do all the things:
func main() {

	// Set the vars from the command-line args:
	flag.Parse()

	// Check our connection to MX4J:
	if checkConnection(cassandraHost, cassandraPort) != nil {
		fmt.Printf("Can't connect to stats-provider (%s)! Trying localhost before bailing...\n", cassandraHost)
		if checkConnection("localhost", cassandraPort) != nil {
			fmt.Println("Can't even connect to localhost! Are you running C* with MX4J?")
			os.Exit(2)
		} else {
			fmt.Println("Proceeding with localhost..")
			cassandraHost = "localhost"
		}
	}

	// Initialise "termbox" (console interface):
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Get the initial window-size:
	termWidth, termHeight = termbox.Size()

	// Get the display running in the right mode:
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	// Render the initial "UI":
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawBorder(termWidth, termHeight)
	termbox.Flush()

	// Run the metrics-collector:
	go MetricsCollector(cassandraHost)
	go handleMetrics()
	go refreshScreen()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		// Key pressed:
		case termbox.EventKey:

			// Handle keypresses:
			if ev.Ch == 113 {
				// "q" (quit):
				printfTb(2, 1, messageForeGroundColour, termbox.ColorBlack, "Goodbye!: %s", ev.Ch)
				break loop
			} else if ev.Ch == 0 { // "Space-bar (refresh)"
				showStats()
			} else if ev.Ch == 109 { // "M"
				dataDisplayed = "Metrics"
				showStats()
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
				handleKeypress(&ev)
			}

			// Redraw the display:
			drawBorder(termWidth, termHeight)
			termbox.Flush()

		// Window is re-sized:
		case termbox.EventResize:
			// Remember the new sizes:
			termWidth = ev.Width
			termHeight = ev.Height

			// Redraw the screen:
			drawBorder(termWidth, termHeight)
			termbox.Flush()

		// Error:
		case termbox.EventError:
			panic(ev.Err)

		default:
		}
	}
}
