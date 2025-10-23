package main

import (
	"github.com/xiaoyi510/xbot"

	_ "mian/plugins/xarr-merchant"
)

func main() {
	cfg, err := xbot.LoadConfigFile("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	err = xbot.RunAndListen(cfg)
	if err != nil {
		panic(err)
	}
}
