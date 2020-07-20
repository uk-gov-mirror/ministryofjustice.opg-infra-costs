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
	rm -Rf ./tmp/*.json
	aws-vault exec sandbox-sandbox-breakglass -- ./aws-costs costs-to-file
	aws-vault exec sirius-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec sirius-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec sirius-prod-breakglass -- ./aws-costs costs-to-file	
	aws-vault exec sirius-backup-breakglass -- ./aws-costs costs-to-file
	aws-vault exec lpa-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec lpa-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec lpa-prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec digideps-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec digideps-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec digideps-prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec refunds-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec refunds-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec refunds-prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec ual-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec ual-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec ual-prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec org-mangagement-breakglass -- ./aws-costs costs-to-file
	aws-vault exec jenkins-dev-accountwrite -- ./aws-costs costs-to-file	
	aws-vault exec jenkins-prod-accountwrite -- ./aws-costs costs-to-file
	aws-vault exec lpa-LEGACY_prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec serve-dev-breakglass -- ./aws-costs costs-to-file
	aws-vault exec serve-preprod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec serve-prod-breakglass -- ./aws-costs costs-to-file
	aws-vault exec refunds-LEGACY_prod-accountwrite -- ./aws-costs costs-to-file
	aws-vault exec refunds-LEGACY_dev-accountwrite -- ./aws-costs costs-to-file
	aws-vault exec identity-identity-breakglass -- ./aws-costs costs-to-file
	

.PHONY: create-spreadsheet
create-spreadsheet:
	[ -f ./costs.xlsx ] && rm ./costs.xlsx || echo "Costs spreadsheet does not exist"
	./aws-costs spreadsheet