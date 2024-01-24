# Minimized

This system has been minimized by removing packages and content that are
not required on a system that users do not log into.

To restore this content, including manpages, you can run the 'unminimize'
command. You will still need to ensure the 'man-db' package is installed.

```bash
yes | unminimize
```


# policy-rc.d denied execution of start

invoke-rc.d: policy-rc.d denied execution of start

在Dockerfile中添加
```bash
RUN echo "#!/bin/sh\nexit 0" > /usr/sbin/policy-rc.d
```
