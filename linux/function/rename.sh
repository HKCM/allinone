#!/usr/bin/env bash
# 批量重命名当前目录下的图片文件

count=1;
for img in $(find . -iname '*.png' -o -iname '*.jpg' -type f -maxdepth 1)
do
    new=image-${count}.${img##*.}
    echo "Renaming ${img} to ${new}"
    mv "${img}" "${new}"
    let count++;
done