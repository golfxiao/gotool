package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

func GetDownFile(releaseId int) (map[string]*TReleaseFile, error) {
	sql := fmt.Sprintf(`
		SELECT
			id,
			path,
			version,
			checksum,
			storage_path,
			size,
			release_date
		FROM
			us_file_element
		WHERE
			release_id = %d
	`, releaseId)

	var files []*TReleaseFile
	num, err := orm.NewOrm().Raw(sql).QueryRows(&files)
	if err != nil {
		return nil, err
	}
	if num < 1 {
		return nil, fmt.Errorf("Unexpected release file empty")
	}

	result := make(map[string]*TReleaseFile, num)
	for _, fileInfo := range files {
		result[fileInfo.Path] = fileInfo
	}
	return result, nil
}
