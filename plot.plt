set pm3d map
set pm3d ftriangles interpolate 10,1
set terminal png enhanced size 1024,1024 font 'Verdana,16'
set output 'osm_coverage.png'
set xrange[-6:10]
set yrange[41:51.5]
unset xtics
unset ytics
unset cbtics
set title "Couverture de la France par OpenStreetMap"
set cblabel "Densité de noeuds répertoriés"
splot "analysis.ssv" using 4:3:5

set output "wheelmap_raw.png"
set title "Lieux répertoriés sur Wheelmap.org"
set cblabel "Densité de lieux répertoriés"
splot "analysis.ssv" using 4:3:($6 == 0 ? 0 : log10($6))

set output "wheelmap_ratio.png"
set title "Couverture (ratio) de Wheelmap.org"
set cblabel "Densité relative de lieux répertoriés"
splot "analysis.ssv" using 4:3:($7 == 0 ? -5.5 : log10($7))
