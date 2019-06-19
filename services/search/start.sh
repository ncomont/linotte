#!/bin/bash

export LINOTTE_SEARCH_SERVICE_ENDPOINT=":50055"
export LINOTTE_SEARCH_ELSATIC_HOST="localhost"
export LINOTTE_SEARCH_ELASTIC_PORT="9200"
export LINOTTE_SEARCH_ELASTIC_INDEX="taxref"

make && go build && ./search
