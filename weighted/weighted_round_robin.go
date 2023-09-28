package weighted

// 支持加权轮询的选取器
type WeightedSelector struct {
	servers     []*WeightedServer // 带权重的服务器列表
	totalWeight int               // 所有服务器加起来的总权重，用于限定实际服务器权重的上下限
}

// 带权重的服务器结构定义
type WeightedServer struct {
	server        string // 服务器标识，可以是ID、域名或其它标识
	weight        int    // 预设权重
	currentWeight int    // 当前权重，会随分配轮次动态变化
}

func NewWeightedSelector(servers map[string]int) *WeightedSelector {
	ss := make([]*WeightedServer, 0, len(servers))
	totalWeight := 0
	for s, weight := range servers {
		ss = append(ss, &WeightedServer{
			server:        s,
			weight:        weight,
			currentWeight: 0,
		})
		totalWeight += weight
	}
	return &WeightedSelector{
		servers: ss,
	}
}

func (this *WeightedSelector) RoundRobin() string {
	var maxWeight = -this.totalWeight // 定义maxWeight表示实时权重最大值，初始给个最小值-totalWeight
	var cur *WeightedServer           // 与maxWeight对应的权重最大的server
	for _, e := range this.servers {
		e.currentWeight += e.weight      // 每轮分配前的加权操作
		if e.currentWeight > maxWeight { // 比较大小,找到实时权重最大的server
			cur, maxWeight = e, e.currentWeight
		}
	}
	cur.currentWeight -= this.totalWeight // 被选中的服务器，做减权操作
	return cur.server
}
