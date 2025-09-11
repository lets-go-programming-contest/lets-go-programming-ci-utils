DOCKER_CONTAINER_NAME?=spbstu-2024-lgp-ci
DOCKER_FILE_PATH?=deployment/dev/Dockerfile

DOCKER_CACHE_ARGS?= \
	-v ${HOME}/.docker_cache:${HOME}/.cache:rw

DOCKER_USER_ARGS?= \
	--user="$(shell id -u):$(shell id -g)" \
	-v /etc/passwd:/etc/passwd:ro \
	-v /etc/group:/etc/group:ro

DOCKER_USER?= --rm \
	-w ${PWD} -v ${PWD}:${PWD}:rw \

docker-build:
	docker build -t $(DOCKER_CONTAINER_NAME) -f $(DOCKER_FILE_PATH) .

docker-bash: docker-build
	docker run -it ${DOCKER_USER} -e COMMON_REPO_URL=${COMMON_REPO_URL} ${DOCKER_CONTAINER_NAME} bash
