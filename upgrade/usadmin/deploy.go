package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"gotool/upgrade/models"
	"gotool/upgrade/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type USDeploy struct {
	rootPath string
	release  *models.TRelease
}

func NewUSDeploy(rootPath string) *USDeploy {
	return &USDeploy{
		rootPath: rootPath,
	}
}

func (self *USDeploy) Deploy(deployFile, deployXml string) error {
	releaseInfo, err := self.readXmlContent(deployXml)
	if err != nil {
		return err
	}
	self.release = releaseInfo
	filePathMap, err := self.unzip(deployFile)
	if err != nil {
		return err
	}

	log.Printf("filePathMap:%v", filePathMap)
	self.release.SetFileStoragePath(filePathMap)
	err = models.SaveRelease(self.release)
	if err != nil {
		return err
	}
	// err1 := os.Remove(deployFile)
	// err2 := os.Remove(deployXml)
	// if err1 != nil || err2 != nil {
	// 	log.Printf("delete deployfile error: %v, %v", err1, err2)
	// }
	return nil
}

func (self *USDeploy) unzip(zfile string) (filePathMap map[string]string, err error) {
	rc, err := zip.OpenReader(zfile) // 打开压缩包的读取句柄
	if err != nil {
		return
	}
	defer rc.Close()

	storageDir, err := self.buildStorageDir() // 构造当前软件包的文件存储目录
	if err != nil {
		return nil, err
	}

	filePathMap = make(map[string]string) // 相对目录与存储目录的映射，发布配置时需要使用
	for i, f := range rc.File {
		if ok := self.isDir(f); ok { // 只存储文件，文件夹不处理
			continue
		}
		filePath := fmt.Sprintf("%s/f_%d", storageDir, i) // 文件存储路径
		if err = self.saveFile(f, filePath); err != nil { // 保存文件
			return
		}
		// 维护相对目录与存储目录的映射, 为便于存储目录迁移，这里只存rootPath以外的部分
		filePathMap[f.Name] = strings.TrimPrefix(filePath, self.rootPath)
	}
	return filePathMap, nil
}

func (self *USDeploy) saveFile(file *zip.File, storagePath string) (err error) {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("open file[%s] error: %s", file.Name, err.Error())
	}
	defer src.Close()

	dst, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	n, err := io.CopyN(dst, src, int64(file.UncompressedSize64))
	log.Printf("copy file [%s] to path [%s], copy size :%d", file.Name, storagePath, n)
	if err != nil {
		return err
	}
	return nil
}

// 构造当前软件包的文件存储目录，并创建该目录
// - rootPath: 可以理解为共享存储为升级业务分配的根存储目录，由外面指定
func (self *USDeploy) buildStorageDir() (string, error) {
	storageDir := fmt.Sprintf("%s/files/%d_%s", self.rootPath, self.release.Application, self.release.Version)
	if ok := util.CreateDir(storageDir); !ok {
		return "", fmt.Errorf("create storage dir error")
	}
	return storageDir, nil
}

func (self *USDeploy) isDir(file *zip.File) bool {
	return file.FileInfo().IsDir()
}

func (self *USDeploy) readXmlContent(xmlConfig string) (*models.TRelease, error) {
	file, err := os.Open(xmlConfig)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	xmlFileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	releaseInfo := &models.TRelease{}
	if err := xml.Unmarshal(xmlFileBytes, releaseInfo); err != nil {
		return nil, err
	}

	return releaseInfo, nil
}

func (self *USDeploy) Revoke(version string) error {
	// TODO
	return nil
}
