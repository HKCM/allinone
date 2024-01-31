# enable_protection

## Introduce

该脚本用于一次性为现有的EC2添加stop-protection和terminate-protection

脚本可以重复运行

## Command

```bash
cd scripts/EnableEC2Protection
# 创建虚拟环境
python3 -m venv .venv
# 启动虚拟环境
source .venv/bin/activate
# 安装依赖
pip3 install -r requirements.txt

# 仅检查所有机器的protection
python3 enable_protection.py --profile staging --check
# 对所有机器启用terminate-protect
python3 enable_protection.py --profile staging --terminate-protect
# 对所有机器启用stop-protect
python3 enable_protection.py --profile staging --stop-protect


# 仅检查指定机器的protection
python3 enable_protection.py --profile staging --check --instance-ids i-08be869c220bba555,i-02a298b96d8f6308f
# 对指定机器启用terminate_protect
python3 enable_protection.py --profile staging --instance-ids i-00c4e4b8597279400 --terminate-protect
# 对机器机器启用stop_protect
python3 enable_protection.py --profile staging --instance-ids i-00c4e4b8597279400 --stop-protect




# 退出虚拟环境
deactivate
```


