#!/bin/bash

dir=$(dirname "$0")
. "$dir/constants.sh"

get_diff () {
  local head="$1"
  local files=$(git diff --name-only -z "$BASE_BRANCH"..."$head" | tr '\000' '\n')
  echo "$files"
}

get_labs_files() {
  local head="$1"
  local files=$(get_diff "$head" | grep -E "$LAB_FILES_REGEXP_PATTERN")
  echo "$files"
}

get_students () {
  local head="$1"
  local students=$(get_labs_files "$head" | grep -o -E "^$STUDENT_REGEXP_PATTERN" | sort | uniq)
  echo "$students"
}

get_tasks () {
  local head="$1"
  local student_name="$2"
  local tasks=$(get_labs_files "$head" | grep -E "^$student_name" | grep -o -E "$TASK_REGEXP_PATTERN" | sort | uniq)
  echo "$tasks"
}

get_cfg_value() {
    local file="$1"
    local key="$2"
    if test -f "$PATH_TO_COMMON/$task/ci.cfg" ; then
      local value=$(grep -E "$key" "$file" | awk '{print $2}')
    fi
    echo "$value"
}

if test -z "$WORKDIR"; then
  WORKDIR=$(PWD)
fi
export WORKDIR

if test "$(git rev-parse) --is-shallow-repository)" == true; then
  git fetch -a --unshallow
else
  git fetch -a
fi

if test "$CI_MERGE_REQUEST_PROJECT_URL"
then
    if git remote | grep 'main' > /dev/null 2>&1
    then
        git remote remove main || exit 1
    fi
    git remote add main "https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.com/$CI_MERGE_REQUEST_PROJECT_PATH" || exit 1
    git fetch main || exit 1
    BASE_BRANCH=main/master
fi

if test -z "$BASE_BRANCH"
then
    BASE_BRANCH=origin/master
fi
export BASE_BRANCH

if test -z "$HEAD"
then
  HEAD=HEAD
fi
export HEAD