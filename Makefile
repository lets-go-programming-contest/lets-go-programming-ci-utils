DOCKER_CONTAINER_NAME?=spbstu-2024-lgp-ci
DOCKER_LOGIN?=andrianovartemii
DOCKER_FILE_PATH?=deployment/docker/Dockerfile

build:
	docker build -t $(DOCKER_CONTAINER_NAME) -f $(DOCKER_FILE_PATH) .

DOCKER_CONTAINER_TAG?=latest
tag: build
	docker tag $(DOCKER_CONTAINER_NAME) $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):$(DOCKER_CONTAINER_TAG)
	docker tag $(DOCKER_CONTAINER_NAME) $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):latest

push: build tag
	docker push --disable-content-trust $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):$(DOCKER_CONTAINER_TAG)
	docker push --disable-content-trust $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):latest


bash: build tag
	docker run -it $(DOCKER_LOGIN)/$(DOCKER_CONTAINER_NAME):$(DOCKER_CONTAINER_TAG) bash