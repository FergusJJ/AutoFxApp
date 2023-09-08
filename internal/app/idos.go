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
	Program            AppProgram
	UiPositionsDataMap map[string]uiPositionData
}

type AppProgram struct {
	Program *tea.Program
}

type uiPositionData struct {
	symbolName     string
	entryPrice     string
	currentPrice   string
	copyPositionId string
	positionId     string
	volume         string
	grossProfit    string
	timestamp      string
	side           string
	symbol         string
	isProfit       bool
}

type orderInProgressData struct {
	orderID     string
	status      string
	avgPrice    string
	unfilledQty int
}
