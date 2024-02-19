run_account:
	cd cmd/account & go run main.go
run_product:
	cd services/product & go run cmd/api/main.go
run_order:
	cd cmd/order & go run main.go
run_payment:
	cd services/payment & go run cmd/api/main.go