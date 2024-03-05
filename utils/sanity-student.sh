#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

students=$(get_students "${HEAD}")

if test -z "${students}"; then
  printf "Changes for at least one student must be submitted in commits!\n" >&2
  exit 0
fi

if test $(echo "${students}" | wc -l) -ne 1; then
  cat >&2 <<EOF
The work of the following students has been submitted for review:
$(echo "${students}" | awk '{print "\t- " $1}')
However, commits should only contain changes for one student!
EOF
  exit 1
else
  printf "Lab works by %s are accepted for automatic review!\n" $(echo "${students}" | awk '{print $1}')
fi

