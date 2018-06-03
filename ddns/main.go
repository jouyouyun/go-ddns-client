package main

import (
	"ddns/dnspod"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Manager struct {
	config *configInfo
}

var (
	_m *Manager
)

const (
	configPath = "/etc/ddns-client/config.json"
)

func main() {
	if _m != nil {
		fmt.Println("There has a manager running, exit!")
		return
	}

	_m = new(Manager)
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}
	_m.config = config
	ip, err := getIP()
	if err != nil {
		fmt.Println("Failed to get ip:", err)
		return
	}
	fmt.Println("IP:", ip)

	if _m.config.Service != "dnspod" {
		fmt.Println("Unsupported")
		return
	}

	info := _m.toDNSPod(ip)
	info.DDNS()
}

func (m *Manager) toDNSPod(ip string) *dnspod.DNSPodInfo {
	var info dnspod.DNSPodInfo
	info.Token = m.config.Token
	for _, d := range m.config.Domains {
		var rcds dnspod.RecordInfos
		for _, v := range d.SubDomains {
			rcds = append(rcds, &dnspod.RecordInfo{
				Name:  v,
				Value: ip,
			})
		}
		info.Domains = append(info.Domains, &dnspod.DomainInfo{
			Name:    d.Domain,
			Records: rcds,
		})
	}
	return &info
}

func getIP() (string, error) {
	resp, err := http.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	l := len(data)
	v := string(data)
	if data[l-1] == '\n' {
		v = string(data[0 : l-1])
	}
	return v, nil
}
