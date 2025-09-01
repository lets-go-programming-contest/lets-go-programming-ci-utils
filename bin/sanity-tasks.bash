#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"

task_files=$(get_tasks_files "${TARGET_BRANCH}" "${HEAD}")

if [[ -z "${task_files}" ]]; then
  printf "Tasks are not represented.\n" >&2
  exit 1
fi

students=$(get_students "${TARGET_BRANCH}" "${HEAD}")
tasks=$(get_tasks "${TARGET_BRANCH}" "${HEAD}" "${student}")

if test $(echo "${tasks}" | wc -l) -ne 1; then
  printf "The next tasks has been presented for review:\n" >&2
  printf -- "\t- %s\n" ${tasks} >&2
  printf "However, commits should only contain changes for one task!\n" >&2

  exit 1
fi

printf "Task ${tasks} are accepted for automatic review!\n"
exit 0

