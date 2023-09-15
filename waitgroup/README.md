## waitgroup
此包用于实现一个带并发数量限制的waitgroup，适用于需要一组任务并发运行来提高性能，但又不能并发量太高的场景。

### API说明
- NewLimitedWaitGroup(max int) ： 创建一个最大并发数量为max的WaitGroup实例
- Add(n int) : 同原生WaitGroup的Add方法
- Done(): 同原生WaitGroup的Done方法
- Wait(): 同原生WaitGroup的Wait方法

使用示例：请参考 [limited_wait_group_test.go](./limited_wait_group_test.go)

### 设计文档
参考：https://blog.csdn.net/xiaojia1001/article/details/132842641