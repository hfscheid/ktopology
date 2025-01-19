#!/usr/bin/env bash
kubectl delete -f ktmonitor/
kubectl delete -f ktsim/nwsim/manifests -f ktsim/fluxgen
