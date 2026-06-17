#!/bin/sh

set -e

echo "Starting sync-worker..."
./sync-worker &

echo "Starting API..."
exec ./api
