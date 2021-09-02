package main

import (
	"fmt"

	"github.com/DeniesKresna/beinventaris/Configs"
	"github.com/DeniesKresna/beinventaris/Routers"
	check "github.com/asaskevich/govalidator"
)

func main() {
	check.SetFieldsRequiredByDefault(true)
	if err := Configs.DatabaseInit(); err != nil {
		fmt.Println("status ", err)
	}

	Configs.DatabaseMigrate()

	r := Routers.SetupRouter()
	r.Run(":8090")
}
