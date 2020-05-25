package Bitcoke

type OrderMode struct {
	// 下单订单结构体
	Qty string //传入交易数量
	Side string //传入交易方向 Buy/Sell
}

//////////////////////////////////////////交易执行后的返回信息//////////////////////////////////////////
type Exec1 struct {
	AvgPx float64`json:"avgPx"`
	CumQty float64 `json:"cumQty"`
	Qty float64 `json:"qty"`
	Side string `json:"side"`
}

type Exec struct {
	//交易返回的状态信息
	ExecType string`json:"execType"`
	Order Exec1`json:"order"`
	Csclass string `json:"csclass"`
}

func (a *Exec)Valid()bool{
	if a.Csclass == "org.cyanspring.exbusiness.event.order.OrderUpdateEvent" && a.ExecType != ""{
		return true
	}else{
		return false
	}
}
////////////////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////仓位主动更新和仓位被动更新信息////////////////////////////////
type PositionData1 struct {
	Price float64`json:"price"`
	Qty float64`json:"qty"`
	Side string`json:"side"`
}

type PositionData struct {
	//订正仓位信息
	Positions []PositionData1`json:"positions"`
	Csclass  string `json:"csclass"`
}

func (a *PositionData)Valid()bool{
	if a.Csclass=="org.cyanspring.exbusiness.event.position.PositionSnapshotReplyEvent"{
		return true
	}else {
		return false
	}
}

type PositionUpdateData struct {
	//每次交易后的重新更新的仓位信息
	Position PositionData1`json:"position"`
	Csclass  string `json:"csclass"`
}

func (a *PositionUpdateData)Valid()bool{
	if a.Csclass == "org.cyanspring.exbusiness.event.position.PositionUpdateEvent"{
		return true
	}else {
		return false
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////

