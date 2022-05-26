package PrivatePkg

import (
	"BitcokeApi/PublicPkg"
	"encoding/json"
	"fmt"
	goLog "github.com/RandomEstimate/go-log"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WsUrl         string = "wss://www.bitcoke.com/prod1/trade/"
	ReconnectTime        = time.Hour * 6
)

type BcWsP2 struct {
	conn *websocket.Conn

	accountInfo map[string]accountInfo

	tradingDataChan chan positionData
	orderUpdateChan chan execData
	orderChan       chan OrderMode

	positionSearchCD     atomic.Value // time.Time
	positionUpdateSwitch atomic.Value // int

	//配置个人信息 F12抓取
	AgentNum     string
	PasswordNum  string
	AccountIdNum string

	BcWsP *PublicPkg.BcWsP

	log *goLog.FileLog
	// deal 单次启动
	once *sync.Once

	m          *sync.RWMutex
	safeTradeM *sync.RWMutex

	execInfoChan chan int
}

func NewBcPrivate(AgentNum string, PasswordNum string, AccountIdNum string, BcWsP *PublicPkg.BcWsP, log *goLog.FileLog) *BcWsP2 {

	tmp := BcWsP2{
		accountInfo:          make(map[string]accountInfo),
		tradingDataChan:      make(chan positionData, 100),
		orderUpdateChan:      make(chan execData, 100),
		orderChan:            make(chan OrderMode, 10),
		AgentNum:             AgentNum,
		PasswordNum:          PasswordNum,
		AccountIdNum:         AccountIdNum,
		BcWsP:                BcWsP,
		once:                 new(sync.Once),
		m:                    new(sync.RWMutex),
		safeTradeM:           new(sync.RWMutex),
		positionSearchCD:     atomic.Value{},
		positionUpdateSwitch: atomic.Value{},
		execInfoChan:         make(chan int, 10),
	}

	tmp.positionSearchCD.Store(time.Now())
	tmp.positionUpdateSwitch.Store(true)

	if log != nil {
		tmp.log = log
	}

	return &tmp
}

func (a *BcWsP2) logIn() error {
	sendMessage := loginMode{
		Agent:     a.AgentNum,
		Channel:   "Windows,chrome/69.0.3947.100",
		Encrypted: "true",
		Token:     a.PasswordNum,
		Priority:  "HIGH",
		TxId:      "TX20200323-112235-253-12", // 时间设置不影响使用
		User:      a.AgentNum,
		Csclass:   "org.cyanspring.event.RemoteUserLoginRequestEvent",
	}

	buf, _ := json.Marshal(sendMessage)
	send2 := strings.Replace(string(buf), "csclass", "_csclass", -1)
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(send2)); err != nil {
		if a.log != nil {
			a.log.E("login error : %v. ", err)
		}
		return err
	}

	return nil
}

func (a *BcWsP2) sub() error {

	messageList := make([]string, 0)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.user.UserSettingRequestEvent","key":null,"txId":"TX20200417-094350-826-3"}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.user.UserSettingUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.user.UserNoticeUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.user.UserKickOutUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.wallet.WalletListRequestEvent","key":null,"txId":"TX20200417-094350-826-4"}`)
	messageList = append(messageList, fmt.Sprintf(`{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.account.UserAccountRequestEvent","key":null,"txId":"TX20200417-094350-827-5","forUser":"%v"}`, a.AgentNum))
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.notice.TargetNoticeUserUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.account.AccountUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.symbol.SymbolSettingUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.symbol.SymbolSettingRequestEvent","key":null,"txId":"TX20200417-094351-205-17"}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.account.AccountSettingRequestEvent","key":null,"txId":"TX20200417-094351-206-18"}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.account.AccountSettingUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.order.OrderUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.position.PositionUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.order.ExecutionUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.position.PositionSnapshotRequestEvent","key":null,"txId":"TX20200417-094351-209-19"}`)
	messageList = append(messageList, `{"priority":"HIGH","_csclass":"org.cyanspring.exbusiness.event.order.ActiveOrderSnapshotRequestEvent","key":null,"txId":"TX20200417-094351-210-20"}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.risk.OrderRiskSettingUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.risk.PositionRiskSettingUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.wallet.WalletUpdateEvent","subKey":null}`)
	messageList = append(messageList, `{"priority":"NORMAL","_csclass":"org.cyanspring.event.RemoteSubscribeEvent","clazz":"org.cyanspring.exbusiness.event.data.TradingDataUpdateEvent","subKey":null}`)

	for _, v := range messageList {
		if err := a.conn.WriteMessage(websocket.TextMessage, []byte(v)); err != nil {
			if a.log != nil {
				a.log.E("send message error : %v .", err)
			}
			return err
		}
		time.Sleep(time.Millisecond * 30)
	}

	return nil
}

func (a *BcWsP2) loadMessage(message string) interface{} {
	Position2struct := positionData{}
	err := json.Unmarshal([]byte(message), &Position2struct)

	if err == nil && Position2struct.valid() {
		return &Position2struct
	}

	Exec2struct := execData{}
	err = json.Unmarshal([]byte(message), &Exec2struct)
	if err == nil && Exec2struct.valid() {
		return &Exec2struct
	}

	return nil
}

func (a *BcWsP2) Start(w *sync.WaitGroup) {

	defer func() {
		if err := recover(); err != nil {
			err = a.conn.Close()
			if err != nil {
				if a.log != nil {
					a.log.E("close error : %v", err)
				}
				//panic(err)
			}
			go a.Start(nil)

		}
	}()

	for {
		c, _, err := websocket.DefaultDialer.Dial(WsUrl, nil)
		if err != nil {
			if a.log != nil {
				a.log.E("connect error : %v", err)
			}
			time.Sleep(time.Second * 3)
			continue
		}
		a.conn = c
		break

	}

	a.once.Do(func() {
		go a.deal()
	})

	err := a.logIn()
	if err != nil {
		panic("logIn panic")
	}

	err = a.sub()
	if err != nil {
		panic("sub panic")
	}

	if w != nil {
		w.Done()
	}

	t := time.NewTimer(ReconnectTime)
	for {
		select {
		case <-t.C:
			if a.log != nil {
				a.log.I("ReconnectTime .")
			}
			panic("Reconnect")
		default:

		}
		_, buf, err := a.conn.ReadMessage()
		if err != nil {
			if a.log != nil {
				a.log.E("receive error : %v", err)
			}
			panic("receive error")
		}

		message := string(buf)
		message = strings.Replace(message, "_csclass", "csclass", -1)
		parse := a.loadMessage(message)

		//fmt.Println(message)
		switch parse.(type) {
		case *positionData:
			a.tradingDataChan <- *parse.(*positionData)
		case *execData:
			a.orderUpdateChan <- *parse.(*execData)

		}

	}

}

func (a *BcWsP2) deal() {

	t := time.NewTicker(time.Second * 30)
	for {
		select {
		case data := <-a.orderChan:
			a.handleTrade(data)
		case data := <-a.tradingDataChan:
			a.handlePositionData(data)
		case data := <-a.orderUpdateChan:
			a.handleExec(data)
		case <-t.C:
			a.pong()
		}
	}

}
