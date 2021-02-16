flyway := \
	docker run --rm -v $${LOCAL_WORKSPACE_FOLDER}/store/postgres/migrations:/flyway/sql --network api_dev \
		flyway/flyway \
		-user=$${POSTGRES_USER} \
		-password=$${POSTGRES_PASSWORD} \
		-url=jdbc:postgresql://$${POSTGRES_HOST}:$${POSTGRES_PORT}/$${POSTGRES_DB}?ssl=false

migrate:
	$(flyway) migrate

fresh:
	$(flyway) clean && $(flyway) migrate

test:
	go test -v -cover ./...

.PHONY: migrate fresh test
