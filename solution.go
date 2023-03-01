package main

import (
	"fmt"
	"sort"
)

type Interval struct {
	Start int
	End   int
}

func merge(intervals []Interval) []Interval {
	// Step 1: Sort the intervals by start time
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	// Step 2: Initialize a stack to store the merged intervals
	stack := []Interval{}

	// Step 3: Loop through the sorted intervals and merge
	for _, interval := range intervals {
		if len(stack) == 0 || interval.Start > stack[len(stack)-1].End {
			// If the stack is empty or the current interval does not overlap with the top interval on the stack,
			// push the current interval onto the stack.
			stack = append(stack, interval)
		} else {
			// Otherwise, merge the current interval with the top interval on the stack and update the top interval.
			top := &stack[len(stack)-1]
			if interval.End > top.End {
				top.End = interval.End
			}
		}
	}

	// Step 4: Return the merged intervals from the stack
	return stack
}

func main() {
	intervals := []Interval{
		{1, 3},
		{2, 6},
		{8, 10},
		{15, 18},
	}
	mergedIntervals := merge(intervals)
	fmt.Println(mergedIntervals)
}
