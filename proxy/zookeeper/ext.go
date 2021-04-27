package zookeeper

/*import "github.com/samuel/go-zookeeper/zk"

//注册服务
func (z *ZkManager) RegisterServer(module, host string) error {
	nodePath := z.pathPrefix + module

	return z.RegisterServerPath(nodePath, host)
}

func (z *ZkManager) GetServerList(module string) ([]string, error) {
	return z.GetServerListByPath(z.pathPrefix + module)
}

func (z *ZkManager) WatchServerList(module string) (chan []string, chan error) {
	return z.WatchServerListByPath(z.pathPrefix + module)
}

//watch机制,监听节点变化
func (z *ZkManager) WatchGetData(module string) (chan []byte, chan error) {
	nodePath := z.pathPrefix + "config_" + module

	return z.WatchPathData(nodePath)
}

//获取配置
func (z *ZkManager) GetData(module string) ([]byte, *zk.Stat, error) {
	nodePath := z.pathPrefix + "config_" + module

	return z.GetPathData(nodePath)
}

//更新配置
func (z *ZkManager) SetData(module string, config []byte) error {
	nodePath := z.pathPrefix + "config_" + module

	return z.SetPathData(nodePath, config)
}*/
