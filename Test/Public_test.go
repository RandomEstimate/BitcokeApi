package Test

import (
	"Bitcoke/PublicPkg"
	"fmt"
	"testing"
	"time"
)

func Test_PublicPriceMap(t *testing.T) {
	Bc, _ := PublicPkg.NewPublic("XBTCUSD")
	go Bc.ConnAndReceive()
	fmt.Println("行情接收已经开启，等待接收")
	time.Sleep(1 * time.Second)
	fmt.Println("进行接口测试，目前测试接口 GetPriceMap ")
	fmt.Println("此接口提供行情数据返回")
	var DataFinal PublicPkg.DepthMap
	for i := 0; i <= 10; i++ {
		Data, err := Bc.GetPriceMap()
		if err != nil {
			fmt.Println(err)
		} else {
			DataFinal = Data
			fmt.Println("接收到Data")
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("打印最后一次行情请求")
	fmt.Println(DataFinal)
	fmt.Println("测试结束")
}
