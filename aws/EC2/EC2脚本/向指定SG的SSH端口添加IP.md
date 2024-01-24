

## 描述: 向指定SG的SSH端口添加IP

### Code
```shell
#!/usr/bin/env bash

# Example1:
# It will add the your current IP address to specified SG

# ./Mac/OneClickTools/check_connect.sh -p int-developer

# Example2:
# It will add the specified IP address and add describe to specified SG
# Valid descriptions are strings less than 256 characters from the following set:  a-zA-Z0-9. _-:/()#,@[]+=&;{}!$*

# ./Mac/OneClickTools/check_connect.sh -p int-developer -i 10.20.30.40 -d "home IP" 

REGION=eu-west-1
SECURITYGROUP=sg-056fb9a9bbd306400
USERNAME=$USER
CREATEDATE=`date +%Y-%m-%d`

set -e

function usage() {
  echo "Usage:

$0 -p <aws_profile> [-i Target IP] [-d Describe]

Example:

$0 -p int-developer -i 10.10.10.10 -d HomeIP

Example2:

It will add your current IP to SG

$0 -p int-developer
"
  exit 0
}

while getopts "p:i:d:h" opt; do
  case "$opt" in
  p) PROFILE="$OPTARG" ;;
  i) USERIP="$OPTARG" ;;
  d) DESCRIBE="$OPTARG" ;;
  h) usage ;;
  [?]) usage ;;
  esac
done

if [ -z "${PROFILE}" ]; then
    usage
fi

if [ -z "${USERIP}" ]; then
    echo -e "\nThere is no Target IP address, it will add your current IP to SG\n"
    USERIP=$(curl http://checkip.amazonaws.com)
    echo -e "\nYour current IP is ${USERIP}"
fi


# Valid descriptions are strings less than 256 characters from the following set:  a-zA-Z0-9. _-:/()#,@[]+=&;{}!$*
if [ -z "${DESCRIBE}" ]; then
    echo -e "\nThere is no Describe, it will add your USER Name to SG describe\n"
    DESCRIBE=$USERNAME
fi

IPS=$(aws ec2 describe-security-groups \
    --profile ${PROFILE} \
    --region ${REGION} \
    --group-ids ${SECURITYGROUP} \
    --query 'SecurityGroups[*].IpPermissions[*].IpRanges[*].CidrIp' \
    --output text)

if ! [[ ${IPS} == *${USERIP}* ]]; then
    aws ec2 authorize-security-group-ingress \
        --profile ${PROFILE} \
        --region ${REGION} \
        --group-id ${SECURITYGROUP} \
        --ip-permissions IpProtocol=tcp,FromPort=22,ToPort=22,IpRanges="[{CidrIp=${USERIP}/32,Description=\"${DESCRIBE}-${CREATEDATE}\"}]"
    if [ $? = 0 ]; then 
        echo "${USERIP}/32 add into security group successful"
    else
        echo "${USERIP}/32 add into security group failed. Please check"
    fi
else
    echo "${USERIP}/32 already exist in security group"
fi
```
