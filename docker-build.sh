#!/bin/bash

sh docker-clean.sh


# echo "Building backoffice ..."
# pushd backoffice > /dev/null
# mkdir ../docker/backoffice/www
# yarn
# NODE_ENV='devel' yarn run build
# cp -rip dist/* ../docker/backoffice/www/
# popd > /dev/null

echo "Building job ..."
pushd services/job > /dev/null
make
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../docker/job/job
popd > /dev/null

echo "Building search ..."
pushd services/search > /dev/null
make
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../docker/search/search
popd > /dev/null

echo "Building taxref ..."
pushd services/taxref > /dev/null
make
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../docker/taxref/taxref
popd > /dev/null

echo "Building user ..."
pushd services/user > /dev/null
make
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../docker/user/user
popd > /dev/null

echo "Building api ..."
pushd services/api > /dev/null
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../docker/api/api
popd > /dev/null

echo "Building job runner ..."
pushd job_runner > /dev/null
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../docker/runner/job_runner
popd > /dev/null

echo "Generating keys ..."
openssl genrsa -out docker/user/app.rsa 1024
openssl rsa -in docker/user/app.rsa -pubout > docker/user/app.rsa.pub


echo "Docker package ready !"
