SERVICE_NAME ?= 242617/notify:0.0.10

# Debug
.PHONY: proto
proto:
	@protoc \
		--proto_path=server/proto/ \
		--go_out=server/ \
		--go_opt=plugins=grpc \
		server/proto/*.proto
	@protoc \
		--proto_path=server/proto/ \
		--go_out=server/ \
		--go_opt=plugins=grpc \
		server/proto/grpc/health/v1/health.proto

.PHONY: build
build:
	go build \
		-o bin/app \
		cmd/main/main.go

.PHONY: run
run: build
	CONSUL_HTTP_TOKEN="" \
		./bin/app \
			--address=localhost:8080 \
			--consul=engine:8080


# Docker
docker\:build:
	docker build \
		-t ${SERVICE_NAME} \
		-f Dockerfile \
		.

docker\:run:
	docker run \
		-it --rm \
		${SERVICE_NAME}

docker\:push:
	docker push ${SERVICE_NAME}
