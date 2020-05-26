package PublicPkg

import "time"

/////////////////////////////////////////价格深度返回数据///////////////////////////////////
type DepthMap struct {
	// 所有行情数据统一返回的数据结构Bid1表示bid的委托一,Bid1Amount表示张数,Bid1Count表示委托数量,若没有数量则忽略
	Bid1       float64
	Bid1Amount float64
	Bid1Count  float64
	Bid2       float64
	Bid2Amount float64
	Bid2Count  float64
	Bid3       float64
	Bid3Amount float64
	Bid3Count  float64
	Ask1       float64
	Ask1Amount float64
	Ask1Count  float64
	Ask2       float64
	Ask2Amount float64
	Ask2Count  float64
	Ask3       float64
	Ask3Amount float64
	Ask3Count  float64

	Time      time.Time
	TimeStamp float64
	TimeStr   string
}

///////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////行情数据内部请求///////////////////////////////////
type PriceStruct struct {
	buy        []float64
	buyAmount  []float64
	buyCount   []float64
	sell       []float64
	sellAmount []float64
	sellCount  []float64
}

///////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////深度信息//////////////////////////////////////////
type DepthData struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"Qty"`
	Count float64 `json:"count"`
}

type OrderBookData struct {
	BuyDepth  []DepthData `json:"buyDepth"`
	SellDepth []DepthData `json:"sellDepth"`
	Csclass   string      `json:"csclass"`
}

func (a *OrderBookData) FullValid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.marketdata.DepthReplyEvent" {
		return true
	} else {
		return false
	}
}

func (a *OrderBookData) UpdateValid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.marketdata.DepthUpdateEvent" {
		return true
	} else {
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
