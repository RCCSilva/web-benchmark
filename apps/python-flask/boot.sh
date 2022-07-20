#!/bin/bash

set -e

source venv/bin/activate
exec gunicorn --workers=4 -b :3000 app:app
