# grafana-wecaht-alert
grafana 企业微信 报警系统
### 介绍
+ 前提你已经在 [企业微信文档](https://work.weixin.qq.com/api/doc/90000/90136/91770) 申请过账号，也已经配置好配置好key
### 打包文件
```
    1. export GOPROXY=https://goproxy.io,direct
    2. go mod tidy 
    3. go build main.go -o wechat-alert
    4 ./wechat-alert
```
### 配置请求方式
```
 1. 在grafana中配置webhook
 2. Webhook settings里添写 
    http://你的ip:88/send?key=微信的key
 3. 点击Send Test 
 
 备注：Username和Password 无关紧要，可以随意填写。
```
### 备注
代码里实现了 发送企业微信的两种类型的消息。
+ markdown（默认）
+ news
### 结束
