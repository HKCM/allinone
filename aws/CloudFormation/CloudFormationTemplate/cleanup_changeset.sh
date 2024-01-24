# cleanup_changeset (Profile, Stack-name, [Region])
# It can clear the failed changeset in the stack, but need to change env/environment.yml, 
# add delete ChangeSet Action "DeleteChangeSet"
# Solve: An error occurred (LimitExceededException) when calling the CreateChangeSet operation: 
#       ChangeSet limit exceeded for stack ...
cleanup_changeset () {
    local profile=$1
    local region=$2
    local stackname=$3
    i=0
    echo "Cleaning up failed change sets"
    changesets=$(aws cloudformation list-change-sets \
        --profile ${profile} \
        --region ${region} \
        --stack-name ${stackname} --query 'Summaries[?Status==`FAILED`].ChangeSetId' --output text)
    echo changesets
    for changeset in $changesets; do
      ((i++))
      echo "${stackname}: deleting change set ${i}: ${changeset:0-36}"
      aws cloudformation delete-change-set \
        --profile ${profile} \
        --region ${region} \
        --stack-name ${stackname} \
        --change-set-name ${changeset}
    done
}

function usage() {
  echo "Usage:
./Mac/OneClickTools/cleanup_changeset -p <aws_profile> -r <region> -s <stack_name>

Example1:
It will delete all failed changeset 

./Mac/OneClickTools/cleanup_changeset -p int-developer -r ap-southeast-1 -s int-auth-interface-test-pre-infra
"
  exit 0
}

while getopts "p:r:s:" opt; do
  case "$opt" in
  p) profile="$OPTARG" ;;
  r) region="$OPTARG";;
  s) stackname="$OPTARG" ;;
  [?]) usage ;;
  esac
done

echo "profile: ${profile}"
echo "region: ${region}"
echo "stackname: ${stackname}"

if [ -z "${profile}" ] || [ -z "${region}" ] || [ -z "${stackname}" ]; then
    usage
fi

cleanup_changeset ${profile} ${region} ${stackname}