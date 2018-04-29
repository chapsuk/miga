postgres_up:
	docker run -d \
		-p 5432:5432 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_DB=miga	\
		-e POSTGRES_PASSWORD=postgres \
		--name=miga-pg postgres:9.6.5-alpine

postgres_down:
	docker rm -f miga-pg
