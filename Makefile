DOCKER_CONTAINER_NAME=andrianovartemii/lgp-ci-utils
DOCKER_CONTAINER_TAG?=latest
DOCKER_FILE_PATH=deployment/docker/Dockerfile
DOCKER_PLATFORM_ARGS?=linux/amd64

docker-container-build:
	docker build --platform=${DOCKER_PLATFORM_ARGS} -t ${DOCKER_CONTAINER_NAME} -f ${DOCKER_FILE_PATH} .

docker-container-tag: docker-container-build
	docker tag ${DOCKER_CONTAINER_NAME} ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG}
    docker tag ${DOCKER_CONTAINER_NAME} ${DOCKER_CONTAINER_NAME}:latest"	

docker-container-push: docker-container-tag
	docker push --disable-content-trust ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG}
    docker push --disable-content-trust ${DOCKER_CONTAINER_NAME}:latest
