localName = douban.out

local: clean
	go build -o localName *.go
clean:
	rm -rvf *.out *.log
