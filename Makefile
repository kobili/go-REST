build:
	docker compose build

run:
	docker compose up

make-migration:
	docker compose -f compose.yaml -f compose.migrate.yaml run --rm migrate create -ext sql -dir ./migrations -seq ${migration_name}
mm: make-migration

migrate:
	docker compose -f compose.yaml -f compose.migrate.yaml run --rm migrate -path=/migrations/ -database postgres://postgres:password@postgres:5432/go_test?sslmode=disable up
m: migrate

count ?= 1 # use -all to go all the way down
migrate-down:
	docker compose -f compose.yaml -f compose.migrate.yaml run --rm migrate -path=/migrations/ -database postgres://postgres:password@postgres:5432/go_test?sslmode=disable down ${count}
md: migrate-down

destroy:
	docker compose stop
	docker compose down -v --rmi local
