package main

import (
	"log"
	"os"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var dataToken = cache.New(time.Hour, 24*time.Hour)

func init() {
	if readConfigFromEnv("AK", "ALIBABA_CLOUD_ACCESS_KEY_ID") != nil {
		log.Panicln("Get Access key fail, crash now.")
	}
	if readConfigFromEnv("SK", "ALIBABA_CLOUD_SECRET_KEY_ID") != nil {
		log.Panicln("Get Secret key fail, crash now.")
	}
	if readConfigFromEnv("STS", "ALIBABA_CLOUD_SECURITY_TOKEN") != nil && gin.Mode() == gin.ReleaseMode {
		log.Panicln("Get SecurityToken fail, crash now.")
	}
	if readConfigFromEnv("OTS_NAME", "TARGET_DB") != nil {
		log.Panicln("Get database name fail, crash now.")
	}
	if readConfigFromEnv("GH_CLIENT_ID", "GH_BASIC_CLIENT_ID") != nil {
		log.Panicln("read github oauth client id failed, crash now")
	}
	if readConfigFromEnv("GH_CLIENT_SECRET", "GH_BASIC_SECRET_SECRET") != nil {
		log.Panicln("read github OAuth client secret failed, crash now")
	}
}

func readConfigFromEnv(cfgIndex, env string) error {
	return dataToken.Add(cfgIndex, os.Getenv(env), cache.DefaultExpiration)
}
func newDbConn() *tablestore.TableStoreClient {
	client := func() *tablestore.TableStoreClient {
		otsName, cached := dataToken.Get("OTS_NAME")
		crashWithNoCache(cached)
		accessKey, cached := dataToken.Get("AK")
		crashWithNoCache(cached)
		secretKey, cached := dataToken.Get("SK")
		crashWithNoCache(cached)
		securityKey, cached := dataToken.Get("STS")
		crashWithNoCache(cached)
		setEndpoint := func() string {
			switch gin.Mode() {
			case gin.DebugMode:
				return "https://" + otsName.(string) + os.Getenv("OTS_REGION_DEBUG") + ".ots.aliyuncs.com"
			case gin.TestMode:
				return "https://" + otsName.(string) + os.Getenv("OTS_REGION_TEST") + ".ots.aliyuncs.com"
			case gin.ReleaseMode:
				return "https://" + otsName.(string) + os.Getenv("FC_REGION") + ".ots-internal.aliyuncs.com"
			default:
				return "https://" + otsName.(string) + os.Getenv("OTS_REGION_DEBUG") + ".ots.aliyuncs.com"
			}
		}
		return tablestore.NewClientWithConfig(otsName.(string)+setEndpoint(), otsName.(string), accessKey.(string), secretKey.(string), securityKey.(string), tablestore.NewDefaultTableStoreConfig())
	}
	return client()
}

func crashWithNoCache(cached bool) {
	if !cached {
		log.Panicln("config no cached")
	}
}
