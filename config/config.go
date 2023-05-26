package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"pollo/pkg/fix"
	"strings"

	"github.com/joho/godotenv"
)

func Config(key string) (envVar string, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	envVar = os.Getenv(key)
	if envVar == "" {
		err = fmt.Errorf("%s does not exist", key)
	}
	return envVar, err
}

func LoadDataFromJson() (readUserData *fix.FxUser, err error) {
	jsonFile, err := os.Open("user/test.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var userData fix.FxUser

	err = json.Unmarshal(jsonBytes, &userData)
	if err != nil {
		return nil, err
	}

	readUserData = &userData
	return readUserData, nil
}

func LoadLicenseKeyFromJson() (string, error) {
	jsonFile, err := os.Open("user/settings.json")
	if err != nil {
		return "", err
	}

	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	type License struct {
		LicenseKey string `json:"licenseKey"`
	}
	var license License
	err = json.Unmarshal(jsonBytes, &license)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(license.LicenseKey), nil
}
