init-collection:
	flow transactions send ./transactions/create_account.cdc --signer emulator-account

check-collection:
	flow scripts execute ./scripts/items_in_accounts.cdc 0xf8d6e0586b0a20c7

account-exists:
	flow scripts execute ./scripts/account_exists.cdc 0xf8d6e0586b0a20c7

set-image:
		flow transactions send ./transactions/set_image.cdc --signer emulator-account "display name" "https://ipfs.io/ipfs/QmU4SndLCLdbGtnthGm7ZbwKSw6cUs6WRUPdJ18pUrePQc?filename=piece-1.png"

mint:
	flow transactions send ./transactions/mint.cdc --signer emulator-account 0xf8d6e0586b0a20c7 "display name" 1 1

generate-qr:
	go run ./cmd/qrcodes