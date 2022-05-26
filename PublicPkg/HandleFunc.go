package PublicPkg

import (
	"fmt"
	"time"
)

var connectErrThreshold = 3. // scan 时针对 sendErr 错误的容忍上线
var connectErrTimes = 0.     // scan 连续错误的次数

//
func (a *BcWsP) handleOrder(data *orderBookData) {

	a.m.Lock()
	defer a.m.Unlock()
	// 获取symbol
	symbol := data.Key[:len(data.Key)-3][1:]

	if _, ok := a.orderBook[symbol]; !ok {
		a.orderBook[symbol] = orderBook{
			ask:        make(map[float64]float64),
			bid:        make(map[float64]float64),
			updateTime: time.Now(),
		}
	}

	d := a.orderBook[symbol]
	askMap := d.ask
	bidMap := d.bid
	for _, v := range data.SellDepth {
		if v.Qty == 0 {
			delete(askMap, v.Price)
		} else {
			askMap[v.Price] = v.Qty
		}
	}

	for _, v := range data.BuyDepth {
		if v.Qty == 0 {
			delete(bidMap, v.Price)
		} else {
			bidMap[v.Price] = v.Qty
		}
	}

	d.ask = askMap
	d.bid = bidMap
	d.updateTime = time.Now()

	a.orderBook[symbol] = d

}

func (a *BcWsP) scan() {
	a.m.Lock()
	a.m.Unlock()
	for symbol := range a.symbolList {

		if d := a.orderBook[symbol]; time.Since(d.updateTime) > 10*time.Second {

			delete(a.symbolList, symbol)
			if err := a.OrderSymbol(symbol); err != nil {
				connectErrTimes += 1
				if connectErrThreshold == connectErrTimes {
					panic(fmt.Sprintf("orderSymbol error : %v", err))
				}
				if a.log != nil {
					a.log.E("connectErr is %v", connectErrTimes)
				}
				break

			}
		}
	}
}

// 重启时需要清空order book
func (a *BcWsP) reconnectUpdate() {
	a.m.Lock()
	defer a.m.Unlock()

	for k := range a.orderBook {
		delete(a.orderBook, k)
	}

}
