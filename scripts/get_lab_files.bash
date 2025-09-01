#!/bin/bash

if [[ -z "$LAB_FILES_REGEXP_PATTERN" ]]; then
    echo "Regexp for lab files not set." >&2 
    exit 1
fi

if [[ -z "$DIFF" ]]; then
    echo "No files are represented. Can't calculate lab files." >&2
    exit 1
fi

echo "$(echo "$DIFF" | grep -E "$LAB_FILES_REGEXP_PATTERN")"