ROOT=./../../
SUDO=sudo
ACT=act
DIRECTORY=--directory "$(ROOT)"
ENVIRONMENT=--env-file environment
SECRETS=--secret-file secrets
BOILERPLATE=$(DIRECTORY) $(ENVIRONMENT) $(SECRETS) --use-gitignore

# https://github.com/nektos/act

graph:
	$(SUDO) $(ACT) --graph $(DIRECTORY)

list:
	$(SUDO) $(ACT) --list $(DIRECTORY)

on-push:
	$(SUDO) $(ACT) $(BOILERPLATE)

on-pullrequest:
	$(SUDO) $(ACT) pull_request $(BOILERPLATE)

# Eg. "make go-linting", "make go-test"
%:
	$(SUDO) $(ACT) --job "$@" $(BOILERPLATE)
