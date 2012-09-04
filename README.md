wheelmap_analysis
=================

Description
-----------
Go program to analyze Wheelmap.org data, and associated Gnuplot script.

How to use
----------
* Download OpenStreetMap XML data in the same directory:

> wget http://download.geofabrik.de/osm/europe/france.osm.bz2

* Compile the GO first program

> go build osm.go parser.go analysis.go

* Run the program

> ./osm

* Use Gnuplot to draw nice pictures

> gnuplot < plot.plt

* Compile the second go program

> go build counters.go parser.go

* Run it and see statistics on standard output
> ./counters|less
