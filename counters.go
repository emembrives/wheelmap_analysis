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

type Sequence []vc

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
        fmt.Println("Rated places read from gob")
        
        tags = make(map[pair]int)
        for _, place := range places {
            for key, value := range place.Tags {
                p := pair{key, value}
                tags[p]++
            }
        }
        tag_counts := make(Sequence, 0)
        for key, value := range tags {
            tag_counts = append(tag_counts, vc{key, value})
        }
        sort.Sort(tag_counts)
        for _, tag_count := range tag_counts {
            fmt.Println(tag_count.value.key, tag_count.value.value, tag_count.count)
        }
}
