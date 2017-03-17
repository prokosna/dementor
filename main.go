package main

import (
	"fmt"
	"github.com/prokosna/dementor/dementor"
)

func main() {
	id, err := api.Authenticate("http://192.168.0.102:8081/", "azkaban", "azkaban")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(id)
	//err = api.CreateProject("http://192.168.0.102:8081/", id, "test1", "desc1")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//err = api.UploadProjectZip("http://192.168.0.102:8081/", id, "test1", "./hoge.zip")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	err = api.ScheduleFlow("http://192.168.0.102:8081/", id, "t", "1", `0 23/30 5,7-10 ? * 6#3`)
	if err != nil {
		fmt.Println(err)
		return
	}
}