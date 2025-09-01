#!/bin/bash

if [[ -z "$STUDENT_REGEXP_PATTERN" ]]; then
    echo "Regexp for task number not set." >&2 
    exit 1
fi

if [[ -z "$TASK_REGEXP_PATTERN" ]]; then
    echo "Regexp for task number not set." >&2 
    exit 1
fi

if [[ -z "$LAB_FILES" ]]; then
    echo ""
    exit 0
fi

echo "$(echo "$LAB_FILES" | grep -E "^$STUDENT_REGEXP_PATTERN" | grep -o -E "$TASK_REGEXP_PATTERN" | sort | uniq)"
