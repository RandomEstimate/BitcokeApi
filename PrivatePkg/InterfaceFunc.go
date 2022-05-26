package PrivatePkg

import (
	"fmt"
	"time"
)

func (a *BcWsP2) SafeTrade(TradeInfo OrderMode) error {
	a.safeTradeM.Lock()
	defer a.safeTradeM.Unlock()

	for len(a.execInfoChan) != 0 {
		<-a.execInfoChan
	}

	a.orderChan <- TradeInfo

	select {
	case <-time.After(time.Second * 3):
		if a.log != nil {
			a.log.E("safeTrade timeout")
		}
		a.positionUpdateSwitch.Store(true)
		return fmt.Errorf("safeTrade timeout")
	case <-a.execInfoChan:
		return nil
	}
}

func (a *BcWsP2) GetPosition(symbol string) (float64, float64, error) {




	flag := false
	for i := 0; i < 3; i++ {
		a.m.Lock()
		if a.positionUpdateSwitch.Load().(bool) && time.Since(a.positionSearchCD.Load().(time.Time)) > 4*time.Second {
			flag = true
			a.m.Unlock()
			break
		}
		a.m.Unlock()
		time.Sleep(time.Second * 2)
	}
	if !flag {
		if a.log != nil {
			a.log.E("positionUpdateSwitch or positionSearchCD error | %v  %v", a.positionUpdateSwitch.Load().(bool), time.Since(a.positionSearchCD.Load().(time.Time)))
		}
		return -1, -1, fmt.Errorf("positionUpdateSwitch or positionSearchCD error")
	}

	a.m.Lock()
	defer a.m.Unlock()
	if d, ok := a.accountInfo[symbol]; !ok {
		return 0, 0, nil
	} else {
		return d.position, d.openPrice, nil
	}

}
