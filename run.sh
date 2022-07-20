#!/bin/bash

set -e 

eval $(minikube docker-env)
exec npm run start
