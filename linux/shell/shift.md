# shift
```shell
$ cat shift.sh
#!/usr/bin/env bash

echo 
count=1
while [ -n $1 ];do
  echo "Parameter #$count = $1"
  count=$[ $count + 1 ]
  shift
done

./shift.sh rich barbara katie jessica
Parameter #1 = rich
Parameter #2 = barbara
Parameter #3 = katie
Parameter #4 = jessica
```