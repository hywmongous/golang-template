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
	sudo docker-compose build --force-rm --compress --parallel

create: build
	sudo docker-compose up --force-recreate --build --no-start --always-recreate-deps

up:
	sudo docker-compose up --build --remove-orphans
	@echo ''
	sudo docker ps

down:
	sudo docker-compose down --volumes

deploy:
	sudo docker stack deploy --compose-file docker-compose.yml ${STACK}
	@echo ''
	sudo docker ps

rm:
	sudo docker stack rm ${STACK}
