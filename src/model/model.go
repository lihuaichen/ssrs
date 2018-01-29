package model

import (
	"encoding/base64"
	"strconv"
	"net"
	"os/exec"
)

type SSR struct {
	Server_port    int
	Password       string
	Method         string
	Protocol       string
	Protocol_param string
	Obfs           string
	Obfs_param     string
}

type SSR_URL struct {
	Remarks string
	Url     string
}

type Json_proxy struct {
	State int
	Token string
	Host  string
	Url   []SSR_URL
	Err   string
}

type Json_root struct {
	State   int
	Token   string
	Remarks string
}

func urlencode(str string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(str))
}

func (s SSR) ToUrl(host string, remarks string, group string) SSR_URL {
	var protocol string
	if s.Protocol != "" {
		protocol = s.Protocol
	} else {
		protocol = "origin"
	}
	var obfs string
	if s.Obfs != "" {
		obfs = s.Obfs
	} else {
		obfs = "plain"
	}
	main_part := host + ":" + strconv.Itoa(s.Server_port) + ":" + protocol + ":" + s.Method + ":" + obfs + ":" + urlencode(s.Password)
	param_str := "obfsparam=" + urlencode(s.Obfs_param)
	if s.Protocol_param != "" {
		param_str += "&protoparam=" + urlencode(s.Protocol_param)
	}
	if remarks != "" {
		param_str += "&remarks=" + urlencode(remarks)
	}
	param_str += "&group=" + urlencode(group)
	url := "ssr://" + urlencode(main_part+"/?"+param_str) + "\n"
	ssrurl := SSR_URL{Remarks: remarks, Url: url}
	return ssrurl
}

func (s SSR) Port_open() bool {
	addr := "127.0.0.1:" + strconv.Itoa(s.Server_port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (s SSR) Restart(commit string) (bool, error) {
	cmd := exec.Command(commit)
	_, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	return s.Port_open(), nil
}
