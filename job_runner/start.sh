#!/bin/bash

export LINOTTE_JOB_RUNNER_RABBIT_ENDPOINT="amqp://guest:guest@localhost:5672/"
export LINOTTE_JOB_RUNNER_RABBIT_TASK_QUEUE_ID="job_task_queue"
export LINOTTE_JOB_RUNNER_RABBIT_RESULT_QUEUE_ID="job_result_queue"
export LINOTTE_JOB_RUNNER_RABBIT_STATUS_QUEUE_ID="job_status_queue"
export LINOTTE_JOB_RUNNER_INDEXER_BATCH_SIZE=250
export LINOTTE_JOB_RUNNER_TAXREF_SERVICE_ENDPOINT=":50052"
export LINOTTE_JOB_RUNNER_STORAGE_PATH="../docker/volumes/ingest"

go build && ./job_runner
