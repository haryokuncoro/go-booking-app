DB_URL=postgres://root:root@localhost:5432/go_bookingdb?sslmode=disable

migrate-up:
	migrate -path migrations \
	-database "$(DB_URL)" up

migrate-down:
	migrate -path migrations \
	-database "$(DB_URL)" down 1