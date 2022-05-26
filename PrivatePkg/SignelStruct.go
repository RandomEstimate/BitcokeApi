package PrivatePkg

////////////////////////////////////////////////登录结构体////////////////////////////////////////////////
type loginMode struct {
	Agent     string `json:"agent"`
	Channel   string `json:"channel"`
	Encrypted string `json:"encrypted"`
	Key       string `json:"key"`
	Token     string `json:"token"`
	Priority  string `json:"priority"`
	TxId      string `json:"txId"`
	User      string `json:"user"`
	Csclass   string `json:"csclass"`
}

////////////////////////////////////////////////订阅结构体////////////////////////////////////////////////
type subMode struct {
	Key      string `json:"key"`
	Priority string `json:"priority"`
	TxId     string `json:"txId"`
	Csclass  string `json:"csclass"`
}

////////////////////////////////////////////////账户信息存储////////////////////////////////////////////////
type accountInfo struct {
	position  float64
	openPrice float64
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
type orderMode struct {
	AccountId    string  `json:"accountId"`
	Channel      string  `json:"channel"`
	Iceberg      string  `json:"iceberg"`
	Key          string  `json:"key"`
	OpenPosition string  `json:"openPosition"`
	OrderType    string  `json:"orderType"`
	Price        float64 `json:"price"`
	Priority     string  `json:"priority"`
	Qty          int     `json:"qty"`
	ShowQty      int     `json:"showQty"`
	Side         string  `json:"side"`
	Symbol       string  `json:"symbol"`
	Tif          string  `json:"tif"`
	TxId         string  `json:"txId"`
	Csclass      string  `json:"csclass"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type OrderMode struct {
	// 下单订单结构体
	Symbol     string
	Qty        float64 //传入交易数量
	Side       string  //传入交易方向 Buy/Sell
	Type       string  //传入交易方式 Limit/Market
	Threshold  float64 //传入当前距离挂单口的amount数量
	LimitNum   int64   //传入最大Limit次数
	MarketSign string  //为了防止市价单子产生巨大的滑点，当MarketSign = Market时进行大范围市价成交，但是是以Limit形式进行，不保证一定完全成交
}

//////////////////////////////////////////交易执行后的返回信息//////////////////////////////////////////
type execData2 struct {
	AvgPx  float64 `json:"avgPx"`
	CumQty float64 `json:"cumQty"`
	Qty    float64 `json:"qty"`
	Side   string  `json:"side"`
}

type execData struct {
	//交易返回的状态信息
	ExecType string    `json:"execType"`
	Order    execData2 `json:"order"`
	Csclass  string    `json:"csclass"`
}

func (a *execData) valid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.order.OrderUpdateEvent" && a.ExecType != "" {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////仓位主动更新和仓位被动更新信息////////////////////////////////
type positionData2 struct {
	Volume float64 `json:"closableQty"`
	Qty    float64 `json:"qty"`
	Side   string  `json:"side"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
type positionData1 struct {
	Pl []positionData2 `json:"pl"`
}

type positionData struct {
	//订正仓位信息
	TradingDataUpdate positionData1 `json:"tradingDataUpdate"`
	Csclass           string        `json:"csclass"`
}

func (a *positionData) valid() bool {
	if a.Csclass == "org.cyanspring.exbusiness.event.data.TradingDataReplyEvent" ||
		a.Csclass == "org.cyanspring.exbusiness.event.data.TradingDataUpdateEvent" {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////
