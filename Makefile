all: build

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags '-s -w -extldflags "-static"' -o bin/fy cli/main.go
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags '-s -w -extldflags "-static"' -o bin/fy.exe cli/main.go

clean:
	rm -rf ./bin
