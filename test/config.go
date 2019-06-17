package main

import (
	"encoding/json"
	"fmt"
	"github.com/zdnscloud/zke/types"
)

var zkeConfig = types.ZKEConfig{}

func main() {
	config, _ := json.Marshal(zkeConfig)
	fmt.Println(string(config))
}
