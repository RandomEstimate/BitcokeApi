package PublicPkg

import (
	"errors"
	"time"
)

func (a *BcWsP) GetPriceMap() (DepthMap, error) {

	if a.LiveBool == true {
		a.Channel2 <- 1
		select {
		case dataMap := <-a.Channel3:
			//fmt.Println("dataMap",dataMap)
			tmp := DepthMap{
				Bid1:       dataMap.buy[0],
				Bid1Amount: dataMap.buyAmount[0],
				Bid1Count:  dataMap.buyCount[0],
				Bid2:       dataMap.buy[1],
				Bid2Amount: dataMap.buyAmount[1],
				Bid2Count:  dataMap.buyCount[1],
				Bid3:       dataMap.buy[2],
				Bid3Amount: dataMap.buyAmount[2],
				Bid3Count:  dataMap.buyCount[2],
				Ask1:       dataMap.sell[0],
				Ask1Amount: dataMap.sellAmount[0],
				Ask1Count:  dataMap.sellCount[0],
				Ask2:       dataMap.sell[1],
				Ask2Amount: dataMap.sellAmount[1],
				Ask2Count:  dataMap.sellCount[1],
				Ask3:       dataMap.sell[2],
				Ask3Amount: dataMap.sellAmount[2],
				Ask3Count:  dataMap.sellCount[2],
				Time:       time.Now(),
				TimeStamp:  float64(time.Now().UnixNano()) / 1000000000,
				TimeStr:    time.Now().Format("2006-01-02 15:04:05.0000"),
			}
			return tmp, nil
		case <-time.After(3 * time.Second):
			err := errors.New("请求错误，未能接受到行情数据")
			return DepthMap{}, err
		}
	}

	err := errors.New("行情接口正在重启")
	return DepthMap{}, err
}
