package dnspod

import (
	"encoding/json"
)

// API Address: http://www.dnspod.cn/docs/index.html
//
// API Common Args:
//     login_token: token_id,token
//     format: [json, xml] default: xml, recommend: json
//     lang: [en, cn] default: en, recommend: cn
//     error_on_empry: [yes, no] default: yes, recommend: no
//     user_id: uid optional
//     login_code: D code only for enable D
//     login_remember: [yes, no] default yes. The result contains a cookie(t+uid), the following requests with the cookie.
//
// API Domain Create:
//     domain: no 'www', such as: 'dnspod.com'
//     group_id: optional
//     is_mark: [yes, no] optional
//
// API Domain List
//     type: [all,mine,share,ismark,pause,vip,recent,share_out] default: all optional
//     offset: default: 0 optional
//     length: the domain number optional
//     group_id: optional
//     keyword: optional
//
// API Record Create
//     domain_id: from Domain.List
//     sub_domain: default '@'
//     record_type: from Record.List, such as: 'A', 'CNAME', 'MX'
//     record_line: default '默认'
//     value: Such as 'A' is 'ip', 'CNAME' is 'domain', 'MX' is 'mail'
//     mx: only for 'MX'
//     ttl: [1 - 604800] optional
//     status: [enable, disable] optional
//     weight: [0 - 100] only for vip optional
//
// API Record List
//     domain_id
//     offset: default: 0 optional
//     length: the domain number optional
//     sub_domain: optional
//     keyword: optional
//
// API Record Modify
//     domain_id
//     record_id
//     sub_domain: optional
//     record_type
//     record_line
//     value
//     ... same as create

const (
	apiRoot = "https://dnsapi.cn"

	// data format: 'login_token=xxx&format=json&lang=en'
	apiVersion = apiRoot + "/Info.Version"
	// data format: 'login_token=xxx&format=json&lang=en'
	apiUserInfo = apiRoot + "/User.Detail"
	// data format: 'login_token=xx&format=json&lang=en&domain=xx'
	apiDomainCreate = apiRoot + "/Domain.Create"
	// data format: 'login_token=xxx&format=json&lang=en'
	apiDomainList = apiRoot + "/Domain.List"
	// data format: 'login_token=xxx&format=json&lang=en&domain_id=xx&sub_domain=xx&record_type=A@record_line=默认&value=xx.xx.xx.xx'
	apiRecordCreate = apiRoot + "/Record.Create"
	// data format: 'login_token=xxx&format=json&lang=en&domain_id=xx'
	apiRecordList = apiRoot + "/Record.List"
	// data format: 'login_token=xxx&format=json&lang=en&domain_id=xx&record_id=xx&record_type=A@record_line=默认&value=xx.xx.xx.xx'
	apiRecordModify = apiRoot + "/Record.Modify"
)

type statusResponse struct {
	Status struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		CreatedAt string `json:"created_at"`
	} `json:"status"`
}

type domainResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Punycode string `json:"punycode"`
	Grade    string `json:"grade"`
	Owner    string `json:"owner"`
	TTL      string `json:"ttl"`
}

type domainListResponse struct {
	statusResponse
	Info struct {
		DomainTotal string `json:"domain_total"`
		AllDomain   string `json:"all_domain"`
	} `json:"info"`
	Domains []domainResponse `json:"domains"`
}

type recordCommonInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Value   string `json:"value"`
	Line    string `json:"line"`
	Type    string `json:"A"`
	TTL     string `json:"ttl"`
	Enabled string `json:"enabled"`
}

type recordAddInfo struct {
	ID string `json:"id"`
	recordCommonInfo
}
type recordAddInfos []*recordAddInfo

type recordModifyInfo struct {
	ID int `json:"id"`
	recordAddInfo
}

type recordResponse struct {
	statusResponse
	Record recordAddInfo `json:"record"`
}

type recordMResponse struct {
	statusResponse
	Record recordModifyInfo `json:"record"`
}

type recordListResponse struct {
	statusResponse
	domainResponse
	Info struct {
		SubDomains string         `json:"sub_domains"`
		Records    recordAddInfos `json:"records"`
	} `json:"info"`
	Records recordAddInfos `json:"records"`
}

func (info *DNSPodInfo) getVersion() (*statusResponse, error) {
	params := info.genCommonParams()
	data, err := info.doPost(apiVersion, params)
	if err != nil {
		return nil, err
	}
	var resp statusResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (info *DNSPodInfo) listDomain() (*domainListResponse, error) {
	params := info.genCommonParams()
	data, err := info.doPost(apiDomainList, params)
	if err != nil {
		return nil, err
	}
	var resp domainListResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (info *DNSPodInfo) createRecord(name string,
	record *RecordInfo) (*recordResponse, error) {
	params := info.genCommonParams()
	params["domain"] = []string{name}
	params["sub_domain"] = []string{record.Name}
	params["record_type"] = []string{"A"}
	params["record_line"] = []string{"默认"}
	params["value"] = []string{record.Value}
	data, err := info.doPost(apiRecordCreate, params)
	if err != nil {
		return nil, err
	}
	var resp recordResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (info *DNSPodInfo) modifyRecord(name, id string,
	record *RecordInfo) (*recordMResponse, error) {
	params := info.genCommonParams()
	params["domain"] = []string{name}
	params["record_id"] = []string{id}
	params["sub_domain"] = []string{record.Name}
	params["record_type"] = []string{"A"}
	params["record_line"] = []string{"默认"}
	params["value"] = []string{record.Value}
	data, err := info.doPost(apiRecordModify, params)
	if err != nil {
		return nil, err
	}
	var resp recordMResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (info *DNSPodInfo) listRecord(name string) (*recordListResponse, error) {
	params := info.genCommonParams()
	params["domain"] = []string{name}
	data, err := info.doPost(apiRecordList, params)
	if err != nil {
		return nil, err
	}
	var resp recordListResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (records recordAddInfos) Exists(rcd *RecordInfo) (bool, bool, *recordAddInfo) {
	var (
		exists = false
		equal  = false
		resp   *recordAddInfo
	)
	for _, v := range records {
		if v.Name != rcd.Name {
			continue
		}
		exists = true
		resp = v
		if v.Value == rcd.Value {
			equal = true
		}
		break
	}
	return exists, equal, resp
}
