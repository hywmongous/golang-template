ROOTDIRECTORY=./../../
DOCKERFILEPATH=./Dockerfile

HOSTPORT=5000
CONTAINERPORT=5000
PORTMAPPING=$(HOSTPORT):$(CONTAINERPORT)

TAG=ghcr.io/hywmongous/golang-template:latest
DOCKER=docker

.PHONY:

help:
	@echo  'Targets:'
	@echo  '  build                         - Builds image with tag $(TAG)'
	@echo  '  create                        - Builds the image and creates the container'
	@echo  '  run                           - Builds the image and runs it'
	@echo  '  runit                         - Builds if necessary and runs the images with an interactive shell'
	@echo  '  rm                            - Removes the image from the local docker registry'
	@echo  '  ps                            - Lists all docker containers'
	@echo  '  networks                      - Lists all docker networks'
	@echo  '  volumes                       - Lists all docker volumes'
	@echo  '  kill-all                      - Kill all containers'
	@echo 'Examples:'
	@echo '  "make run"                     - Runs the docker image in a container with ports $(PORTMAPPING)'

build:
	docker build -t $(TAG) -f $(DOCKERFILEPATH) $(ROOTDIRECTORY) --compress

create: build
	docker create $(TAG)

run: create
	docker run -p $(PORTMAPPING) $(TAG)

runit: create
	docker run -p $(PORTMAPPING) -it $(TAG)

rm:
	docker rm -f $(TAG) 2> /dev/null

ps:
	sudo $(DOCKER) ps

networks:
	sudo $(DOCKER) network ls

volumes:
	sudo $(DOCKER) volume ls

kill-all:
	sudo $(DOCKER) kill $$(sudo $(DOCKER) ps -q)
