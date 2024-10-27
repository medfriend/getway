package consulRegister

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"log"
	"os"
)

func handleJson(data map[string]string) string {

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Error converting map to JSON: %v", err)
	}

	return string(jsonData)
}

func RegisterConstants(consulClient *api.Client) {
	getwayData := map[string]string{
		"SERVICE_ID":   os.Getenv("SERVICE_ID"),
		"SERVICE_NAME": os.Getenv("SERVICE_NAME"),
		"SERVICE_PORT": os.Getenv("SERVICE_PORT"),
		"SERVICE_PATH": os.Getenv("SERVICE_ADDRESS"),
	}

	dbData := map[string]string{
		"DB_HOST":           os.Getenv("DB_HOST"),
		"DB_PORT":           os.Getenv("DB_PORT"),
		"DB_USER":           os.Getenv("DB_USER"),
		"DB_PASSWORD":       os.Getenv("DB_PASSWORD"),
		"DB_NAME":           os.Getenv("DB_NAME"),
		"MAX_OPEN_CONN":     os.Getenv("MAX_OPEN_CONN"),
		"MAX_IDLE_CONN":     os.Getenv("MAX_IDLE_CONN"),
		"MAX_LIFETIME_CONN": os.Getenv("MAX_LIFETIME_CONN"),
	}

	cacheData := map[string]string{
		"SERVICE_ID":     os.Getenv("CACHE_SERVICE_ID"),
		"SERVICE_NAME":   os.Getenv("CACHE_SERVICE_NAME"),
		"SERVICE_PORT":   os.Getenv("CACHE_SERVICE_PORT"),
		"SERVICE_PATH":   os.Getenv("CACHE_SERVICE_ADDRESS"),
		"REDIS_ADDRESS":  os.Getenv("REDIS_ADDRESS"),
		"REDIS_PORT":     os.Getenv("REDIS_PORT"),
		"REDIS_PASSWORD": os.Getenv("REDIS_PASSWORD"),
	}

	securityData := map[string]string{
		"SERVICE_ID":   os.Getenv("SECURITY_SERVICE_ID"),
		"SERVICE_NAME": os.Getenv("SECURITY_SERVICE_NAME"),
		"SERVICE_PORT": os.Getenv("SECURITY_SERVICE_PORT"),
		"SERVICE_PATH": os.Getenv("SECURITY_SERVICE_ADDRESS"),
	}

	consul.StoreKeyValue(consulClient, "SERVICE_CACHE", os.Getenv("SERVICE_CACHE"))
	consul.StoreKeyValue(consulClient, "JWT", os.Getenv("JWT_KEY"))
	consul.StoreKeyValue(consulClient, "GETWAY", handleJson(getwayData))
	consul.StoreKeyValue(consulClient, "DB", handleJson(dbData))
	consul.StoreKeyValue(consulClient, "CACHE", handleJson(cacheData))
	consul.StoreKeyValue(consulClient, "SECURITY", handleJson(securityData))
}
