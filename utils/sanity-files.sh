#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"

  cat >&1 <<EOF
Affected files lists:
$(get_diff "$HEAD" | awk '{print "\t- " $1}')
EOF

no_lab_files=$(get_diff "$HEAD" | grep -v -E "$LAB_FILES_REGEXP_PATTERN")

for file in $no_lab_files; do
    users=$(get_log "$HEAD" "$file" '%aE')
    for user in $users; do
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
if test -z "${lab_files}"; then
  printf 'Warning: Lab files is not represented in this commit...\n'
fi