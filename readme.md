# <center>ReadMe<center>

-----------------------------

# 使用连接
[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

[Go语言开发CLI实用工具](https://blog.csdn.net/Yanzu_Wu/article/details/83021193)

-------------------------

# 使用注意事项
> 使用与`selpg`类似，但是用了`pflag`而不是`flag`，需要注意以下格式

```
USAGE: ./selpg [--s start] [--e end] [--l lines | --f ] [ --d dest ] [ in_filename ]

 selpg --s start    : start page
 selpg --e end      : end page
 selpg --l lines    : lines/page
 selpg --f          : check page with '\f'
 selpg --d dest     : pipe destination
```

> 演示实例

```
selpg --s 1 --e 2 --l 3 --d lp1 input >output 2>err
```
