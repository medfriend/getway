package util

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/medfriend/shared-commons-go/util/global"
)

func ConnectRabbit(consulClient *api.Client) {
	rabbitInfo, _ := consul.GetKeyValue(consulClient, "RABBIT")

	var resultRabbitmq map[string]string

	err := json.Unmarshal([]byte(rabbitInfo), &resultRabbitmq)

	if err != nil {
		fmt.Errorf("error while unmarshalling RABBIT")
	}

	s := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		resultRabbitmq["RABBIT_USER"],
		resultRabbitmq["RABBIT_PASSWORD"],
		resultRabbitmq["RABBIT_HOST"],
		resultRabbitmq["RABBIT_PORT"])

	global.SetRabbitConn(s)
}
