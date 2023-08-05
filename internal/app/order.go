package app

import (
	"fmt"
	"log"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"strings"
)

// may need to make sure that connection is not being written to by the update position part?
// haven't decided whether I want to always return errors here or not, probably will, but will just print twice now for time being
func (app *FxApp) OpenPosition(copyPosition *api.ApiMonitorMessage) (string, error) {
	//get oid
	orderData := &fix.OrderData{
		Symbol:    fmt.Sprint(copyPosition.SymbolID),
		Volume:    float64(copyPosition.Volume),
		Direction: strings.ToLower(copyPosition.Direction),
		OrderType: "market",
	}

	executionReport, fxErr := app.FxSession.CtraderNewOrderSingle(app.FxUser, *orderData)
	if fxErr != nil {
		switch fxErr.ErrorCause {
		case fix.ProgramError:
			log.Fatal(fxErr.ErrorMessage)
		case fix.ConnectionError:
			app.Progam.Send(FeedUpdate(fmt.Sprintf("connection error whilst attempting to open position: %s", fxErr.ErrorMessage)))
		case fix.MarketError:
			app.Progam.Send(FeedUpdate(fmt.Sprintf(fxErr.ErrorMessage)))
		case fix.UserDataError:
			log.Fatal(fxErr.ErrorMessage)
		}
	}
	log.Fatalf("%+v", executionReport)
	return "", nil
}

func (app *FxApp) ClosePosition(copyPosition *api.ApiMonitorMessage) (string, error) {
	return "", nil
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
