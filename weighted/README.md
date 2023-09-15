## weighted
描述：一种加权随机法的实现，适用于需要在随机抽取基础上加入权重的场景，例如：抽奖、流量调度等。

### API使用说明：
- NewWeightedRandom(defaultWeight int): 创建一个实例
   - defaultWeight： 所有用户的初始权重值，这个参数决定了中奖次数对权重影响的细腻程度；
- Draw(winCount []int, num int) []int: 从参数指定的用户列表中随机抽取几个用户
   - winCount: 用户列表及中奖次数
        - 数组下标：表示用户，中奖的用户也是返回的数组下标，使用方需要自己做对应； 
        - 数组元素：用户的中奖次数，未中奖填充0； 
   - num: 抽取的用户数； 
   - return []int: 中奖用户列表，每个元素表示一个中奖者，值来自入参winCount的下标，与业务场景中的真实用户对应需要使用方来完成。
- DrawOne(winCount []int) int: 从用户列表中随机抽取一个用户 
     
使用示例：请参考 [weighted_random_test.go](./weighted_random_test.go)

### 设计文档
请参考：https://blog.csdn.net/xiaojia1001/article/details/132914175