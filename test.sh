#!/bin/bash

# Define the URL
url="https://localhost:8090/api/endpoint"

# Number of parallel requests
num_requests=1000

# Use seq to generate a sequence of numbers from 1 to num_requests
seq 1 $num_requests | xargs -n 1 -P 10 -I {} sh -c "curl -s $url && echo 'Request {} completed'"
