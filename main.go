// cloudforge
//
// CloudForge API - Enterprise Cloud Provisioning Platform
//
//	Schemes: http
//	Host: localhost:8080
//	Version: 0.1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"github.com/Luiz-Cruz/cloudforge/cmd"
	"github.com/Luiz-Cruz/cloudforge/platform/config"
)

//go:generate swagger generate spec -m -o ./docs/specs/swagger.yaml
func main() {
	config.InitConfiguration()
	cmd.ProvideRunner().Run()
}
