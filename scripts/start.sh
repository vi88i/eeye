#!/bin/bash

docker ps --filter "name=eeye-db" --filter "status=running" | grep eeye-db || docker start eeye-db

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
GO_ENTRY_PATH="$PROJECT_ROOT/src/main.go"

go run $GO_ENTRY_PATH
