#!/usr/bin/bash 

aws ec2 start-instances --instance-ids i-1234567890abcdef0

aws ec2 wait instance-status-ok \
    --instance-ids i-1234567890abcdef0

aws ec2 stop-instances --instance-ids i-1234567890abcdef0

aws ec2 wait instance-stopped \
    --instance-ids i-1234567890abcdef0