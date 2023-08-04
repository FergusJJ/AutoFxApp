package app

import (
	"pollo/pkg/api"
	"pollo/pkg/fix"

	tea "github.com/charmbracelet/bubbletea"
)

type FxApp struct {
	FxSession      fix.FxSession
	FxSecurityList map[string]string
	FxUser         fix.FxUser
	ApiSession     api.ApiSession //done
	LicenseKey     string         `json:"licenseKey"` //
	Progam         *tea.Program
}
