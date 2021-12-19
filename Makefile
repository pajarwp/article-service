migrateup:
	migrate -path migrations -database "mysql://user:password@tcp(127.0.0.1:4306)/article" -verbose up

migrateuptest:
	migrate -path migrations -database "mysql://user:password@tcp(127.0.0.1:4306)/article_test" -verbose up

migratedown:
	migrate -path migrations -database "mysql://user:password@tcp(127.0.0.1:4306)/article" -verbose down

migratedowntest:
	migrate -path migrations -database "mysql://user:password@tcp(127.0.0.1:4306)/article_test" -verbose down

test:
	ENV=testing go test ./...

.PHONY: migrateup migrateuptest migratedown migratedowntest test