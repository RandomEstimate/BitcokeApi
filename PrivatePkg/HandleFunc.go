package PrivatePkg

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"

	"strings"
	"time"
)

//////////////////////////////////////////////处理交易函数/////////////////////////////////////////////
func (a *BcWsP2) handleTrade(data OrderMode) {

	PriceMap, err := a.BcWsP.GetPriceMap(data.Symbol)
	if err != nil {
		if a.log != nil {
			a.log.E("handleTrade GetPriceMap error : %v", err)
		}
		return
	}
	// 有可能不存在 bug
	BidPrice := PriceMap["BidPrice5"]
	AskPrice := PriceMap["AskPrice5"]

	threshold := 30000.
	if data.Threshold != 0 {
		threshold = data.Threshold
	}

	tmp := 0.0
	for i := 1; i <= 10; i++ {
		if s, ok := PriceMap["AskAmount"+fmt.Sprint(i)]; ok {
			tmp += s
			if tmp > threshold {
				AskPrice = PriceMap["AskPrice"+fmt.Sprint(i)]
				break
			}
		}
	}

	tmp = 0.0
	for i := 1; i <= 10; i++ {
		if s, ok := PriceMap["BidAmount"+fmt.Sprint(i)]; ok {
			tmp += s
			if tmp > threshold {
				BidPrice = PriceMap["BidPrice"+fmt.Sprint(i)]
				break
			}
		}
	}

	var price float64
	if data.Side == "Buy" {
		price = AskPrice
	} else if data.Side == "Sell" {
		price = BidPrice
	}

	//下单交易
	a.tradeOrder(data.Symbol, int(data.Qty), data.Side, price, "Limit")

	a.positionUpdateSwitch.Store(false)
	a.positionSearchCD.Store(time.Now())
	//避免交易被忽略
	time.Sleep(10 * time.Millisecond)
}

func (a *BcWsP2) tradeOrder(Symbol string, Qty int, Side string, Price float64, OrderType string) {
	// Qty表示下单数量 Side表示下单方向  Price 表示下单价格 OrderType表示下单类型

	sendData := orderMode{
		AccountId:    a.AccountIdNum,
		Channel:      "Windows,chrome/69.0.3947.100",
		Iceberg:      "false",
		OpenPosition: "true",
		Priority:     "HIGH",
		Qty:          Qty,
		ShowQty:      0,
		Side:         Side,
		Symbol:       fmt.Sprintf("X%vUSD", Symbol),
		TxId:         "TX20200324-080054-692-70",
		Csclass:      "org.cyanspring.exbusiness.event.order.EnterOrderRequestEvent",
	}

	if OrderType == "Market" {
		sendData.OrderType = "Market"
		sendData.Price = 0
		sendData.Tif = "GOOD_TILL_CANCEL"
	} else if OrderType == "Limit" {
		sendData.OrderType = "Limit"
		sendData.Price = Price
		//sendData.Tif = "FILL_OR_KILL"
		sendData.Tif = "IMMEDIATE_OR_CANCEL"
	}

	buf, _ := json.Marshal(sendData)
	send2 := strings.Replace(string(buf), "csclass", "_csclass", -1)

	if a.log != nil {
		a.log.I("order: %v", send2)
	}
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(send2)); err != nil {
		if a.log != nil {
			a.log.E("tradeOrder error : %v", err)
		}
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////仓位主动更新////////////////////////////////////////
func (a *BcWsP2) positionUpdate() {

	sendMessage := subMode{
		Key:      "null",
		Priority: "HIGH",
		TxId:     "TX20200323-112235-410-19",
		Csclass:  "org.cyanspring.exbusiness.event.data.TradingDataRequestEvent",
	}
	buf, _ := json.Marshal(sendMessage)
	send2 := strings.Replace(string(buf), "csclass", "_csclass", -1)
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(send2)); err != nil {
		if a.log != nil {
			a.log.E("positionUpdate error : %v", err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////心跳保持/////////////////////////////////////////////
func (a *BcWsP2) pong() {
	defer func() {
		if err := recover(); err != nil {
			a.log.E("pong panic : %v", err)
		}
	}()

	sendMessage := subMode{
		Key:      "null",
		Priority: "HIGH",
		TxId:     "TX20200323-112235-410-19",
		Csclass:  "org.cyanspring.event.RequestSystemStateEvent",
	}
	buf, _ := json.Marshal(sendMessage)
	send2 := strings.Replace(string(buf), "csclass", "_csclass", -1)
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(send2)); err != nil {
		if a.log != nil {
			a.log.E("pong error : %v", err)
		}
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (a *BcWsP2) handlePositionData(data positionData) {

	a.m.Lock()
	defer a.m.Unlock()

	// delete all
	for k := range a.accountInfo {
		delete(a.accountInfo, k)
	}

	for _, v := range data.TradingDataUpdate.Pl {
		symbol := v.Symbol[:len(v.Symbol)-3][1:]

		var sideNumber float64
		if v.Side == "Long" {
			sideNumber = 1
		} else if v.Side == "Short" {
			sideNumber = -1
		}

		d := accountInfo{}
		d.openPrice = v.Price
		d.position = v.Qty * sideNumber

		a.accountInfo[symbol] = d

	}

	a.positionUpdateSwitch.Store(true)

}

/////////////////////////////////////////////处理执行函数//////////////////////////////////////////

func (a *BcWsP2) handleExec(data execData) {

	if data.ExecType == "FILLED" || data.ExecType == "CANCELED" {

		a.execInfoChan <- 1
		a.positionUpdateSwitch.Store(true)
	}

}
