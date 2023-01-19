#!/usr/bin/env sh

FOLDERS=`echo "clientApi createroomsubscriber createusersubscriber roomservice userservice"`

for FOLDER in ${FOLDERS}; do
    cd ${FOLDER}
    go mod tidy
    go mod vendor
    printf "tidy and vendored %s\n" "${FOLDER}"
    cd ..
done