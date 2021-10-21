FEATURES=$(shell ls -d */)

.ONESHELL:
.SILENT:
.PHONY:

help:
	@echo 'Targets:'
	@echo '  features                      - Lists all the subdirectories which should be all features'
	@echo '  %                             - Wildcard which either targets a specific feature or all'
	@echo 'Examples:'
	@echo '  "make smoke"                  - Runs the smoke_all.js k6 test'
	@echo '  "make identity-login_smoke"   - Runs the identity-login feature smoke test with k6'


features:
	@(for feature in $(subst _, , ${FEATURES}); do \
		echo $${feature}
	done)

%:
	$(eval argv=$(subst _, , ${MAKECMDGOALS})) \
	$(eval type=$(word 1, $(argv))) \
	$(eval feature=$(word 2, $(argv))) \

	$(if $(feature), \
		k6 run ./$(type)/$(feature)_test.js \
	, \
		k6 run ./$(type)_all.js \
	)
