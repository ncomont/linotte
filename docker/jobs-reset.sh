#!/bin/bash

echo "Connecting to th db container ..."

docker-compose -f docker-compose-dev.yml exec db /bin/bash -c "mysql -uroot -pl1n0tt3passwd linotte -e \"truncate table job_results;truncate table job_reports;update jobs set status = 'NEW';\""

docker-compose -f docker-compose-dev.yml exec rabbitmq /bin/bash -c "rabbitmqctl purge_queue job_result_queue && rabbitmqctl purge_queue job_status_queue && rabbitmqctl purge_queue job_task_queue"

docker-compose -f docker-compose-dev.yml restart runner

echo "Jobs stuff resetted !"
