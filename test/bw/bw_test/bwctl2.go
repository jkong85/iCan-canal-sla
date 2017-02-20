//bandwidth control interface for EyeQ
//author: Yan Sun

package main

import (
	//"os"
	"encoding/json"
	//"io/ioutil"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"golang.org/x/net/context"
	"log"
	"net"
	"time"
)

type ContainerBW struct {
	NodeIP          string
	PodID           string
	VlanID          string
	VxlanID         string
	PodIP           string
	Action          string
	InBandWidthMin  string // unit is Mbps
	InBandWidthMax  string // unit is Mbps
	OutBandWidthMin string // unit is Mbps
	OutBandWidthMax string // unit is Mbps
	PodPriority     string // 0-7, 0 is the highest priority, 7 is the lowest priority.
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func SaveAsJson(bw []ContainerBW, path string) {
	b, err := json.Marshal(bw)
	check(err)
	ioutil.WriteFile(path, b, 0644)
}
*/
func main() {

	bw := []ContainerBW{
		{"192.0.0.1", "1", "100", "1", "all", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "default", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.2", "delete", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.3", "delete", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.4", "delete", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.5", "delete", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.6", "delete", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.7", "delete", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.8", "delete", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.9", "delete", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.10", "delete", "666", "100", "100", "100", "2"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.11", "change", "666", "777", "10", "777", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.12", "change", "666", "777", "700", "777", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.13", "change", "666", "777", "200", "777", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.14", "change", "666", "777", "200", "777", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.15", "change", "666", "777", "1000", "777", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.16", "change", "666", "777", "10", "777", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.17", "change", "666", "777", "700", "777", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.18", "change", "666", "777", "200", "777", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.19", "change", "666", "777", "200", "777", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.20", "change", "666", "777", "1000", "777", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.21", "change", "666", "777", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.22", "change", "666", "777", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.23", "change", "666", "777", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.24", "change", "666", "777", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.25", "change", "666", "777", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.26", "change", "666", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.27", "change", "666", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.28", "change", "666", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.29", "change", "666", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.30", "change", "666", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.31", "change", "666", "777", "10", "777", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.32", "change", "666", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.33", "change", "666", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.34", "change", "666", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.35", "change", "666", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.36", "change", "666", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.37", "change", "666", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.38", "change", "666", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.39", "change", "666", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.40", "change", "666", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.41", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.42", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.43", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.44", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.45", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.46", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.47", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.48", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.49", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.50", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.51", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.52", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.53", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.54", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.55", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.56", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.57", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.58", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.59", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.60", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.61", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.62", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.63", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.64", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.65", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.66", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.67", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.68", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.69", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.70", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.71", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.72", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.73", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.74", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.75", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.76", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.77", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.78", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.79", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.80", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.81", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.82", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.83", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.84", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.85", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.86", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.87", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.88", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.89", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.90", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.91", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.92", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.93", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.94", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.95", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.96", "", "10", "100", "10", "100", "5"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.97", "", "500", "500", "700", "700", "0"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.98", "", "200", "200", "200", "200", "5"},
		{"192.0.0.2", "2", "102", "2", "172.17.0.99", "", "200", "200", "200", "200", "7"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.100", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "172.17.0.101", "", "10", "100", "10", "100", "5"},
	}

	// We can register as many backends that are supported by libkv
	etcd.Register()

	server := "127.0.0.1:4001"

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{server},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store", kv)
	}

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:4001"},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	qos_encode, err := json.Marshal(bw)
	if err != nil {
		log.Fatal(err)
	}

	intf, err := net.InterfaceByName("wlp3s0")

	if err != nil {
		log.Fatal("Cannot find interface by name eth0")
	}

	mac := intf.HardwareAddr

	key := "/" + string(mac)

	resp, err := kapi.Set(context.Background(), key, string(qos_encode), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		//log.Printf("Set is done. Metadata is %q\n", resp)
		fmt.Println("Set is done. Metadata is %q\n", resp)
	}

	/*
		path := "./qos.json"
		f, err := os.Create(path)
		check(err)

		defer f.Close()
		SaveAsJson(bw, path)
		f.Sync()
	*/
}
