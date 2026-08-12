# Build the three release binaries. Go cross-compiles from any host, so `make` on a Mac
# produces the Windows .exe too. Flags: -s -w strips debug info; -H windowsgui keeps Windows
# from opening a console window for a background daemon.
LDFLAGS := -s -w

GOBIN := $(shell go env GOPATH)/bin
GOVERSIONINFO := $(GOBIN)/goversioninfo

# 윈도우 리소스(버전 정보·아이콘·매니페스트). 파일명이 _windows_amd64 로 끝나므로 Go 가
# 윈도우 amd64 빌드에서만 링크한다 — 맥 빌드에는 아무 영향이 없다.
WINRES := resource_windows_amd64.syso

.PHONY: all clean check mac-arm64 mac-x64 windows-x64 winres

all: check mac-arm64 mac-x64 windows-x64

# Formatting + the built-in static analyser. Run before building.
#   windows 쪽 코드는 빌드 태그로 갈려 있어 GOOS 를 바꿔 한 번 더 봐야 검사가 된다.
check:
	gofmt -l .
	go vet ./...
	GOOS=windows go vet ./...

$(GOVERSIONINFO):
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.7.0

winres: $(WINRES)

$(WINRES): versioninfo.json assets/icon.ico assets/claude-awake.manifest | $(GOVERSIONINFO)
	$(GOVERSIONINFO) -64 -o $@ versioninfo.json

mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o dist/claude-awake-macos-arm64 .

mac-x64:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o dist/claude-awake-macos-x64 .

windows-x64: $(WINRES)
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS) -H windowsgui" -o dist/claude-awake-windows-x64.exe .

# 아이콘 원본은 SVG 다. ImageMagick 이 있어야 하고, 아이콘을 고칠 때만 돌리면 된다.
#   256 크기는 넣지 않는다 — ICO 안에서 비압축이라 그것 하나로 실행파일이 260 KB 더 뚱뚱해진다.
assets/icon.ico: assets/icon.svg
	magick -background none $< -define icon:auto-resize=128,64,48,32,16 $@

clean:
	rm -rf dist $(WINRES)
