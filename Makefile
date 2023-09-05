start-docker-compose:
	docker-compose up

fw-migrate:
	docker-compose run flyway -locations=filesystem:/flyway/sql -connectRetries=60 migrate

.PHONY: start-docker-compose fw-migrate