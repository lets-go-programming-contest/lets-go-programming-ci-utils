DOCKER_CONTAINER_NAME?=andrianovartemii/lgp-ci-utils
DOCKER_CONTAINER_TAG?=latest
DOCKER_FILE_PATH?=deployment/docker/Dockerfile

DOCKER_WORKDIR?=${PWD}

TEST_DIR_UTILS?=${PWD}/.ci_utils
TEST_DIR_COMMON?=${PWD}/.ci_common

docker-build-container:
	docker build --progress=plain --platform=linux/amd64 -t ${DOCKER_CONTAINER_NAME} -f ${DOCKER_FILE_PATH} .

docker-tag: docker-build-container
	docker tag ${DOCKER_CONTAINER_NAME} ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG}
	docker tag ${DOCKER_CONTAINER_NAME} ${DOCKER_CONTAINER_NAME}:latest

docker-push: docker-tag
	docker push --disable-content-trust ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG}
	docker push --disable-content-trust ${DOCKER_CONTAINER_NAME}:latest

DOCKER_CI_ARGS?= \
	-e TEST_DIR_UTILS=${TEST_DIR_UTILS} \
	-e TEST_DIR_COMMON=${TEST_DIR_COMMON}

DOCKER_ARGS?= --rm -w ${DOCKER_WORKDIR} \
	-v ${PWD}:${DOCKER_WORKDIR}:rw \
	$(DOCKER_CI_ARGS)

docker-%:
	docker run ${DOCKER_ARGS} ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG} make -f ${TEST_DIR_UTILS}/Makefile ${*}

docker-bash:
	docker run -it ${DOCKER_ARGS} ${DOCKER_CONTAINER_NAME}:${DOCKER_CONTAINER_TAG} bash

sanity-files:
	${TEST_DIR_UTILS}/bin/sanity-files.bash

sanity-student:
	${TEST_DIR_UTILS}/bin/sanity-student.bash

sanity-tasks:
	${TEST_DIR_UTILS}/bin/sanity-tasks.bash

build:
	${TEST_DIR_UTILS}/bin/build.bash

lint:
	${TEST_DIR_UTILS}/bin/lint.bash

tests:
	${TEST_DIR_UTILS}/bin/tests.bash
