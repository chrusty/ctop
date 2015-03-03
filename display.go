package main

import (
	// "flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
	// "os"
)

// Reads log-messages out of the logMessage chan and displays them to screen:
func show_logs() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	printf_tb(2, 1, messageForeGroundColour, termbox.ColorBlack, "         |")
	printf_tb(2, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "Severity")
	printf_tb(13, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "Message")

	for y := 2; y < termHeight; y++ {
		select {
		// attempt to receive from channel:
		case logMessage := <-messageChannel:
			printf_tb(2, y, messageForeGroundColour, termbox.ColorBlack, "%s", logMessage.Severity)
			printf_tb(13, y, messageForeGroundColour, termbox.ColorBlack, "%s", logMessage.Message)
		default:
			printf_tb(2, y, messageForeGroundColour, termbox.ColorBlack, "No more logs")
			return
		}
	}
}

// Draws stats on the screen:
func show_stats() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Positions:                                                                  2                   22                  42        52         63             79              94
	printf_tb(2, 1, messageForeGroundColour|termbox.AttrBold, termbox.ColorBlack, "KeySpace            ColumnFamily        Reads/s   Writes/s   LiveSpace(B)   R-Latency(ms)   W-Latency(ms)")
	printf_tb(20, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printf_tb(40, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printf_tb(50, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printf_tb(61, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printf_tb(76, 1, messageForeGroundColour, termbox.ColorBlack, "|")
	printf_tb(92, 1, messageForeGroundColour, termbox.ColorBlack, "|")

	y := 2

	// Get a lock on stats, then make a sorted map of the stats:
	statsMutex.Lock()
	sortedStats := sortedKeys(stats)
	statsMutex.Unlock()

	for _, cfStatsKey := range sortedStats {
		if y < termHeight {
			// printf_tb(2, y, messageForeGroundColour, termbox.ColorBlack, "(%s:%s) r:%d, w:%d", cfStats.KeySpace, cfStats.ColumnFamily, cfStats.ReadCount, cfStats.WriteCount)
			printf_tb(2, y, messageForeGroundColour, termbox.ColorBlack, "%s", stats[cfStatsKey].KeySpace)
			printf_tb(20, y, messageForeGroundColour, termbox.ColorBlack, "  %s", stats[cfStatsKey].ColumnFamily)
			printf_tb(40, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].ReadRate)
			printf_tb(50, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].WriteRate)
			printf_tb(61, y, messageForeGroundColour, termbox.ColorBlack, "  %d", stats[cfStatsKey].LiveDiskSpaceUsed)
			printf_tb(76, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].ReadLatency)
			printf_tb(92, y, messageForeGroundColour, termbox.ColorBlack, "  %f", stats[cfStatsKey].WriteLatency)
			y++
		}
	}
}

// Refreshes the on-screen data:
func refresh_screen() {
	for {

		if dataDisplayed == "Metrics" {
			show_stats()
		}

		if dataDisplayed == "Logs" {
			show_logs()
		}

		// Sleep:
		time.Sleep(refreshTime)
	}
}

// Print function for TermBox:
func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// PrintF function for TermBox:
func printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	print_tb(x, y, fg, bg, s)
}

// Draw the border around the edge of the screen:
func draw_border(width int, height int) {
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
	print_tb(1, 0, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, " C-top ")
	print_tb(8, 0, termbox.ColorBlue, defaultBackGroundColour, "(top for Cassandra) ")

	// Menu:
	// Positions:                                                       2            15           28            42              58                76                94         105     113
	print_tb(1, height-1, termbox.ColorBlue, defaultBackGroundColour, " Organise by (1)Reads/s / (2)Writes/s / (3)Space-used / (4)Read-latency / (5)Write-latency, (M)etrics, (L)ogs, (Q)uit ")
	print_tb(15, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "1")
	print_tb(28, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "2")
	print_tb(42, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "3")
	print_tb(58, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "4")
	print_tb(76, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "5")
	print_tb(94, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "M")
	print_tb(105, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "L")
	print_tb(113, height-1, termbox.ColorBlue|termbox.AttrBold, defaultBackGroundColour, "Q")

	// Highlight the sorting mode:
	if dataSortedBy == "Reads" {
		print_tb(15, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "1")
	}
	if dataSortedBy == "Writes" {
		print_tb(28, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "2")
	}
	if dataSortedBy == "Space" {
		print_tb(42, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "3")
	}
	if dataSortedBy == "ReadLatency" {
		print_tb(58, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "4")
	}
	if dataSortedBy == "WriteLatency" {
		print_tb(76, height-1, termbox.ColorWhite|termbox.AttrBold, defaultBackGroundColour, "5")
	}

	// Show what mode we're in:
	if dataDisplayed == "Metrics" {
		printf_tb(termWidth-10, 0, termbox.ColorBlue|termbox.AttrBold, termbox.ColorBlack, " Metrics ")
	}
	if dataDisplayed == "Logs" {
		printf_tb(termWidth-7, 0, termbox.ColorBlue|termbox.AttrBold, termbox.ColorBlack, " Logs ")
	}
}
