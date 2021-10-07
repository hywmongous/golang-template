DOCKERCOMPOSE=docker-compose
DOCKER=docker
DOCKERSTACK=$(DOCKER) stack
STACK=hywmongous

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
	@echo ''
	sudo $(DOCKER) ps

.PHONY:
down:
	sudo $(DOCKERCOMPOSE) down --volumes

.PHONY:
deploy:
	sudo $(DOCKERSTACK) deploy --compose-file docker-compose.yml ${STACK}
	@echo ''
	sudo $(DOCKER) ps

.PHONY:
rm:
	sudo $(DOCKERSTACK) rm ${STACK}

.PHONY:
ps:
	sudo $(DOCKER) ps

.PHONY:
networks:
	sudo $(DOCKER) network ls

.PHONY:
volumes:
	sudo $(DOCKER) volume ls

.PHONY:
kill-all:
	sudo $(DOCKER) kill $$(sudo $(DOCKER) ps -q)

.PHONY:
dist_clean: build_clean
.PHONY:
build_clean:
	sudo $(DOCKERCOMPOSE) rm -s -f
