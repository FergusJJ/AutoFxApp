package app

import (
	"pollo/pkg/api"
	"pollo/pkg/fix"

	"github.com/apex/log"
)

type FxApp struct {
	FxSession      fix.FxSession
	FxSecurityList map[string]string
	FxUser         fix.FxUser
	ApiSession     api.ApiSession //done
	AppLogger      *log.Logger    //done
	LicenseKey     string         `json:"licenseKey"` //
}
