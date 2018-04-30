.PHONY: postgres_up
postgres_up:
	docker run -d \
		-p 5432:5432 \
		-e POSTGRES_DB=miga	\
		-e POSTGRES_USER=user \
		-e POSTGRES_PASSWORD=password \
		--name=miga-pg postgres:9.6.5-alpine

.PHONY: postgres_down
postgres_down:
	docker rm -f miga-pg

.PHONY: mysql_up
mysql_up:
	docker run -d \
		-p 3306:3306 \
		-e MYSQL_DATABASE=miga \
		-e MYSQL_USER=user \
		-e MYSQL_PASSWORD=password \
		-e MYSQL_ROOT_PASSWORD=mysql \
		--name=miga-mysql mysql:5.7

.PHONY: mysql_down
mysql_down:
	docker rm -f miga-mysql