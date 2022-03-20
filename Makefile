.PHONY : build run start execute debug dlv

build:
	go build -o cmd/main

execute:
	./cmd/main ~/torrent/in/spider.torrent  ~/torrent/out udp

dlv:
	dlv debug ./main.go --  ~/torrent/in/sabaton.torrent ~/torrent/out udp

run: build execute

debug: build dlv
