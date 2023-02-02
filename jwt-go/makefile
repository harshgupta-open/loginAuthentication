migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mock -destination db/mock/store.go golang-boiler-plate/db Store
migrate:
	migrate create -ext sql -dir db/migration -seq ${FILE_NAME}

server:
	go run main.go