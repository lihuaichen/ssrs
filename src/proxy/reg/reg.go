package reg

import (
	"proxy/base"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type Json_reg struct {
	Host  string
	Port  int
	Token string
}

func Reg(wg *sync.WaitGroup) bool {
	defer wg.Done()
	time.Sleep(5 *time.Second)
	config, _ := base.CONFIG()
	var r Json_reg
	r.Host = config.Server.Host
	r.Port = config.Server.Port
	r.Token = config.Token
	b, err := json.Marshal(r)
	if err != nil {
		base.LOGs("proxy/reg Reg() json encode error")
		base.LOG(err)
		return false
	}
	body := bytes.NewBuffer(b)
	res, err := http.Post(config.Url + "/api/reg.php", "application/json;charset=utf-8", body)
	if err != nil {
		base.LOGs("proxy/reg Reg() Post error")
		base.LOG(err)
		return false
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		base.LOGs("proxy/reg Reg() ReadAll error")
		base.LOG(err)
		return false
	}
	data := string(result)
	if strings.Contains(data, "601") {
		return true
	} else if strings.Contains(data, "401") {
		base.LOGs("proxy/reg reg() Token error")
	} else {
		base.LOGs("proxy/reg reg() UnKnow error:[" + data + "]")
	}
	return false
}
