package service

import (
	"acmweb/app/authorization/model"
	"acmweb/constants"
	"acmweb/system"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"github.com/pkg/errors"
	"strings"
	"time"
)

///*
//#cgo CXXFLAGS: -std=c++11
//#cgo  LDFLAGS: -L ../../middleware/lic_cgo/ -llicense
//#include "../../middleware/lic_cgo/lic_interface.hpp"
//*/
//import "C"

type GenLicense struct {
	qryId       int64
	srcData     string
	destData    string
	destDataMd5 string
	xmlInfo     *model.ACMLicense
}

func NewGenLic(id int64) *GenLicense {
	return &GenLicense{qryId: id}
}

func (o *GenLicense) GetLicData() string {
	return o.destData
}

/*func (o *GenLicense) GenSign() error {
	//生成签名
	retData := C.AcmGenSignToStr(C.CString(o.srcData))
	dataLen := C.AcmLicStrLen()
	o.destData = string(((*[1 << 31]byte)(unsafe.Pointer(retData)))[0:int(dataLen):int(dataLen)])
	//fmt.Println(o.srcData)
	//system.Log.Info("Generator license sign success,len:", len(o.destData))
	var err error
	if o.destDataMd5, err = system.Secret.MD5.MD5(o.destData); err != nil {
		system.Log.Error("License data convert to MD5 error:", err)
		return err
	}
	return nil
}*/

func (o *GenLicense) GenSign() (err error) {
	//生成签名
	for idx, grp := range o.xmlInfo.LicenseGroup {
		if grp.Component == model.LIC_GRP_AUTH {
			if o.xmlInfo.LicenseGroup[idx].Signature, err = o.SignAuth(grp.License); err != nil {
				return err
			}
			//system.Log.Debugf("----auth sign[%s]", licGrp.Signature)
		} else {
			if o.xmlInfo.LicenseGroup[idx].Signature, err = o.Sign(grp.License); err != nil {
				return err
			}
		}
	}
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(&o.xmlInfo); err != nil {
		system.Log.Error("marshal xml to string err :%v\n", err)
		return err
	}
	o.destData = base64.StdEncoding.EncodeToString(buf.Bytes())
	if o.destDataMd5, err = system.Secret.MD5.MD5(o.destData); err != nil {
		system.Log.Error("License data convert to MD5 error:", err)
		return err
	}
	return nil
}

func (o *GenLicense) GenLicData() (string, error) {
	//查询申请信息
	licInfo, err := o.FindLicInfo(o.qryId)
	if err != nil {
		system.Log.Error("Find license information failed", err.Error())
	}
	//存储数据文件名
	hostIps := strings.Split(licInfo.Auth.Host, ",")
	filenamekey := system.Text.Conv.String(licInfo.AuthId) + "-" + system.Text.Conv.String(licInfo.Auth.Company) + "-" + hostIps[0]

	o.xmlInfo = model.NewLicXml(licInfo.Scope)
	setCompentVal := func(k, v string) model.License {
		return model.License{Component: k, Value: v}
	}
	var licGrp model.LicenseGroup
	//Cert
	licGrp.Component = model.LIC_GRP_CERT
	licInfo.Cert.CertTime = time.Now().Format(constants.DBDATEFORMAT)
	licGrp.License = append(licGrp.License,
		setCompentVal(model.LIC_COMPONENT_EFFECTIVE, licInfo.Cert.CertTime),
		setCompentVal(model.LIC_COMPONENT_EXPIRATION, licInfo.Cert.ExpireDate))
	o.xmlInfo.LicenseGroup = append(o.xmlInfo.LicenseGroup, licGrp)

	//Auth
	licGrp.Clear()
	licGrp.Component = model.LIC_GRP_AUTH
	licGrp.License = append(licGrp.License,
		setCompentVal(model.LIC_COMPONENT_CODE, licInfo.Auth.Code),
		setCompentVal(model.LIC_COMPONENT_COMPANY, licInfo.Auth.Company),
		setCompentVal(model.LIC_COMPONENT_HOST, licInfo.Auth.Host),
		setCompentVal(model.LIC_COMPONENT_MAC, licInfo.Auth.Mac),
		setCompentVal(model.LIC_COMPONENT_EXPIRATION, licInfo.Auth.Expiration))
	o.xmlInfo.LicenseGroup = append(o.xmlInfo.LicenseGroup, licGrp)

	//Info
	licGrp.Clear()
	licGrp.Component = model.LIC_GRP_INFO
	water := setCompentVal(model.LIC_COMPONET_WATER, system.Text.Conv.String(len(licInfo.InfoList.WaterMark)))
	water.Text = licInfo.InfoList.WaterMark
	licGrp.License = append(licGrp.License, water)
	o.xmlInfo.LicenseGroup = append(o.xmlInfo.LicenseGroup, licGrp)

	//Serv
	licGrp.Clear()
	licGrp.Component = model.LIC_GRP_SERV
	for _, v := range licInfo.ServList.Modules {
		lic := model.License{Component: v}
		licGrp.License = append(licGrp.License, lic)
	}
	o.xmlInfo.LicenseGroup = append(o.xmlInfo.LicenseGroup, licGrp)
	//User
	licGrp.Clear()
	licGrp.Component = model.LIC_GRP_USER
	o.xmlInfo.LicenseGroup = append(o.xmlInfo.LicenseGroup, licGrp)

	//resXML, err := xml.MarshalIndent(xmlInfo, " ", " ")
	//if err != nil {
	//	system.Log.Error("marshal xml err :%v\n", err)
	//	return err
	//}
	//_, _ = osFile.Write([]byte(xml.Header))
	//_, _ = osFile.Write(resXML)
	//o.licFile = osFile.Name()
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(&o.xmlInfo); err != nil {
		system.Log.Error("marshal xml to string err :%v\n", err)
		return filenamekey, err
	}
	o.srcData = buf.String()
	//system.Log.Debugf("License data[%s]", o.srcData)
	//	system.Log.Info("Generator license data success,len:", len(o.srcData))
	return filenamekey, nil
}
func (o *GenLicense) GetLicFile() string {
	return ""
}
func (o *GenLicense) GetLicMd5() string {
	return o.destDataMd5
}

func (o *GenLicense) FindLicInfo(id int64) (model.LicInfo, error) {
	var licInfo model.LicInfo
	sql := "SELECT a.auth_id,a.trial_days,a.unit_id,b.watermark,b.inst_name,c.type_code,a.host_ip,a.host_mac,group_concat(e.app_code) app_code,a.app_list " +
		",date_format(a.appv_time,'%Y-%m-%d') effective from auth_req a,scene b LEFT JOIN scene_app d on b.inst_id = d.inst_id LEFT JOIN app e on d.app_id = e.app_id,auth_mode c " +
		"  where a.inst_id = b.inst_id and a.auth_mode = c.type_id and (auth_id =? or wf_inst_id = ?) "
	ds, err := system.Data.DB.QueryRows(sql, id, id)
	if err != nil {
		return licInfo, err
	}
	if len(ds) == 0 {
		return licInfo, errors.New("empty")
	}
	licInfo.AuthId = system.Text.Conv.Int64(ds[0]["auth_id"])
	licInfo.AppName = ds[0]["app_name"]
	licInfo.Scope = ds[0]["type_code"]
	licInfo.Auth.Host = ds[0]["host_ip"]
	licInfo.Auth.Mac = ds[0]["host_mac"]
	//licInfo.Auth.Code = ds[0]["app_code"]
	licInfo.Auth.Code = ds[0]["app_list"]
	licInfo.Cert.CertTime = ds[0]["effective"]
	trialDays := system.Text.Conv.Int(ds[0]["trial_days"])
	// 如果申请的是试用版，且应用的试用天数大于0，则计算证书到期日期
	if trialDays > 0 {
		certDate, _ := time.Parse(constants.DBDATEFORMAT, licInfo.Cert.CertTime)
		certDate = certDate.AddDate(0, trialDays, 0)
		licInfo.Cert.ExpireDate = certDate.Format(constants.DBDATEFORMAT)
	} else {
		licInfo.Cert.ExpireDate = ""
	}
	//licInfo.Auth.Expiration = licInfo.Cert.ExpireDate
	licInfo.InfoList.WaterMark = ds[0]["watermark"]
	licInfo.Auth.Company = ds[0]["unit_id"]
	return licInfo, nil
}

func (o *GenLicense) FindLicData() (string, error) {
	sql := "SELECT a.cert from auth_req a where auth_id=?"
	ds, err := system.Data.DB.QueryRows(sql, o.qryId)
	if err != nil {
		return "", err
	}
	if len(ds) == 0 {
		return "", nil
	}
	zipFile := system.Config.Attachment.FilePath + "/" + ds[0]["cert"]
	return zipFile, nil
}

func (o *GenLicense) Sign(params []model.License) (result string, err error) {
	var origSign string
	for _, value := range params {
		origSign += "|" + value.Component + ":" + value.Value
	}

	if encodeByte, err := system.Secret.DES.DesEncryopt(system.Text.Conv.Bytes(origSign), system.Text.Conv.Bytes("lyy20201207")); err != nil {
		system.Log.Error("License data convert to des error:", err)
		return result, err
	} else {
		if result, err = system.Secret.MD5.MD5(system.Text.Conv.String(encodeByte)); err != nil {
			system.Log.Error("Sign des data encryption md5 failed,err:", err.Error())
			return "", err
		}
		return result, nil
	}
}

func (o *GenLicense) SignAuth(params []model.License) (result string, err error) {
	var orginSign string
	for _, value := range params {
		orginSign += "|" + value.Component + ":" + value.Value
	}
	orginSign += "|"
	orginSign = "ZOC#@666//" + orginSign
	s1, _ := system.Secret.MD5.MD5(orginSign)
	orginSign = "Z1573Q2486D0901//" + orginSign
	s2, _ := system.Secret.MD5.MD5(orginSign)
	orginSign = "//X888@#" + orginSign
	s3, _ := system.Secret.MD5.MD5(orginSign)
	return s1 + s2 + s3, nil
}
