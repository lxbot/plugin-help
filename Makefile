.PHONY: build

build:
	go build -buildmode=plugin -o plugin-help.so plugin.go
