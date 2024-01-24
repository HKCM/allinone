### AWS CLI JMESPath查询技巧

https://opensourceconnections.com/blog/2015/07/27/advanced-aws-cli-jmespath-query/

```
$ aws ec2 describe-images --owner amazon --output text | grep -c ami
977
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[].[ImageId,Name]' --output text | grep -m5 "ami-"
ami-0048c968    .NET Beanstalk Cfn Container v2.0.2.1 on Windows 2012
ami-005daf69    ElasticBeanstalk-Tomcat6-64bit-20110322-2041
ami-0078da69    amzn-ami-pv-2012.03.2.x86_64-s3
ami-00c17768    aws-elasticbeanstalk-amzn-2014.09.0.x86_64-php55-gpu-201409291824
ami-013aca6a    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121

[Errno 32] Broken pipe
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[0:5].[ImageId,Name]' --output text
aki-0251b36b    None
aki-0a4aa863    None
aki-12f0127b    None
aki-1a946e73    None
aki-1c669375    vmlinuz-2.6.21.7-2.ec2.v1.3.fc8xen.manifest.xml
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[?starts_with(ImageId, `ami-`) == `true`]|[0:5].[ImageId,Name]' --output text
ami-0048c968    .NET Beanstalk Cfn Container v2.0.2.1 on Windows 2012
ami-005daf69    ElasticBeanstalk-Tomcat6-64bit-20110322-2041
ami-0078da69    amzn-ami-pv-2012.03.2.x86_64-s3
ami-00c17768    aws-elasticbeanstalk-amzn-2014.09.0.x86_64-php55-gpu-201409291824
ami-013aca6a    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[?Name!=`null`]|[?starts_with(Name, `aws-elasticbeanstalk`) == `true`]|[?contains(Name, `tomcat7`) == `true`]|[0:5].[ImageId,Name]' --output text
ami-013aca6a    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121
ami-033aca68    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-hvm-201506152121
ami-143b257c    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201505181913
ami-248dff4c    aws-elasticbeanstalk-amzn-2014.09.1.i386-tomcat7java7-pv-201501140057
ami-34c1775c    aws-elasticbeanstalk-amzn-2014.09.0.x86_64-tomcat7java7-gpu-201409291824
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[?Name!=`null`]|[?starts_with(Name, `aws-elasticbeanstalk`) == `true`]|[?contains(Name, `x86_64-tomcat7java6-pv`) == `true`].[CreationDate,ImageId,Name]' --output text | sort -k1
2014-09-29T18:32:37.000Z        ami-bec177d6    aws-elasticbeanstalk-amzn-2014.09.0.x86_64-tomcat7java6-pv-201409291829
2014-10-15T22:26:52.000Z        ami-721aa21a    aws-elasticbeanstalk-amzn-2014.09.0.x86_64-tomcat7java6-pv-201410152224
2015-01-27T23:14:12.000Z        ami-a06d29c8    aws-elasticbeanstalk-amzn-2014.09.1.x86_64-tomcat7java6-pv-201501272310
2015-04-03T20:07:00.000Z        ami-6252620a    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201504032003
2015-04-20T17:45:22.000Z        ami-58959230    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201504201742
2015-05-18T19:18:26.000Z        ami-143b257c    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201505181913
2015-06-15T21:35:00.000Z        ami-013aca6a    aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121
```
```
$ aws ec2 describe-images --owner amazon --query 'Images[?Name!=`null`]|[?starts_with(Name, `aws-elasticbeanstalk`) == `true`]|[?contains(Name, `x86_64-tomcat7java6-pv`) == `true`].[CreationDate,ImageId,Name]' --output text | sort -k1 | tail -n1 | gawk '{print $2}' | xargs aws ec2 describe-images --image-ids "$@"
{
    "Images": [
        {
            "Hypervisor": "xen",
            "Public": true,
            "RootDeviceType": "ebs",
            "KernelId": "aki-919dcaf8",
            "State": "available",
            "ImageLocation": "amazon/aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121",
            "ImageOwnerAlias": "amazon",
            "Name": "aws-elasticbeanstalk-amzn-2015.03.0.x86_64-tomcat7java6-pv-201506152121",
            "VirtualizationType": "paravirtual",
            "Architecture": "x86_64",
            "CreationDate": "2015-06-15T21:35:00.000Z",
            "ImageId": "ami-013aca6a",
            "BlockDeviceMappings": [
                {
                    "DeviceName": "/dev/sda1",
                    "Ebs": {
                        "Encrypted": false,
                        "DeleteOnTermination": true,
                        "VolumeType": "standard",
                        "VolumeSize": 8,
                        "SnapshotId": "snap-521f631c"
                    }
                }
            ],
            "RootDeviceName": "/dev/sda1",
            "ImageType": "machine",
            "OwnerId": "102837901569"
        }
    ]
}
```