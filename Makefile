# https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents

BUILDPATH=./builds
BUILDDIRS=golang docker

DEPLOYMENTPATH=./deployments
DOCKERCOMPOSE=docker-compose
DOCKERSTACK=docker-stack
DEPLOYMENTDIRS=$(DOCKERCOMPOSE) $(DOCKERSTACK)
TARGETREQURIESDOCKERBUILD=up create

TESTPATH=./test
TESTMETHODS=k6

GITHUBPATH=./.github
WORKFLOWS=workflow workflows
WORKFLOWSPATH=$(GITHUBPATH)/workflows

.ONESHELL:
.SILENT:
.PHONY:

help:
	@echo 'Project targets:'
	@echo '  lint                          - Applies golangci and misspell, staticcheck, vet, and gofmt'
	@echo '  misspell                      - Applies client9 misspell with UK'
	@echo '  vet                           - Runs go vet with all checks'
	@echo '  gofmt                         - Applies go formatting'
	@echo '  gotest                        - runs all go tests with race check'
	@echo '  install                       - Installs all dedpendencies, eg. golangci-lint, and misspell'
	@echo '  protoc                        - Codegen proto files matching ./protos/*.proto'
	@echo '  %                             - Wildcard which constructs a bubild/deployment/act target'
	@echo 'Help targets:'
	@echo '  help-builds                   - Prints all the help targets from the builds'
	@echo '  help-deployments              - Prints all the help targets from the deployment'
	@echo 'Examples:'
	@echo '  "make vet"                    - Golang vet project'
	@echo '  "make golang_build"           - Build golang binary'
	@echo '  "make docker_run"             - Run docker'
	@echo '  "make docker_hellp"           - Get docker help'
	@echo '  "make docker-compose_up"      - Deploy docker-compose'
	@echo '  "make k6_smoke"               - Runs k6 smoke_all test'
	@echo '  "make k6_smoke_identity-login"- Runs k6 smoke test for the feature identity-login'

lint: misspell staticcheck vet gofmt
# The disabled linters are deprecated
	golangci-lint run --enable-all --verbose --sort-results --tests --fix \
		--disable maligned --disable interfacer --disable scopelint --disable golint ./...

misspell:
	misspell -locale UK .

staticcheck:
	staticcheck ./...

vet:
	go vet -all ./...

fmt:
	gofmt -s -d -w .
	gofumpt -l -w .

gotest:
	go test -v -race -covermode=atomic ./...

install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install mvdan.cc/gofumpt@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	sudo apt install -y protobuf-compiler

protoc:
	@(cd ./protos/ ; protoc --go_out=. *.proto)

%:
	$(eval argv=$(subst _, , ${MAKECMDGOALS})) \
	$(eval target=$(word 2, $(argv))) \

# No target? Then the default is "help"
	$(if $(target),, \
		$(eval target=help) \
	)

# Test targets
	$(if $(filter $(firstword $(argv)),$(TESTMETHODS)), \
		$(eval test=$(word 1, $(argv))) \

		$(info Test: $(test)) \
		$(info Target: $(target)) \

		$(MAKE) -C $(TESTPATH)/$(test) $(target)
	)

# Workflow targets
	$(if $(filter $(firstword $(argv)),$(WORKFLOWS)), \
		$(info Workflow: $(target)) \

		$(MAKE) -C $(WORKFLOWSPATH) $(target)
	)

# Build targets
	$(if $(filter $(firstword $(argv)),$(BUILDDIRS)), \
		$(eval build=$(word 1, $(argv))) \

		$(info Build: $(build)) \
		$(info Target: $(target)) \

		$(MAKE) -C $(BUILDPATH)/$(build) $(target)
	)

# Deployment targets
	$(if $(filter $(firstword $(argv)),$(DEPLOYMENTDIRS)), \
		$(eval deployment=$(word 1, $(argv))) \
		$(info Deployment: $(deployment))

# The docker stack targets are in the same folder as docker-compose
		$(if $(filter $(deployment),$(DOCKERSTACK)), \
			$(eval deployment=$(DOCKERCOMPOSE)) \
		)

		$(info Target: $(target)) \

# Only build the docker image if we are deploying it
		$(if $(filter $(target),$(TARGETREQURIESDOCKERBUILD)), \
			$(MAKE) docker_build \
		)

		$(MAKE) -C $(DEPLOYMENTPATH)/$(deployment) $(target)
	)
