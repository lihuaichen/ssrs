package sock

import (
	"sync"
	"net"
	"proxy/base"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"model"
	"strconv"
)

func Sock(wg *sync.WaitGroup) {
	defer wg.Done()
	var ch = make(chan int, 20)
	config, _ := base.CONFIG()
	addr := "0.0.0.0:" + strconv.Itoa(config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		base.LOGs("proxy/sock Sock() Listen error")
		base.LOG(err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			base.LOGs("proxy/sock Sock() Accept error")
			base.LOG(err)
			return
		}
		ch <- 1
		go handel(conn, ch)
	}
}

func int2bytes(i int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes()
}

func bytes2int(i []byte) int32 {
	buf := bytes.NewBuffer(i)
	var x int32
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

func chan_done(ch chan int) {
	<- ch
}

func handel(conn net.Conn, ch chan int) {
	defer chan_done(ch)
	defer base.Groupflush()
	defer conn.Close()
	var config, _ = base.CONFIG()
	var resp = model.Json_proxy{Token: config.Token, Host: config.Server.Host}
	var len int32
	var data model.Json_root
	buf := make([]byte, 4)
	l, err := conn.Read(buf)
	if l != 4 {
		base.LOGs("proxy/sock handel() Read error:len error")
		resp.State = 500
		resp.Err = "proxy/sock handel() Read error:len error"
		goto send
	}
	if err != nil {
		base.LOGs("proxy/sock handel() Read error")
		base.LOG(err)
		resp.State = 500
		resp.Err = "proxy/sock handel() Read error"
		goto send
	}
	len = bytes2int(buf)
	buf = make([]byte, len)
	data = model.Json_root{}
	conn.Read(buf)
	err = json.Unmarshal(buf, &data)
	if err != nil {
		base.LOGs("proxy/sock handel() json decode error")
		base.LOG(err)
		resp.State = 500
		resp.Err = "proxy/sock handel() json decode error"
		goto send
	}
	if data.Token != config.Token {
		base.LOGs("proxy/sock handel() token error")
		resp.State = 500
		resp.Err = "proxy/sock handel() token error"
		goto send
	}
	if data.State == 101 {
		ssr_url_list, err := base.All()
		if err != nil {
			resp.State = 500
			resp.Err = "proxy/sock handel() All error"
			goto send
		}
		resp.State = 200
		resp.Url = ssr_url_list
		goto send
	}
	if data.State == 102 {
		b, err := base.Only(data.Remarks)
		if err != nil {
			resp.State = 500
			resp.Err = "proxy/sock handel() Only error"
			goto send
		}
		if b == true {
			resp.State = 200
		} else {
			resp.State = 500
			resp.Err = "proxy/sock handel() Only error: false"
		}
		goto send
	}
send:
	buf, err = json.Marshal(resp)
	if err != nil {
		base.LOGs("proxy/sock handel() json encode error")
		base.LOG(err)
		return
	}
	len = int32(bytes.Count(buf, nil) - 1)
	conn.Write(int2bytes(len))
	conn.Write(buf)
}
