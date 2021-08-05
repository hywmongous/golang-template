DOCKERCOMPOSE=docker-compose

.PHONY:
help:
	@echo 'Cleaning targets:'
	@echo '  build_clean                   - Removes the image from the local docker library'
	@echo '  dist_clean                    - Removes the image from the local docker library'
	@echo 'Configuration targets:'
	@echo '  build                         - builds images in the docker-compose. yml file'
	@echo '  create                        - Builds, (re)creates, starts, and attaches to containers for a service without starting it'
	@echo '  up                            - Builds, (re)creates, starts, and attaches to containers for a service'
	@echo '  down                          - Stops containers and removes containers, networks, volumes, and images created by up'

.PHONY:
build:
	sudo $(DOCKERCOMPOSE) build --force-rm --compress --parallel

.PHONY:
create: build
	sudo $(DOCKERCOMPOSE) up --force-recreate --build --no-start --always-recreate-deps

.PHONY:
up:
	sudo $(DOCKERCOMPOSE) up --build --remove-orphans

.PHONY:
down:
	sudo $(DOCKERCOMPOSE) down --volumes

.PHONY:
dist_clean: build_clean
.PHONY:
build_clean:
	sudo $(DOCKERCOMPOSE) rm -s -f
