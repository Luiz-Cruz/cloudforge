package main

import (
	"github.com/Luiz-Cruz/cloudforge/cmd"
	"github.com/Luiz-Cruz/cloudforge/platform/config"
)

func main() {
	config.InitConfiguration()
	cmd.ProvideRunner().Run()
}
