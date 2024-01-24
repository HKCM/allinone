

## 描述: 这是一个通过EBS的snapshot还原Instance的脚本

```shell
#!/bin/bash

# 这是一个通过EBS的snapshot还原的脚本
# 通过寻找具有指定tag的最新的snapshot创建volume
# 并将旧的volume从EC2上卸下,挂载新的volume,以达到还原的目的
# 也可以指定snapshot ID和instance ID
set -e

# source ./env.sh
# env.sh中包含默认的Name变量,用于指定Name

function usage() {
  echo "Usage:

./recovery.sh -p <profile> [-r <region>] [-s <snapshotID>] [-i <instanceID>]

Example:

  ./recovery.sh -p <profile>
"
  exit 0
}

while getopts "p:r:s:i:h" opt; do
  case "$opt" in
  p) PROFILE="$OPTARG" ;;
  r) REGION="$OPTARG" ;;
  s) SnapshotId="$OPTARG" ;;
  i) InstanceId="$OPTARG" ;;
  [?]) usage ;;
  esac
done

# Check the parameters and configure the necessary informations
function checkParameters(){
    if [ -z "$PROFILE" ]; then
        echo "No Profile!!!"
        usage
    fi
    if [ -z "$REGION" ]; then
        echo "No Region parameters"
        REGION=us-east-1
        echo "Set default region: ${REGION}"
    fi
    if [ -z "$SnapshotId" ]; then
        echo "No SnapshotId parameters"
        getLatestSnapshot
        echo "Get SnapshotId: ${SnapshotId}"
    fi
    if [ -z "$InstanceId" ]; then
        echo "No InstanceId parameters"
        getInstanceINFO
        echo "Get InstanceId: ${InstanceId}"
    else
        echo "InstanceId: ${InstanceId}"
        getInstanceINFOById
    fi
    echo "------------------------"
    echo ""
    echo "InstanceId: ${InstanceId}"
    echo "OldVolumeId: ${OldVolumeId}"
    echo "AvailabilityZone: ${AvailabilityZone}"
    echo "SnapshotId: ${SnapshotId}"

}

# Get Snapshot ID
function getLatestSnapshot()
{
    SnapshotId=$(aws ec2 describe-snapshots \
    --profile ${PROFILE} \
    --region ${REGION} \
    --filters "Name=tag:Name,Values=${Name}" \
    --query 'reverse(sort_by(Snapshots,&StartTime))[0].[SnapshotId]' \
    --output text)
}

# Get instance ID, OldVolumeId and AvailabilityZone
function getInstanceINFO()
{
    InstanceINFO=($(aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --filters "Name=tag:Name,Values=${Name}" \
    --query 'Reservations[*].Instances[*].[InstanceId,BlockDeviceMappings[*].Ebs.VolumeId,Placement.AvailabilityZone]' \
    --output text))

    InstanceId=${InstanceINFO[0]}
    AvailabilityZone=${InstanceINFO[1]}
    OldVolumeId=${InstanceINFO[2]}
}

# Get OldVolumeId and AvailabilityZone by instance ID
function getInstanceINFOById()
{
    InstanceINFO=($(aws ec2 describe-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${InstanceID}
    --query 'Reservations[*].Instances[*].[InstanceId,BlockDeviceMappings[*].Ebs.VolumeId,Placement.AvailabilityZone]' \
    --output text))

    InstanceId=${InstanceINFO[0]}
    AvailabilityZone=${InstanceINFO[1]}
    OldVolumeId=${InstanceINFO[2]}
}

checkParameters

# Create new volume from latest snapshot
echo "Creating new volume..."
NewVolumeId=$(aws ec2 create-volume \
    --profile ${PROFILE} \
    --region ${REGION} \
    --snapshot-id ${SnapshotId} \
    --availability-zone ${AvailabilityZone} \
    --query 'VolumeId' \
    --tag-specifications "ResourceType=volume,Tags=[ \
        {Key=Team,Value=${Team}},{Key=Department,Value=${Department}}, \
        {Key=Name,Value=${Name}},{Key=Environment,Value=${Environment}}]" \
    --output text)

echo "NewVolumeId: ${NewVolumeId}"

# Stop sentry server
echo "Stop sentry server..."
aws ec2 stop-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${InstanceId}

echo "Waiting instance server stop..."
aws ec2 wait instance-stopped \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${InstanceId}

echo "Waiting volume available..."
aws ec2 wait volume-available \
    --profile ${PROFILE} \
    --region ${REGION} \
    --volume-ids ${NewVolumeId}

# Deatch old volume from instance
echo "Deatch old volume..."
aws ec2 detach-volume \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-id ${InstanceId} \
    --volume-id ${OldVolumeId}

# attatch new volume from instance
echo "Attach new volume..."
aws ec2 attach-volume \
    --profile ${PROFILE} \
    --region ${REGION} \
    --volume-id ${NewVolumeId} \
    --instance-id ${InstanceId} \
    --device /dev/sda1

echo "Waiting volume in-use..."
aws ec2 wait volume-in-use \
    --profile ${PROFILE} \
    --region ${REGION} \
    --volume-ids ${NewVolumeId}

echo "Start instance server..."
aws ec2 start-instances \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${InstanceId}

echo "Waiting instance server status-ok..."
aws ec2 wait instance-status-ok \
    --profile ${PROFILE} \
    --region ${REGION} \
    --instance-ids ${InstanceId}

read -r -p "Do you want to delete old volume?(Default No) [Y/n] " input
echo 
case $input in
    [yY][eE][sS]|[yY])
        aws ec2 delete-volume \
            --profile ${PROFILE} \
            --region ${REGION} \
            --volume-id ${OldVolumeId}
        echo "Deleted old volume"
        ;;

    [nN][oO]|[nN])
        echo "Keep the old volume: ${OldVolumeId}"	       	
        ;;

    *)
        echo "Keep the old volume: ${OldVolumeId}"
        ;;
esac

echo "Everything is OK!!"


```
