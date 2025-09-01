#!/bin/bash

if [[ -z ${STUDENT_REGEXP_PATTERN} ]]; then
    echo "Regexp for student names not set." >&2 
    exit 1
fi

if [[ -z "$LAB_FILES" ]]; then
    echo ""
    exit 0
fi

echo "$(echo "$LAB_FILES" | grep -oE "^${STUDENT_REGEXP_PATTERN}" | sort | uniq)"
