package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type configInfo struct {
	Service string      `json:"service"`
	Name    string      `json:"name"`
	Token   string      `json:"token"`
	Domains domainInfos `json:"domains"`
}

// only for 'A' record
type domainInfo struct {
	Domain     string   `json:"domain"`
	SubDomains []string `json:"sub_domains"`
}
type domainInfos []*domainInfo

func loadConfig(file string) (*configInfo, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config configInfo
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Config info: %#v\n", config.Domains[0])
	return &config, nil
}
