#!/usr/bin/env bash
docker build -t test-poddiscovery ../ -f Dockerfile
kubectl apply -f ./manifest.yaml
