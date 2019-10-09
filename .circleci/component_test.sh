#!/bin/bash

set -xe

mkdir -p /tmp/bin
PATH=/tmp/bin:$PATH

# Install KinD
KIND_VERSION="0.5.1"
cd /tmp
curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v${KIND_VERSION}/kind-$(uname)-amd64
chmod +x ./kind
mv ./kind /tmp/bin/

# start a cluster
kind create cluster

# delete the cluster
kind delete cluster
