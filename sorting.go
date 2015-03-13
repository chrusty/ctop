package main

import (
	types "./types"
	"sort"
)

type sortedMap struct {
	m map[string]types.CFStats
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

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
	} else {
		// Default to "Reads":
		return sm.m[sm.s[i]].ReadRate > sm.m[sm.s[j]].ReadRate
	}
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]types.CFStats) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
