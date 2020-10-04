.PHONY: build
build:
	rm -rf _build && mkdir _build && go build -o _build/img_generator -v ./cmd

.PHONY: build
run:
	go run ./cmd

.PHONY: clear
clear:
	rm -rf _build


.DEFAULT_GOAL := build
