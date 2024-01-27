此示例是做dynamodb的rollback全过程
1. 首先以cloudformation template部署dynamodb,并启用时间点恢复
```shell
bash 1.deploy.sh
```
2. 使用CLI以时间点恢复,创建一个新的dynamodb table
```shell
bash 2.restoreBackup.sh
```
3. 将新的dynamodb table导入cloudformation template
```shell
bash 3.importResource.sh
```