package dnspod

import (
	"ddns/utils"
	"fmt"
	"net/url"
)

// DNSPodInfo dnspod account and domains info
type DNSPodInfo struct {
	Token   string
	Domains DomainInfos
}

// DomainInfo domain and records info
type DomainInfo struct {
	Name    string
	Records RecordInfos
}
type DomainInfos []*DomainInfo

// RecordInfo records info
type RecordInfo struct {
	Name  string
	Value string
}
type RecordInfos []*RecordInfo

// DDNS dnspod ddns query and modify
func (info *DNSPodInfo) DDNS() {
	for _, d := range info.Domains {
		info.parseDomain(d)
	}
}

func (info *DNSPodInfo) parseDomain(domain *DomainInfo) {
	rcdResp, err := info.listRecord(domain.Name)
	if err != nil {
		fmt.Printf("Failed to list records for (%s): %v\n", domain.Name, err)
		// add records
		info.addDomain(domain)
		return
	}
	for _, rcd := range domain.Records {
		exists, equal, cur := rcdResp.Records.Exists(rcd)
		fmt.Println("Exists:", exists, equal, cur)
		if !exists {
			_, err := info.createRecord(domain.Name, rcd)
			if err != nil {
				fmt.Println("Failed to add record:", err)
			}
			continue
		}
		if equal {
			continue
		}
		_, err := info.modifyRecord(domain.Name, cur.ID, rcd)
		if err != nil {
			fmt.Println("Failed to modify record:", err, rcd)
		}
	}
}

func (info *DNSPodInfo) addDomain(d *DomainInfo) {
	for _, v := range d.Records {
		_, err := info.createRecord(d.Name, v)
		if err != nil {
			fmt.Println("Failed to add record:", d.Name, v)
			continue
		}
	}
}

func (info *DNSPodInfo) genCommonParams() url.Values {
	var params = make(url.Values)
	params["login_token"] = []string{info.Token}
	params["format"] = []string{"json"}
	params["lang"] = []string{"en"}
	return params
}

func (info *DNSPodInfo) doPost(url string, params url.Values) ([]byte, error) {
	return utils.PostForm(url, params)
}
