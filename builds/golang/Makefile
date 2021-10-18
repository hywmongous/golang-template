ROOTPATH=./../..
CMDNAME=service
BASEPATH=$(ROOTPATH)/cmd/$(CMDNAME)
MAINPATH=$(BASEPATH)/main.go

BUILDPATH=.
DISTPATH=$(ROOTPATH)/dist/golang/main

.PHONY:

help:
	@echo 'Targets:'
	@echo '  build                         - Builds the rest-server into binary and places it inside "./bin"'
	@echo '  run                           - Runs the rest-server'

build:
	go build $(BFLAGS) -o $(DISTPATH) $(MAINPATH)

run:
	go run $(MAINPATH)
