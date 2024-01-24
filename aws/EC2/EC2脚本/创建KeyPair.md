## 描述: 创建EC2 Keypair



```shell
#!/usr/bin/env bash
set -x

DEFAULT_REGION=us-east-1
DEFAULT_KEYNAME=mykey

function usage() {
  echo "Usage:

./createKeypair.sh -p <profile> [-r <region>] -k <keypair name>
Example:

  ./createKeypair.sh -p admin -r us-east-1 -k mykey
"
}

while getopts "p:r:k:h" opt; do
  case "$opt" in
  p) PROFILE="$OPTARG" ;;
  r) REGION="$OPTARG" 
     echo 1;;
  k) KeyName="$OPTARG" ;;
  [?]) usage 
     exit 0;;
  esac
done

# Check the parameters and configure the necessary information
function checkParameters(){
    if [ -z "$PROFILE" ]; then
        echo "No Profile!!!"
        usage
        exit 1
    fi
    if [ -z "$REGION" ]; then
        echo "No Region parameters"
        REGION=${DEFAULT_REGION}
        echo "Set default region: ${REGION}"
    fi
    if [ -z "$KeyName" ]; then
        echo "Checking keypair.."
        KeyName=${DEFAULT_KEYNAME}
        echo "Set default keypair name: ${KeyName}"
    fi
}

checkParameters

aws ec2 create-key-pair \
    --profile ${PROFILE} \
    --region ${REGION} \
    --key-name ${KeyName} \
    --query 'KeyMaterial' \
    --output text > ${KeyName}.pem
```

