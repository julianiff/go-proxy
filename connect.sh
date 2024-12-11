#!/bin/bash

# Array of API endpoints
declare -a endpoints=(
    "/api/users"
    "/api/products"
    "/api/orders"
    "/api/auth"
    "/api/payments"
    "/api/inventory"
    "/api/shipping"
    "/api/notifications"
    "/api/analytics"
    "/api/admin"
)

# Initialize counter
count=0
total_requests=1000

# Function to send request to proxy
send_request() {
    local endpoint=$1
    echo "eel $endpoint test another one here also" | nc localhost 3000
    count=$((count + 1))
    echo "Request $count of $total_requests: Sent to $endpoint"
    sleep 0.001
}

send_wrong_request() {
    local endpoint=$1
    echo "kwel $endpoint thats just to be wrong, lolz" | nc localhost 3000
    count=$((count + 1))
    echo "Request $count of $total_requests: FAAAKE SENT to $endpoint"
    sleep 0.001
}

# Outer loop for 1000 iterations
for ((i=1; i<=100; i++)); do
    # Inner loop through endpoints
    for endpoint in "${endpoints[@]}"; do
        if [ $((i % 10)) -eq 0 ]; then
            send_wrong_request "$endpoint"
            else
            send_request "$endpoint"
        fi
        # Break if we've reached 1000 requests
        if [ $count -eq $total_requests ]; then
            echo "Completed $total_requests requests"
            exit 0
        fi
    done
done
