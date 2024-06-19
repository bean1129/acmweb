package model

import (
	"acmweb/system"
	"encoding/xml"
)

const (
	LIC_SCOPE_TYPE_TV = "TV"
	LIC_SCOPE_TYPE_CV = "CV"
	LIC_SCOPE_TYPE_GV = "GV"
)
const (
	LIC_GRP_CERT = "Cert"
	LIC_GRP_AUTH = "Auth"
	LIC_GRP_INFO = "Info"
	LIC_GRP_SERV = "Serv"
	LIC_GRP_USER = "User"

	LIC_COMPONENT_CODE       = "Code"
	LIC_COMPONENT_COMPANY    = "Company"
	LIC_COMPONENT_PROJECT    = "Project"
	LIC_COMPONENT_HOST       = "Host"
	LIC_COMPONENT_MAC        = "MAC"
	LIC_COMPONENT_EFFECTIVE  = "Effective"
	LIC_COMPONENT_EXPIRATION = "Expiration"
	LIC_COMPONENT_FIFD       = "FIFD"

	LIC_COMPONET_TITLE = "Title"
	LIC_COMPONET_ABOUT = "About"
	LIC_COMPONET_WATER = "Watemark"
)

func NewLicXml(scope string) *ACMLicense {
	root := xml.Name{Local: "ACMLicense"}
	return &ACMLicense{XMLName: root, Scope: scope}
}

type ACMLicense struct {
	XMLName      xml.Name       `xml:"ACMLicense"`
	Scope        string         `xml:"scope,attr"`
	LicenseGroup []LicenseGroup `xml:"License-Group"`
}

type LicenseGroup struct {
	Text      string    `xml:",chardata"`
	Component string    `xml:"component,attr"`
	Signature string    `xml:"signature,attr"`
	License   []License `xml:"License"`
}

func (lg *LicenseGroup) Clear() {
	lg.Component = ""
	lg.Signature = ""
	lg.License = []License{}
}

type License struct {
	Text      string `xml:",chardata"`
	Component string `xml:"component,attr"`
	Value     string `xml:"value,attr"`
	Desc      string `xml:"desc,attr"`
	Max       string `xml:"max,attr,omitempty"`
}

type LicInfo struct {
	AuthId   int64
	AppName  string
	Scope    string
	Cert     CertGrp
	Auth     AuthGrp
	InfoList InfoGrp
	ServList ServGrp
}

type CertGrp struct {
	CertTime    string //控制证书24H内导入
	CertVersion string
	ExpireDate  string
}

type AuthGrp struct {
	Code       string
	Company    string
	Host       string
	Mac        string
	Effective  string //license生成时间
	Expiration string //试用期到期日（TV）
}
type InfoGrp struct {
	WaterMark string
}
type ServGrp struct {
	Modules []string
}
type UserGrp struct {
}

type LicDesc struct {
	XMLName   xml.Name `xml:"acm-lic"`
	Id        string   `xml:"id"`
	Code      string   `xml:"code"`
	Name      string   `xml:"name"`
	Issuedate string   `xml:"issuedate"`
	Proposer  string   `xml:"proposer"`
	OwnerIp   string   `xml:"ownerip"`
	OwnerMac  string   `xml:"ownermac"`
	FileName  string   `xml:"filename"`
	UnitCode  string   `xml:"unitcode"`
}

func (o *LicDesc) GetLicFileInfo(qryId int64) error {
	sql := "SELECT b.auth_id,a.app_code,a.app_name,c.unit_code,b.proposer,DATE_FORMAT(b.appv_time,'%Y-%m-%d')  appv_time ,b.host_ip,b.host_mac from app a,auth_req b,unit c " +
		"WHERE a.app_id = b.app_id and c.unit_id = b.unit_id and (wf_inst_id = ? or auth_id = ?)"
	ds, err := system.Data.DB.QueryRows(sql, qryId, qryId)
	if err != nil {
		return err
	}
	for _, rs := range ds {
		o.Id = rs["auth_id"]
		o.Code = rs["app_code"]
		o.Name = rs["app_name"]
		o.Issuedate = rs["release_date"]
		o.Proposer = rs["proposer"]
		o.OwnerMac = rs["host_mac"]
		o.OwnerIp = rs["host_ip"]
		o.UnitCode = system.Text.Conv.String(rs["unit_code"])
	}
	return nil
}
