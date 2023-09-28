package upgrade

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gotool/upgrade/models"
	"gotool/upgrade/util"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type USDownload struct {
	ReleaseId int      // 发布ID
	Path      []string // 客户端请求的差异化文件列表，以相对路径形式提供
	RootPath  string   // 升级业务的磁盘存储目录
}

func (self *USDownload) Download() (*models.TDiffPatch, error) {
	// 查询文件列表
	allFiles, err := models.GetDownFile(self.ReleaseId)
	if err != nil {
		return nil, fmt.Errorf("query release file error: %s", err.Error())
	}
	// 过滤出差异化文件
	diffFiles := self.filterDiffFiles(allFiles)
	if len(diffFiles) == 0 {
		return nil, fmt.Errorf("not match client upgrage file")
	}
	// 检查patch文件是否已经存在
	fileHash := self.getFileMd5(diffFiles)
	patch := models.GetDiffPatch(fileHash)
	if patch != nil && self.FileExists(patch.StoragePath) {
		log.Printf("patch file already exists, [%s] %s", fileHash, patch.StoragePath)
		return patch, nil
	}
	// 如果不存在，则生成一个新的patch
	patch = &models.TDiffPatch{FileHash: fileHash}
	patch.StoragePath, err = self.singleBuildPatch(fileHash, diffFiles)
	if err != nil {
		return nil, err
	}
	// 计算文件大小和内容校验和
	patch.CheckSum, patch.Size = self.calculateChecksumAndSize(patch.StoragePath)
	if err := models.SaveDiffPatch(patch); err != nil {
		return nil, err
	}
	return patch, nil
}

func (self *USDownload) filterDiffFiles(allFiles map[string]*models.TReleaseFile) []*models.TReleaseFile {
	diffFiles := make([]*models.TReleaseFile, 0, len(allFiles))
	for _, v := range self.Path { // self.Path为客户端请求的文件集, 相对路径的集合
		if fileInfo, ok := allFiles[v]; ok { // allFiles为从数据库查询出来的全部文件列表，以相对路径为Key
			diffFiles = append(diffFiles, fileInfo)
		}
	}
	return diffFiles
}

func (self *USDownload) getFileMd5(diffFiles []*models.TReleaseFile) string {
	pathList := make([]string, len(diffFiles)) // 获取差异文件存储路径列表
	for i, v := range diffFiles {
		pathList[i] = v.StoragePath
	}
	sort.StringSlice(pathList).Sort() // 给差异文件的路径排序

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(strings.Join(pathList, ","))) // 计算MD5值
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func (self *USDownload) calculateChecksumAndSize(storagePath string) (checksum string, size int64) {
	absPath := filepath.Join(self.RootPath, storagePath)
	file, err := os.Open(absPath)
	if err != nil {
		log.Printf("open file[%s] error: %s", absPath, err.Error())
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Printf("stat file size error: %s", err.Error())
		return
	}

	hash := sha1.New()
	io.Copy(hash, file)
	checksum = hex.EncodeToString(hash.Sum(nil))
	size = fileStat.Size()
	log.Printf("patchfile:%s, checksum:%s, size: %d", storagePath, checksum, size)
	return
}

func (self *USDownload) buildPatch(files []*models.TReleaseFile) (string, error) {
	packageFile, path, err := self.createPatchFile() // 创建压缩包，并打开文件
	if err != nil {
		return "", err
	}
	defer packageFile.Close()

	w := zip.NewWriter(packageFile) // 构造一个压缩格式的写操作句柄
	defer w.Close()

	for _, file := range files { // 逐个遍历文件，将文件内容写入压缩包
		if err := self.writeFile(w, file); err != nil {
			return "", err
		}
	}
	return strings.TrimPrefix(path, self.RootPath), nil // 返回压缩包路径
}

func (self *USDownload) createPatchFile() (*os.File, string, error) {
	packageDir := filepath.Join(self.RootPath, "packages")
	if !util.CreateDir(packageDir) {
		return nil, "", fmt.Errorf("create dir %s error", packageDir)
	}
	filePath := fmt.Sprintf("%s/%s.zip", packageDir, util.UUID())
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	return file, filePath, err
}

func (self *USDownload) writeFile(w *zip.Writer, file *models.TReleaseFile) error {
	if len(file.Path) < 2 || len(file.StoragePath) < 2 {
		log.Printf("invalid path or storage path: %s, %s", file.Path, file.StoragePath)
		return nil
	}

	srcPath := filepath.Join(self.RootPath, file.StoragePath)
	src, err := os.Open(srcPath) // 打开源文件句柄
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := w.Create(file.Path[1:]) // 打开目标文件句柄
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src) // 从源文件向目标文件拷贝内容
	if err != nil {
		return err
	}
	return nil
}
func (self *USDownload) FileExists(storagePath string) bool {
	packageFile := filepath.Join(self.RootPath, storagePath)
	if _, err := os.Stat(packageFile); os.IsNotExist(err) {
		log.Printf("storage path [%s] not exits", packageFile)
		return false
	} else {
		return true
	}
}

func (self *USDownload) singleBuildPatch(patchKey string, diffFiles []*models.TReleaseFile) (string, error) {
	result, err, _ := sg.Do(patchKey, func() (interface{}, error) {
		return self.buildPatch(diffFiles)
	})
	if err != nil {
		return "", err
	}
	patchPath := result.(string)
	return patchPath, nil
}
