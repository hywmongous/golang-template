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

build_clean:
	@(for dir in $(BUILDDIRS); do \
		$(MAKE) -C $(BUILDPATH)/$$dir build_clean; \
	done)

dist_clean:
	@(for dir in $(BUILDDIRS); do \
		$(MAKE) -C $(BUILDPATH)/$$dir dist_clean; \
	done)

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

		$(info Deploying: $(deployment)) \
		$(info Target: $(target)) \

		$(MAKE) -C $(DEPLOYMENTPATH)/$(deployment) $(target)
	)
