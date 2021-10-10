#!/bin/sh
sudo docker kill $(sudo docker ps -q)
