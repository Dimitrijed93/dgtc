.PHONY : clean build run start execute debug dlv

clean:
	rm -r cmd/main

build:
	go build -o cmd/main

execute:
	./cmd/main ~/torrent/in/spider.torrent  ~/torrent/out

dlv:
	dlv debug ./main.go --  ~/torrent/in/debian.iso.torrent ~/torrent/out

run: clean build execute

debug: clean build dlv
