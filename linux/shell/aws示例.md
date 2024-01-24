描述: 基于AWS EC2的Bash脚本示例


本文基于AWS EC2的CLI命令,实现EC2 instance type的转化.代码源于https://docs.amazonaws.cn/cli/latest/userguide/cli-services-ec2-instance-type-script.html

主要是学习Bash shell脚本的格式和规范.例如,函数声明使用,提示信息,测试以及说明文档

### change_ec2_instance_type.sh

该主脚本文件包含执行以下任务的 change_ec2_instance_type() 函数：
* 验证指定的 EC2 实例是否存在。
* 除非选择了 -f，否则函数会在停止实例之前向用户发出警告。
* 更改实例类型
* 如果您设置了 -r，请重启实例并确认实例正在运行

```shell
###############################################################################
#
# function change_ec2_instance_type
#
# This function changes the instance type of the specified Amazon EC2 instance.
# 
# Parameters:
#   -i   [string, mandatory] The instance ID of the instance whose type you
#                            want to change.
#   -t   [string, mandatory] The instance type to switch the instance to.
#   -f   [switch, optional]  If set, the function doesn't pause and ask before
#                            stopping the instance.
#   -r   [switch, optional]  If set, the function restarts the instance after
#                            changing the type.
#   -h   [switch, optional]  Displays this help.
#
# Example:
#      The following example converts the specified instance to type "t2.micro"
#      without pausing to ask permission. It automatically restarts the
#      instance after changing the type.
#
#      change-instance-type -i i-123456789012 -f -r
#
# Returns:
#      0 if successful
#      1 if it fails
###############################################################################

# Import the general_purpose functions.
source awsdocs_general.sh 

###############################################################################
# function instance-exists
#
# This function checks to see if the specified instance already exists. If it
# does, it sets two global parameters to return the running state and the 
# instance type.
#
# Input parameters:
#       $1 - The id of the instance to check
# 
# Returns:
#       0 if the instance already exists
#       1 if the instance doesn't exist
#     AND: 
#       Sets two global variables:
#            EXISTING_STATE - Contains the running/stopped state of the instance.
#            EXISTING_TYPE  - Contains the current type of the instance.
###############################################################################
function get_instance_info {

    # Declare local variables.
    local INSTANCE_ID RESPONSE
    
    # This function accepts a single parameter.
    INSTANCE_ID=$1

    # The following --filters parameter causes server-side filtering to limit
    # results to only the records that match the specified ID. The --query 
    # parameter causes CLI client-side filtering to include only the values of
    # the InstanceType and State.Code fields. 
    
    RESPONSE=$(aws ec2 describe-instances \
                   --query 'Reservations[*].Instances[*].[State.Name, InstanceType]' \
                   --filters Name=instance-id,Values="$INSTANCE_ID" \
                   --output text \
               )

    if [[ $? -ne 0 ]] || [[ -z "$RESPONSE" ]]; then
        # There was no response, so no such instance.
        return 1        # 1 in Bash script means error/false
    fi

    # If we got a response, the instance exists.
    # Retrieve the values of interest and set them as global variables.
    EXISTING_STATE=$(echo "$RESPONSE" | cut -f 1 )        
    EXISTING_TYPE=$(echo "$RESPONSE" | cut -f 2 )

    return 0        # 0 in Bash script means no error/true
}

######################################
#
#  See header at top of this file
#
######################################

function change_ec2_instance_type {

    function usage() (
        echo ""
        echo "This function changes the instance type of the specified instance."
        echo "Parameter:"
        echo "  -i  Specify the instance ID whose type you want to modify."
        echo "  -t  Specify the instance type to convert the instance to."
        echo "  -d  If the instance was originally running, this option prevents"
        echo "      automatically restarting the instance."
        echo "  -f  If the instance was originally running, this option prevents"
        echo "      the script from asking permission before stopping the instance."
        echo ""
    )

    local FORCE RESTART REQUESTED_TYPE INSTANCE_ID VERBOSE OPTION RESPONSE ANSWER 
    local OPTIND OPTARG # Required to use getopts command in a function.
    
    # Set default values.
    FORCE=false
    RESTART=false
    REQUESTED_TYPE=""
    INSTANCE_ID=""
    VERBOSE=false

    # Retrieve the calling parameters.
    while getopts "i:t:frvh" OPTION; do
        case "${OPTION}"
        in
            i)  INSTANCE_ID="${OPTARG}";;
            t)  REQUESTED_TYPE="${OPTARG}";;
            f)  FORCE=true;;
            r)  RESTART=true;;
            v)  VERBOSE=true;;
            h)  usage; return 0;;
            \?) echo "Invalid parameter"; usage; return 1;; 
        esac
    done

    if [[ -z "$INSTANCE_ID" ]]; then
        errecho "ERROR: You must provide an instance ID with the -i parameter."
        usage
        return 1
    fi

    if [[ -z "$REQUESTED_TYPE" ]]; then
        errecho "ERROR: You must provide an instance type with the -t parameter."
        usage
        return 1
    fi

    iecho "Parameters:\n"
    iecho "    Instance ID:   $INSTANCE_ID"
    iecho "    Requests type: $REQUESTED_TYPE"
    iecho "    Force stop:    $FORCE"
    iecho "    Restart:       $RESTART"
    iecho "    Verbose:       $VERBOSE"
    iecho ""
    
    # Check that the specified instance exists.
    iecho -n "Confirming that instance $INSTANCE_ID exists..."
    get_instance_info "$INSTANCE_ID"
    # If the instance doesn't exist, the function returns an error code <> 0.
    if [[ ${?} -ne 0 ]]; then
        errecho "ERROR: I can't find the instance \"$INSTANCE_ID\" in the current AWS account."
        return 1
    fi
    # Function get_instance_info has returned two global values:
    #   $EXISTING_TYPE  -- The instance type of the specified instance
    #   $EXISTING_STATE -- Whether the specified instance is running

    iecho "confirmed $INSTANCE_ID exists."
    iecho "      Current type: $EXISTING_TYPE"
    iecho "      Current state code: $EXISTING_STATE"

    # Are we trying to change the instance to the same type?
    if [[ "$EXISTING_TYPE" == "$REQUESTED_TYPE" ]]; then
        errecho "ERROR: Can't change instance type to the same type: $REQUESTED_TYPE."
        return 1
    fi

    # Check if the instance is currently running.
    # 16="running"
    if [[ "$EXISTING_STATE" == "running" ]]; then
        # If it is, we need to stop it.
        # Do we have permission to stop it? 
        # If -f (FORCE) was set, we do.
        # If not, we need to ask the user.
        if [[ $FORCE == false ]]; then
            while true; do
                echo ""
                echo "The instance $INSTANCE_ID is currently running. It must be stopped to change the type."
                read -r -p "ARE YOU SURE YOU WANT TO STOP THE INSTANCE? (Y or N) " ANSWER
                case $ANSWER in
                    [yY]* )
                        break;;
                    [nN]* )
                        echo "Aborting."
                        exit;;
                    * ) 
                        echo "Please answer Y or N."
                        ;;
                esac
            done
        else
            iecho "Forcing stop of instance without prompt because of -f."
        fi

        # stop the instance
        iecho -n "Attempting to stop instance $INSTANCE_ID..."
        RESPONSE=$( aws ec2 stop-instances \
                        --instance-ids "$INSTANCE_ID" )

        if [[ ${?} -ne 0 ]]; then
            echo "ERROR - AWS reports that it's unable to stop instance $INSTANCE_ID.\n$RESPONSE"
            return 1
        fi
        iecho "request accepted."
    else
        iecho "Instance is not in running state, so not requesting a stop."
    fi;

    # Wait until stopped.
    iecho "Waiting for $INSTANCE_ID to report 'stopped' state..."
    aws ec2 wait instance-stopped \
        --instance-ids "$INSTANCE_ID"
    if [[ ${?} -ne 0 ]]; then
        echo "\nERROR - AWS reports that Wait command failed.\n$RESPONSE"
        return 1
    fi 
    iecho "stopped.\n"

    # Change the type - command produces no output.
    iecho "Attempting to change type from $EXISTING_TYPE to $REQUESTED_TYPE..."
    RESPONSE=$(aws ec2 modify-instance-attribute \
                   --instance-id "$INSTANCE_ID" \
                   --instance-type "{\"Value\":\"$REQUESTED_TYPE\"}"
              )
    if [[ ${?} -ne 0 ]]; then
        errecho "ERROR - AWS reports that it's unable to change the instance type for instance $INSTANCE_ID from $EXISTING_TYPE to $REQUESTED_TYPE.\n$RESPONSE"
        return 1
    fi
    iecho "changed.\n"

    # Restart if asked
    if [[ "$RESTART" == "true" ]]; then

        iecho "Requesting to restart instance $INSTANCE_ID..."
        RESPONSE=$(aws ec2 start-instances \
                        --instance-ids "$INSTANCE_ID" \
                   )
        if [[ ${?} -ne 0 ]]; then
            errecho "ERROR - AWS reports that it's unable to restart instance $INSTANCE_ID.\n$RESPONSE"
            return 1
        fi
        iecho "started.\n"
        iecho "Waiting for instance $INSTANCE_ID to report 'running' state..."
        RESPONSE=$(aws ec2 wait instance-running \
                       --instance-ids "$INSTANCE_ID" )
        if [[ ${?} -ne 0 ]]; then
            errecho "ERROR - AWS reports that Wait command failed.\n$RESPONSE"
            return 1
        fi 
        
        iecho "running.\n"
        
    else
        iecho "Restart was not requested with -r.\n"
    fi
    
}
```

### test_change_ec2_instance_type.sh
脚本文件 change_ec2_instance_type_test.sh 会测试 change_ec2_instance_type 函数的各种代码路径。如果测试脚本中的所有步骤都正确完成，测试脚本将会删除它创建的所有资源。

您可以使用以下参数运行该测试脚本：

* -v -（开关）每个测试在运行时都会显示通过/失败状态。默认情况下，测试以静默方式运行，输出仅包含最终的总体通过/失败状态。
* -i -（开关）脚本会在每个测试后暂停，以便您能查看每个步骤的中间结果。让您能够使用 Amazon EC2 控制台检查实例的当前状态。您在系统提示符处按 Enter 后，脚本才会继续执行下一步。

```shell
#!/usr/bin/env bash

###############################################################################
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# This file is licensed under the Apache License, Version 2.0 (the "License").
#
# You may not use this file except in compliance with the License. A copy of
# the License is located at http://aws.amazon.com/apache2.0/.
#
# This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.
###############################################################################

source ./awsdocs_general.sh
source ./change_ec2_instance_type.sh

function usage {
    echo "This script tests the change_ec2_instance_type function by calling the"
    echo "function in a variety of ways and checking the output. It converts the"
    echo "instance between types t2.nano and t2.micro."
    echo ""
    echo "Parameters:"
    echo ""
    echo "  -v  Verbose. Shows diagnostic messages about the tests as they run."
    echo "  -i  Interactive. Pauses the script between steps so you can browse"
    echo "      the results in the AWS Management Console as they occur."
    echo ""
    echo "IMPORTANT: Running this test script creates an Amazon EC2 instance in"
    echo "   your Amazon account that can incur charges. It is your responsibility"
    echo "   to ensure that no resources are left in your account after the script"
    echo "   completes. If an error occurs during the operation of the script,this"
    echo "   instance can remain. Check for the instance and delete it manually to"
    echo "   avoid charges."
}

# Set default values.
INTERACTIVE=false

# Retrieve the calling parameters.
while getopts "ivh" OPTION; do
    case "${OPTION}"
    in
        i)
            INTERACTIVE=true
            VERBOSE=true
            ;;
        v)
            VERBOSE=true
            ;;
        h)
            usage
            return 0
            ;;
        \?)
            echo "Invalid parameter."
            usage
            return 1
            ;;
    esac
done

if [ "$VERBOSE" == "true" ]; then iecho "Tests running in verbose mode."; fi
if [ "$INTERACTIVE" == "true" ]; then iecho "Tests running in interactive mode."; fi

iecho ""
iecho "***************SETUP STEPS******************"
    # First, get the AMI ID for the one running the latest Amazon Linix 2.
    iecho -n "Retrieving the AMI ID for the latest Amazon Linux 2 AMI..."
    AMI_ID=$(aws ec2 describe-images \
                --owners 'amazon' \
                --filters 'Name=name,Values=amzn2-ami-hvm-2.0.????????-x86_64-gp2' 'Name=state,Values=available' \
                --query 'sort_by(Images, &CreationDate)[-1].[ImageId]' \
                --output 'text')
    if [ ${?} -ne 0 ]; then
        echo "ERROR: Unable to retrieve latest Amazon Linux 2 AMI ID: $AMI_ID"
        echo "Tests canceled."
        return 1
    else
        iecho "retrieved $AMI_ID."
    fi

    # Now launch the instance as a t2.micro and capture its instance ID.
    # All other instance settings are left to default.
    iecho -n "Requesting new Amazon EC2 instance of type t2.micro..."
    EC2_INSTANCE_ID=$(aws ec2 run-instances \
                        --image-id "$AMI_ID" \
                        --instance-type t2.micro \
                        --query 'Instances[0].InstanceId' \
                        --output text)
    if [ ${?} -ne 0 ]; then
        echo "ERROR: Unable to launch EC2 instance: $EC2_INSTANCE_ID"
        echo "Tests canceled."
        return 1
    else
        iecho "launched. ID:$EC2_INSTANCE_ID"
    fi

    iecho -n "Waiting for instance $EC2_INSTANCE_ID to exist..."
    aws ec2 wait instance-exists \
        --instance-id "$EC2_INSTANCE_ID"
    iecho "confirmed."

iecho "***************END OF SETUP*****************"
iecho ""


run_test "1. Missing mandatory -i parameter" \
         "change_ec2_instance_type" \
         1 \
         "ERROR: You must provide an instance id."

run_test "2. Missing mandatory -t parameter" \
         "change_ec2_instance_type -i abc" \
         1 \
         "ERROR: You must provide an instance type."

run_test "3. Using an instance ID that doesn't exist" \
         "change_ec2_instance_type -i abc -t t2.micro" \
         1 \
         "ERROR: I can't find the instance."

# Test changing to the same type. We can do this while the instance is starting up.
run_test "4. Trying to change to same type" \
         "change_ec2_instance_type -v -i $EC2_INSTANCE_ID -t t2.micro" \
         1 \
         "ERROR: Can't change instance type to the same type."

iecho -n "Waiting for instance $EC2_INSTANCE_ID to reach running state..."
RESPONSE=$(aws ec2 wait instance-running --instance-id "$EC2_INSTANCE_ID")
if [[ ${?} -ne 0 ]]; then
    errecho "\nERROR: AWS reports that the Wait command failed.\n$RESPONSE"
    return 1
fi 
iecho "running."

# Test changing to t2.micro without -r : should still be in stopped state.
run_test "5. Changing to type t2.nano without restart" \
         "change_ec2_instance_type -f -i $EC2_INSTANCE_ID -t t2.nano" \
         0

         # Validate result was "t2.nano" and that it's in "stopped" state.
         get_instance_info "$EC2_INSTANCE_ID"
         if [ "$EXISTING_TYPE" != "t2.nano" ]; then
             test_failed "Unable to validate change. Should be t2.nano. Found $EXISTING_TYPE."
         fi
         if [ "$EXISTING_STATE" != "stopped" ]; then
             test_failed "Unable to validate state. Should be stopped. Found $EXISTING_STATE."
         fi

# Test changing back to t2.micro with -r. Should now be in running state
run_test "6. Changing to type t2.micro with restart" \
         "change_ec2_instance_type -f -r -i $EC2_INSTANCE_ID -t t2.micro" \
         0

         # Validate result was "t2.micro" and that it's in "running" state.
         get_instance_info "$EC2_INSTANCE_ID"
         if [ "$EXISTING_TYPE" != "t2.micro" ]; then
             test_failed "Unable to validate change. Should be t2.micro. Found $EXISTING_TYPE."
         fi
         if [ "$EXISTING_STATE" != "running" ]; then
             test_failed "Unable to validate state. Should be running. Found $EXISTING_STATE."
         fi

iecho ""
iecho "*************TEAR DOWN STEPS****************"
    iecho -n "Requesting termination of instance $EC2_INSTANCE_ID..."
    # Delete and terminate the instance.
    RESPONSE=$(aws ec2 terminate-instances \
                        --instance-ids "$EC2_INSTANCE_ID"
              )
    if [ ${?} -ne 0 ]; then
        errecho "**** ERROR ****"
        errecho "AWS reported a failure to terminate EC2 instance: $EC2_INSTANCE_ID"
        errecho "You must terminate the instance using the AWS Management Console"
        errecho "or CLI commands. Failure to terminate the instance can result in"
        errecho "charges to your AWS account.\n"
    else
        iecho "request accepted."
    fi

    iecho -n "Waiting for instance $EC2_INSTANCE_ID to terminate..."
    aws ec2 wait instance-terminated \
        --instance-id "$EC2_INSTANCE_ID"
    iecho "confirmed."
    if [[ ${?} -ne 0 ]]; then
        errecho "ERROR - AWS reports that Wait command failed."
        errecho "You must ensure that the instance terminated successfully yourself using the"
        errecho "AWS Management Console or CLI commands. Failure to terminate the instance can" 
        errecho "result in charges to your AWS account.\n"
        return 1
    fi 


iecho "************END OF TEAR DOWN****************"
iecho ""

echo "Tests completed successfully."
```

### awsdocs_general.sh
脚本文件 awsdocs_general.sh 中包含了在 AWS CLI 的高级示例中使用的通用函数。
```shell
#!/usr/bin/env bash
###############################################################################
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# This file is licensed under the Apache License, Version 2.0 (the "License").
#
# You may not use this file except in compliance with the License. A copy of
# the License is located at http://aws.amazon.com/apache2.0/.
#
# This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied. See the License for the
# specific language governing permissions and limitations under the License.
###############################################################################
#
# This script contains general-purpose functions that are used throughout 
# the AWS Command Line Interface (AWS CLI) code examples that are maintained
# in the repo at https://github.com/awsdocs/aws-doc-sdk-examples.
#
# They are intended to abstract functionality that is required for the tests
# to work without cluttering up the code. The intent is to ensure that the 
# purpose of the code is clear.

# Set global defaults:
VERBOSE=false

###############################################################################
# function run_test
#
# This function is used to perform a command and compare its output to both
# the expected error code and the expected output string. If there isn't a 
# match, then the function invokes the test_failed function.
###############################################################################
function run_test {
    local DESCRIPTION COMMAND EXPECTED_ERR_CODE EXPECTED_OUTPUT RESPONSE
    
    DESCRIPTION="$1"
    COMMAND="$2"
    EXPECTED_ERR_CODE="$3"
    if [[ -z "$4" ]]; then EXPECTED_OUTPUT="$4"; else EXPECTED_OUTPUT=""; fi
    
    iecho -n "Running test: $DESCRIPTION..."
    RESPONSE="$($COMMAND)"
    ERR="${?}"

    # Check to see if we got the expected error code.
    if [[ "$EXPECTED_ERR_CODE" -ne "$ERR" ]]; then
        test_failed "The test \"$DESCRIPTION\" returned an unexpected error code: $ERR"
    fi

    # Now check the error message, if we provided other than "".
    if [[ -n "$EXPECTED_OUTPUT" ]]; then
        MATCH=$(echo "$RESPONSE" | grep "$EXPECTED_OUTPUT")
        # If there was no match (it's an empty string), then fail.
        if [[ -z "$MATCH" ]]; then
            test_failed "The test \"$DESCRIPTION\" returned an unexpected output: $RESPONSE"
        fi
    fi
    
    iecho "OK"
    ipause
}

###############################################################################
# function test_failed
#
# This function is used to terminate a failed test and to warn the customer
# about possible undeleted resources that could incur costs to their account.
###############################################################################

function test_failed {
    
    errecho ""
    errecho "===TEST FAILED==="
    errecho "$@"
    errecho ""
    errecho "    One or more of the tests failed to complete successfully. This means that"
    errecho "    any tests after the one that failed didn't run and might have left resources"
    errecho "    still active in your account."
    errecho ""
    errecho "IMPORTANT:"
    errecho "    Resources created by this script can incur charges to your AWS account. If the"
    errecho "    script didn't complete successfully, then you must review and manually delete"
    errecho "    any resources created by this script that were not automatically removed."
    errecho ""
    exit 1 
}


###############################################################################
# function errecho
#
# This function outputs everything sent to it to STDERR (standard error output).
###############################################################################
function errecho {
    printf "%s\n" "$*" 2>&1
}

###############################################################################
# function iecho
#
# This function enables the script to display the specified text only if 
# the global variable $VERBOSE is set to true.
###############################################################################
function iecho {
    if [[ $VERBOSE == true ]]; then
        echo "$@"
    fi
}

###############################################################################
# function ipause
#
# This function enables the script to pause after each command if interactive
# mode is set (by including -i on the script invocation command).
###############################################################################
function ipause {
    if [[ $INTERACTIVE == true ]]; then
        read -r -p "Press ENTER to continue..."
    fi
}

# Initialize the shell's RANDOM variable.
RANDOM=$$
###############################################################################
# function generate_random_name
#
# This function generates a random file name with using the specified root,
# followed by 4 groups that each have 4 digits.
# The default root name is "test".
###############################################################################
function generate_random_name {

    ROOTNAME="test"
    if [[ -n $1 ]]; then
        ROOTNAME=$1
    fi

    # Initialize the FILENAME variable
    FILENAME="$ROOTNAME"
    # Configure random number generator to issue numbers between 1000 and 9999, 
    # inclusive.
    DIFF=$((9999-1000+1))

    for _ in {1..4}
    do
        rnd=$(($((RANDOM%DIFF))+X))
        # Make sure that the number is 4 digits long.
        while [ "${#rnd}" -lt 4 ]; do rnd="0$rnd"; done
        FILENAME+="-$rnd"
    done
    echo $FILENAME
}
```