package PublicPkg

import (
	"fmt"
	"github.com/RandomEstimate/util_go/common"
	"github.com/gorilla/websocket"
	"sort"
	"time"
)

func (a *BcWsP) OrderSymbol(symbol string) error {

	messageList := make([]string, 0)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.marketdata.DepthUpdateEvent","subKey":"X%vUSD"}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.marketdata.DepthRequestEvent","key":"X%vUSD,"txId":"TX20200413-213011-836-5"}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.marketdata.DepthFullUpdateEvent","subKey":"X%vUSD"}"`)

	if common.Index(symbol, a.symbolList) != -1 {
		return nil
	}

	if a.conn == nil {
		if a.log != nil {
			a.log.E("orderSymbol and connection is nil ")
		}

		return fmt.Errorf("orderSymbol and connection is nil ")
	}

	for _, v := range messageList {
		if err := a.conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(v, symbol))); err != nil {
			if a.log != nil {
				a.log.E("symbol %v. orderSymbol error %v ", symbol, err)
			}
			a.symbolList[symbol] = struct{}{}
			return err
		}
		time.Sleep(time.Millisecond * 50)
	}

	a.symbolList[symbol] = struct{}{}

	return nil

}

func (a *BcWsP) GetUpdateTime(symbol string) (*time.Time, error) {
	a.m.Lock()
	defer a.m.Unlock()

	if d, ok := a.orderBook[symbol]; !ok {
		return nil, fmt.Errorf("symbol not exist. ")
	} else {
		t := d.updateTime
		return &t, nil
	}
}

func (a *BcWsP) GetPriceMap(symbol string) (map[string]float64, error) {
	a.m.Lock()
	defer a.m.Unlock()

	//
	if d, ok := a.orderBook[symbol]; !ok {
		return nil, fmt.Errorf("symbol not exist. ")
	} else {
		// sort
		askKey := make([]float64, len(d.ask))
		j := 0
		for k := range d.ask {
			askKey[j] = k
			j++
		}
		sort.Float64s(askKey)

		bidKey := make([]float64, len(d.bid))
		j = 0
		for k := range d.bid {
			bidKey[j] = k
			j++
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(bidKey)))

		m := make(map[string]float64)
		for i := 1; (i < 11) && (i < len(askKey)+1); i++ {
			m["AskPrice"+fmt.Sprint(i)] = askKey[i-1]
			m["AskAmount"+fmt.Sprint(i)] = d.ask[askKey[i-1]]
		}

		for i := 1; (i < 11) && (i < len(bidKey)+1); i++ {
			m["BidPrice"+fmt.Sprint(i)] = bidKey[i-1]
			m["BidAmount"+fmt.Sprint(i)] = d.bid[bidKey[i-1]]
		}
		return m, nil
	}

}
