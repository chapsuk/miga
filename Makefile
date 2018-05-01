NAME = miga
VERSION ?= develop
PG_CONTAINER_NAME = miga-pg
MYSQL_CONTAINER_NAME = miga-mysql

.PHONY: build
build:
	go build -o bin/$(NAME) -ldflags "-X main.Version=${VERSION}" main.go

.PHONY: postgres_up
postgres_up: postgres_down
	docker run -d \
		-p 5432:5432 \
		-e POSTGRES_DB=miga	\
		-e POSTGRES_USER=user \
		-e POSTGRES_PASSWORD=password \
		--name=$(PG_CONTAINER_NAME) postgres:9.6.5-alpine

.PHONY: postgres_down
postgres_down:
	-docker rm -f $(PG_CONTAINER_NAME)

.PHONY: mysql_up
mysql_up: mysql_down
	docker run -d \
		-p 3306:3306 \
		-e MYSQL_DATABASE=miga \
		-e MYSQL_USER=user \
		-e MYSQL_PASSWORD=password \
		-e MYSQL_ROOT_PASSWORD=mysql \
		--name=$(MYSQL_CONTAINER_NAME) mysql:5.7

.PHONY: mysql_down
mysql_down:
	-docker rm -f $(MYSQL_CONTAINER_NAME)

.PHONY: db_up
db_up: postgres_up mysql_up

.PHONY: db_down
db_down: postgres_down mysql_down