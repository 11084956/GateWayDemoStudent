package load_balance

import (
	"fmt"
	"testing"
)

//随机负载均衡
func TestRandomBalance(t *testing.T) {
	rb := &RandomBalance{}

	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}

//轮询负载均衡
func TestPollingBalance(t *testing.T) {
	rb := &RoundRobinBalance{}

	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}

//加权负载均衡
func TestWeightedBalance(t *testing.T) {
	rb := &WeightRoundRobinBalance{}
	_ = rb.Add("127.0.0.1:2003", "4") //0
	_ = rb.Add("127.0.0.1:2004", "3") //1
	_ = rb.Add("127.0.0.1:2005", "2") //2

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
}

//hash一致性,负载均衡
func TestNewConsistentHashBanlance(t *testing.T) {
	rb := NewConsistentHashBanlance(10, nil)
	_ = rb.Add("127.0.0.1:2003") //0
	_ = rb.Add("127.0.0.1:2004") //1
	_ = rb.Add("127.0.0.1:2005") //2
	_ = rb.Add("127.0.0.1:2006") //3
	_ = rb.Add("127.0.0.1:2007") //4

	//url hash
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/error"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/changepwd"))

	//ip hash
	fmt.Println(rb.Get("127.0.0.1"))
	fmt.Println(rb.Get("192.168.0.1"))
	fmt.Println(rb.Get("127.0.0.1"))
}
