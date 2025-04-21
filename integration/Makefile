SERVICE_NAME = integration
DOCKER_IMAGE_NAME = lib_manager/${SERVICE_NAME}
DOCKER_IMAGE_TAGGED = lib_manager/${SERVICE_NAME}:0.0.1

buildDockerWithTag:
	docker build -f Dockerfile -t ${DOCKER_IMAGE_TAGGED} .

buildDocker:
	docker build -f Dockerfile -t ${SERVICE_NAME} .

runDocker:
	docker run \
	-it \
	--name ${SERVICE_NAME} \
	--network applications \
	${SERVICE_NAME}

tests:
	poetry run pytest