package main

import (
	types "./types"
	"sort"
)

// Define a type for a sorted-map of columnfamily-stats:
type sortedMap struct {
	m map[string]types.CFStats
	s []string
}

// Return the length of a sorted-map:
func (sm *sortedMap) Len() int {
	return len(sm.m)
}

// Handles the different attributes we might sort by:
func (sm *sortedMap) Less(i, j int) bool {
	if dataSortedBy == "Reads" {
		return sm.m[sm.s[i]].ReadRate > sm.m[sm.s[j]].ReadRate
	}
	if dataSortedBy == "Writes" {
		return sm.m[sm.s[i]].WriteRate > sm.m[sm.s[j]].WriteRate
	}
	if dataSortedBy == "Space" {
		return sm.m[sm.s[i]].LiveDiskSpaceUsed > sm.m[sm.s[j]].LiveDiskSpaceUsed
	}
	if dataSortedBy == "ReadLatency" {
		return sm.m[sm.s[i]].ReadLatency > sm.m[sm.s[j]].ReadLatency
	}
	if dataSortedBy == "WriteLatency" {
		return sm.m[sm.s[i]].WriteLatency > sm.m[sm.s[j]].WriteLatency
	}
	// Default to "Reads":
	return sm.m[sm.s[i]].ReadRate > sm.m[sm.s[j]].ReadRate
}

// Replace two values in a list:
func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

// Return keys in order:
func sortedKeys(m map[string]types.CFStats) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
