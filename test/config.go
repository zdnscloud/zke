package main

import (
	"encoding/json"
	"fmt"
	types "github.com/zdnscloud/zke/typesnew"
	"io/ioutil"
	"os"
)

var zkeConfig = types.ZKEConfig{}

var masterNode types.ZKEConfigNode
var workerNode types.ZKEConfigNode

func main() {
	/*
		masterNode.NodeName = "master"
		masterNode.Address = "192.168.1.10"
		masterNode.InternalAddress = ""
		masterNode.Roles = []string{"controlPannel", "etcd"}
		workerNode.NodeName = "worker"
		workerNode.Address = "192.168.1.11"
		workerNode.InternalAddress = ""
		workerNode.Roles = []string{"worker"}

		zkeConfig.Nodes = append(zkeConfig.Nodes, masterNode)
		zkeConfig.Nodes = append(zkeConfig.Nodes, workerNode)
		config, _ := json.Marshal(zkeConfig)
		fmt.Println(string(config))
	*/

	f, _ := os.Open("config1.json")
	configContent, err := ioutil.ReadAll(f)
	err = json.Unmarshal(configContent, &zkeConfig)
	if err == nil {
		newConfig, _ := json.Marshal(zkeConfig)
		fmt.Println(string(newConfig))
	}
	fmt.Println(err)
}
