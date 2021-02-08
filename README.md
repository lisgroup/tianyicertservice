## 天翼翼证-区块链 验签服务

### 一、启动方法
1. run
```shell script
go run main.go
```
2. build
```shell script
go build main.go
./main
```

### 二、 客户端获取sign

- 请求地址设定为 http://127.0.0.1:8080/

- 请求参数列表

参数名称 | 描述
---|---
body | 接口参数
nonce | 签名内容里面的随机数
timeStamp | 签名内容里面的时间戳
apiSecret | 秘钥


完整请求示例：
```
http://127.0.0.1:8080/sign?body={"Data":"test","Remark":"test"}&nonce=100&timeStamp=1608626661&apiSecret=apiSecret*********
```
响应JSON示例：

```json
{
    "code": 0, 
    "data": {
        "sign": "40677972351786153192903865501230545775316569974084416717902566734041634855574+24262690419852203592204540437043224196359361711282558324821312146479811246097"
    },
    "reason": "success"
}
```
