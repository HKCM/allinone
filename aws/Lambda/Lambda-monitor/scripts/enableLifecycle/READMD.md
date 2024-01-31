## Introduce

This script will enable lifecycle for all buckets on the AWS account unless they have a specific tag or the bucket itself already has lifecycle configured.

The specific tag: **EnableLifecycle:false**


## How to execute

!!!! This script is dangerous !!!! 

Go to the project's root folder
```sh
# go run script/enableLifecycle/enableLifecycle.go --profile <profile-name>
# example
# go run script/enableLifecycle/enableLifecycle.go --profile test-account
```