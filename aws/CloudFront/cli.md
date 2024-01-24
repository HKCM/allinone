
清除缓存
```shell
aws cloudfront create-invalidation --distribution-id <id> --paths "/*"

aws cloudfront create-invalidation \
    --distribution-id <id> \
    --paths "/example-path/example-file.jpg" "/example-path/example-file2.png"
```