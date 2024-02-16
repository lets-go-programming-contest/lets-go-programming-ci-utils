#!/bin/bash

student=$(get_students "$HEAD" | head -n 1)
tasks=$(get_tasks "$HEAD" "$student")

for task in $tasks; do
  task_makefile="$TEST_DIR_COMMON/$task/Makefile"
  if [[ -f "$task_makefile" ]]; then
    if test $(grep -e 'custom:' "$task_makefile"); then
      make --no-print-directory -f "$task_makefile" custom
    else
      printf "No custom checks for %s\n" "$task"
    fi
  fi
done