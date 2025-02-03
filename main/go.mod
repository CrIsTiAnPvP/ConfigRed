module configred/main

go 1.23.4

replace configred/interfacesv4 => ../interfacesv4

require configred/interfacesv4 v0.5.0

replace configred/rainbow => ../rainbow

require configred/rainbow v1.0.0

replace configred/utils => ../utils

require configred/utils v0.0.1

require (
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/schollz/progressbar/v3 v3.17.1 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
)
