.PHONY : build run start execute debug dlv

build:
	go build -o cmd/main

execute:
	./cmd/main ~/torrent/in/debian.iso.torrent  ~/torrent/out http

dlv:
	dlv debug ./main.go --  ~/torrent/in/debian.iso.torrent ~/torrent/out http

run: build execute

debug: build dlv
