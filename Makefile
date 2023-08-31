docker:
	docker build -t docker-postgres-go-rest-api .

run:
	docker compose up --build -d

stop:
	docker compose down

serve:
	go run cmd/main.gow