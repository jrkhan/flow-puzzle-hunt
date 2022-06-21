init-collection:
	flow transactions send ./transactions/create_account.cdc --signer emulator-account

check-collection:
	flow scripts execute ./scripts/items_in_accounts.cdc 