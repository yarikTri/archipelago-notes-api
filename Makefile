.PHONY: swagger

# Generate Swagger documentation
swagger:
	~/go/bin/swag init -g cmd/api/main.go

docker/rebuild:
	docker compose down && docker compose build && docker compose up -d

test/docker/rebuild:
	docker compose -f docker-compose-stateless.yml down && docker compose -f docker-compose-stateless.yml build && docker compose -f docker-compose-stateless.yml up -d
