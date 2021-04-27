package load_balance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Hash func(data []byte) uint32

type Uint32Slice []uint32

//hash一致性 负载均衡
type ConsistentHashBanlance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int               //复制因子
	keys     Uint32Slice       //已排序的节点hash切片
	hashMap  map[uint32]string //节点哈希和key的map, 键是hash值,值是节点key

	//观察主体
	conf LoadBalanceConf
}

func (s Uint32Slice) Len() int {
	return len(s)
}

func (s Uint32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Uint32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewConsistentHashBanlance(replicas int, fn Hash) *ConsistentHashBanlance {
	m := &ConsistentHashBanlance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}

	if m.hash == nil {
		//最多32位,保证是一个 2^32-1 环
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

//验证是否为空
func (c *ConsistentHashBanlance) IsEmpty() bool {
	return len(c.keys) == 0
}

// Add 方法用户增减缓存节点,参数为节点key, 比如使用IP
func (c *ConsistentHashBanlance) Add(params ...string) error {
	if len(params) <= 0 {
		return errors.New("param len 1 at least")
	}

	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()

	//结合复制因子计算所有虚拟节点的hash值,
	//并存入m.keys中,同时在m.hasMap中保存哈希值和key的映射
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))

		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}

	//对所有虚拟节点的hash值进行排序,方便二分查找
	sort.Sort(c.keys)

	return nil
}

func (c *ConsistentHashBanlance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}

	hash := c.hash([]byte(key))

	//通过二分查找最优节点,第一个 "服务器hash" 值大于 "数据hash" 值的就是最优 "服务器节点"
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	//如果查找结果 大于 服务器节点hash数组的最大索引, 标识此时该对象哈希值位于最后一个节点之后,那么放入第一个节点中
	if idx == len(c.keys) {
		idx = 0
	}

	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.hashMap[c.keys[idx]], nil
}

func (c *ConsistentHashBanlance) SetConf(conf LoadBalanceConf) {
	c.conf = conf
}

//通知更新
func (c *ConsistentHashBanlance) Update() {
	if conf, ok := c.conf.(*LoadBalanceZkConf); ok {
		confList := conf.GetConf()
		fmt.Println("Update get conf:", confList)

		c.keys = nil
		c.hashMap = nil

		for _, ip := range confList {
			_ = c.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := c.conf.(*LoadBalanceCheckConf); ok {
		confList := conf.GetConf()
		fmt.Println("Update get conf:", confList)

		c.keys = nil
		c.hashMap = map[uint32]string{}

		for _, ip := range confList {
			_ = c.Add(strings.Split(ip, ",")...)
		}
	}
}
