package PublicPkg

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

const (
	WsUrl = "wss://www.nanotradetech.com/prod1/market/"
)

var (
	//避免在自启动时重复启动Deal函数导致资源竞争
	DealGoBool bool = false
)

type BcWsP struct {
	Symbol    string
	LiveBool  bool
	Conn      *websocket.Conn
	OrderChan chan *OrderBookData

	// 存储行情深度数据数组
	BuyPriceOrderBook  []float64
	BuyQtyOrderBook    []float64
	BuyCount           []float64
	SellPriceOrderBook []float64
	SellQtyOrderBook   []float64
	SellCount          []float64

	//避免同时数据更新和外部接口同时访问深度数据采用channel机制 Chan2是发送int 请求数据， Chan3接收数据
	Channel2 chan int
	Channel3 chan PriceStruct
}

func NewPublic(Symbol string) (BcWsP, error) {
	return BcWsP{
		Symbol:   Symbol,
		LiveBool: false,
	}, nil
}

func (a *BcWsP) Work() {
	a.BuyPriceOrderBook = make([]float64, 0, 50)
	a.BuyQtyOrderBook = make([]float64, 0, 50)
	a.BuyCount = make([]float64, 0, 50)
	a.SellPriceOrderBook = make([]float64, 0, 50)
	a.SellQtyOrderBook = make([]float64, 0, 50)
	a.SellCount = make([]float64, 0, 50)
	a.OrderChan = make(chan *OrderBookData, 5)
	a.Channel2 = make(chan int, 10)
	a.Channel3 = make(chan PriceStruct, 10)
}

func (a *BcWsP) Deal() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Deal线程 err :" + fmt.Sprint(err))
			go a.Deal()
		}
	}()

	for {
		select {
		case data := <-a.OrderChan:
			a.HandleOrder(data)

			// 防止出现重启时访问接口导致数组越界
			if a.LiveBool == false {
				a.LiveBool = true
			}

		case <-a.Channel2:
			a.HandlePriceReq()
		}
	}
}

func (a *BcWsP) LoadMessage(message string) interface{} {
	orderData := OrderBookData{}
	err := json.Unmarshal([]byte(message), &orderData)
	if err == nil && orderData.FullValid() {
		return &orderData
	}

	if err == nil && orderData.UpdateValid() {
		return &orderData
	}
	return nil
}

func (a *BcWsP) LogIn() error {

	sendMessage := "{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.marketdata.DepthUpdateEvent\",\"subKey\":\"" + a.Symbol + "\"}"
	sendMessage2 := "{\"priority\":\"HIGH\",\"_csclass\":\"org.cyanspring.exbusiness.event.marketdata.DepthRequestEvent\",\"key\":\"" + a.Symbol + "\",\"txId\":\"TX20200413-213011-836-5\"}"
	sendMessage3 := "{\"priority\":\"NORMAL\",\"_csclass\":\"org.cyanspring.event.RemoteSubscribeEvent\",\"clazz\":\"org.cyanspring.exbusiness.event.marketdata.DepthFullUpdateEvent\",\"subKey\":\"" + a.Symbol + "\"}"

	if err := a.Conn.WriteMessage(websocket.TextMessage, []byte(sendMessage)); err != nil {
		log.Println("send ERR")
		return err
	}
	if err := a.Conn.WriteMessage(websocket.TextMessage, []byte(sendMessage2)); err != nil {
		log.Println("send ERR2")
		return err
	}
	if err := a.Conn.WriteMessage(websocket.TextMessage, []byte(sendMessage3)); err != nil {
		log.Println("send ERR3")
		return err
	}
	return nil

}

func (a *BcWsP) ConnAndReceive() {
	defer func() {
		if err := recover(); err != nil {
			a.LiveBool = false
			log.Println("行情数据接口重启 err :" + fmt.Sprint(err))
			err = a.Conn.Close()
			if err != nil {
				log.Println("行情数据关闭错误 err :" + fmt.Sprint(err))
				panic("行情数据关闭错误")
			}
			go a.ConnAndReceive()
		}
	}()

	for {
		c, _, err := websocket.DefaultDialer.Dial(WsUrl, nil)
		if err != nil {
			log.Println("connect2 err")
			log.Println(err)
		} else {
			a.Conn = c
			break
		}

	}

	a.Work()

	if DealGoBool == false {
		go a.Deal()
		DealGoBool = true
	}

	err := a.LogIn()
	if err != nil {
		panic("无法登陆错误")
	}

	for {
		_, buf, err := a.Conn.ReadMessage()
		if err != nil {
			log.Println("接收数据出现err , err :" + fmt.Sprint(err))
			panic("接收数据出现err")
		}

		message := string(buf)
		message = strings.Replace(message, "_csclass", "csclass", -1)
		parse := a.LoadMessage(message)

		switch parse.(type) {
		case *OrderBookData:
			a.OrderChan <- parse.(*OrderBookData)
		}
	}

}
