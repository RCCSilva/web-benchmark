#!/bin/bash

set -e

source venv/bin/activate
exec uvicorn app:app --workers 4 --host 0.0.0.0 --port 3000
