#!/bin/bash

STUDENT_REGEXP_PATTERN='[a-z0-9]+\.[a-z0-9]+'
TASK_REGEXP_PATTERN='task-[0-9-]+'
LAB_FILES_REGEXP_PATTERN="^${STUDENT_REGEXP_PATTERN}\/${TASK_REGEXP_PATTERN}\/.+$"

get_diff() {
  local target="$1"
  local head="$2"

  local changes="$(git diff --name-only -z "${target}"..."${head}" | tr '\000' '\n')"

  if test -z "${changes}"
  then
      changes="$(git diff --name-only -z "${head}^" | tr '\000' '\n')"
  fi

  echo "${changes}"
}

get_log() {
    local target="$1"
    local head="$2"
    local files="$3"
    local format="$4"

    if test "${files}"
    then
        files="-- ${files}"
    fi

    if test "${format}"
    then
        format="--format=${format}"
    fi

    local changes="$(git log $format "${target}"..."${head}" ${files})"
    if test -z "${changes}"
    then
        changes="$(git log ${format} "${head}^..${head}" ${files})"
    fi

    echo "${changes}"
}

get_tasks_files() {
  local target="$1"
  local head="$2"

  local files=$(get_diff "${target}" "${head}" | \
    grep -E "${LAB_FILES_REGEXP_PATTERN}")

  echo "${files}"
}

get_no_tasks_files() {
  local target="$1"
  local head="$2"

  local files=$(get_diff "${target}" "${head}" | \
    grep -v -E "${LAB_FILES_REGEXP_PATTERN}")

  echo "${files}"
}

get_students() {
  local target="$1"
  local head="$2"

  local student=$(get_diff "${target}" "${head}" | \
    grep -oE "^${STUDENT_REGEXP_PATTERN}" | sort | uniq)

  echo "${student}"
}

get_tasks() {
  local target="$1"
  local head="$2"
  local student_name="$3"

  local tasks=$(get_tasks_files "${target}" "${head}" | \
    grep -E "^${student_name}" | grep -oE "${TASK_REGEXP_PATTERN}" | \
    sort | uniq)

  echo "${tasks}"
}