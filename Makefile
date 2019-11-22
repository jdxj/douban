localName = douban.out
linuxName = douban_linux.out
macName = douban_mac.out

local: clean
	go build -o $(localName) *.go
linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $(linuxName) *.go
	upx --best $(linuxName)
mac: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $(macName) *.go
	upx --best $(macName)
clean:
	rm -rvf *.out *.log
