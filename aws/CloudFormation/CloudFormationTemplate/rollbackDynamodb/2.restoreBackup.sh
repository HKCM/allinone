#!/usr/bin/env bash
set -e

source ./0.env.sh
# PROFILE=eu
# REGION=eu-west-1

# SourceTableName, Name of the source table that is being restored.
# SourceTableName=Music-rollback

# TargetTableName, The name of the new table to which it must be restored to.
# TargetTableName=Music-restore

# restore-date-time need UTC timezone
# Use the shell to execute the following command to get the current UTC time:
# date -u "+%Y-%m-%dT%H:%M:%SZ"
# 2020-08-19T09:03:21Z
# RestoreDateTime=2020-08-19T09:03:21Z

read -r -p "Are You Sure to restore table ${SourceTableName} from ${RestoreDateTime} time point ? [Y/n] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Yes"
		;;

    [nN][oO]|[nN])
		echo "No"
       	;;

    *)
		echo "Invalid input..."
		exit 1
		;;
esac
# restore
echo "Start restore table ${SourceTableName} from ${RestoreDateTime} time point..."
tableArn=$(aws dynamodb restore-table-to-point-in-time \
    --profile ${PROFILE} \
    --region ${REGION} \
    --source-table-name ${SourceTableName} \
    --target-table-name ${TargetTableName} \
    --restore-date-time ${RestoreDateTime} \
    --query "TableDescription.TableArn" \
    --output text)

# --local-secondary-index-override \
#         IndexName=AlbumTitleSongTitle,KeySchema=["{AttributeName=AlbumTitle,KeyType=HASH},{AttributeName=SongTitle,KeyType=RANGE}"],Projection="{ProjectionType=ALL}" \
    
# wait table-exists this command looks will wait until table available, so hardcode sleep 60 second
# aws dynamodb wait table-exists \
#     --profile ${PROFILE} \
#     --region ${REGION} \
#     --table-name ${TargetTableName}

echo "Wait table exists..."
sleep 60

echo "Tag restore dynamodb..."
# add tags
aws dynamodb tag-resource \
    --profile ${PROFILE} \
    --region ${REGION} \
    --resource-arn ${tableArn} \
    --tags \
        Key=Department,Value=dev \
        Key=Environment,Value=int-xmn \
        Key=Name,Value=Music-test-restore \
        Key=Team,Value=int-xmn \
        Key=RestoreTimePoint,Value=${RestoreDateTime} \
        Key=StartRestoreActionTime,Value=$(date -u "+%Y-%m-%dT%H:%M:%SZ")

echo "Everything is good..."
