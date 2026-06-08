#!/bin/bash

set -e

BASE_URL=http://localhost:8080

echo "Testing health"

curl -f \
  "$BASE_URL/health"

echo

echo "Testing ready"

curl -f \
  "$BASE_URL/ready"

echo

echo "Smoke tests passed"