package database

import (
	"log"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

var dataToken = cache.New(time.Hour, 24*time.Hour)

func init() {
	if readConfig("AK", "ALIBABA_CLOUD_ACCESS_KEY_ID") != nil {
		log.Panicln("Get Access key fail, crash now.")
	}
	if readConfig("SK", "ALIBABA_CLOUD_SECRET_KEY_ID") != nil {
		log.Panicln("Get Secret key fail, crash now.")
	}
	if readConfig("STS", "ALIBABA_CLOUD_SECURITY_TOKEN") != nil {
		log.Panicln("Get SecurityToken fail, crash now.")
	}
	if readConfig("OTS_NAME", "TARGET_DB") != nil {
		log.Panicln("Get database name fail, crash now.")
	}
	if readConfig("GH_CLIENT_ID", "GH_BASIC_CLIENT_ID") != nil {
		log.Panicln("read github oauth client id failed, crash now")
	}
	if readConfig("GH_CLIENT_SECRET", "GH_BASIC_SECRET_SECRET") != nil {
		log.Panicln("read github OAuth client secret failed, crash now")
	}
}

func readConfig(cfgIndex, env string) error {
	return dataToken.Add(cfgIndex, os.Getenv(env), cache.DefaultExpiration)
}
