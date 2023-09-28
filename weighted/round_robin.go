package weighted

type Selector struct {
	servers []string // 服务器环境列表
	seq     int      // 序列号自增
}

// 构造方法，传入服务器列表servers，服务器标识与权重组成的键值对
func NewSelector(servers map[string]int) *Selector {
	weightedServers := make([]string, 0, len(servers))
	for server, weighted := range servers {
		for i := 0; i < weighted; i++ {
			weightedServers = append(weightedServers, server)
		}
	}
	return &Selector{
		servers: weightedServers,
		seq:     0,
	}
}

func (this *Selector) RoundRobin() string {
	// this.servers := []string{"server1", "server2", "server3"}
	curIndex := this.seq % len(this.servers)
	this.seq++
	return this.servers[curIndex]
}
