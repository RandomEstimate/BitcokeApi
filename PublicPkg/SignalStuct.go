package PublicPkg

import "time"

//////////////////////////////////////////深度信息//////////////////////////////////////////
type DepthData struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"Qty"`
	Count float64 `json:"count"`
}

type orderBookData struct {
	BuyDepth  []DepthData `json:"buyDepth"`
	SellDepth []DepthData `json:"sellDepth"`
	Key       string      `json:"key"`
	Csclass   string      `json:"csclass"`
}

func (a *orderBookData) fullValid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.marketdata.DepthReplyEvent" {
		return true
	} else {
		return false
	}
}

func (a *orderBookData) updateValid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.marketdata.DepthUpdateEvent" {
		return true
	} else {
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
// 本地维护order book
type orderBook struct {
	ask        map[float64]float64
	bid        map[float64]float64
	updateTime time.Time
}
