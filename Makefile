.PHONY: build run test clean lint cross-compile

build:
	go build -o jorbites-scripts main.go

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f jorbites-scripts jorbites-scripts-pi32 jorbites-scripts-pi64

cross-compile:
	# Cross-compile for Raspberry Pi 3B (ARMv7 32-bit)
	GOOS=linux GOARCH=arm GOARM=7 go build -o jorbites-scripts-pi32 main.go
	# Cross-compile for Raspberry Pi 3B (ARM64 64-bit)
	GOOS=linux GOARCH=arm64 go build -o jorbites-scripts-pi64 main.go
