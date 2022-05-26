### Bitcoke交易所api 

bitcoke交易所对开放api的限制频率，目前这一套api利用web端进行行情接收，同时也可以通过web端进行部分功能的下单。

***

##### 目录结构

PublicPkg模块不需要额外的用户信息，PrivatePkg模块需要额外的用户信息。

> **BitcokeApi**
>
> **├── PrivatePkg**
> **│   ├── BitcokeTrade.go**
> **│   ├── HandleFunc.go**
> **│   ├── InterfaceFunc.go**
> **│   └── SignelStruct.go**
> **├── PublicPkg**
> **│   ├── BitcokePrice.go**
> **│   ├── HandleFunc.go**
> **│   ├── InterfaceFunc.go**
> **│   └── SignalStuct.go**

***

##### 行情接收模块支持接口

```go
// 行情订阅
// BcWsP.OrderSymbol("BTC")
func (a *BcWsP) OrderSymbol(symbol string) error

// 行情更新时间
// BcWsP.GetUpdateTime("BTC")
func (a *BcWsP) GetUpdateTime(symbol string) (*time.Time, error) 

// 行情orderbook获取
// BcWsP.GetPriceMap("BTC")
func (a *BcWsP) GetPriceMap(symbol string) (map[string]float64, error)

```

##### 行情接收模块example

```go
package main

import (
	"BitcokeApi/PublicPkg"
	"fmt"
	"sync"
	"time"
)

func main() {
    
	obj, _ := PublicPkg.NewPublic(nil)
	w := &sync.WaitGroup{}
	w.Add(1)
	go obj.Start(w)
	w.Wait()
	err := obj.OrderSymbol("BTC")
    if err != nil {
        panic(err)
    }
    
	fmt.Println("订阅成功")
    
    go func() {
		for {
			time.Sleep(time.Second)
            fmt.Println(obj.GetUpdateTime("BTC"))
			fmt.Println(obj.GetPriceMap("BTC"))
		}
	}()
    
    select{}
    
}
```

