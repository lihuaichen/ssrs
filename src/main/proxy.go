package main

import (
	"proxy/reg"
	"sync"
	"proxy/sock"
	"proxy/base"
)

func main() {
	wg := &sync.WaitGroup{}
	_, err := base.CONFIG()
	if err != nil {
		base.LOGs("proxy/main main() CONFIG error")
		base.LOG(err)
		return
	}
	for {
		wg.Add(2)
		go reg.Reg(wg)
		go sock.Sock(wg)
		wg.Wait()
	}
}
