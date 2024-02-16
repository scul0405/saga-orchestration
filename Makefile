run_account:
	go run services/account/cmd/api/main.go
run_product:
	go run services/product/cmd/api/main.go
run_all:
	go run services/account/cmd/api/main.go &
	go run services/product/cmd/api/main.go
