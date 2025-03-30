.PHONY: swagger

# Generate Swagger documentation
swagger:
	~/go/bin/swag init -g cmd/api/main.go 

docker/rebuild:
	docker compose down && docker compose build && docker compose up -d