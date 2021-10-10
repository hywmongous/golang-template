#!/bin/sh

# Source: https://phoenixnap.com/kb/install-minikube-on-ubuntu
# Source: https://minikube.sigs.k8s.io/docs/start/

# 1. Install pre-requists
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install curl -y
sudo apt-get install apt-transport-https -y

# 2. Install VirtualBox Hypervisor
sudo apt install -y virtualbox virtualbox-ext-pack

# 3. Install Minikube
sudo curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# 4. Verify install
sudo minikube version

# 5. Install kubectl
sudo minikube kubectl -- get po -A