## freetry
描述：一个功能试用模块的逻辑抽取，但不完整，推荐拷贝到项目中加上DB支持后对外提供访问。

### API使用说明：
- ApplyFreeUse(userConfig map[string]*UserFreeItem,
	feature string, defaultConfig map[string]FreeItem) error: 申请功能试用
    - userConfig：用户级已有的试用配置
    - feature: 要试用的功能名称
    - defaultConfig: 系统级的试用配置，作为默认配置提供给用户级作配置初始化

     
使用示例：请参考 [freetry_test.go](./freetry_test.go)

配置示例：请参考 [free_use.yaml](./free_use.yaml)

### 设计文档
请参考：https://blog.csdn.net/xiaojia1001/article/details/132959395