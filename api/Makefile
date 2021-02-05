flyway := \
	docker run --rm -v $${LOCAL_WORKSPACE_FOLDER}/store/postgres/migrations:/flyway/sql --network syncedvideo_dev \
		flyway/flyway \
		-user=$${API_POSTGRES_USER} \
		-password=$${API_POSTGRES_PASSWORD} \
		-url=jdbc:postgresql://$${API_POSTGRES_HOST}:$${API_POSTGRES_PORT}/$${API_POSTGRES_DB}?ssl=false

migrate:
	$(flyway) migrate

fresh:
	$(flyway) clean && $(flyway) migrate

test:
	go test -v -cover ./...

.PHONY: migrate fresh test
