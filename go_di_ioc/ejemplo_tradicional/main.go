package main

import (
	"fmt"

	"github.com/nmarsollier/go_di_ioc/ejemplo_tradicional/model/hello/dao"
	"github.com/nmarsollier/go_di_ioc/ejemplo_tradicional/model/hello/service"
)

func main() {

	srv := service.NewService(dao.NewDao())

	fmt.Println(srv.SayHello())
}
