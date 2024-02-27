#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

base_branch=$(git rev-list --max-parents=0 --date-order --reverse "$BASE_BRANCH" --)
head_branch=$(git rev-list --max-parents=0 --date-order --reverse "$HEAD" --)

if test $(echo "${head_branch}" | wc -l) -ne 1; then
  cat >&2 <<EOF
The branch $(git rev-parse --abbrev-ref ${HEAD}) has multiple initial commits:
$(echo "${head_branch}" | awk '{print "\t- " $1}')
Fix your branch structure before proceeding!
EOF
  exit 1
fi

if test $(echo "${BASE_BRANCH}" | wc -l) -ne 1; then
  cat >&2 <<EOF
The branch ${BASE_BRANCH} has multiple initial commits:
$(echo "${BASE_BRANCH}" | awk '{print "\t- " $1}')
Fix your branch structure before proceeding!
EOF
  exit 1
fi

base_branch_root=$(echo "${base_branch}" | head -1)
head_branch_root=$(echo "${head_branch}" | head -1)

if test "${base_branch_root}" != "${head_branch_root}"; then
  cat >&2 <<EOF
Initial commits for $(git rev-parse --abbrev-ref HEAD) and ${BASE_BRANCH} not equal:
$(printf "\t- %s\t%s\n"  $(git rev-parse --abbrev-ref HEAD) "${head_branch_root}")
$(printf "\t- %s\t%s\n" "${BASE_BRANCH}" "${base_branch_root}")
Fix your branch structure before proceeding!
EOF
  exit 1
else
  printf "No discrepancies were found between branch %s and branch %s!\n"  $(git rev-parse --abbrev-ref HEAD) "${BASE_BRANCH}"
fi