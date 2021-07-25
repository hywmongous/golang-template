FEATURESPATH=./features
K6=k6
TESTS=$(K6)

FEATURES=$(shell ls $(FEATURESPATH))

help:
	@echo 'Features:'
	$(MAKE) help-features
	@echo '-- make help-k6 --'
	$(MAKE) help-k6

help-k6:
	@echo '  eg. k6_account-login_load'

help-features:
	echo -n '  ' $(FEATURES)
	echo ''

.ONESHELL:
.SILENT:
.PHONY:
%:
	$(eval argv=$(subst _, , ${MAKECMDGOALS})) \
	$(if $(filter $(word 1, $(argv)),$(TESTS)), \
		$(eval test=$(word 1, $(argv))) \
		$(eval feature=$(word 2, $(argv))) \
		$(eval type=$(word 3, $(argv))) \

		$(info Testing: $(test)) \
		$(info Feature: $(feature)) \
		$(info Type: $(type)) \

		$(if $(filter $(word 1, $(test)),$(K6)),
			$(K6) run $(FEATURESPATH)/$(feature)/$(type)_test.js
		)
	)
