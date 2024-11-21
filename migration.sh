goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/restWallet?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:password@localhost:5432/restWallet?sslmode=disable" up