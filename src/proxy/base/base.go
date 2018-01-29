package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"log"
	"model"
	"encoding/json"
	"net/http"
)

type config struct {
	Url      string
	Token    string
	Server   server
	Ssr_list []service
}

type server struct {
	Host string
	Port int
}

type service struct {
	Config_file  string
	Remarks string
	Restart string
}

var conf *config = nil
var logger *log.Logger = nil
var group string

func CONFIG() (config, error) {
	if conf == nil {
		yamlFile, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			LOGs("proxy/base CONFIG() ReadFile error")
			LOG(err)
			return *conf, err
		}
		var con config
		err = yaml.Unmarshal(yamlFile, &con)
		if err != nil {
			LOGs("proxy/base CONFIG() yaml decode error")
			LOG(err)
			return *conf, err
		}
		conf = &con
	}
	return *conf, nil
}

func LOG(err error) {
	if logger == nil {
		newLOG()
	}
	logger.Println(err.Error())
}

func LOGs(err string) {
	if logger == nil {
		newLOG()
	}
	logger.Println(err)
}

func newLOG() {
	if logger == nil {
		file, err := os.Create("err.log")
		if err != nil {
			log.Fatal("proxy/base.LOG() os.Create ERROR:" + err.Error())
		}
		logger = log.New(file, "[ERROR]", log.LstdFlags)
	}
}

func Group() string {
	if group == "" {
		res, err := http.Get(conf.Url + "/api/group.php")
		if err != nil {
			LOGs("proxy/base Group() Get error")
			LOG(err)
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			LOGs("proxy/base Group() ReadAll error")
			LOG(err)
		}
		group = string(result)
	}
	return group
}

func Groupflush() {
	group = ""
}

func All() ([]model.SSR_URL, error) {
	service_len := len(conf.Ssr_list)
	ssr_url_list := make([]model.SSR_URL, service_len)
	var s service
	for _, s = range conf.Ssr_list {
		str_conf, err := ioutil.ReadFile(s.Config_file)
		if err != nil {
			LOGs("proxy/base All() ReadFile error")
			LOG(err)
			return ssr_url_list, err
		}
		ssr := model.SSR{}
		err = json.Unmarshal(str_conf, &ssr)
		if err != nil {
			LOGs("proxy/base All() json decode error")
			LOG(err)
			return ssr_url_list, err
		}
		if !ssr.Port_open() {
			b, err := ssr.Restart(s.Restart)
			if err != nil {
				LOGs("proxy/base All() Restart error")
				LOG(err)
				return ssr_url_list, err
			}
			if !b {
				LOGs("proxy/base All() Restart error: Can not restart " + s.Remarks)
				continue
			}
		}
		ssr_url := ssr.ToUrl(conf.Server.Host, s.Remarks, Group())
		ssr_url_list = append(ssr_url_list, ssr_url)
	}
	return ssr_url_list, nil
}

func Only(remarks string) (bool, error) {
	var s service
	for _, s = range conf.Ssr_list {
		if s.Remarks == remarks {
			str_conf, err := ioutil.ReadFile(s.Config_file)
			if err != nil {
				LOGs("proxy/base Only() ReadFile error")
				LOG(err)
				return false, err
			}
			ssr := model.SSR{}
			err = json.Unmarshal(str_conf, &ssr)
			if err != nil {
				LOGs("proxy/base Only() json decode error")
				LOG(err)
				return false, err
			}
			if !ssr.Port_open() {
				b, err := ssr.Restart(s.Restart)
				if err!= nil {
					LOGs("proxy/base Only() Restart error")
					LOG(err)
					return false, err
				}
				return b, nil
			} else {
				return true, err
			}
		}
	}
	return false, nil
}
