flyway := \
	docker run --rm -v  $${LOCAL_WORKSPACE_FOLDER}/postgres/migrations:/flyway/sql --network api_dev \
		flyway/flyway \
		-user=$${APP_POSTGRES_USER} \
		-password=$${APP_POSTGRES_PASSWORD} \
		-url=jdbc:postgresql://$${APP_POSTGRES_HOST}:$${APP_POSTGRES_PORT}/$${APP_POSTGRES_DB}?ssl=false

migrate:
	$(flyway) migrate

fresh:
	$(flyway) clean && $(flyway) migrate

test:
	go test -v -cover ./...

.PHONY: migrate fresh test
