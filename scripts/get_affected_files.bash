#!/bin/bash

if [[ -z "$TARGET" ]]; then
    echo "Target not set. Set TARGET env." >&2
    exit 1
fi

if [[ -z "$HEAD" ]]; then 
    echo "Head not set. Set HEAD env." >&2
    exit 1
fi 

changes=$(git diff --name-only -z "$TARGET"..."$HEAD" | tr '\000' '\n')
if [[ -z "${changes}" ]]; then 
    changes=$(git diff --name-only -z "$HEAD^".."$HEAD" | tr '\000' '\n')
fi

echo "$changes"