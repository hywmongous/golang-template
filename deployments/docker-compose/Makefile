STACK=hywmongous

.PHONY:

help:
	@echo 'Targets (docker-compose):'
	@echo '  build                         - Builds images in the docker-compose. yml file'
	@echo '  up                            - Builds, (re)creates, starts, and attaches to containers for a service'
	@echo '  create                        - Builds, (re)creates, starts, and attaches to containers for a service without starting it'
	@echo '  down                          - Stops containers and removes containers, networks, volumes, and images created by up'
	@echo 'Targets (docker stack):'
	@echo '  deploy                        - Deploys the stack $(STACK) as a new stack or updates the existing stack'
	@echo '  rm                            - Removes the $(STACK) stack'

build:
	docker-compose build --force-rm --compress --parallel

create: build
	docker-compose up --force-recreate --build --no-start --always-recreate-deps

up:
	docker-compose up --build --remove-orphans
	@echo ''
	docker ps

down:
	docker-compose down --volumes

deploy:
	docker stack deploy --compose-file docker-compose.yml ${STACK}
	@echo ''
	docker ps

rm:
	docker stack rm ${STACK}
