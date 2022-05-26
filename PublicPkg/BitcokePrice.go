package PublicPkg

import (
	"encoding/json"
	goLog "github.com/RandomEstimate/go-log"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

const (
	WsUrl         = "wss://www.bitcoke.com/prod1/market/"
	ReconnectTime = time.Hour * 6
)

type BcWsP struct {
	conn *websocket.Conn

	orderBookDataChan chan *orderBookData

	symbolList map[string]struct{}
	// 存储行情
	orderBook map[string]orderBook

	log  *goLog.FileLog
	once *sync.Once
	m    *sync.RWMutex
}

func NewPublic(log *goLog.FileLog) (*BcWsP, error) {
	return &BcWsP{
		orderBookDataChan: make(chan *orderBookData, 1000),
		symbolList:        make(map[string]struct{}),
		orderBook:         make(map[string]orderBook),
		once:              &sync.Once{},
		m:                 &sync.RWMutex{},
		log:               log,
	}, nil
}

func (a *BcWsP) deal() {

	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case data := <-a.orderBookDataChan:
			a.handleOrder(data)

		case <-t.C:
			a.scan()


		}
	}
}

func (a *BcWsP) LoadMessage(message string) interface{} {
	orderData := orderBookData{}
	err := json.Unmarshal([]byte(message), &orderData)
	if err == nil && orderData.fullValid() {
		return &orderData
	}

	if err == nil && orderData.updateValid() {
		return &orderData
	}
	return nil
}

func (a *BcWsP) Start(w *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {

			err = a.conn.Close()
			if err != nil {
				a.log.E("conn.Close error : %v", err)
				//panic("conn.Close error")
			}

			// clear all
			for len(a.orderBookDataChan) > 0 {
				<-a.orderBookDataChan
			}
			a.reconnectUpdate()
			go a.Start(nil)
		}
	}()

	for {
		c, _, err := websocket.DefaultDialer.Dial(WsUrl, nil)
		if err != nil {
			if a.log != nil {
				a.log.E("connect error: %v", err)
			}
		} else {
			a.conn = c
			break
		}
		time.Sleep(time.Second * 3)

	}

	a.once.Do(func() {
		go a.deal()
	})

	for v := range a.symbolList {
		if a.OrderSymbol(v) != nil {
			panic("orderSymbol error")
		}

	}

	if w != nil {
		w.Done()
	}

	t := time.NewTimer(ReconnectTime)
	for {
		select {
		case <-t.C:
			if a.log != nil {
				a.log.I("reconnect time. ")
			}
			panic("reconnect")
		default:
			_, buf, err := a.conn.ReadMessage()
			if err != nil {
				if a.log != nil {
					a.log.E("receive data error : %v", err)
				}
				panic("receive data error")
			}

			message := string(buf)
			message = strings.Replace(message, "_csclass", "csclass", -1)
			parse := a.LoadMessage(message)
			//fmt.Println(message)
			switch parse.(type) {
			case *orderBookData:
				a.orderBookDataChan <- parse.(*orderBookData)
			}

		}

	}

}
