### 统一更新交易所接口 -- Bitcoke

> 1. 交易API说明，由于Bitcoke目前只支持只读API无法进行交易，因此没有REST接口，目前Bitcoke采用Ws发送消息指令进行下单成交。
> 2. 由于Bitcoke采用Ws进行下单的机制，因此其私有Ws接口与交易接口为同一 类 class（struct）
> 3. 目前Bitcoke可以进行挂单Limit以及市价单Market进行交易（就目前使用情况来看，直接进行市价交易会出现巨大滑点，可能交易所的做市交易是对手盘，为了避免这一情况，目前统一采用对手限价ask1、bid1单进行成交，效果还不错。）
> 4. 目前所设计的接口针对的是单向开仓，没有平仓的接口，如果希望平仓只需要反向开仓即可，所以需要在Bitcoke的Web上进行设置为单向持仓否则接口报错。

#### 私有接口说明

```go

//下单接口,非挂单，目前暂时不支持做事（挂单）交易
func Trade(TradeInfo OrderMode) (error){
    //返回err 目前err返回当前格式是否正确且当前下单Order传入进行执行的队列，不保证一定能够触发交易
    //经过测试 当快速对交易所进行发送连续交易指令，可能会有交易指令被默认取消
}

//交易品种的当前仓位（方向 、仓位 、均价）
func GetPostion() (string,float,float,error){
	//返回开仓方向 Buy/Sell 仓位数 单位：张 = 1USD  开仓均价 err
}

//获取最后更改仓位时间
func GetUpdateTime() (time.Time){
    //返回time.Time 
}

//通过Ws更新当前仓位
func RepairPosition() (error){
    //返回err
}

```

#### 私有接口出现交易事件返回

```go
//当触发事件时 通过结构体内部Channel进行返回基本成交信息 、 当然在初始时可以设置不开启这个功能
```

