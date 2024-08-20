NAME = miga
VERSION ?= dev
PG_CONTAINER_NAME = miga-pg
MYSQL_CONTAINER_NAME = miga-mysql
CLICKHOUSE_CONTAINER_NAME = miga-clickhouse
VERTICA_CONTAINER_NAME = miga-vertica
IMAGE_NAME = chapsuk/$(NAME)

TRAVIS_POSTGRES = postgres://postgres:@127.0.0.1:5432/miga?sslmode=disable
TRAVIS_MYSQL = travis:@tcp(127.0.0.1:3306)/miga

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -o bin/$(NAME) -ldflags "-w -v -extldflags \"-static\" -X main.Version=$(VERSION)" main.go

.PHONY: docker_build
docker_build:
	docker build -t $(IMAGE_NAME):$(VERSION) .

.PHONY: release
release: docker_build
	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest

.PHONY: test
test:
	go clean -testcache
	go test -v -race ./...

.PHONY: db_up
db_up: vertica_up postgres_up mysql_up clickhouse_up

.PHONY: db_down
db_down: postgres_down mysql_down clickhouse_down vertica_down

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

.PHONY: clickhouse_down
clickhouse_down:
	-docker rm -f $(CLICKHOUSE_CONTAINER_NAME)

.PHONY: clickhouse_up
clickhouse_up: clickhouse_down
	docker compose up -d

.PHONY: vertica_down
vertica_down:
	-docker rm -f $(VERTICA_CONTAINER_NAME)

.PHONY: vertica_up
vertica_up: vertica_down
	docker run -d \
		-p 5433:5433 \
		--name=$(VERTICA_CONTAINER_NAME) \
		dataplatform/docker-vertica
