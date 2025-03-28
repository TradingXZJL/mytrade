package mytrade

import (
	"strconv"
	"strings"

	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
)

type GateExchangeInfo struct {
	ExchangeBase
	isLoaded                   bool
	spotSymbolMap              *MySyncMap[string, *mygateapi.PublicRestSpotCurrencyPairCommon]
	futuresSymbolMap           *MySyncMap[string, *mygateapi.ContractCommon]
	deliverySymbolMap          *MySyncMap[string, *mygateapi.DeliveryContractCommon]
	spotIsolatedMarginLeverage *MySyncMap[string, decimal.Decimal]

	spotExchangeInfoMap     *MySyncMap[string, TradeSymbolInfo]
	futuresExchangeInfoMap  *MySyncMap[string, TradeSymbolInfo]
	deliveryExchangeInfoMap *MySyncMap[string, TradeSymbolInfo]
}

type gateSymbolInfo struct {
	symbolInfo
}

func (e *GateExchangeInfo) loadExchangeInfo() error {
	e.spotSymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.PublicRestSpotCurrencyPairCommon]())
	e.futuresSymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.ContractCommon]())
	e.deliverySymbolMap = GetPointer(NewMySyncMap[string, *mygateapi.DeliveryContractCommon]())
	e.spotIsolatedMarginLeverage = GetPointer(NewMySyncMap[string, decimal.Decimal]())

	var err error
	spotRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestCurrencyPairsAll().Do()
	if err != nil {
		return err
	}

	spotMarginRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestMarginUniCurrencyPairs().Do()
	if err != nil {
		return err
	}

	for _, v := range spotMarginRes.Data {
		e.spotIsolatedMarginLeverage.Store(v.CurrencyPair, decimal.RequireFromString(v.Leverage))
	}

	futuresRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestFuturesSettleContracts().Settle("usdt").Do()
	if err != nil {
		return err
	}

	deliveryRes, err := mygateapi.NewRestClient("", "").PublicRestClient().
		NewPublicRestDeliverySettleContracts().Settle("usdt").Do()
	if err != nil {
		return err
	}

	for _, v := range spotRes.Data {
		newSymbol := v
		e.spotSymbolMap.Store(v.ID, &newSymbol)
	}
	for _, v := range deliveryRes.Data {
		newSymbol := v
		e.deliverySymbolMap.Store(v.Name, &newSymbol)
	}
	for _, v := range futuresRes.Data {
		newSymbol := v
		e.futuresSymbolMap.Store(v.Name, &newSymbol)
	}
	e.isLoaded = true
	e.spotExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.deliveryExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	e.futuresExchangeInfoMap = GetPointer(NewMySyncMap[string, TradeSymbolInfo]())
	return nil
}

func (e *GateExchangeInfo) Refresh() error {
	e.isLoaded = false
	return e.loadExchangeInfo()
}

func (e *GateExchangeInfo) GetSymbolInfo(accountType string, symbol string) (TradeSymbolInfo, error) {
	if !e.isLoaded {
		err := e.loadExchangeInfo()
		if err != nil {
			return nil, err
		}
	}

	var pricePrecision, amtPrecision, maxOrderNum int
	var tickSize, minPrice, maxPrice = "0", "0", "0"
	var lotSize, minAmt, maxLmtAmt, maxMktAmt = "0", "0", "0", "0"
	var minNotional = "0"
	var baseCoin, quoteCoin string
	var isContract, isContractAmt bool
	var contractSize, contractCoin, contractType = "0", "", ""
	var minLeverage, maxLeverage, stepLeverage = "0", "0", "0"

	var isTrading bool

	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		v, ok := e.spotSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = v.Base, v.Quote
		isContract, isContractAmt = false, false
		contractSize, contractCoin = "0", ""
		if isolatedMarginLeverage, ok := e.spotIsolatedMarginLeverage.Load(symbol); ok {
			minLeverage, maxLeverage = isolatedMarginLeverage.String(), isolatedMarginLeverage.String()
		}

		pricePrecision = v.Precision
		tickSize = getSizeFromPrecision(v.Precision)
		amtPrecision = v.AmountPrecision
		lotSize = getSizeFromPrecision(v.AmountPrecision)
		if v.MinBaseAmount != "null" && v.MinBaseAmount != "" {
			minAmt = v.MinBaseAmount
		}
		if v.MinQuoteAmount != "null" && v.MinQuoteAmount != "" {
			minNotional = v.MinQuoteAmount
		}

		if v.MaxBaseAmount != "null" && v.MaxBaseAmount != "" {
			maxMktAmt = v.MaxBaseAmount
			maxLmtAmt = v.MaxBaseAmount
		}
		if v.TradeStatus == "tradable" {
			isTrading = true
		}

	case GATE_ACCOUNT_TYPE_FUTURES:
		v, ok := e.futuresSymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		sp := strings.Split(v.Name, "_")
		if len(sp) != 2 {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = sp[0], sp[1]
		isContract, isContractAmt = true, true
		contractSize = v.QuantoMultiplier
		if v.Type == "direct" {
			//正向合约
			contractCoin = baseCoin
		} else {
			//反向合约
			contractCoin = quoteCoin
		}
		minLeverage, maxLeverage, stepLeverage = v.LeverageMin, v.LeverageMax, "1"
		contractType = v.Type

		markPrice, _ := decimal.NewFromString(v.MarkPrice)
		priceDeviate, _ := decimal.NewFromString(v.OrderPriceDeviate)

		minPrice = markPrice.Sub(markPrice.Mul(priceDeviate)).String()
		maxPrice = markPrice.Add(markPrice.Mul(priceDeviate)).String()

		pricePrecision = countDecimalPlaces(v.OrderPriceRound)
		tickSize = v.OrderPriceRound
		amtPrecision = countDecimalPlaces(strconv.FormatInt(v.OrderSizeMin, 10))
		lotSize = strconv.FormatInt(v.OrderSizeMin, 10)

		minAmt = strconv.FormatInt(v.OrderSizeMin, 10)

		maxMktAmt = strconv.FormatInt(v.OrderSizeMax, 10)
		maxLmtAmt = strconv.FormatInt(v.OrderSizeMax, 10)

		isTrading = true
	case GATE_ACCOUNT_TYPE_DELIVERY:
		v, ok := e.deliverySymbolMap.Load(symbol)
		if !ok {
			return nil, ErrorSymbolNotFound
		}
		sp := strings.Split(v.Underlying, "_")
		if len(sp) != 2 {
			return nil, ErrorSymbolNotFound
		}
		baseCoin, quoteCoin = sp[0], sp[1]
		isContract, isContractAmt = true, true
		contractSize = v.QuantoMultiplier
		if v.Type == "direct" {
			//正向合约
			contractCoin = baseCoin
		} else {
			//反向合约
			contractCoin = quoteCoin
		}
		minLeverage, maxLeverage, stepLeverage = v.LeverageMin, v.LeverageMax, "1"
		contractType = v.Type

		markPrice, _ := decimal.NewFromString(v.MarkPrice)
		priceDeviate, _ := decimal.NewFromString(v.OrderPriceDeviate)

		minPrice = markPrice.Sub(markPrice.Mul(priceDeviate)).String()
		maxPrice = markPrice.Add(markPrice.Mul(priceDeviate)).String()

		pricePrecision = countDecimalPlaces(v.OrderPriceRound)
		tickSize = v.OrderPriceRound
		amtPrecision = countDecimalPlaces(strconv.FormatInt(v.OrderSizeMin, 10))
		lotSize = strconv.FormatInt(v.OrderSizeMin, 10)

		minAmt = strconv.FormatInt(v.OrderSizeMin, 10)

		maxMktAmt = strconv.FormatInt(v.OrderSizeMax, 10)
		maxLmtAmt = strconv.FormatInt(v.OrderSizeMax, 10)

		isTrading = true
	default:
		return nil, ErrorAccountType
	}

	return &gateSymbolInfo{symbolInfo: symbolInfo{
		symbolInfoStruct: symbolInfoStruct{
			Exchange:      GATE_NAME.String(),
			AccountType:   accountType,
			Symbol:        symbol,
			BaseCoin:      baseCoin,
			QuoteCoin:     quoteCoin,
			IsTrading:     isTrading,
			IsContract:    isContract,
			IsContractAmt: isContractAmt,
			ContractSize:  contractSize,
			ContractCoin:  contractCoin,
			ContractType:  contractType,

			PricePrecision: pricePrecision,
			AmtPrecision:   amtPrecision,

			TickSize: tickSize,
			MinPrice: minPrice,
			MaxPrice: maxPrice,

			LotSize:   lotSize,
			MinAmt:    minAmt,
			MaxLmtAmt: maxLmtAmt,
			MaxMktAmt: maxMktAmt,

			MaxLeverage:  maxLeverage,
			MinLeverage:  minLeverage,
			StepLeverage: stepLeverage,
			MaxOrderNum:  maxOrderNum,
			MinNotional:  minNotional,
		},
	}}, nil
}

func (e *GateExchangeInfo) GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error) {
	if !e.isLoaded {
		err := e.Refresh()
		if err != nil {
			return nil, err
		}
	}
	var symbolInfoList []TradeSymbolInfo
	switch GateAccountType(accountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		e.spotSymbolMap.Range(func(key string, value *mygateapi.PublicRestSpotCurrencyPairCommon) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case GATE_ACCOUNT_TYPE_FUTURES:
		e.futuresSymbolMap.Range(func(key string, value *mygateapi.ContractCommon) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	case GATE_ACCOUNT_TYPE_DELIVERY:
		e.deliverySymbolMap.Range(func(key string, value *mygateapi.DeliveryContractCommon) bool {
			symbolInfo, err := e.GetSymbolInfo(accountType, key)
			if err != nil {
				return false
			}
			symbolInfoList = append(symbolInfoList, symbolInfo)
			return true
		})
	default:
		return nil, ErrorAccountType
	}
	return symbolInfoList, nil
}
