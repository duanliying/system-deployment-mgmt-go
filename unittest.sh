#!/bin/bash

export GOPATH=$PWD

go get github.com/golang/mock/gomock

pkg_list=("api" "commons/errors" "commons/logger" "commons/url" "db" "db/mongo" "manager/agent" "manager/group" "messenger")

count=0
for pkg in "${pkg_list[@]}"; do
    go test -c -v -gcflags "-N -l" $pkg
    go test -coverprofile=$count.cover.out $pkg
    if [ $? -ne 0 ]; then
	echo "Unittest is failed."
	rm *.out *.test
	exit 1
    fi
    count=$count.0
done

echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out

go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverall.html

rm *.out *.test
