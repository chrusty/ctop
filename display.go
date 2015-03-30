package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

// Reads log-messages out of the logMessage chan and displays them to screen:
func showLogs() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printfTb(2, 1, messageForeGroundColour, termbox.ColorBlack, "         |")
	printfTb(2, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "Severity")
	printfTb(13, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "Message")

	for y := 2; y < termHeight; y++ {
		select {
		// attempt to receive from channel:
		case logMessage := <-messageChannel:
			printfTb(2, y, messageForeGroundColour, termbox.ColorBlack, "%s", logMessage.Severity)
			printfTb(13, y, messageForeGroundColour, termbox.ColorBlack, "%s", logMessage.Message)
		default:
			printfTb(2, y, messageForeGroundColour, termbox.ColorBlack, "No more logs")
			return
		}
	}
}

// Draws stats on the screen:
func showStats() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Positions:                                                                  2                   22                  42        52         63             79              94
	printfTb(2, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "KeySpace            ColumnFamily        Reads/s   Writes/s   LiveSpace(B)   R-Latency(ms)   W-Latency(ms)")
	printfTb(20, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printfTb(40, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printfTb(50, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printfTb(61, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printfTb(76, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printfTb(92, 1, messageForeGroundColour, termbox.ColorBlack, "|")

	y := 2

	// Get a lock on stats, then make a sorted map of the stats:
	statsMutex.Lock()
	sortedStats := sortedKeys(stats)
	statsMutex.Unlock()

	for _, cfStatsKey := range sortedStats {
		if y < termHeight {
			// printfTb(2, y, messageForeGroundColour, termbox.ColorBlack, "(%s:%s) r:%d, w:%d", cfStats.KeySpace, cfStats.ColumnFamily, cfStats.ReadCount, cfStats.WriteCount)
			printfTb(2, y, messageForeGroundColour, termbox.ColorBlack, "%s", stats[cfStatsKey].KeySpace)
			printfTb(20, y, messageForeGroundColour, termbox.ColorBlack, "  %s", stats[cfStatsKey].ColumnFamily)
			printfTb(40, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].ReadRate)
			printfTb(50, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].WriteRate)
			printfTb(61, y, messageForeGroundColour, termbox.ColorBlack, "  %d", stats[cfStatsKey].LiveDiskSpaceUsed)
			printfTb(76, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].ReadLatency)
			printfTb(92, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].WriteLatency)
			y++
		}
	}
}

// Refreshes the on-screen data:
func refreshScreen() {
	for {

		if dataDisplayed == "Metrics" {
			showStats()
		}

		if dataDisplayed == "Logs" {
			showLogs()
		}

		// Sleep:
		time.Sleep(refreshTime)
	}
}

// Print function for TermBox:
func printTb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// PrintF function for TermBox:
func printfTb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	printTb(x, y, fg, bg, s)
}

// Draw the border around the edge of the screen:
func drawBorder(width int, height int) {
	// Sides:
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '-', defaultForeGroundColour, defaultBackGroundColour)
		termbox.SetCell(x, height-1, '-', defaultForeGroundColour, defaultBackGroundColour)
	}

	// Top and bottom:
	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, '|', defaultForeGroundColour, defaultBackGroundColour)
		termbox.SetCell(width-1, y, '|', defaultForeGroundColour, defaultBackGroundColour)
	}

	// Corners:
	termbox.SetCell(0, 0, '+', defaultForeGroundColour, defaultBackGroundColour)
	termbox.SetCell(width-1, 0, '+', defaultForeGroundColour, defaultBackGroundColour)
	termbox.SetCell(0, height-1, '+', defaultForeGroundColour, defaultBackGroundColour)
	termbox.SetCell(width-1, height-1, '+', defaultForeGroundColour, defaultBackGroundColour)

	// Title:
	printTb(1, 0, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, " C-top ")
	printTb(8, 0, termbox.ColorBlue, defaultBackGroundColour, "(top for Cassandra) connected to ")
	printTb(41, 0, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, *cassandraHost)

	// Menu:
	// Positions:                                                      2            15           28            42              58                76                94         105     113
	printTb(1, height-1, termbox.ColorBlue, defaultBackGroundColour, " Organise by (1)Reads/s / (2)Writes/s / (3)Space-used / (4)Read-latency / (5)Write-latency, (M)etrics, (L)ogs, (Q)uit ")
	printTb(15, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "1")
	printTb(28, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "2")
	printTb(42, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "3")
	printTb(58, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "4")
	printTb(76, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "5")
	printTb(94, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "M")
	printTb(105, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "L")
	printTb(113, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "Q")

	// Highlight the sorting mode:
	if dataSortedBy == "Reads" {
		printTb(15, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "1")
	}
	if dataSortedBy == "Writes" {
		printTb(28, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "2")
	}
	if dataSortedBy == "Space" {
		printTb(42, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "3")
	}
	if dataSortedBy == "ReadLatency" {
		printTb(58, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "4")
	}
	if dataSortedBy == "WriteLatency" {
		printTb(76, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "5")
	}

	// Show what mode we're in:
	if dataDisplayed == "Metrics" {
		printfTb(termWidth-10, 0, termbox.ColorBlue|termbox.AttrBold, termbox.ColorBlack, " Metrics ")
	}
	if dataDisplayed == "Logs" {
		printfTb(termWidth-7, 0, termbox.ColorBlue|termbox.AttrBold, termbox.ColorBlack, " Logs ")
	}
}
