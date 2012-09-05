package main

import (
	"encoding/gob"
        "fmt"
	"os"
        "runtime"
        "sort"
)

type pair struct {
    key, value string
}

type vc struct {
    value pair
    count int
}

type sc struct {
    value string
    count int
}

type Sequence []vc
type StringCount []sc

// Methods required by sort.Interface.
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i].count > s[j].count
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Methods required by sort.Interface.
func (s StringCount) Len() int {
    return len(s)
}
func (s StringCount) Less(i, j int) bool {
    return s[i].count > s[j].count
}
func (s StringCount) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func main() {
	var (
		file *os.File
		err  error
                places  []*RatedPlace
                tags map[pair]int
	)
        runtime.GOMAXPROCS(runtime.NumCPU() - 2)

        if file, err = os.Open("places.gob"); err != nil {
		panic(err.Error())
	}
	decoder := gob.NewDecoder(file)
	decoder.Decode(&places)
        
        tags = make(map[pair]int)
        keys := make(map[string]int)
        for _, place := range places {
            for key, value := range place.Tags {
                p := pair{key, value}
                tags[p]++
                keys[key]++
            }
        }

        tag_counts := make(Sequence, 0)
        for key, value := range tags {
            tag_counts = append(tag_counts, vc{key, value})
        }
        sort.Sort(tag_counts)
        if file, err = os.Create("counters_all.txt"); err != nil {
		panic(err.Error())
	}
        defer file.Close()

        for _, tag_count := range tag_counts {
            s := fmt.Sprintln(tag_count.value.key, tag_count.value.value, tag_count.count)
            file.WriteString(s)
        }
        
        tag_aggregates := make(StringCount, 0)
        for key, value := range keys {
            tag_aggregates = append(tag_aggregates, sc{key, value})
        }
        sort.Sort(tag_aggregates)
        if file, err = os.Create("counters_aggregates.txt"); err != nil {
		panic(err.Error())
	}
        defer file.Close()
        for _, tag_count := range tag_aggregates {
            s := fmt.Sprintln(tag_count.value, tag_count.count)
            file.WriteString(s)
        }
}
