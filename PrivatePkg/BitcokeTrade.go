package PrivatePkg

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

const (
	//此线路为Bitcoke的国际线路，可以通过ping的延时进行选择
	WsUrl string = "wss://www.ghcyd.top/prod1/trade/"
)

var (
	PanicBool bool = false

	//避免在自启动时重复启动Deal函数导致资源竞争
	DealGoBool bool = false
)

type Bc struct {
	Conn *websocket.Conn

	PositionNow float64

	PositionSell  float64
	PositionBuy   float64
	OpenSellPrice float64
	OpenBuyPrice  float64

	OpenPrice          float64
	OpenSide           string
	PositionChan       chan *PositionData
	ExecChan           chan *Exec
	PositionUpdateChan chan *PositionUpdateData

	// 由于采用ws接口进行交易避免同时使用ws导致无法下单因此采用队列形式
	TradeChan chan OrderMode

	//配置个人信息 F12抓取
	//{"priority":"HIGH","_csclass":"org.cyanspring.event.RemoteUserLoginRequestEvent","key":null,"txId":"TX89862982-20200525-101342-685-13","channel":"Windows,chrome/69.0.3947.100","password":"","encrypted":true,"user":"","agent":""}
	AgentNum     string
	PasswordNum  string
	AccountIdNum string

	//交易品种 目前交易所支持品种XBTCUSD XETHUSD XEOSUSD XBCHUSD XLTCUSD
	Symbol string
}

func NewBcPrivate()Bc{
	return Bc{}
}

func (bc *Bc) LogIn() error {
	type tmp struct {
		Agent     string `json:"agent"`
		Channel   string `json:"channel"`
		Encrypted string `json:"encrypted"`
		Key       string `json:"key"`
		Password  string `json:"password"`
		Priority  string `json:"priority"`
		TxId      string `json:"txId"`
		User      string `json:"user"`
		Csclass   string `json:"csclass"`
	}
	sendMessage := tmp{
		Agent:     bc.AgentNum,
		Channel:   "Windows,chrome/69.0.3947.100",
		Encrypted: "true",
		Password:  bc.PasswordNum,
		Priority:  "HIGH",
		TxId:      "TX20200323-112235-253-12", // 时间设置不影响使用
		User:      bc.AgentNum,
		Csclass:   "org.cyanspring.event.RemoteUserLoginRequestEvent",
	}

	send, _ := json.Marshal(sendMessage)
	send2 := strings.Replace(string(send), "csclass", "_csclass", -1)
	err := bc.Conn.WriteMessage(websocket.TextMessage, []byte(send2))

	if err != nil {
		return err
	}
	return nil
}

func (bc *Bc) Work() error {

	var err error
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.user.UserSettingRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094350-826-3\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.user.UserSettingUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.user.UserNoticeUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.user.UserKickOutUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.wallet.WalletListRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094350-826-4\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.account.UserAccountRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094350-827-5\",\"forUser\":\"89862982\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.notice.TargetNoticeUserUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.account.AccountUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.symbol.SymbolSettingUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.symbol.SymbolSettingRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094351-205-17\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.account.AccountSettingRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094351-206-18\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.account.AccountSettingUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.order.OrderUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.position.PositionUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.order.ExecutionUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.position.PositionSnapshotRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094351-209-19\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.order.ActiveOrderSnapshotRequestEvent\",\"key\":null,\"txId\":\"TX20200417-094351-210-20\"}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.risk.OrderRiskSettingUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.risk.PositionRiskSettingUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)
	err = bc.Conn.WriteMessage(websocket.TextMessage, []byte("{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.wallet.WalletUpdateEvent\",\"subKey\":null}"))
	if err != nil {
		return err
	}

	return nil
}

func (bc *Bc) LoadMessage(message string) interface{} {
	Position2struct := PositionData{}
	err := json.Unmarshal([]byte(message), &Position2struct)
	if err == nil && Position2struct.Valid() {
		return &Position2struct
	}

	PositionUpdate2struct := PositionUpdateData{}
	err = json.Unmarshal([]byte(message), &PositionUpdate2struct)
	if err == nil && PositionUpdate2struct.Valid() {
		return &PositionUpdate2struct
	}

	Exec2struct := Exec{}
	err = json.Unmarshal([]byte(message), &Exec2struct)
	if err == nil && Exec2struct.Valid() {
		return &Exec2struct
	}

	return nil
}

func (bc *Bc) ConnAndReceive() {

	defer func() {
		if err := recover(); err != nil {

			//标识此时进行重启避免在重启过程中接口出现panic
			PanicBool = true

			log.Println("Bitcoke私有接口进行重启,err :" + fmt.Sprint(err))
			err = bc.Conn.Close()
			if err != nil {
				panic("关闭Bitcoke私有接口出现错误,err :" + fmt.Sprint(err))
			}
			go bc.ConnAndReceive() //回调自启动

		}
	}()

	symbol = bc.Symbol

	bc.PositionChan = make(chan *PositionData, 10)
	bc.ExecChan = make(chan *Exec, 10)
	bc.PositionUpdateChan = make(chan *PositionUpdateData, 10)

	for {
		c, _, err := websocket.DefaultDialer.Dial(WsUrl, nil)
		if err != nil {
			fmt.Println("connect err")
			fmt.Println(err)
		} else {
			bc.Conn = c
			break
		}

	}

	if DealGoBool == false {
		go bc.Deal()
		DealGoBool = true
	}

	err := bc.LogIn()
	if err != nil {
		panic("登录发生消息出现错误 , err :" + fmt.Sprint(err))
	}

	err = bc.Work()
	if err != nil {
		panic("登录后订阅消息出现错误 , err :" + fmt.Sprint(err))
	}

	for {

		_, buf, err := bc.Conn.ReadMessage()
		if err != nil {
			log.Println("接收数据出现err , err :" + fmt.Sprint(err))
			panic("接收数据出现err")
		}

		message := string(buf)
		message = strings.Replace(message, "_csclass", "csclass", -1)
		parse := bc.LoadMessage(message)

		//log.Println(message)
		switch parse.(type) {
		case *PositionData:
			bc.PositionChan <- parse.(*PositionData)
		case *Exec:
			bc.ExecChan <- parse.(*Exec)
		case *PositionUpdateData:
			bc.PositionUpdateChan <- parse.(*PositionUpdateData)
		}

		PanicBool = false

	}

}
