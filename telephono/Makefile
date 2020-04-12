.PHONY: all clean

all: test

clean:
	rm -f telephono

telephono: *.go
	go build -o $@ cmd/telephono/main.go

test:
	go test
