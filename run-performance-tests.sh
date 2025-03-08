#!/bin/bash

# Create performance test reports directory if it doesn't exist
mkdir -p performance-test-reports

# Define services and their endpoints to test - using separate arrays instead of associative array
# services=("account" "transaction")
# endpoints=("http://localhost:8080" "http://localhost:8081")
# test transaction service only for now
services=( "transaction")
endpoints=("http://localhost:8081")

# Function to check if k6 is installed
check_k6() {
    if ! command -v k6 &> /dev/null; then
        echo "k6 is not installed. Please install it first:"
        echo "For macOS: brew install k6"
        echo "For Ubuntu: sudo apt-get install k6"
        exit 1
    fi
}

# Function to wait for service to be ready
wait_for_service() {
    local service=$1
    local url=$2
    echo "Waiting for $service service to be ready..."
    
    for i in {1..30}; do
        if curl -s "$url/health-check" &> /dev/null; then
            echo "$service service is ready!"
            return 0
        fi
        echo "Waiting... ($i/30)"
        sleep 1
    done
    
    echo "$service service did not become ready in time"
    return 1
}

# Main execution
echo "Starting performance tests..."

# Check if k6 is installed
check_k6

# Wait for services to be ready
for i in "${!services[@]}"; do
    service="${services[$i]}"
    endpoint="${endpoints[$i]}"
    if ! wait_for_service "$service" "$endpoint"; then
        echo "Error: Unable to connect to $service service"
        exit 1
    fi
done

# Create timestamp for this test run
timestamp=$(date +"%Y%m%d_%H%M%S")
report_dir="performance-test-reports/$timestamp"
mkdir -p "$report_dir"

# Run tests for account service endpoints
# Run tests for each service
services=("transaction")
for service in "${services[@]}"; do
    echo "Testing $service Service endpoints..."
    test_file="performance-tests/${service}-service.js"
    if [ -f "$test_file" ]; then
        k6 run \
            --out json="$report_dir/${service}-service.json" \
            --summary-export="$report_dir/${service}-service-summary.json" \
            "$test_file"
    else
        echo "Warning: ${service}-service.js test file not found"
    fi
done

# Generate summary report
echo "Generating summary report..."
cat > "$report_dir/summary.md" << EOF
# Performance Test Results - $(date)

## Test Environment
- Date: $(date)
- Duration: 10 seconds per endpoint
- Virtual Users: 100 concurrent users
- Environment: Docker containers on local machine

## Results Summary
Results are available in JSON format in this directory:
- account-service.json
- transaction-service.json

For detailed metrics, please check the JSON files.
EOF

echo "Performance tests completed!"
echo "Reports available in: $report_dir"
