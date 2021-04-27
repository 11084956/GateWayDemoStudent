package load_balance

import (
	"fmt"
	"testing"
)

func TestB(t *testing.T) {
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
