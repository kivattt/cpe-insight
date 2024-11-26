#!/usr/bin/env bash
bin=./bin

if [ ! -d $bin ]; then
	mkdir $bin
fi

ldflags="-s -w"

# amd64
echo "linux   amd64"; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-linux-amd64
echo "freebsd amd64"; GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-freebsd-amd64
echo "windows amd64"; GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-windows-amd64.exe

# i386
echo "linux   386"; GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-linux-i386
echo "freebsd 386"; GOOS=freebsd GOARCH=386 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-freebsd-i386
echo "windows 386"; GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-windows-i386.exe

# arm64
echo "linux   arm64"; GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-linux-arm64
echo "freebsd arm64"; GOOS=freebsd GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-freebsd-arm64
echo "windows arm64"; GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-windows-arm64.exe

# arm
echo "linux   arm"; GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-linux-arm
echo "freebsd arm"; GOOS=freebsd GOARCH=arm CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-freebsd-arm
echo "windows arm"; GOOS=windows GOARCH=arm CGO_ENABLED=0 go build -ldflags="$ldflags" -o $bin/cpe-insight-windows-arm.exe
