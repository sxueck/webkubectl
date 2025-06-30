#!/bin/bash

SHARED_DIR="/shared"
RETENTION_DAYS=7

if [ -d "$SHARED_DIR" ]; then
    find "$SHARED_DIR" -type f -mtime +$RETENTION_DAYS -delete
    find "$SHARED_DIR" -type d -empty -delete
fi 