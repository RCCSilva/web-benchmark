#!/bin/bash

set -e

source venv/bin/activate
exec gunicorn -b :3000 app:app
