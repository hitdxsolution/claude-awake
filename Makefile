# Build the three release binaries. Go cross-compiles from any host, so `make` on a Mac
# produces the Windows .exe too. Flags: -s -w strips debug info; -H windowsgui keeps Windows
# from opening a console window for a background daemon.
LDFLAGS := -s -w

.PHONY: all clean check mac-arm64 mac-x64 windows-x64

all: check mac-arm64 mac-x64 windows-x64

# Formatting + the built-in static analyser. Run before building.
check:
	gofmt -l .
	go vet ./...

mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o dist/claude-awake-macos-arm64 .

mac-x64:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o dist/claude-awake-macos-x64 .

windows-x64:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS) -H windowsgui" -o dist/claude-awake-windows-x64.exe .

clean:
	rm -rf dist
