package cfg

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type ApiCfg struct {
	XMLName   xml.Name `xml:"ApiTstCfg"` // 指定最外层的标签为ApiTstCfg
	DirUrl    string   `xml:"dir_addr"`
	AppId     uint64   `xml:"app_id"`
	ZoneId    uint32   `xml:"zone_id"`
	Signature string   `xml:"signature"`
}

var ApiConfig *ApiCfg = nil

func ReadApiCfg(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadFile " + filePath + " err:" + err.Error())
		return err
	}

	ApiConfig = new(ApiCfg)
	err = xml.Unmarshal(data, ApiConfig)
	if err != nil {
		fmt.Println("xml Unmarshal " + filePath + " err:" + err.Error())
		return err
	}

	fmt.Printf("API CONF: %+v\n", *ApiConfig)
	return nil
}
