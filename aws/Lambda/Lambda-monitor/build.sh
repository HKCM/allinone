#!/bin/bash

buildGo() {

    echo -e "Start build $1..."
    env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -tags lambda.norpc -o ./bootstrap main.go

    zip bootstrap.zip bootstrap
}

cd lambda

echo "********** Start Build Lambda **********"

for dir in $(ls); do
    cd $dir
    buildGo $dir
    echo -e "Build $dir finished..\n"
    cd ..
done

echo -e "********** Build Lambda Finished **********\n"



