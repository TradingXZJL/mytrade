package mytrade

import (
	"errors"
	"fmt"
	"github.com/Hongssd/myokxapi"
	"github.com/shopspring/decimal"
)

// 查询订单返回结果处理
func (o *OkxTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrdersPendingRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	orders := make([]*Order, 0, len(res.Data))
	for _, r := range res.Data {
		orderType, timeInForce := o.okxConverter.FromOKXOrderType(r.OrdType)
		var isMargin, isIsolated bool
		if r.InstType == OKX_AC_MARGIN.String() {
			r.InstType = OKX_AC_SPOT.String()
			isMargin = true
			if r.TdMode == OKX_MARGIN_MODE_ISOLATED {
				isIsolated = true
			}
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   r.InstType,
			Symbol:        r.InstId,
			IsMargin:      isMargin,
			IsIsolated:    isIsolated,
			IsAlgo:        req.IsAlgo,
			Price:         r.Px,
			Quantity:      r.Sz,
			ExecutedQty:   r.FillSz,
			AvgPrice:      r.AvgPx,
			Status:        o.okxConverter.FromOKXOrderStatus(r.State, false),
			Type:          orderType,
			Side:          o.okxConverter.FromOKXOrderSide(r.Side),
			PositionSide:  o.okxConverter.FromOKXPositionSide(r.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    stringToBool(r.ReduceOnly),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
			RealizedPnl:   r.Pnl,
		}
		if r.AttachAlgoOrds != nil && len(r.AttachAlgoOrds) > 0 {
			order.AttachTpOrdPrice = r.AttachAlgoOrds[0].TpOrdPx
			order.AttachTpTriggerPrice = r.AttachAlgoOrds[0].TpTriggerPx
			order.AttachSlOrdPrice = r.AttachAlgoOrds[0].SlOrdPx
			order.AttachSlTriggerPrice = r.AttachAlgoOrds[0].SlTriggerPx
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromQueryOrderGet(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderGetRes]) (*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	r := res.Data[0]

	orderType, timeInForce := o.okxConverter.FromOKXOrderType(r.OrdType)
	var isMargin, isIsolated bool
	if r.InstType == OKX_AC_MARGIN.String() {
		r.InstType = OKX_AC_SPOT.String()
		isMargin = true
		if r.TdMode == OKX_MARGIN_MODE_ISOLATED {
			isIsolated = true
		}
	}
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   r.InstType,
		Symbol:        r.InstId,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		IsAlgo:        req.IsAlgo,
		OrderAlgoType: req.OrderAlgoType,
		Price:         r.Px,
		Quantity:      r.Sz,
		ExecutedQty:   r.FillSz,
		AvgPrice:      r.AvgPx,
		Status:        o.okxConverter.FromOKXOrderStatus(r.State, false),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(r.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(r.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(r.ReduceOnly),
		CreateTime:    stringToInt64(r.CTime),
		UpdateTime:    stringToInt64(r.UTime),
	}
	if r.AttachAlgoOrds != nil && len(r.AttachAlgoOrds) > 0 {
		order.AttachTpOrdPrice = r.AttachAlgoOrds[0].TpOrdPx
		order.AttachTpTriggerPrice = r.AttachAlgoOrds[0].TpTriggerPx
		order.AttachSlOrdPrice = r.AttachAlgoOrds[0].SlOrdPx
		order.AttachSlTriggerPrice = r.AttachAlgoOrds[0].SlTriggerPx
	}
	return order, nil
}
func (o *OkxTradeEngine) handleOrdersFromQueryOrderGet(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderHistoryRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}

	var orders []*Order
	for _, r := range res.Data {
		orderType, timeInForce := o.okxConverter.FromOKXOrderType(r.OrdType)
		var isMargin, isIsolated bool
		if r.InstType == OKX_AC_MARGIN.String() {
			r.InstType = OKX_AC_SPOT.String()
			isMargin = true
			if r.TdMode == OKX_MARGIN_MODE_ISOLATED {
				isIsolated = true
			}
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   r.InstType,
			Symbol:        r.InstId,
			IsMargin:      isMargin,
			IsIsolated:    isIsolated,
			IsAlgo:        req.IsAlgo,
			OrderAlgoType: req.OrderAlgoType,
			Price:         r.Px,
			Quantity:      r.Sz,
			ExecutedQty:   r.FillSz,
			AvgPrice:      r.AvgPx,
			Status:        o.okxConverter.FromOKXOrderStatus(r.State, false),
			Type:          orderType,
			Side:          o.okxConverter.FromOKXOrderSide(r.Side),
			PositionSide:  o.okxConverter.FromOKXPositionSide(r.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    stringToBool(r.ReduceOnly),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
			RealizedPnl:   r.Pnl,
		}
		if r.AttachAlgoOrds != nil && len(r.AttachAlgoOrds) > 0 {
			order.AttachTpOrdPrice = r.AttachAlgoOrds[0].TpOrdPx
			order.AttachTpTriggerPrice = r.AttachAlgoOrds[0].TpTriggerPx
			order.AttachSlOrdPrice = r.AttachAlgoOrds[0].SlOrdPx
			order.AttachSlTriggerPrice = r.AttachAlgoOrds[0].SlTriggerPx
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeFillsRes]) ([]*Trade, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	trades := make([]*Trade, 0, len(res.Data))

	for _, r := range res.Data {
		quoteQty := decimal.RequireFromString(r.FillPx).Mul(decimal.RequireFromString(r.FillSz))
		isMaker := r.ExecType == "M"
		if r.InstType == OKX_AC_MARGIN.String() {
			r.InstType = OKX_AC_SPOT.String()
		}
		trade := &Trade{
			Exchange:     OKX_NAME.String(),
			AccountType:  r.InstType,
			Symbol:       r.InstId,
			TradeId:      r.TradeId,
			OrderId:      r.OrdId,
			Price:        r.FillPx,
			Quantity:     r.FillSz,
			QuoteQty:     quoteQty.String(),
			Side:         o.okxConverter.FromOKXOrderSide(r.Side),
			PositionSide: o.okxConverter.FromOKXPositionSide(r.PosSide),
			FeeAmount:    r.Fee,
			FeeCcy:       r.FeeCcy,
			RealizedPnl:  r.FillPnl,
			IsMaker:      isMaker,
			Timestamp:    stringToInt64(r.FillTime),
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// 单订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromOrderCreate(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderPostRes]) (*Order, error) {

	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderAmend(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeAmendOrderRes]) (*Order, error) {

	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderCancel(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeCancelOrderRes]) (*Order, error) {

	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

// 批量订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromBatchOrderCreate(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeBatchOrdersRes]) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClOrdId, r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	if errStr != "" {
		return orders, fmt.Errorf("[%s]%s: [%s]", res.Code, res.Msg, errStr)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromBatchOrderAmend(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeAmendBatchOrdersRes]) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClOrdId, r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	if errStr != "" {
		return orders, fmt.Errorf("[%s]%s: [%s]", res.Code, res.Msg, errStr)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromBatchOrderCancel(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeCancelBatchOrdersRes]) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClOrdId, r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	if errStr != "" {
		return orders, fmt.Errorf("[%s]%s: [%s]", res.Code, res.Msg, errStr)
	}
	return orders, nil
}

// 订单推送处理
func (o *OkxTradeEngine) handleOrderFromWsOrder(order myokxapi.WsOrders) *Order {
	orderType, timeInForce := o.okxConverter.FromOKXOrderType(order.OrdType)
	var IsMargin, IsIsolated bool

	//accountType := order.Orders.InstType

	if order.Orders.InstType == OKX_AC_MARGIN.String() {
		//accountType = OKX_AC_SPOT.String()
		IsMargin = true
		if order.Orders.TdMode == OKX_MARGIN_MODE_ISOLATED {
			IsIsolated = true
		}
	}

	return &Order{
		Exchange: OKX_NAME.String(),
		//AccountType:   order.Orders.InstType,
		Symbol:        order.Orders.InstId,
		IsMargin:      IsMargin,
		IsIsolated:    IsIsolated,
		OrderId:       order.OrdId,
		ClientOrderId: order.ClOrdId,
		Price:         order.Px,
		Quantity:      order.Sz,
		ExecutedQty:   order.FillSz,
		AvgPrice:      order.AvgPx,
		Status:        o.okxConverter.FromOKXOrderStatus(order.State, false),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(order.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(order.ReduceOnly),
		CreateTime:    stringToInt64(order.CTime),
		UpdateTime:    stringToInt64(order.UTime),
		RealizedPnl:   order.FillPnl,

		ErrorMsg:  order.Msg,
		ErrorCode: order.Code,
	}

}

// 策略订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromOrderAlgoCreate(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderAlgoPostRes]) (*Order, error) {
	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]

	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.AlgoId,
		ClientOrderId: r.AlgoClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		IsAlgo:        true,
	}

	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderAlgoAmend(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeAmendOrderAlgoRes]) (*Order, error) {
	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]

	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.AlgoId,
		ClientOrderId: r.AlgoClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		IsAlgo:        true,
	}

	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderAlgoCancel(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeCancelOrderAlgoRes]) (*Order, error) {
	if res == nil || len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		//return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
		return nil, fmt.Errorf("[%s]%s: {[%s]:}", res.Code, res.Msg, res.Data[0].SCode)
	}
	r := res.Data[0]

	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.AlgoId,
		ClientOrderId: r.AlgoClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}

	return order, nil
}

// 查询策略订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromQueryOrderAlgo(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderAlgoGetRes]) (*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	order := res.Data[0]
	var orderType OrderType
	timeInForce := TIME_IN_FORCE_GTC
	var px decimal.Decimal
	var triggerPx decimal.Decimal
	var triggerType OrderTriggerType
	var triggerConditionType OrderTriggerConditionType

	if order.TpTriggerPx != "" {
		px, _ = decimal.NewFromString(order.TpOrdPx)
		triggerPx, _ = decimal.NewFromString(order.TpTriggerPx)
		triggerType = ORDER_TRIGGER_TYPE_TAKE_PROFIT

		if order.Side == OKX_ORDER_SIDE_BUY {
			//止盈买入 价格下穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		} else {
			//止盈卖出 价格上穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		}
	} else if order.SlTriggerPx != "" {
		px, _ = decimal.NewFromString(order.SlOrdPx)
		triggerPx, _ = decimal.NewFromString(order.SlTriggerPx)
		triggerType = ORDER_TRIGGER_TYPE_STOP_LOSS

		if order.Side == OKX_ORDER_SIDE_BUY {
			//止损买入 价格上穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		} else {
			//止损卖出 价格下穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		}
	}

	if px.Equal(decimal.NewFromInt(-1)) {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var IsMargin, IsIsolated bool
	if order.TdMode != "cash" {
		switch order.TdMode {
		case "cross":
			IsIsolated = false
		case "isolated":
			IsIsolated = true
		}
	}
	if order.InstType == OKX_AC_MARGIN.String() {
		IsMargin = true
		order.InstType = OKX_AC_SPOT.String()
	}

	retOrder := &Order{
		Exchange:      OKX_NAME.String(),
		AccountType:   order.InstType,
		Symbol:        order.InstId,
		IsMargin:      IsMargin,
		IsIsolated:    IsIsolated,
		OrderId:       order.AlgoId,
		ClientOrderId: order.AlgoClOrdId,
		Price:         px.String(),
		Quantity:      order.Sz,
		ExecutedQty:   decimal.Zero.String(),
		AvgPrice:      decimal.Zero.String(),
		Status:        o.okxConverter.FromOKXOrderStatus(order.State, true),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(order.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(order.ReduceOnly),
		CreateTime:    stringToInt64(order.CTime),
		UpdateTime:    stringToInt64(order.UTime),
		RealizedPnl:   decimal.Zero.String(),

		IsAlgo:        true,
		OrderAlgoType: OrderAlgoType(order.OrdType),
	}

	switch order.OrdType {
	case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
		retOrder.TriggerPrice = triggerPx.String()
		retOrder.TriggerType = triggerType
		retOrder.TriggerConditionType = triggerConditionType
	case OKX_ORDER_ALGO_TYPE_OCO:
		retOrder.OcoTpTriggerPrice = order.TpTriggerPx
		retOrder.OcoTpOrdPrice = order.TpOrdPx
		retOrder.OcoSlTriggerPrice = order.SlTriggerPx
		retOrder.OcoSlOrdPrice = order.SlOrdPx
		if order.TpOrdPx == "-1" {
			retOrder.OcoTpOrdType = ORDER_TYPE_MARKET
		} else {
			retOrder.OcoTpOrdType = ORDER_TYPE_LIMIT
		}
		if order.SlOrdPx == "-1" {
			retOrder.OcoSlOrdType = ORDER_TYPE_MARKET
		} else {
			retOrder.OcoSlOrdType = ORDER_TYPE_LIMIT
		}
	}

	return retOrder, nil
}
func (o *OkxTradeEngine) handleOrdersFromQueryOrderAlgo(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderAlgoHistoryRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	var orders []*Order
	for _, order := range res.Data {
		var orderType OrderType
		timeInForce := TIME_IN_FORCE_GTC
		var px decimal.Decimal
		var triggerPx decimal.Decimal
		var triggerType OrderTriggerType
		var triggerConditionType OrderTriggerConditionType

		if order.TpTriggerPx != "" {
			px, _ = decimal.NewFromString(order.TpOrdPx)
			triggerPx, _ = decimal.NewFromString(order.TpTriggerPx)
			triggerType = ORDER_TRIGGER_TYPE_TAKE_PROFIT

			if order.Side == OKX_ORDER_SIDE_BUY {
				//止盈买入 价格下穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
			} else {
				//止盈卖出 价格上穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
			}
		} else if order.SlTriggerPx != "" {
			px, _ = decimal.NewFromString(order.SlOrdPx)
			triggerPx, _ = decimal.NewFromString(order.SlTriggerPx)
			triggerType = ORDER_TRIGGER_TYPE_STOP_LOSS

			if order.Side == OKX_ORDER_SIDE_BUY {
				//止损买入 价格上穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
			} else {
				//止损卖出 价格下穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
			}
		}

		if px.Equal(decimal.NewFromInt(-1)) {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var IsMargin, IsIsolated bool
		if order.TdMode != "cash" {
			switch order.TdMode {
			case "cross":
				IsIsolated = false
			case "isolated":
				IsIsolated = true
			}
		}
		if order.InstType == OKX_AC_MARGIN.String() {
			IsMargin = true
			order.InstType = OKX_AC_SPOT.String()
		}

		appendOrder := &Order{
			Exchange:      OKX_NAME.String(),
			Symbol:        order.InstId,
			AccountType:   order.InstType,
			IsMargin:      IsMargin,
			IsIsolated:    IsIsolated,
			OrderId:       order.AlgoId,
			ClientOrderId: order.AlgoClOrdId,
			Price:         px.String(),
			Quantity:      order.Sz,
			ExecutedQty:   decimal.Zero.String(),
			AvgPrice:      decimal.Zero.String(),
			Status:        o.okxConverter.FromOKXOrderStatus(order.State, true),
			Type:          orderType,
			Side:          o.okxConverter.FromOKXOrderSide(order.Side),
			PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    stringToBool(order.ReduceOnly),
			CreateTime:    stringToInt64(order.CTime),
			UpdateTime:    stringToInt64(order.UTime),
			RealizedPnl:   decimal.Zero.String(),

			IsAlgo:        true,
			OrderAlgoType: OrderAlgoType(order.OrdType),
		}

		switch order.OrdType {
		case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
			appendOrder.TriggerPrice = triggerPx.String()
			appendOrder.TriggerType = triggerType
			appendOrder.TriggerConditionType = triggerConditionType
		case OKX_ORDER_ALGO_TYPE_OCO:
			appendOrder.OcoTpTriggerPrice = order.TpTriggerPx
			appendOrder.OcoTpOrdPrice = order.TpOrdPx
			appendOrder.OcoSlTriggerPrice = order.SlTriggerPx
			appendOrder.OcoSlOrdPrice = order.SlOrdPx
			if order.TpOrdPx == "-1" {
				appendOrder.OcoTpOrdType = ORDER_TYPE_MARKET
			} else {
				appendOrder.OcoTpOrdType = ORDER_TYPE_LIMIT
			}
			if order.SlOrdPx == "-1" {
				appendOrder.OcoSlOrdType = ORDER_TYPE_MARKET
			} else {
				appendOrder.OcoSlOrdType = ORDER_TYPE_LIMIT
			}
		}
		orders = append(orders, appendOrder)
	}

	return orders, nil
}

// 策略订单推送处理
func (o *OkxTradeEngine) handleOrderFromWsOrderAlgo(order myokxapi.WsOrdersAlgo) *Order {

	var orderType OrderType
	timeInForce := TIME_IN_FORCE_GTC
	var px decimal.Decimal
	var triggerPx decimal.Decimal
	var triggerType OrderTriggerType
	var triggerConditionType OrderTriggerConditionType

	if order.TpTriggerPx != "" {
		px, _ = decimal.NewFromString(order.TpOrdPx)
		triggerPx, _ = decimal.NewFromString(order.TpTriggerPx)
		triggerType = ORDER_TRIGGER_TYPE_TAKE_PROFIT

		if order.Side == OKX_ORDER_SIDE_BUY {
			//止盈买入 价格下穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		} else {
			//止盈卖出 价格上穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		}
	} else if order.SlTriggerPx != "" {
		px, _ = decimal.NewFromString(order.SlOrdPx)
		triggerPx, _ = decimal.NewFromString(order.SlTriggerPx)
		triggerType = ORDER_TRIGGER_TYPE_STOP_LOSS

		if order.Side == OKX_ORDER_SIDE_BUY {
			//止损买入 价格上穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		} else {
			//止损卖出 价格下穿触发
			triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		}
	}

	if px.Equal(decimal.NewFromInt(-1)) {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var IsMargin, IsIsolated bool
	if order.TdMode != "cash" {
		switch order.TdMode {
		case "cross":
			IsIsolated = false
		case "isolated":
			IsIsolated = true
		}
	}
	if order.OrdersAlgo.InstType == OKX_AC_MARGIN.String() {
		IsMargin = true
		order.OrdersAlgo.InstType = OKX_AC_SPOT.String()
	}

	retOrder := &Order{
		Exchange:      OKX_NAME.String(),
		AccountType:   order.OrdersAlgo.InstType,
		Symbol:        order.OrdersAlgo.InstId,
		IsMargin:      IsMargin,
		IsIsolated:    IsIsolated,
		OrderId:       order.AlgoId,
		ClientOrderId: order.AlgoClOrdId,
		Price:         px.String(),
		Quantity:      order.Sz,
		ExecutedQty:   decimal.Zero.String(),
		AvgPrice:      decimal.Zero.String(),
		Status:        o.okxConverter.FromOKXOrderStatus(order.State, true),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(order.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(order.ReduceOnly),
		CreateTime:    stringToInt64(order.CTime),
		UpdateTime:    stringToInt64(order.UTime),
		RealizedPnl:   decimal.Zero.String(),

		IsAlgo:        true,
		OrderAlgoType: OrderAlgoType(order.OrdersAlgo.OrdType),
	}
	switch order.OrdersAlgo.OrdType {
	case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
		retOrder.TriggerPrice = triggerPx.String()
		retOrder.TriggerType = triggerType
		retOrder.TriggerConditionType = triggerConditionType
	case OKX_ORDER_ALGO_TYPE_OCO:
		retOrder.OcoTpTriggerPrice = order.OrdersAlgo.TpTriggerPx
		retOrder.OcoTpOrdType = OrderType(order.OrdersAlgo.TpTriggerPxType)
		retOrder.OcoTpOrdPrice = order.OrdersAlgo.TpOrdPx
		retOrder.OcoSlTriggerPrice = order.OrdersAlgo.SlTriggerPx
		retOrder.OcoSlOrdType = OrderType(order.OrdersAlgo.SlTriggerPxType)
		retOrder.OcoSlOrdPrice = order.OrdersAlgo.SlOrdPx
	}

	return retOrder
}
func (o *OkxTradeEngine) handleOrdersFromQueryOpenOrderAlgo(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderAlgoPendingRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	orders := make([]*Order, 0, len(res.Data))
	for _, order := range res.Data {
		var orderType OrderType
		timeInForce := TIME_IN_FORCE_GTC
		var px decimal.Decimal
		var triggerPx decimal.Decimal
		var triggerType OrderTriggerType
		var triggerConditionType OrderTriggerConditionType

		if order.TpTriggerPx != "" {
			px, _ = decimal.NewFromString(order.TpOrdPx)
			triggerPx, _ = decimal.NewFromString(order.TpTriggerPx)
			triggerType = ORDER_TRIGGER_TYPE_TAKE_PROFIT

			if order.Side == OKX_ORDER_SIDE_BUY {
				//止盈买入 价格下穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
			} else {
				//止盈卖出 价格上穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
			}
		} else if order.SlTriggerPx != "" {
			px, _ = decimal.NewFromString(order.SlOrdPx)
			triggerPx, _ = decimal.NewFromString(order.SlTriggerPx)
			triggerType = ORDER_TRIGGER_TYPE_STOP_LOSS

			if order.Side == OKX_ORDER_SIDE_BUY {
				//止损买入 价格上穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
			} else {
				//止损卖出 价格下穿触发
				triggerConditionType = ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
			}
		}

		if px.Equal(decimal.NewFromInt(-1)) {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var IsMargin, IsIsolated bool
		if order.TdMode != "cash" {
			switch order.TdMode {
			case "cross":
				IsIsolated = false
			case "isolated":
				IsIsolated = true
			}
		}
		if order.InstType == OKX_AC_MARGIN.String() {
			IsMargin = true
			order.InstType = OKX_AC_SPOT.String()
		}

		appendOrder := &Order{
			Exchange:      OKX_NAME.String(),
			AccountType:   order.InstType,
			Symbol:        order.InstId,
			IsMargin:      IsMargin,
			IsIsolated:    IsIsolated,
			OrderId:       order.AlgoId,
			ClientOrderId: order.AlgoClOrdId,
			Price:         px.String(),
			Quantity:      order.Sz,
			ExecutedQty:   decimal.Zero.String(),
			AvgPrice:      decimal.Zero.String(),
			Status:        o.okxConverter.FromOKXOrderStatus(order.State, true),
			Type:          orderType,
			Side:          o.okxConverter.FromOKXOrderSide(order.Side),
			PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    stringToBool(order.ReduceOnly),
			CreateTime:    stringToInt64(order.CTime),
			UpdateTime:    stringToInt64(order.UTime),
			RealizedPnl:   decimal.Zero.String(),

			IsAlgo:        true,
			OrderAlgoType: OrderAlgoType(order.OrdType),
		}

		switch order.OrdType {
		case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
			appendOrder.TriggerPrice = triggerPx.String()
			appendOrder.TriggerType = triggerType
			appendOrder.TriggerConditionType = triggerConditionType
		case OKX_ORDER_ALGO_TYPE_OCO:
			appendOrder.OcoTpTriggerPrice = order.TpTriggerPx
			appendOrder.OcoTpOrdPrice = order.TpOrdPx
			appendOrder.OcoSlTriggerPrice = order.SlTriggerPx
			appendOrder.OcoSlOrdPrice = order.SlOrdPx
			if order.TpOrdPx == "-1" {
				appendOrder.OcoTpOrdType = ORDER_TYPE_MARKET
			} else {
				appendOrder.OcoTpOrdType = ORDER_TYPE_LIMIT
			}
			if order.SlOrdPx == "-1" {
				appendOrder.OcoSlOrdType = ORDER_TYPE_MARKET
			} else {
				appendOrder.OcoSlOrdType = ORDER_TYPE_LIMIT
			}
		}
		orders = append(orders, appendOrder)
	}
	return orders, nil
}

// ws单订单请求相关
func (o *OkxTradeEngine) handleWsOrderCreateFromOrderParam(req *OrderParam) myokxapi.WsOrderArgData {
	tdMode := o.okxConverter.getTdModeFromAccountType(OkxAccountType(req.AccountType),
		o.okxConverter.ToOKXAccountMode(req.AccountMode), req.IsIsolated, req.IsMargin)
	return myokxapi.WsOrderArgData{
		InstId:     req.Symbol,
		TdMode:     tdMode,
		Px:         req.Price.String(),
		Sz:         req.Quantity.String(),
		Side:       o.okxConverter.ToOKXOrderSide(req.OrderSide),
		PosSide:    o.okxConverter.ToOKXPositionSide(req.PositionSide),
		OrdType:    o.okxConverter.ToOKXOrderType(req.OrderType, req.TimeInForce),
		ClOrdId:    req.ClientOrderId,
		ReduceOnly: req.ReduceOnly,
	}
}
func (o *OkxTradeEngine) handleWsOrderAmendFromOrderParam(req *OrderParam) myokxapi.WsAmendOrderArgData {
	return myokxapi.WsAmendOrderArgData{
		InstId:    req.Symbol,
		CxlOnFail: false,
		OrdId:     req.OrderId,
		ClOrdId:   req.ClientOrderId,
		NewPx:     req.Price.String(),
		NewSz:     req.Quantity.String(),
	}
}
func (o *OkxTradeEngine) handleWsOrderCancelFromOrderParam(req *OrderParam) myokxapi.WsCancelOrderArgData {
	return myokxapi.WsCancelOrderArgData{
		InstId:  req.Symbol,
		OrdId:   req.OrderId,
		ClOrdId: req.ClientOrderId,
	}
}

// ws批量订单请求相关
func (o *OkxTradeEngine) handleBatchWsOrderCreateFromOrderParams(reqs []*OrderParam) []myokxapi.WsOrderArgData {
	args := make([]myokxapi.WsOrderArgData, 0, len(reqs))
	for _, req := range reqs {
		args = append(args, o.handleWsOrderCreateFromOrderParam(req))
	}
	return args
}
func (o *OkxTradeEngine) handleBatchWsOrderAmendFromOrderParams(reqs []*OrderParam) []myokxapi.WsAmendOrderArgData {
	args := make([]myokxapi.WsAmendOrderArgData, 0, len(reqs))
	for _, req := range reqs {
		args = append(args, o.handleWsOrderAmendFromOrderParam(req))
	}
	return args
}
func (o *OkxTradeEngine) handleBatchWsOrderCancelFromOrderParams(reqs []*OrderParam) []myokxapi.WsCancelOrderArgData {
	args := make([]myokxapi.WsCancelOrderArgData, 0, len(reqs))
	for _, req := range reqs {
		args = append(args, o.handleWsOrderCancelFromOrderParam(req))
	}
	return args
}

// ws单订单结果返回
func (o *OkxTradeEngine) handleOrderFromWsOrderResult(req *OrderParam, res *myokxapi.WsOrderResult) (*Order, error) {
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]%s: {[%s]:%s}", res.Code, res.Msg, res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   `1`,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

// ws批量订单结果返回
func (o *OkxTradeEngine) handleOrdersFromWsBatchOrderResult(reqs []*OrderParam, res *myokxapi.WsOrderResult) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClOrdId, r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	if errStr != "" {
		return orders, fmt.Errorf("[%s]%s: [%s]", res.Code, res.Msg, errStr)
	}
	return orders, nil
}
