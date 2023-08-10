package app

import (
	"pollo/pkg/api"
	"pollo/pkg/fix"

	tea "github.com/charmbracelet/bubbletea"
)

type FxApp struct {
	FxSession          fix.FxSession
	FxSecurityList     map[string]string
	FxUser             fix.FxUser
	ApiSession         api.ApiSession //done
	LicenseKey         string         `json:"licenseKey"` //
	Progam             *tea.Program
	UiPositionsDataMap map[string]uiPositionData
}

type uiPositionData struct {
	positionId  string
	volume      string
	grossProfit string
	netProfit   string
}

type orderInProgressData struct {
	orderID     string
	status      string
	avgPrice    string
	unfilledQty int
}
