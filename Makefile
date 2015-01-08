checkDB:
ifndef DB_URL
	@echo "DB_URL not set"
	@exit 1
endif
.PHONY: checkDB

migrate: checkDB
	migrate -url $(DB_URL) -path ./data/migrations up
.PHONY: migrate
