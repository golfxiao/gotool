## ukey
此包主要用于生成业务短链中常见的唯一标识key，可以通过一定的算法来将业务标识构造成一个短链ukey，也支持对ukey作解析来还原业务标识信息。

### API说明
- func BuildUkey(data1, data2, data3, randomLen int) string ： 构造并返回ukey
    - 参数data1, data2, data3: 分别为3个整数形式的业务标识, 如果不够用可以自行改造； 
    - 参数randomLen: 支持指定几位随机字符来增强安全性； 
- func ParseUkey(ukey string) (data1, data2, data3, randomLen int) : 解析并返回组成ukey的业务标识
    - 参数ukey：要解析的ukey串； 
    - 返回值 data1, data2, data3为解析出来的三段业务标识， 不足三段时，例如只有第一段有效，则data1有值，data2和data3为0； 
    - randomLen： ukey中包含的随机字符长度； 

使用示例：请参考 [ukey_test.go](./ukey_test.go)

### 设计文档
参考：https://blog.csdn.net/xiaojia1001/article/details/132649727



