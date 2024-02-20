package main

import (
	"astralstack-hyperv/boostrap"
	"astralstack-hyperv/web"
	"log"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {
	if err := boostrap.Init(); err != nil {
		log.Panicln(err)
		return
	}

	app := web.NewRoute()

	if err := app.Run(":8080"); err != nil {
		log.Panicln(err)
		return
	}

	//v, _ := mem.VirtualMemory()
	//fmt.Printf("Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)
	//fmt.Println(v)

	//name := "cn_windows_server_2019_datacenter_x64"
	//name := "1111"
	//_, err := hyperv.CreateVM(name, 2, 4096, hyperv.Generation1)
	//if err != nil {
	//	return
	//}
	//new, _ := compute.NewInstance(name)
	//
	//instance, err := compute.GetByName(name)
	//if err != nil {
	//	return
	//}
	//
	//_, err = instance.Stop()
	//if err != nil {
	//	return
	//}

	//instanceList, err := compute.GetAllInstance()
	//if err != nil {
	//	panic(err)
	//}
	//for _, instance := range *instanceList {
	//	fmt.Println(instance.Name)
	//}
}
