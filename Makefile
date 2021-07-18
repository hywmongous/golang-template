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
	@echo 'Configuration targets:'
	@echo '  format: "make BUIILD TARGET"  - eg. "make golang build_rest-server'
	@echo ''
	@echo '--Help for builds--'
	@(for dir in $(BUILDDIRS); do \
		echo $${dir}
		echo -n '  '
		$(MAKE) -s -C $(BUILDPATH)/$$dir help; \
		echo ''
	done)
	@echo ''
	@echo '--Help for deployments--'
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
	$(info Building: $(filter-out ${MAKECMDGOALS}, $(firstword $(subst _, , ${MAKECMDGOALS}))))
	$(info Target: $(subst ${SPACE},_,$(filter-out $(firstword $(subst _, , ${MAKECMDGOALS})), $(subst _, , ${MAKECMDGOALS}))))
	$(MAKE) -C $(BUILDPATH)/$(firstword $(subst _, , ${MAKECMDGOALS})) \
		$(subst ${SPACE},_,$(filter-out $(firstword $(subst _, , ${MAKECMDGOALS})), $(subst _, , ${MAKECMDGOALS})))
