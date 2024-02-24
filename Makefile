build:
	docker-compose build

run:
	docker-compose up

swagger:
	swag init -g cmd/main.go

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go