package main

import (
    "fmt"
    "sync"
    "math"
    "os"
)

type PlaceList []Location
type Buckets [][]Bucket

const degToRad float64 = 180/math.Pi
const divisions int = 1000

func distance(latitude1, longitude1, latitude2, longitude2 float64) (d float64) {
    R := 6371.0
    dLat := (latitude2-latitude1) * degToRad
    dLon := (longitude2-longitude1) * degToRad
    lat1 := latitude1 * degToRad
    lat2 := latitude2 * degToRad

    a := math.Sin(dLat/2) * math.Sin(dLat/2) + math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    d = R * c
    return
}

type Bucket struct {
    Lat, Lon float64
    FullCount, RatedCount int
    fullCountMutex, ratedCountMutex sync.Mutex
}

func (b Bucket) Ratio() (r float64) {
    if b.FullCount != 0 {
        r = float64(b.RatedCount)/float64(b.FullCount)
    } else {
        r = 0.0
    }
    return
}

func bucketize(list []*SimpleNode, rated []*RatedPlace) (Buckets) {
    var minLat, maxLat, minLon, maxLon float64
    minLat = 99
    minLon = 99
    maxLat = -99
    maxLon = -99
    for _, place := range list {
        minLat = math.Min(place.Latitude, minLat)
        minLon = math.Min(place.Longitude, minLon)
        maxLat = math.Max(place.Latitude, maxLat)
        maxLon = math.Max(place.Longitude, maxLon)
    }
    maxLat = maxLat + 0.001
    maxLon = maxLon + 0.001
    latStep := (maxLat - minLat) / float64(divisions)
    lonStep := (maxLon - minLon) / float64(divisions)
    var buckets Buckets = make([][]Bucket, divisions)
    for i, _ := range buckets {
        buckets[i] = make([]Bucket, divisions)
        for j, _ := range buckets[i] {
            bucket := &buckets[i][j]
            bucket.Lat = minLat + float64(i) * latStep + latStep / 2
            bucket.Lon = minLon + float64(j) * lonStep + lonStep / 2
        }
    }
    for _, place := range list {
        var indexLat int = int((place.Latitude - minLat) / latStep)
        var indexLon int = int((place.Longitude - minLon) / lonStep)
        buckets[indexLat][indexLon].fullCountMutex.Lock()
        buckets[indexLat][indexLon].FullCount++
        buckets[indexLat][indexLon].fullCountMutex.Unlock()
    }
    for _, place := range rated {
        var indexLat int = int((place.Latitude - minLat) / latStep)
        var indexLon int = int((place.Longitude - minLon) / lonStep)
        buckets[indexLat][indexLon].RatedCount++
    }
    return buckets
}

func printBuckets(buckets Buckets, file * os.File) {
    for i, _ := range buckets {
        for j, _ := range buckets[i] {
            b := buckets[i][j]
            s := fmt.Sprintln(i, j, b.Lat, b.Lon, b.FullCount, b.RatedCount, b.Ratio())
            file.WriteString(s)
        }
        file.WriteString("\n")
    }
}
