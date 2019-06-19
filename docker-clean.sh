#!/bin/bash

echo "Cleaning backoffice ..."
rm -rf docker/backoffice/www

echo "Cleaning old executable ..."
rm -f docker/api/api
rm -f docker/user/user
rm -f docker/job/job
rm -f docker/taxref/taxref
rm -f docker/search/search
rm -f docker/runner/job_runner

echo "Cleaning old keys ..."
rm -f docker/user/app.rsa*
