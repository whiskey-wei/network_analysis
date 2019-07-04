# 网络流量分析工具
## 运行维护
### 安装
> 配置好golang环境后先执行 go get github.com/google/gopacket 获取gopacket包，然后执行 go get github.com/whiskey-wei/network_analysis 获取项目源码，进入项目目录执行go build编译
### 使用
* 命令行选项
    * -h  
        查看命令行选项  
    * -d
        选择网卡接口，默认为en0  
    * -i  
        过滤的端口号，默认不过滤（即为0）
    * -p
        过滤的协议，只支持tcp与udp，默认不过滤
    * -t  
        每隔t时间统计该时间段内，连接中数据包的大小，单位是s，默认为10s
    * tp  
        tcp数据包的文件存储路径，默认为./info/tcpinfo
    * dp  
        udp数据包的文件存储路径，默认为./info/udpinfo

## 项目功能
* 实时抓取网络数据  
* 网络协议分析与显示  
* 将网络数据包聚合成数据流，以源 IP、目的 IP、源端口、目的端口及协议等五元组的形式存储
* 计算并显示固定时间间隔内网络连接

## 项目结构
* conf/conf.go  
  存储相关配置
* handlepacket
  * hashkey.go  
    计算哈希值  
  * hashlist.go  
    哈希表的增删操作  
  * info.go  
    对网络包的处理，包括提取相关信息，然后存储到哈希表和磁盘文件
  * model.go
    用到的哈希存储模型

##