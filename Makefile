migrateschema:
		migrate create -ext sql -dir db/migration -seq init_schema

createdb:
		docker exec -it mysql-container-test mysql -u root -proot -e "CREATE DATABASE IF NOT EXISTS service_pro;"

migrateup:
		migrate -path db/migration -database "mysql://root:root@tcp(127.0.0.1:3308)/service_pro" -verbose up 

migratedown:
		migrate -path db/migration -database "mysql://root:root@tcp(127.0.0.1:3308)/service_pro" -verbose down 

sqlc:
		sqlc generate

.PHONY: migrateschema createdb addtable
