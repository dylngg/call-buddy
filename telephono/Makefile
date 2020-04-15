.PHONY: all clean telephono call-buddy

all: call-buddy telephono

clean:
	rm -f telephono call-buddy

telephono: *.go
	go build -o $@ cmd/telephono/main.go

call-buddy: *.go
	go build -o $@ cmd/call-buddy/main.go

test:
	go test -v
