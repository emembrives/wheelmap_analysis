package main

import (
	"compress/bzip2"
	"encoding/gob"
	"encoding/xml"
        "fmt"
	"os"
        "runtime"
)

func main() {
	var (
		file *os.File
		err  error
	)
        runtime.GOMAXPROCS(runtime.NumCPU() - 2)
	if file, err = os.Open("france.osm.bz2"); err != nil {
		panic(err.Error())
	}
	decompressor := bzip2.NewReader(file)
	decoder := xml.NewDecoder(decompressor)
	parser := NewParser(decoder)
	parser.Parse()
        fmt.Println("XML file parsed")

        if file, err = os.Create("places.gob"); err != nil {
		panic(err.Error())
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(parser.Places)
        fmt.Println("Rated places written to gob")

	buckets := bucketize(parser.Nodes, parser.Places)
        fmt.Println("Buckets computed")

        if file, err = os.Create("buckets.gob"); err != nil {
		panic(err.Error())
	}
	encoder = gob.NewEncoder(file)
	encoder.Encode(buckets)
        fmt.Println("Buckets written to gob")
	
        if file, err = os.Create("analysis.ssv"); err != nil {
		panic(err.Error())
	}
	printBuckets(buckets, file)
        fmt.Println("Buckets to SSV file")
}
