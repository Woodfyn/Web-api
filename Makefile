build:
	docker-compose build

run:
	docker-compose up

restart:
	docker-compose restart

swagger:
	swag init -g cmd/main.go