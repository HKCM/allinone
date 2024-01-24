限制IAM的权限
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:UpdateUser",
                "iam:UpdateLoginProfile",
                "iam:PutUserPolicy",
                "iam:PutUserPermissionsBoundary",
                "iam:DetachUserPolicy",
                "iam:DeleteUserPolicy",
                "iam:DeleteUserPermissionsBoundary",
                "iam:DeleteServiceSpecificCredential",
                "iam:DeactivateMFADevice",
                "iam:CreateUser",
                "iam:CreateServiceSpecificCredential",
                "iam:CreateLoginProfile",
                "iam:AttachUserPolicy"
            ],
            "Resource": "arn:aws:iam::${AWS::AccountId}:user/*"
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:UpdateAccessKey",
                "iam:DeleteAccessKey",
                "iam:CreateAccessKey"
            ],
            "NotResource": [
                "arn:aws:iam::${AWS::AccountId}:user/special_user"
            ]
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "iam:UpdateAccessKey",
                "iam:DeleteAccessKey",
                "iam:CreateAccessKey"
            ],
            "Resource": [
                "arn:aws:iam::${AWS::AccountId}:user/special_user"
            ]
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": "iam:AddUserToGroup",
            "Resource": "arn:aws:iam::${AWS::AccountId}:group/*"
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:Update*",
                "iam:Remove*",
                "iam:Put*",
                "iam:Pass*",
                "iam:Detach*",
                "iam:Delete*",
                "iam:Create*",
                "iam:Attach*",
                "iam:Add*"
            ],
            "Resource": "arn:aws:iam::${AWS::AccountId}:role/okta*"
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:Update*",
                "iam:Resync*",
                "iam:Reset*",
                "iam:Put*",
                "iam:Enable*",
                "iam:Detach*",
                "iam:Delete*",
                "iam:*User*",
                "iam:*Policy*",
                "iam:*Password*",
                "iam:*MFA*",
                "iam:*LoginProfile",
                "iam:*Key*",
                "iam:*Credential*"
            ],
            "Resource": [
                "arn:aws:iam::${AWS::AccountId}:user/admin"
                "arn:aws:iam::${AWS::AccountId}:user/OktaSSO"
            ]
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:Update*",
                "iam:Remove*",
                "iam:Put*",
                "iam:Detach*",
                "iam:Delete*",
                "iam:Create*",
                "iam:Attach*",
                "iam:Add*"
            ],
            "Resource": "arn:aws:iam::${AWS::AccountId}:group/EmergencyAdministrators"
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:Update*",
                "iam:Put*",
                "iam:Detach*",
                "iam:Delete*",
                "iam:Create*",
                "iam:Attach*"
            ],
            "Resource": [
                "arn:aws:iam::${AWS::AccountId}:policy/okta*",
                "arn:aws:iam::${AWS::AccountId}:policy/OKTA*"
            ]
        },
        {
            "Sid": "",
            "Effect": "Deny",
            "Action": [
                "iam:Update*",
                "iam:Delete*",
                "iam:Create*"
            ],
            "Resource": "arn:aws:iam::${AWS::AccountId}:saml-provider/OktaSSO"
        }
    ]
}
```

写日志的权限
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
```

限制region
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "*",
            "Effect": "Deny",
            "Resource": "*",
            "Condition": {
                "ForAnyValue:StringNotEquals": {
                    "aws:RequestedRegion": [
                        "us-east-1",
                        "ap-southeast-1",
                        "eu-west-1"
                    ]
                }
            }
        }
    ]
}
```