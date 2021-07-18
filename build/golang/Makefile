BASEPATH=./../../cmd/rest-server
MAINPATH=$(BASEPATH)/main.go

BUILDPATH=.
DISTPATH=./../../dist/golang
BINPATH=bin
BINMAINPATH=$(BINPATH)/main

BFLAGS=-o

.PHONY:
help:
	@echo  'Cleaning targets:'
	@echo  '  build_clean                   - Removes golang "./bin" files in the "build" folder'
	@echo  '  dist_clean                    - Removes golang "./bin" files in the "dist" folder'
	@echo  'Configuration targets:'
	@echo  '  build                         - Builds the rest-server into binary and places it inside "./bin"'
	@echo  '  run                           - Runs the rest-server'
	@echo  '  compile                       - Compiles the source into the dist directory'

.PHONY:
build:
	go build $(BFLAGS) $(BUILDPATH)/$(BINMAINPATH) $(MAINPATH)

.PHONY:
compile:
	@echo 'Dist compilation not implemented'

.PHONY:
run:
	go run $(MAINPATH)

.PHONY:
build_clean:
	rm -rf $(BUILDPATH)/$(BINPATH)

.PHONY:
dist_clean:
	rm -rf $(DISTPATH)/$(BINPATH)
