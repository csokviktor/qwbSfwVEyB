SERVICE_NAME = manager
DOCKER_IMAGE_NAME = lib_manager/${SERVICE_NAME}
DOCKER_IMAGE_TAGGED = lib_manager/${SERVICE_NAME}:0.0.1

buildDockerWithTag:
	docker build -f Dockerfile -t ${DOCKER_IMAGE_TAGGED} .

run:
	go run ./cmd

mocks:
	go generate -x -run="mockgen" ./...

tests:
	mkdir -p test_output
	chmod 777 test_output
	go test ./... -covermode=set -coverprofile=test_output/coverage.out fmt
	go tool cover -html=test_output/coverage.out -o test_output/cover.html

lint:
	golangci-lint run --fix -c ./.golangci.yaml

migration:
	cd db-migrations && \
	dbmate new "$(name)"

buildDocker:
	docker build -f Dockerfile -t ${SERVICE_NAME} .

runDocker:
	docker run \
	--rm -it \
	--name ${SERVICE_NAME} \
	${SERVICE_NAME}