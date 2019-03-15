# snmp_webhook
如何测试？
----------
在trapdebug目录下有一个net-snmp的目录，这里面包含了构建用的dockerfile。net-snmp可以用来测试你的snmp协议报文是否正确。
需要修改时注意挂载进去的MIB文件要跟你报文的格式一致