.PHONY: build debug

build:
	go build -buildmode=plugin -o plugin-help.so plugin.go

debug:
	go build -gcflags="all=-N -l" -buildmode=plugin -o plugin-help.so plugin.go
