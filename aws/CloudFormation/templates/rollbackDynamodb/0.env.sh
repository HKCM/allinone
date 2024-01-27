#!/usr/bin/env bash

StackName=Music-test
PROFILE=eu
REGION=eu-west-1



# For 2.restoreBackup.sh===============================================
# SourceTableName, Name of the source table that is being restored.
SourceTableName=${StackName}

# TargetTableName, The name of the new table to which it must be restored to.
TargetTableName=${StackName}-restore

# restore-date-time need UTC timezone
# Use the shell to execute the following command to get the current UTC time:
# date -u "+%Y-%m-%dT%H:%M:%SZ"
# 2020-08-19T09:03:21Z
RestoreDateTime=2020-08-19T09:03:21Z

# For 3.importResource==================================================
LogicalResourceId=restoreTable      # This variable needs to be the same as the logical ID in the template. file://dynamoRollback.yml       # The table name
ChangeSetName=ImportChangeSet
