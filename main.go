package main

import (
	"github.com/hailocab/ctop/types"
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"sync"
	"time"
)

var metricsChannel = make(chan types.CFMetric, 100)
var messageChannel = make(chan types.LogMessage, 100)
var stats = make(map[string]types.CFStats)
var statsMutex sync.Mutex
var dataDisplayed = "Metrics"
var dataSortedBy = "Reads"
var termWidth = 80
var termHeight = 25
var refreshTime = 1 * time.Second
var localHostName, _ = os.Hostname()
var printVersion = flag.Bool("version", false, "Print version number and exit")
var cassandraHost = flag.String("host", localHostName, "IP address of the Cassandra host to run against")
var cassandraMX4jPort = flag.String("port", "8081", "TCP port of MX4J on the Cassandra host")

const (
	defaultForeGroundColour = termbox.ColorWhite
	defaultBackGroundColour = termbox.ColorBlack
	messageForeGroundColour = termbox.ColorMagenta
	releaseVersion          = 1.5
)

func init() {
	// Set the vars from the command-line args:
	flag.Parse()

	// Print the version and quit (if we've been asked to):
	if *printVersion == true {
		fmt.Printf("CTOP version %v\n", releaseVersion)
		os.Exit(0)
	}
}

// Do all the things:
func main() {

	// Check our connection to MX4J:
	if checkConnection(*cassandraHost, *cassandraMX4jPort) != nil {
		fmt.Printf("Can't connect to stats-provider (%s)! Trying localhost before bailing...\n", *cassandraHost)
		if checkConnection("localhost", *cassandraMX4jPort) != nil {
			fmt.Println("Can't even connect to localhost! Check your destination host and port and try again.")
			os.Exit(2)
		} else {
			fmt.Println("Proceeding with localhost..")
			*cassandraHost = "localhost"
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
	go MetricsCollector()
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
