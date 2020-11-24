package ttp

import (
	"fmt"
	"net"
	"github.com/xela07ax/toolsXela/noRegx101"
	"strings"
)

func GetAllipLocal()(ips []string, err error)  {
	ifaces, err := net.Interfaces()
	if err != nil{
		return ips,err
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil{
			return ips,err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ips = append(ips,fmt.Sprintf("%v",ip))
			// process IP address
		}
	}
	return
}

func GetTryCleanIp()(tryip string, err error)  {
	//fmt.Sprintf("Не удалось определить ваш ip адрес: %v",err)
	// Старается найти рабочий ip адрес, если фильров не хватило, то возвращает через разделитель " | "
	ips, err := GetAllipLocal()
	if err != nil{
		return tryip,err
	}
	var validIps []string
	for _,ip := range ips{
		// Если начинается со 127, скорее всего это локальный адрес который не интересен
		if noRegx101.FindStrinTemplateNoRx(ip, `127`) {
			continue
		}
		// Если присутствуют двоеточия, то скорее всего это ipv6 или что то точно не нужное
		if noRegx101.FindStrinTemplateNoRx(ip, ":") {
			continue
		}
		validIps  = append(validIps,ip)
	}

	return strings.Join(validIps," | "), nil
}
