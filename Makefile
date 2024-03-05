DOCKER_CONTAINER_NAME?=spbstu-2024-lgp-ci
DOCKER_LOGIN?=andrianovartemii
DOCKER_FILE_PATH?=deployment/docker/Dockerfile
COMMON_REPO_URL?=https://gitlab.com/lets-go-programming/lgp-ci-common.git
TAG?=


build:
	docker build -t $(DOCKER_CONTAINER_NAME) -f $(DOCKER_FILE_PATH) .

tag: build
	docker tag $(DOCKER_CONTAINER_NAME) $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):$(TAG)
	docker tag $(DOCKER_CONTAINER_NAME) $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):latest

push: build tag
	docker push --disable-content-trust $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):$(TAG)
	docker push --disable-content-trust $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):latest

docker_workdir=/spbstu-2024-lgp-ci
docker_args = --rm \
    -v ${PWD}/utils:${docker_workdir}/utils \

bash: build tag
	docker run -it ${docker_args} -e COMMON_REPO_URL=${COMMON_REPO_URL}  $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):latest bash