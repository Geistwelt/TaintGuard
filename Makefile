#!/bin/bash

.PHONY: install
install:
	@go build -o tguard main.go && mv tguard ${GOPATH}/bin/