SHELL = '/bin/bash'

.PHONY: all
all: ##
	@${MAKE} install
	@${MAKE} get-costs
	@${MAKE} create-spreadsheet

.PHONY: install
install:
	composer install && composer dump -o

.PHONY: get-costs
get-costs:
	mkdir -p ./tmp
	rm -Rf ./tmp/*.json
	aws-vault exec identity -- ./aws-costs costs-to-file

.PHONY: create-spreadsheet
create-spreadsheet:
	[ -f ./costs.xlsx ] && rm ./costs.xlsx || echo "Costs spreadsheet does not exist"
	./aws-costs spreadsheet
