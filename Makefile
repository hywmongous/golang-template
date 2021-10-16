BUILDPATH=./build
BUILDDIRS=golang docker

DEPLOYMENTPATH=./deployment
DEPLOYMENTDIRS=docker-compose

SPACE :=
SPACE +=

help:
	@echo 'Cleaning targets:'
	@echo '  build_clean                   - Invokes all build_clean in supported build makefiles'
	@echo '  dist_clean                    - Invokes all dist_clean in supported build makefiles'
	@echo 'Examples:'
	@echo '  Build golang binary: "make golang_build'
	@echo '  Run docker: "make docker_run'
	@echo '  Deploy docker-compose: "make docker-compose-postgres_up'
	@echo ''
	@echo '-- make help-build --'
	$(MAKE) help-build
	@echo ''
	@echo '-- make help-deployment --'
	$(MAKE) help-deployment
	@echo ''
	@echo '-- make others --'
	@echo '  protoc                        - Codegen proto files matching ./protos/*.proto'

help-build:
	@(for dir in $(BUILDDIRS); do \
		echo $${dir}
		echo -n '  '
		$(MAKE) -s -C $(BUILDPATH)/$$dir help; \
		echo ''
	done)

help-deployment:
	@(for dir in $(DEPLOYMENTDIRS); do \
		echo $${dir}
		echo -n '  '
		$(MAKE) -s -C $(DEPLOYMENTPATH)/$$dir help; \
		echo ''
	done)

.ONESHELL:
lint: misspell staticcheck vet gofmt
	golangci-lint run --verbose ./...

.ONESHELL:
misspell:
	misspell -locale UK .

.ONESHELL:
staticcheck:
	staticcheck ./...

.ONESHELL:
vet:
	go vet ...

.ONESHELL:
gofmt:
	gofmt -s -d -w .

.ONESHELL:
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	|| go install honnef.co/go/tools/cmd/staticcheck@latest

.ONESHELL:
protoc:
	@(cd ./protos/ ; protoc --go_out=. *.proto)

.ONESHELL:
.SILENT:
.PHONY:
%:
	$(eval argv=$(subst _, , ${MAKECMDGOALS})) \
	$(if $(filter $(firstword $(argv)),$(BUILDDIRS)), \
		$(eval build=$(word 1, $(argv))) \
		$(eval target=$(word 2, $(argv))) \

		$(info Building: $(build)) \
		$(info Target: $(target)) \

		$(MAKE) -C $(BUILDPATH)/$(build) $(target)
	)
	$(if $(filter $(firstword $(argv)),$(DEPLOYMENTDIRS)), \
		$(eval deployment=$(word 1, $(argv))) \
		$(eval target=$(word 2, $(argv))) \

# TODO: Only build docker image if we deploy.
#   Right now we always do it, even on "docker-compose_down"
		$(MAKE) docker_build
		$(info Deploying: $(deployment)) \
		$(info Target: $(target)) \

		$(MAKE) -C $(DEPLOYMENTPATH)/$(deployment) $(target)
	)
