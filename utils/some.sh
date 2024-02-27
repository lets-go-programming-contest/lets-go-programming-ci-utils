#!/bin/bash

. "$TEST_DIR_UTILS/common.sh"


# Получаем список файлов, не относящихся к лабораторным
no_lab_files=$(get_diff "$HEAD" | grep -v -E "$LAB_FILES_REGEXP_PATTERN")

# Проверяем каждый файл
for file in $no_lab_files; do

    # Проверяем, что изменения сделаны пользователями из списка
    user_changes=$(git log --format="%aE" "$file")
    for user in $user_changes; do
        if ! grep -q -e "^$user$" MAINTAINERS; then
            echo "Пользователь $user не найден в файле MAINTAINERS."
            exit 1
        fi
    done
done

echo "Все изменения прошли проверку."
