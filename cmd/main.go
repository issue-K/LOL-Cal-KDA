package main

import (
	"log"
	"lolapi/App"
	admin "lolapi/pkg/window"
)

func init(){
	admin.MustRunWithAdmin()  //改成以管理员方式运行
}

func main(){
	admin.GetLolClientApiInfoV2()

	app := App.NewProphet()   //创建应用
	if err := app.Run(); err != nil{   //开始运行
		log.Fatal( err )
	}
}