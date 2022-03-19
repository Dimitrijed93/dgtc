.PHONY : build run start execute debug dlv

build:
	go build -o cmd/main

execute:
	./cmd/main ~/torrent/in/debian.iso.torrent  ~/torrent/out

dlv:
	dlv debug ./main.go --  ~/torrent/in/debian.iso.torrent ~/torrent/out

run: build execute

debug: build dlv
