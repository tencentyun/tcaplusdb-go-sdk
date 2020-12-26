package cfg

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type TcaplusCase struct {
	XMLName xml.Name `xml:"TcaplusCase"` // 指定最外层的标签为TcaplusCase
	Head    CaseHead `xml:"Head"`
	Body    CaseBody `xml:"Body"`
}

type CaseHead struct {
	Cmd                    CaseCmd          `xml:"Cmd"`
	AppID                  int              `xml:"AppID"`
	ZoneID                 int              `xml:"ZoneID"`
	ShardID                int              `xml:"ShardID"`
	TableName              string           `xml:"TableName"`
	DirUrl                 string           `xml:"DirUrl"`
	AppSignUp              string           `xml:"AppSignUp"`
	ResultFlag             int              `xml:"ResultFlag"`
	AsyncID                uint64           `xml:"AsyncID"`
	CheckDataVersionPolicy int              `xml:"CheckDataVersionPolicy"`
	DataVersion            int              `xml:"DataVersion"`
	SpeedControl           CaseSpeedControl `xml:"SpeedControl"`
}

type CaseCmd struct {
	CmdType    string   `xml:"Type"`
	RandCmd    []string `xml:"RandCmd"`
	CmdTypeInt int
	RandCmdInt []int
}

type CaseSpeedControl struct {
	StartSpeed        int   `xml:"StartSpeed"`
	Step              int64 `xml:"Step"`
	SpeedLimit        int   `xml:"SpeedLimit"`
	AllowErrorNum     int64 `xml:"AllowErrorNum"`
	SpeedChangePeriod int   `xml:"SpeedChangePeriod"`
	ErrSleepSec       int   `xml:"ErrSleepSec"`
	MaxAllowAvgTime   int   `xml:"MaxAllowAvgTime"`
}

type CaseBody struct {
	KeyInfoList   []CaseKeyInfo   `xml:"KeyInfo"`
	ValueInfoList []CaseValueInfo `xml:"ValueInfo"`
}

type CaseKeyInfo struct {
	FieldName    string `xml:"FieldName,attr"`
	FieldType    string `xml:"FieldType,attr"`
	FieldTypeInt int
	FieldBuff    string `xml:"FieldBuff,attr"`
	KeyStep      int    `xml:"KeyStep,attr"`
	KeyRange     int    `xml:"KeyRange,attr"`
}

type CaseValueInfo struct {
	FieldName    string `xml:"FieldName,attr"`
	FieldType    string `xml:"FieldType,attr"`
	FieldTypeInt int
	FieldBuff    string `xml:"FieldBuff,attr"`
	ValueStep    int    `xml:"ValueStep,attr"`
	ValueRange   int    `xml:"ValueRange,attr"`
}

func ParseCase(filePath string) (*TcaplusCase, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadFile " + filePath + " err:" + err.Error())
		return nil, err
	}

	tCase := new(TcaplusCase)
	err = xml.Unmarshal(data, tCase)
	if err != nil {
		fmt.Println("xml Unmarshal " + filePath + " err:" + err.Error())
		return nil, err
	}

	return tCase, nil
}
