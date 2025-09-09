package core_util

import (
	"Badminton-Hub/util"
	"fmt"
)

func HashAuth(rawHash, key string) string {
	data := fmt.Sprint(rawHash + key)
	hashAuth := util.Sha256(data)
	return hashAuth
}

func RandomGoogleState() (string, error) {
	config := util.LoadConfig()

	if config.Mode == "DEVERLOP" {
		return config.DefaultGoogleState, nil
	}

	state := util.RandomString(32, true, true, false)
	return string(state), nil
}
