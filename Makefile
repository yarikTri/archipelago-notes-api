.PHONY: swagger

# Generate Swagger documentation
swagger:
	~/go/bin/swag init -g cmd/api/main.go 