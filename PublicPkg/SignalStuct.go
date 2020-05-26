package PublicPkg

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
