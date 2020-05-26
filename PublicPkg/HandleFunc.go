package PublicPkg

//冒泡排序
func Sort(UpOrDown bool,dataPrice,dataQty,dataCount []float64)([]float64,[]float64,[]float64){
	// true 小的在前 适用于sell ask
	if UpOrDown{
		for i:=0;i<=len(dataPrice)-1;i++{
			for j:=i+1;j<=len(dataPrice)-1;j++{
				if dataPrice[j] < dataPrice[i]{
					t := dataPrice[i]
					dataPrice[i] = dataPrice[j]
					dataPrice[j] = t

					t = dataQty[i]
					dataQty[i] = dataQty[j]
					dataQty[j] = t

					t = dataCount[i]
					dataCount[i] = dataCount[j]
					dataCount[j] = t
				}
			}

		}
	}else{
		for i:=0;i<=len(dataPrice)-1;i++{
			for j:=i+1;j<=len(dataPrice)-1;j++{
				if dataPrice[j] > dataPrice[i]{
					t := dataPrice[i]
					dataPrice[i] = dataPrice[j]
					dataPrice[j] = t

					t = dataQty[i]
					dataQty[i] = dataQty[j]
					dataQty[j] = t

					t = dataCount[i]
					dataCount[i] = dataCount[j]
					dataCount[j] = t
				}
			}

		}

	}
	return dataPrice,dataQty,dataCount
}

func (a *BcWsP)HandleOrder(data *OrderBookData) {

	// Bid价格排序
	for _,v := range(data.BuyDepth){
		isExist := false
		settleNum := 0
		for i,v1 := range(a.BuyPriceOrderBook){
			if v1 == v.Price {
				isExist = true
				settleNum = i
				break
			}
		}
		if isExist {
			if v.Qty == 0 {
				a.BuyPriceOrderBook = append(a.BuyPriceOrderBook[0:settleNum],a.BuyPriceOrderBook[settleNum+1:]...)
				a.BuyQtyOrderBook = append(a.BuyQtyOrderBook[0:settleNum],a.BuyQtyOrderBook[settleNum+1:]...)
				a.BuyCount = append(a.BuyCount[0:settleNum],a.BuyCount[settleNum+1:]...)
			}else {
				a.BuyQtyOrderBook[settleNum] = v.Qty
				a.BuyCount[settleNum] = v.Count
			}
		}else {
			a.BuyPriceOrderBook = append(a.BuyPriceOrderBook, v.Price)
			a.BuyQtyOrderBook = append(a.BuyQtyOrderBook, v.Qty)
			a.BuyCount = append(a.BuyCount, v.Count)
		}
	}

	// Ask价格排序
	for _,v := range(data.SellDepth){
		isExist := false
		settleNum := 0
		for i,v1 := range(a.SellPriceOrderBook){
			if v1 == v.Price {
				isExist = true
				settleNum = i
				break
			}
		}
		if isExist {
			if v.Qty == 0 {
				a.SellPriceOrderBook = append(a.SellPriceOrderBook[0:settleNum],a.SellPriceOrderBook[settleNum+1:]...)
				a.SellQtyOrderBook = append(a.SellQtyOrderBook[0:settleNum],a.SellQtyOrderBook[settleNum+1:]...)
				a.SellCount = append(a.SellCount[0:settleNum],a.SellCount[settleNum+1:]...)
			}else {
				a.SellQtyOrderBook[settleNum] = v.Qty
				a.SellCount[settleNum] = v.Count
			}
		}else {
			a.SellPriceOrderBook = append(a.SellPriceOrderBook, v.Price)
			a.SellQtyOrderBook = append(a.SellQtyOrderBook, v.Qty)
			a.SellCount = append(a.SellCount, v.Count)
		}
	}

	//排序
	a.BuyPriceOrderBook,a.BuyQtyOrderBook,a.BuyCount = Sort(false,a.BuyPriceOrderBook,a.BuyQtyOrderBook,a.BuyCount)
	a.SellPriceOrderBook,a.SellQtyOrderBook,a.SellCount = Sort(true,a.SellPriceOrderBook,a.SellQtyOrderBook,a.SellCount)
}

func (a *BcWsP)HandlePriceReq(){
	tmp := PriceStruct{
		buy:        a.BuyPriceOrderBook,
		buyAmount:  a.BuyQtyOrderBook,
		buyCount:   a.BuyCount,
		sell:       a.SellPriceOrderBook,
		sellAmount: a.SellQtyOrderBook,
		sellCount:  a.SellCount,
	}
	a.Channel3 <- tmp
}
