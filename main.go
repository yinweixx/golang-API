package main

import (
	"e.coding.net/anyun-cloud-api-gateway/config"
	"e.coding.net/anyun-cloud-common/app"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app.RunApplication(config.GatewayInitFlags, config.Bootstrap)
}

// func main() {
// 	val := "test"
// 	change(&val)
// 	fmt.Println(val)
// }

// func change(v *string) {
// 	*v = "abc"
// }

// func main() {
// 	client, err := sql.Open("mysql", "root:1234qwer@/tcp:client.mysql.service.consul:3306/test?charset=utf8")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	stmtIns, err1 := client.Prepare("select 1=1")
// 	if err1 != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(stmtIns)
// }
