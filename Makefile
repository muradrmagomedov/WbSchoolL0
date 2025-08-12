db:
	docker run --name postgress -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
redis:
	docker run --name redis-instance -p 6379:6379 -d redis
migrateup:
	migrate -path storage/migration -database "postgresql://root:secret@localhost:5431/root?sslmode=disable" -verbose up
migratedown:
	migrate -path storage/migration -database "postgresql://root:secret@localhost:5431/root?sslmode=disable" -verbose down
run:
	go run ./cmd/app/.

PHONY: db migrateup migratedown run redis
