## workerPool
描述：管理工作线程的资源池，适用于需要对并发运行的任务数量作控制的场景。

### API使用说明：
- NewWorkerPool(max int) : 创建一个线程数量为max的资源池
- Run(ctx context.Context, run runFunc): 异步运行一个任务，如果资源池中无可用线程，会阻塞等待直到有可用线程。
    - run： 指定待运行的任务
    - ctx： 上下文，可选，预留用于串联上下文日志
- RunWithKey(ctx context.Context, key TaskKey, run runFunc) error : 
    - 在Run函数的基础上，支持对同一个Key重复运行任务的检测，同一时间一个Key只允许运行一个任务； 
    - 如果检测到重复运行，只有第一个调用者能申请到资源来运行，后面的调用者会返回ErrAlreadyRunning错误； 
     
### 设计文档
请参考：https://blog.csdn.net/xiaojia1001/article/details/132757738