IMAGENAME=example-service:latest
DOCKER=docker
ROOTDIRECTORY=./../../
DOCKERFILEPATH=./Dockerfile

.PHONY:
help:
	@echo  'Cleaning targets:'
	@echo  '  build_clean                   - Removes the image from the local docker library'
	@echo  '  dist_clean                    - Removes the image from the local docker library'
	@echo  'Configuration targets:'
	@echo  '  build                         - Builds the rest-server as an image'
	@echo  '  create                        - Builds if nessecary and creates a container with the image'
	@echo  '  run                           - Builds if nessecary and runs the images'
	@echo  '  runit                         - Builds if nessecary and runs the images with an interactive shell'
	@echo  '  push                          - Pushes the image to Dockerhub'
	@echo  '  pull                          - Pulls the image from Dockerhub'

.PHONY:
build:
	sudo $(DOCKER) build -t $(IMAGENAME) -f $(DOCKERFILEPATH) $(ROOTDIRECTORY) --compress --rm

.PHONY:
create: build
	sudo $(DOCKER) create $(IMAGENAME)

.PHONY:
run: create
	sudo $(DOCKER) run -p 8080:8080 $(IMAGENAME) --rm

.PHONY:
runit: create
	sudo $(DOCKER) run -p 8080:8080 -it $(IMAGENAME)

.PHONY:
push: build
	@echo 'Not implemented yet'

.PHONY:
pull:
	@echo 'Not implemented yet'

.PHONY:
dist_clean: build_clean
.PHONY:
build_clean:
	sudo $(DOCKER) rm -f $(IMAGENAME) 2> /dev/null
