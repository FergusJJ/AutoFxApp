package app

import (
	"fmt"
	"log"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"strconv"
	"strings"
	"time"
)

// might want to change from OrderType market, to Limit, user may specify limit
func (app *FxApp) OpenPosition(copyPosition *api.ApiMonitorMessage) (bool, *fix.Position) {
	//get oid
	orderData := fix.OrderData{
		Symbol:    fmt.Sprint(copyPosition.SymbolID),
		Volume:    float64(copyPosition.Volume),
		Direction: strings.ToLower(copyPosition.Direction),
		OrderType: "market",
	}

	executionReport, fxErr := app.FxSession.CtraderNewOrderSingle(app.FxUser, orderData)
	if fxErr != nil {
		switch fxErr.ErrorCause {
		case fix.ProgramError:
			log.Fatalf("program error: %s", fxErr.ErrorMessage)
		case fix.ConnectionError:
			app.Program.SendColor(fmt.Sprintf("client error whilst attempting to open position: %s", fxErr.ErrorMessage), "red")
		case fix.MarketError:
			app.Program.SendColor(fmt.Sprintf(fxErr.ErrorMessage), "red")
		}
		return false, nil
	}
	var pollTimeout = time.Millisecond * 1500
	// The expected status
	for {
		switch executionReport.ExecType {
		case "0":
			app.Program.SendColor("order in progress...", "yellow")
			time.Sleep(pollTimeout)
		case "4":
			app.Program.SendColor("order was cancelled", "red")
			return false, nil
		case "8":
			app.Program.SendColor(fmt.Sprintf("order rejected: %s", executionReport.Text), "red")
			return false, nil
		case "C":
			app.Program.SendColor("order expired", "red")
			return false, nil
		case "F":
			app.Program.SendColor("order was executed", "green")
			if executionReport.OrdStatus == "1" {
				app.Program.SendColor("order could only be partially filled", "yellow")
			}
			avgPx, err := strconv.ParseFloat(executionReport.AvgPx, 64)
			if err != nil {
				log.Fatalf("Error parsing float from :%s", executionReport.AvgPx)
			}
			volumeInt, err := strconv.ParseInt(executionReport.CumQty, 10, 64)
			if err != nil {
				log.Fatalf("Error parsing int from :%s", executionReport.CumQty)
			}
			positionData := &fix.Position{
				PID:       executionReport.PosMaintRptID,
				CopyPID:   copyPosition.CopyPID,
				Side:      copyPosition.Direction,
				Symbol:    fmt.Sprint(copyPosition.SymbolID),
				AvgPx:     avgPx,
				Volume:    volumeInt,
				Timestamp: executionReport.TransactTime,
			}

			return true, positionData

		//this shouldn't happen yet
		//if/when support added, will just need to update the volume of the position i think
		case "5":
			app.Program.SendColor("order was replaced", "yellow")
			return false, nil
			//update positions
		case "I":
			log.Fatalf("in OpenPosition ExecType: %+v", executionReport)
		}
		clOrdID := executionReport.ClOrdID
		executionReport, fxErr = app.FxSession.CtraderOrderStatus(app.FxUser, clOrdID)
		if fxErr != nil {
			switch fxErr.ErrorCause {
			case fix.ProgramError:
				log.Fatalf("program error: %s", fxErr.ErrorMessage)
			case fix.ConnectionError:
				app.Program.SendColor(fmt.Sprintf("client error whilst polling order status: %s", fxErr.ErrorMessage), "red")
			case fix.MarketError:
				app.Program.SendColor(fmt.Sprintf(fxErr.ErrorMessage), "red")
			}
			return false, nil
		}
	}
}

/*
	[11:fa047931-4d8c-4bfa-ba69-7e47db59897f
	14:0
	37:99501203
	38:2000
	39:0
	40:1
	54:2
	55:1
	59:3
	60:20230808-13:26:33.827
	150:0
	151:2000
	721:55019091]}
*/

func (app *FxApp) ClosePosition(copyPosition *api.ApiMonitorMessage) (string, error) {
	return "", nil
}

func (app *FxApp) FetchPositions() ([]string, *fix.ErrorWithCause) {
	positionReports, fxErr := app.FxSession.CtraderRequestForPositions(app.FxUser)
	if fxErr != nil {
		//log error to feed, don't update positions
		return []string{}, fxErr
	}
	//need to get a slice containing the position ids, with their volume, entry price, and symbol
	positionReportString := ""
	for _, v := range positionReports {
		vol := v.LongQty
		if v.LongQty == "0" {
			vol = v.ShortQty
		}
		positionReportString += fmt.Sprintf("Position id: %s | Volume: %s | Symbol: %s \n ", v.PosMaintRptID, vol, v.Symbol)
	}
	if positionReportString == "" {
		app.Program.SendColor("No updates", "yellow")
	}
	app.Program.SendColor(positionReportString, "yellow")
	// fxErr = app.FxSession.CtraderMarketDataRequest(app.FxUser)
	// if fxErr != nil {
	// 	return []string{}, fxErr
	// }
	return []string{}, nil
}

func (app *FxApp) GetPositionInformation() {
	//send all marketDataSubscriptions
}

/*

Errors:
2023/06/30 16:44:52 Got OPEN:&{CopyPID:340483758 SymbolID:3 Price:152.861 Volume:1200000 Direction:SELL MessageType:OPEN}
2023/06/30 16:44:52 &{3 1200 sell market} //EURJPY
volume: 1200
[8=FIX.4.4 9=218 35=j 34=2 49=CServer 50=TRADE 52=20230630-15:44:51.538 56=demo.ctrader.3697899 57=TRADE 58=TRADING_BAD_VOLUME:Order volume 1200.00 must be multiple of stepVolume=1000.00. 379=6be7594a-aa5c-47d1-ab74-4e4cba2beda8 380=0 10=171]
2023/06/30 16:44:52 {RefSeqNum: RefMsgType: BusinessRejectRefID:6be7594a-aa5c-47d1-ab74-4e4cba2beda8 BusinessRejectReason:0 Text:TRADING_BAD_VOLUME:Order volume 1200.00 must be multiple of stepVolume}
2023/06/30 16:44:52 Business Message Reject

Working:
[8=FIX.4.4 9=99 35=A 34=1 49=CServer 50=TRADE 52=20230630-15:48:57.117 56=demo.ctrader.3697899 57=TRADE 98=0 108=0 10=044]
2023/06/30 16:48:57 logged in
2023/06/30 16:49:07 got message of type: 1
2023/06/30 16:49:07 Got OPEN:&{CopyPID:340483758 SymbolID:3 Price:152.861 Volume:1200000 Direction:SELL MessageType:OPEN}
2023/06/30 16:49:07 &{3 12000 sell market}
volume: 12000
[8=FIX.4.4 9=233 35=8 34=2 49=CServer 50=TRADE 52=20230630-15:49:07.165 56=demo.ctrader.3697899 57=TRADE 11=3f357761-3c81-4a50-b777-b3d38fd6bd92 14=0 37=94878953 38=12000 39=0 40=1 54=2 55=3 59=3 60=20230630-15:49:07.160 150=0 151=12000 721=52942663 10=075]
2023/06/30 16:49:07 &{OrderID:94878953 ClOrdID:3f357761-3c81-4a50-b777-b3d38fd6bd92 TotNumReports: ExecType:0 OrdStatus:0 Symbol:3 Side:2 TransactTime:20230630-15:49:07.160 AvgPx: OrderQty:12000 LeavesQty:12000 CumQty:0 LastQty: OrdType:1 Price: StopPx: TimeInForce:3 ExpireTime: Text: OrdRejReason: PosMaintRptID:52942663 Designation: MassStatusReqID: AbsoluteTP: RelativeTP: AbsoluteSL: RelativeSL: TrailingSL: TriggerMethodSL: GuaranteedSL:}


*/
