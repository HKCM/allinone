keyword: list instances, describeInstances, 列出ec2

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install boto3
```

```python
import boto3

# create boto3 client
my_session = boto3.Session(profile_name='prod')
client = my_session.client('ec2')
# client = boto3.client('ec2',)

def main():
    # filter_name := "image-id"
	# filter_value := "ami-003488c9c9e37d326"
	# filter_name := "tag-key" 列出包含指定tag的机器,无论tag的值是什么
	# filter_value := "backup"
	filter_name := "tag:backup" // 列出包含指定tag的机器,且tag的值为weekly
	filter_value := "weekly"

    reservations = client.describe_instances(
        Filters=[ {'Name': filter_name, "Values": [filter_value]}],
    )
    if len(reservations["Reservations"]) == 0:
        print("Nothing")

    for reservation in reservations["Reservations"]:
        for instance in reservation["Instances"]:
            for tag in instance["Tags"]:
                if tag["Key"] == "Name":
                    instance_name = tag["Value"]
                    print(instance_name)

    print(len(reservations["Reservations"]))

if __name__ == "__main__":
    main()
```