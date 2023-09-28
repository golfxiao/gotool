package models

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/orm"
)

func GetDiffPatch(fileHash string) *TDiffPatch {
	sql := fmt.Sprintf(`
		SELECT
			id,
			storage_path,
			checksum,
			size,
			file_hash
		FROM
			us_diff_patch
		WHERE
			file_hash = '%s'
	`, fileHash)

	var patch TDiffPatch
	err := orm.NewOrm().Raw(sql).QueryRow(&patch)
	if err != nil {
		log.Printf("query patch error: %s", err.Error())
		return nil
	}
	return &patch
}

func SaveDiffPatch(patch *TDiffPatch) error {
	sql := fmt.Sprintf(`
		INSERT INTO
  			us_diff_patch (storage_path, checksum, size, file_hash)
		VALUES
  			('%s', '%s', %d, '%s')
		`,
		patch.StoragePath, patch.CheckSum, patch.Size, patch.FileHash)

	_, err := orm.NewOrm().Raw(sql).Exec()
	if err != nil {
		return err
	}

	return nil
}
