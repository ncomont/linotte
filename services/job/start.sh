#!/bin/bash

export LINOTTE_JOB_SERVICE_ENDPOINT=":50053"
export LINOTTE_JOB_DATABASE_HOST="localhost"
export LINOTTE_JOB_DATABASE_USER="linotte"
export LINOTTE_JOB_DATABASE_PASSWORD="l1n0tt3passwd"
export LINOTTE_JOB_DATABASE_NAME="linotte"
export LINOTTE_JOB_DATABASE_PORT=3306
export LINOTTE_JOB_DATABASE_VERBOSE_MODE="false"
export LINOTTE_JOB_RABBIT_ENDPOINT="amqp://guest:guest@localhost:5672/"
export LINOTTE_JOB_RABBIT_TASK_QUEUE_ID="job_task_queue"
export LINOTTE_JOB_RABBIT_RESULT_QUEUE_ID="job_result_queue"
export LINOTTE_JOB_RABBIT_STATUS_QUEUE_ID="job_status_queue"

make && go build && ./job
