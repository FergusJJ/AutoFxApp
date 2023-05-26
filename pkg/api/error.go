package api

import "fmt"

func errorValidatingLicense(status int, message string) error {
	return fmt.Errorf("%d - %s", status, message)
}
