#!/bin/bash

source "${TEST_DIR_UTILS}/bin/setup.bash"
source "${TEST_DIR_UTILS}/includes/lib/git.bash"

students=$(get_students "${TARGET_BRANCH}" "${HEAD}")

if test -z "${students}"; then
  printf "Changes for at least one student must be submitted in commits!\n" >&2
  exit 1
fi

if test $(echo "${students}" | wc -l) -ne 1; then
  printf "The work of the following students has been presented for review:\n" >&2
  printf -- "\t- %s\n" ${students} >&2
  printf "However, commits should only contain changes for one student!\n" >&2
  
  exit 1
fi

printf "Lab works by ${students} are accepted for automatic review!\n"
exit 0

