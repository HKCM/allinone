#!/bin/bash

aclids=(acl-xxxxx acl-xxxxxx)
region=us-east-1
profile=XXX

for aclid in ${aclids[@]}; do 
    echo "${aclid} start"
    aws ec2 --profile ${profile} --region ${region} create-network-acl-entry --network-acl-id ${aclid} \
        --egress --rule-number 1 --protocol tcp --port-range From=9050,To=9050 \
        --cidr-block 0.0.0.0/0 --rule-action deny

    list=(111.111.111.111 222.222.222.222 333.333.333.333 444.444.444.444) 

    # Don't forget change the rule number
    for i in $(seq 2 9)
    do
        aws ec2 --profile ${profile} --region ${region} create-network-acl-entry --network-acl-id ${aclid} \
            --egress --rule-number ${i} --protocol tcp --port-range From=80,To=80 \
            --cidr-block ${list[$i-2]}/32 --rule-action deny
    done
    echo "${aclid}: 80 port finished"

    for i in $(seq 10 17)
    do
        aws ec2 --profile ${profile} --region ${region} create-network-acl-entry --network-acl-id ${aclid} \
            --egress --rule-number ${i} --protocol tcp --port-range From=443,To=443 \
            --cidr-block ${list[$i-10]}/32 --rule-action deny
    done
    echo "${aclid}: 443 port finished"
done




 

