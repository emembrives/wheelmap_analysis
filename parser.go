package main

import (
	"encoding/xml"
	//"io"
        "fmt"
        "math"
	"strconv"
)

type Parser struct {
	Decoder *xml.Decoder
	Places  []*RatedPlace
        Nodes []*SimpleNode
}

type Location interface {
    GetLatLon() (lat, lon float64)
}

type RatedPlace struct {
	NodeId    int
	Latitude  float64
	Longitude float64
	Tags      map[string]string
}

type SimpleNode struct {
    Latitude float64
    Longitude float64
}

func NewParser(decoder *xml.Decoder) *Parser {
	parser := new(Parser)
	parser.Decoder = decoder
	parser.Places = make([]*RatedPlace, 0, 10240)
	parser.Nodes = make([]*SimpleNode, 0, 10240)
	return parser
}

func (this *Parser) Parse() {
	var (
		inNode       bool
		currentPlace *RatedPlace
	)

	for token, err := this.Decoder.Token(); err == nil; token, err = this.Decoder.Token() {
		switch token.(type) {
		case xml.StartElement:
			name := token.(xml.StartElement).Name.Local
			if inNode == false && name != "node" {
				if name == "way" || name == "relationship" {
					return
				} else {
					continue
				}
			} else if inNode == false && name == "node" {
				currentPlace = new(RatedPlace)
				currentPlace.Tags = make(map[string]string)
				inNode = true
				for _, attr := range token.(xml.StartElement).Attr {
					if attr.Name.Local == "id" {
						currentPlace.NodeId, err = strconv.Atoi(attr.Value)
					} else if attr.Name.Local == "lat" {
						currentPlace.Latitude, err = strconv.ParseFloat(attr.Value, 64)
					} else if attr.Name.Local == "lon" {
						currentPlace.Longitude, err = strconv.ParseFloat(attr.Value, 64)
					}
                                        if err != nil {
                                            panic(err.Error())
                                        }
				}
			} else if inNode == true && name == "tag" {
				var k, v string
				for _, attr := range token.(xml.StartElement).Attr {
					switch attr.Name.Local {
					case "k":
						k = attr.Value
					case "v":
						v = attr.Value
					}
				}
				currentPlace.Tags[k] = v
			}
		case xml.EndElement:
			name := token.(xml.EndElement).Name.Local
			if inNode == true && name == "node" {
				inNode = false
				if _, ok := currentPlace.Tags["wheelchair"]; ok {
					this.Places = append(this.Places, currentPlace)
                                        fmt.Println(len(this.Places), "rated places parsed.")
				}
                                storedNode := new(SimpleNode)
                                storedNode.Latitude = currentPlace.Latitude
                                storedNode.Longitude = currentPlace.Longitude
                                this.Nodes = append(this.Nodes, storedNode)
                                if i, _ := math.Modf(math.Mod(float64(len(this.Nodes)), 1000)); i == 0 {
                                    fmt.Print(".")
                                }
			}
		default:
			continue
		}
	}
}
