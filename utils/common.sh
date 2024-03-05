#!/bin/bash

. "$TEST_DIR_UTILS/constants.sh"

get_diff() {
  local head="$1"
  local changes="$(git diff --name-only -z "$BASE_BRANCH"..."$head" | tr '\000' '\n')"

  if test -z "$changes"
  then
      changes="$(git diff --name-only -z "$head^" | tr '\000' '\n')"
  fi

  echo "$changes"
}

get_log() {
    local head="$1"
    local files="$2"
    local format="$3"
    if test "$files"
    then
        files="-- $files"
    fi
    if test "$format"
    then
        format="--format=$format"
    fi
    local changes="$(git log $format "$BASE_BRANCH"..."$head" $files)"
    if test -z "$changes"
    then
        changes="$(git log $format "$head^..$head" $files)"
    fi
    echo "$changes"
}

get_labs_files() {
  local head="$1"
  local files=$(get_diff "$head" | grep -E "$LAB_FILES_REGEXP_PATTERN")
  echo "$files"
}

get_students() {
  local head="$1"
  local students=$(get_labs_files "$head" | grep -o -E "^$STUDENT_REGEXP_PATTERN" | sort | uniq)
  echo "$students"
}

get_tasks() {
  local head="$1"
  local student_name="$2"
  local tasks=$(get_labs_files "$head" | grep -E "^$student_name" | grep -o -E "$TASK_REGEXP_PATTERN" | sort | uniq)
  echo "$tasks"
}

get_cfg_value() {
    local file="$1"
    local key="$2"
    if test -f "$TEST_DIR_COMMON/$task/ci.cfg" ; then
      local value=$(grep -E "$key" "$file" | awk '{print $2}')
    fi
    echo "$value"
}

print_copyright() {
  cat "$TEST_DIR_UTILS/copyright.txt"
}

print_copyright
printf "\033[34mPreparing the environment for execution... \033[0m\n"

if test "$COMMON_REPO_URL"; then
  printf "Download common from %s\n" "$COMMON_REPO_URL"
  git clone "$COMMON_REPO_URL" "$TEST_DIR_COMMON"
fi

if test -z "$WORKDIR"; then
  WORKDIR=$(pwd)
fi
export WORKDIR

if test "$(git rev-parse) --is-shallow-repository)" == true ; then
  git fetch -a --unshallow || exit 1
else
  git fetch -a || exit 1
fi

if test "$CI_MERGE_REQUEST_PROJECT_URL"
then
    if git remote | grep 'main' > /dev/null 2>&1
    then
        git remote remove main || exit 1
    fi
    printf "Add %s as main/master\n" "https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.com/$CI_MERGE_REQUEST_PROJECT_PATH"
    git remote add main "https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.com/$CI_MERGE_REQUEST_PROJECT_PATH" || exit 1
    git fetch main || exit 1
    BASE_BRANCH=main/master
fi

if test -z "$BASE_BRANCH"
then
    BASE_BRANCH=origin/master
fi
export BASE_BRANCH

if test -z "$CI_MERGE_REQUEST_SOURCE_BRANCH_SHA"
then
  HEAD="$CI_MERGE_REQUEST_SOURCE_BRANCH_SHA"
fi
if test -z "$HEAD"
then
  HEAD=HEAD
fi
printf "HEAD=%s\n" "$HEAD"
export HEAD

mkdir -p logs

printf "\033[34mThe preparation of the environment is complete! \033[0m\n"
