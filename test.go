package main

import (
	"BitcokeApi/PrivatePkg"
	"BitcokeApi/PublicPkg"
	"fmt"
	"sync"
	"time"
)

const (
	AgentNum     = ""
	PasswordNum  = ""
	AccountIdNum = ""
)


func main() {
	a, _ := PublicPkg.NewPublic(nil)
	w := &sync.WaitGroup{}
	w.Add(1)
	go a.Start(w)
	w.Wait()
	err := a.OrderSymbol("BTC")
	fmt.Println("订阅成功", err)

	b := PrivatePkg.NewBcPrivate(AgentNum, PasswordNum, AccountIdNum, a, nil)
	w = &sync.WaitGroup{}
	w.Add(1)
	go b.Start(w)
	w.Wait()
	fmt.Println("启动完成")

	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Println(b.GetPosition("BSV"))
		}

	}()
	//time.Sleep(time.Second * 10)
	//err = b.SafeTrade(PrivatePkg.OrderMode{
	//	Symbol:    "BSV",
	//	Qty:       100,
	//	Side:      "Sell",
	//	Type:      "Limit",
	//	Threshold: 3000,
	//})
	//fmt.Println(err)

	select {}

}
