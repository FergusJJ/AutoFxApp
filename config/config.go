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

func LoadSettingsFromJson() (string, error) {
	settingFile, err := os.Open("user/settings.json")
	if err != nil {
		return "", err
	}

	defer settingFile.Close()

	jsonBytes, err := io.ReadAll(settingFile)
	if err != nil {
		return "", err
	}

	type Settings struct {
		LicenseKey string `json:"licenseKey"`
		// Pools      string `json:"pools"`
	}
	var settings Settings
	err = json.Unmarshal(jsonBytes, &settings)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(settings.LicenseKey) /*strings.TrimSpace(settings.Pools),*/, nil
}
