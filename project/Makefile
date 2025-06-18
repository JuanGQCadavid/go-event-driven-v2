up:
	docker compose up --build -d

down:
	docker compose down

run:
	REDIS_ADDR=localhost:6379 GATEWAY_ADDR=http://localhost:8888 go run .

test:
	curl -v -X POST http://localhost:8080/tickets-confirmation \
		-H "Content-Type: application/json" \
		-d '{"tickets": ["ticket-1", "ticket-2"]}'