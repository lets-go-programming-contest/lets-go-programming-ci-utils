#!/bin/bash

if [[ -z "$TARGET" ]]; then
    echo "Target not set. Set TARGET env." >&2
    exit 1
fi

if [[ -z "$HEAD" ]]; then 
    echo "Head not set. Set HEAD env." >&2
    exit 1
fi 

if [[ ! -f "$MAINTAINERS_FILE" ]]; then
    echo "Maintainers file not found. Set MAINTAINERS_FILE env." >&2
    exit 1
fi

if [[ -z "$FILES_FOR_CHECK" ]]; then
    echo "Files for check is empty." 
    exit 0
fi

maintainers=$(cat "$MAINTAINERS_FILE")
commits=$(git log --format="%h" "$TARGET"..."$HEAD")
if [[ -z "$commits" ]]; then
    commits=$(git log --format="%h" "$HEAD^"..."$HEAD")
fi

is_affected=0
for commit in $commits; do
    [[ -z "$commit" ]] && continue

    author=$(git show -s --format="%ae" $commit)
    if echo $maintainers | grep -qw "$author"; then
        continue
    fi  
    
    files=$(git show --name-only --format="" $commit)
    for file in $files; do
        if echo $FILES_FOR_CHECK | grep -qw "$file"; then
            echo "File $file affected by $author in $commit"
            is_affected=1
        fi
    done
done

if [[ $is_affected -ne 0 ]]; then
    exit 1
fi


