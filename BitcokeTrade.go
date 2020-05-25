package Bitcoke

import (
	"github.com/gorilla/websocket"
)


type Bc struct {
	Conn *websocket.Conn


	PositionNow float64

	PositionSell float64
	PositionBuy float64
	OpenSellPrice float64
	OpenBuyPrice float64



	OpenPrice float64
	OpenSide string
	PositionChan chan *PositionData
	ExecChan chan *Exec
	PositionUpdateChan chan *PositionUpdateData

	TradeChan chan OrderMode // 由于采用ws接口进行交易避免同时使用ws导致无法下单因此采用队列形式


	//配置个人信息 F12抓取
	//{"priority":"HIGH","_csclass":"org.cyanspring.event.RemoteUserLoginRequestEvent","key":null,"txId":"TX89862982-20200525-101342-685-13","channel":"Windows,chrome/69.0.3947.100","password":"","encrypted":true,"user":"","agent":""}
	AgentNum string
	PasswordNum string
	AccountIdNum string

	//交易品种 目前交易所支持品种XBTCUSD XETHUSD XEOSUSD XBCHUSD XLTCUSD
	Symbol string

}

