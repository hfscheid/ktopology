#!/usr/bin/env bash
kubectl apply -f ktsim/nwsim/manifests -f ktsim/fluxgen
kubectl apply -f ktmonitor/
