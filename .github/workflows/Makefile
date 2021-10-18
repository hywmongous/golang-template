# https://github.com/nektos/act

ROOT=./../../
DIRECTORY=--directory "$(ROOT)"

ENVIRONMENT=--env-file environment
SECRETS=--secret-file secrets
BOILERPLATE=$(DIRECTORY) $(ENVIRONMENT) $(SECRETS) --use-gitignore

.PHONY:

help:
	@echo 'Targets:'
	@echo '  graph                         - Graphs the workflows and their connections'
	@echo '  list                          - Lists all the workflwos'
	@echo '  on-push                       - Imitates a push and runs the appropriate workflows'
	@echo '  on-pr                         - Imitates a pull--request and runs the appropriate workflows'
	@echo '  %                             - Wildcard which runs an arbitrary job'
	@echo 'Examples:'
	@echo '  "make go-test"                - Runs thte go-test workflow'

graph:
	sudo act --graph $(DIRECTORY)

list:
	sudo act --list $(DIRECTORY)

on-push:
	sudo act $(BOILERPLATE)

on-pr:
	sudo act pull_request $(BOILERPLATE)

%:
	sudo act --job "$@" $(BOILERPLATE)
