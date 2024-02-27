#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

no_lab_files=$(get_diff "$HEAD" | grep -v -E "$LAB_FILES_REGEXP_PATTERN")

for file in $no_lab_files; do
    user_changes=$(git log --format="%aE" "$file")
    for user in $user_changes; do
        if ! grep -q "^$user$" MAINTAINERS; then
          printf "%s affected by %s\n" "$file" "$user" >> logs/sanity-files-error-log.txt
        fi
    done
done

if [ -s "logs/sanity-files-error-log.txt" ]; then
    cat >&2 <<EOF
The following root files are affected:
$(cat logs/sanity-files-error-log.txt | awk '{print "\t- " $0}')
EOF
  exit 1
fi

lab_files=$(get_labs_files "$HEAD")
if test "${lab_files}"; then
  cat >&2 <<EOF
The following lab files are submitted for checking:
$(echo ${lab_files} | awk '{print "\t- " $1}')
EOF
  exit 0
else
  printf 'Warning: Lab files is not represented in this commit...'
fi