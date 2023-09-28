package models

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

type TRelease struct {
	Version     string           `xml:"version,attr"`      //客户端升级包版本
	Time        string           `xml:"time,attr"`         //客户端升级包出包时间
	Application int64            `xml:"application,attr"`  //客户端升级appid，例如：6:pc 7:mac
	Site        string           `xml:"site,attr"`         //升级站点，对应产品标识
	ForceUpdate bool             `xml:"force_update,attr"` //强制升级标记位
	NeedUpdate  bool             `xml:"needupdate,attr"`   //是否升级标记位:1:是 0:否
	MinVersion  string           `xml:"min_version,attr"`  //直接到官网下载安装包最小版本号
	Mv          string           `xml:"mv,attr"`           //云会议升级最小版本配置
	ReleaseId   int64            `xml:"-"`                 //客户端发布标识
	ClientType  int64            `xml:"-"`                 //客户端类型,只用来区分云会议的pc和mac
	Notes       []*TReleaseNotes `xml:"notes,omitempty"`   //升级内容，可能包括多种语言
	File        []*TReleaseFile  `xml:"file"`              //文件列表
}

type TReleaseNotes struct {
	XMLName xml.Name `xml:"notes,omitempty" orm:"-"`
	Id      int64    `xml:"-" orm:"column(id);pk;auto"`            // 主键，自增
	Lang    string   `xml:"lang,attr" orm:"column(lang)"`          // 语言，如zh-cn, en-us
	Notes   string   `xml:",innerxml" orm:"column(release_notes)"` // 升级提示内容
}

type TReleaseFile struct {
	XMLName     xml.Name `xml:"file" orm:"-"`
	FileId      int64    `xml:"-" orm:"column(id);pk;auto"`           // 主键，自增
	Path        string   `xml:"path,attr" orm:"column(path)"`         // 文件在包内的相对路径
	Version     string   `xml:"version,attr" orm:"column(version)"`   // 文件所属包的版本号
	Time        string   `xml:"time,attr" orm:"column(release_date)"` // 文件所属包的发布时间
	CheckSum    string   `xml:"checksum,attr" orm:"column(checksum)"` // 文件的md5
	Size        int64    `xml:"size,attr" orm:"column(size)"`         // 文件的大小
	Url         string   `xml:"url,attr,omitempty" orm:"-"`           // 文件的url，适用于http(s)协议的文件
	StoragePath string   `xml:"-" orm:"column(storage_path)"`         // 存储路径
}

type TDiffPatch struct {
	Id          int64  `orm:"column(id);pk;auto"`
	Size        int64  `orm:"column(size)`         // 差异包大小
	CheckSum    string `orm:"column(checksum)`     // 差异包文件内容校验和
	StoragePath string `orm:"column(storage_path)` // 差异包存储路径
	FileHash    string `orm:"column(files_hash)`   // 差异包标识MD5, 用于判断文件是否已经生成
}

func (r *TRelease) SetFileStoragePath(fileStorageMap map[string]string) {
	for i, file := range r.File {
		//从压缩包中读出的相对路径，不带前面的"/",需要去掉
		if v, ok := fileStorageMap[file.Path[1:]]; ok {
			r.File[i].StoragePath = v
		}
		//对于andriod不需要去掉第一个字符
		if v, ok := fileStorageMap[file.Path]; ok {
			r.File[i].StoragePath = v
		}
	}
}

func (r *TRelease) CheckEmptyStoregeFile() bool {
	f := make([]string, 0, len(r.File))
	for _, fi := range r.File {
		if len(fi.StoragePath) < 1 {
			//过滤掉http(s)的文件
			if strings.HasPrefix(fi.Path, "http") {
				continue
			}
			f = append(f, fi.Path)
		}
	}
	//打印日志
	if len(f) > 0 {
		b, _ := json.Marshal(f)
		log.Printf("config.xml file contains file not in binary zip file,val:%s", string(b))
		return true
	}
	return false
}

func (r *TRelease) SetWSMiddleAttr() (s string) {
	extentAttr := map[string]string{
		"min_version": r.MinVersion,
	}
	b, _ := json.Marshal(extentAttr)
	return string(b)
}

func (r *TRelease) GetReleaseFileInfo(releaseId int64) (sqlval []string) {
	sqlval = make([]string, 0, len(r.File))
	for _, fi := range r.File {
		v := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', %d, %d)",
			fi.Path, fi.Version, fi.StoragePath,
			r.Time, fi.CheckSum, fi.Size, releaseId)
		sqlval = append(sqlval, v)
	}
	return sqlval
}

func (r *TRelease) GetReleaseNotes(releaseId int64) (sqlval []string) {
	sqlval = make([]string, 0, len(r.File))
	for _, fi := range r.Notes {
		v := fmt.Sprintf("(%d, '%s', '%s')", releaseId, fi.Lang, fi.Notes)
		sqlval = append(sqlval, v)
	}
	return sqlval
}
