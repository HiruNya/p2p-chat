#!/bin/bash

# This will only build and deploy them to the AWS instances
# It will not run them

echo "cd src; go build; exit" | docker run -i -v ~/Documents/chat:/go/src golang

while read -p "AWS Instance ID (e.g. ec2-54-252-218-141): " AWS_ID; do
  scp -i "sylo-assesment.pem" ./chat "ubuntu@${AWS_ID}.ap-southeast-2.compute.amazonaws.com:/home/ubuntu/chat"
done
