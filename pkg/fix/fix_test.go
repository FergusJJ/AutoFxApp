package fix_test

import (
	"pollo/pkg/fix"
	"testing"
)

func TestParseFIXResponse(t *testing.T) {
	var executionReportMessage1 = "8=FIX.4.4|9=206|35=8|34=78|49=CSERVER|50=TRADE|52=20170117-10:02:15.045|56=live.theBroker.12345|57=any_string|6=1.0674|11=876316397|14=10000|32=10000|37=101|38=10000|39=2|40=1|54=1|55=1|59=3|60=20170117-10:02:14.963|150=F|151=0|721=101|10=077|"
	var executionReportMessage2 = "8=FIX.4.4|9=197|35=8|34=77|49=CSERVER|50=TRADE|52=20170117-10:02:14.720|56=live.theBroker.12345|57=any_string|11=876316397|14=0|37=101|38=10000|39=0|40=1|54=1|55=1|59=3|60=20170117-10:02:14.591|150=0|151=10000|721=101|10=149|"
	_, err := fix.ParseFIXResponse([]byte(executionReportMessage1), fix.NewOrderSingle)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fix.ParseFIXResponse([]byte(executionReportMessage2), fix.NewOrderSingle)
	if err != nil {
		t.Fatal(err)
	}
}
