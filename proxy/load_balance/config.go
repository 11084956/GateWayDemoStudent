package load_balance

import (
	"GateWayDemoStudent/proxy/zookeeper"
	"fmt"
)

type Observer interface {
	Update()
}

//配置主题
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type LoadBalanceZkConf struct {
	observers    []Observer
	path         string
	zkHosts      []string
	confOpWeight map[string]string
	activeList   []string
	format       string
}

//添加
func (s *LoadBalanceZkConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

//回调通知
func (s *LoadBalanceZkConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceZkConf) GetConf() (confList []string) {
	for _, ip := range s.activeList {
		weight, ok := s.confOpWeight[ip]

		if !ok {
			weight = "50" //默认权重
		}

		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}

	return confList
}

//更新配置时,通知监听者也更新
func (s *LoadBalanceZkConf) WatchConf() {
	zkManager := zookeeper.NewZkManager(s.zkHosts)
	_ = zkManager.GetConnect()

	fmt.Println("watchConf")
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)

	go func() {
		defer zkManager.Close()

		for true {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr", changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed")
				s.UpdateConf(changedList)
			}
		}
	}()
}

//更新配置时,通知监听者也更新
func (s *LoadBalanceZkConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string) (*LoadBalanceZkConf, error) {
	zkManager := zookeeper.NewZkManager(zkHosts) // new
	_ = zkManager.GetConnect()                   //获取连接
	defer zkManager.Close()                      //关闭服务

	//获取服务列表
	zlist, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}

	mConf := &LoadBalanceZkConf{
		format:       format,
		activeList:   zlist,
		confOpWeight: conf,
		zkHosts:      zkHosts,
		path:         path,
	}

	mConf.WatchConf()

	return mConf, nil
}

type LoadBalanceObserver struct {
	ModuleConf *LoadBalanceZkConf
}

func (l *LoadBalanceObserver) Update() {
	fmt.Println("Update get conf:", l.ModuleConf.GetConf())
}

func NewLoadBalanceObserver(conf *LoadBalanceZkConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}
