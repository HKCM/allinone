#!/usr/bin/env python3
import boto3 

def getNon_CompliantSG(ec2_client):
    response = ec2_client.describe_security_groups(
        Filters=[
            {
                'Name': 'ip-permission.cidr',
                'Values': [
                    '0.0.0.0/0',
                ]
            },
        ],
        MaxResults=100
    )

    sgId = []

    # # Only print SecurityGroup ID 
    for i in range(len(response['SecurityGroups'])):
        sgId.append(response['SecurityGroups'][i]['GroupId'])
   
    # # print more info
    # for SecurityGroup in response['SecurityGroups']:
    #     Non_Compliant = False
    #     for permission in SecurityGroup['IpPermissions']:
    #         for cidrIpGroup in permission['IpRanges']:
    #             if '0.0.0.0/0' in cidrIpGroup['CidrIp']:
    #                 if permission['IpProtocol'] == '-1':
    #                     Non_Compliant = True
    #                     print("SecurityGroup: " + SecurityGroup['GroupId'] + " is NON_COMPLIANT, it allow all traffic")
    #                     break
    #                 if not ((permission['FromPort'] == 80 and permission['ToPort'] == 80) or (permission['FromPort'] == 443 and permission['ToPort'] == 443)):
    #                     Non_Compliant = True
    #                     print("SecurityGroup: " + SecurityGroup['GroupId'] + " is NON_COMPLIANT, it allow port from " + str(permission['FromPort']) + " to " + str(permission['ToPort']))
    #                     break
    #     if Non_Compliant:    
    #         sgId.append(SecurityGroup['GroupId'])


    return sgId



if __name__ == '__main__':
    session = boto3.Session(profile_name='hkc')
    ec2_client = session.client('ec2')
    sg = getNon_CompliantSG(ec2_client)
    print()

