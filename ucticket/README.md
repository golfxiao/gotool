# 1. 功能说明
* 此包主要用于提供一个通用的业务ID生成器，通过DB + 分段Cache的方式来保证ID生成的有序性和高性能。

# 2. 使用说明

## 2.1 初始化存储

### 2.1.1 在自己的业务数据库中新创建一张表：ticket
```
    DROP TABLE IF EXISTS `ticket`;
    CREATE TABLE `ticket` (
    `biz_tag` varchar(128) NOT NULL DEFAULT '',
    `max_id` bigint(20) NOT NULL DEFAULT '1',
    `step` int(11) DEFAULT NULL,
    `desc` varchar(256) DEFAULT NULL,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`biz_tag`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

* biz_tag: 需要生成ID的业务名称，一般直接使用表名称
* max_id: 当前已经分配过的ID最大值
* step: 每次批量分配ID的步长，例如：50表示一次分配50个ID并缓存到内存Cache中
* 一个ticket表为一台应用DB中的所有表提供ID存储服务

### 2.1.2 对于线上已经运行的应用，需要将当前数据库中的ID最大值初始化表中
```
    INSERT INTO `ticket`(`biz_tag`, `max_id`, `step`) VALUES ('department', '4066', '20');
    INSERT INTO `ticket`(`biz_tag`, `max_id`, `step`) VALUES ('user', '50038259', '50');
```
## 2.2 应用的配置文件中增加如下配置
```
    # ticket DB连接串, 此处以MySQL为例，也支持Mongo连接串
    ticket_datasrc = root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4
    # ticket db 连接数量
    ticket_cons = 10
    # ticket 表名，可选，未配置时默认为ticket表
    ticket_tb_name = ticket
    # 每个应用需要使用的ID类型，必选，即biz_tag列表，分号分隔
    ticket_allow_tags = user;document;
    # 是否使用预加载
    ticket_use_preload = true

```
## 2.3 服务启动时初始化ticket库
```
	config := ucticket.TicketConfig{
		DataSrc:   beego.AppConfig.String("ticket_datasrc"),
		ConnCount: beego.AppConfig.DefaultInt("ticket_cons", 0),
		TableName: beego.AppConfig.String("ticket_tb_name"),
		ScopeList: beego.AppConfig.Strings("ticket_allow_tags"),
        UsePreload: beego.AppConfig.Bool("ticket_use_preload"),
	}
    // 此初始化函数适用于MySQL数据库
	err = ucticket.InitTicketDB(config)
    // 此初始化函数适用于Mongo数据库
    // err = ucticket.InitTicketMongo(config)
```
## 2.4 使用ucticket生成ID
```
    const (
        BiztagUser ucticket.BiztagType = "user"  // 这里的user需要和ticket_allow_tags中配置的名称一致
        BiztagDepartment ucticket.BiztagType = "department"  
    )
    userId, err := BiztagTest.GetGlobalId()
    // TODO 业务处理
``` 