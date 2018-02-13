package api

import (
	"testing"
	"e.coding.net/anyun-cloud-api-gateway/server"
)

// TestAPIController -- 测试API网关
func TestAPIController(t *testing.T){
	sendHelp(&server.APICONTROLLERPARAMS{
		ID : "123456", 
		Name : "test",
		Version : "1.0.0",
		Dc : "test_dc",
	})
}