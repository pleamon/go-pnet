v0.1.0
1. 添加长度校验机制
    1. block固定为4096
    2. 用读取到的length % block == 0则检验成功,否则放弃数据包
    3. 写入数据,length = length + (block - length % block)