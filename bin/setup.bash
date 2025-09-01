#!/bin/bash

DEFAULT_TARGET_BRANCH="origin/main"

cat "${TEST_DIR_UTILS}/data/copyright.txt"

printf "Setup environment.\n"

if test -z "${WORKDIR}"; then
  WORKDIR=$(pwd)
fi
printf "Using ${WORKDIR} as work directory.\n"

printf "Add current dir into safe directory.\n"
git config --global --add safe.directory ${WORKDIR}

if test "$(git rev-parse --is-shallow-repository)" == true ; then
  git fetch -a --unshallow || exit 1
else
  git fetch -a || exit 1
fi

if test -z "${TARGET_BRANCH}"
then
    TARGET_BRANCH=${DEFAULT_TARGET_BRANCH}
fi

if test -z "$HEAD"
then
  HEAD=HEAD
fi
printf "HEAD=${HEAD}\nTARGET_BRANCH=${TARGET_BRANCH}\n"

printf "Setup environment finished.\n"
