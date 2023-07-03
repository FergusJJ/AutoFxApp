package app

import (
	"pollo/internal/app/ui"
	"pollo/pkg/api"
	"pollo/pkg/fix"
)

type FxApp struct {
	FxSession      fix.FxSession
	FxSecurityList map[string]string
	FxUser         fix.FxUser
	ApiSession     api.ApiSession //done
	LicenseKey     string         `json:"licenseKey"` //
	UI             *ui.AppUi
}
