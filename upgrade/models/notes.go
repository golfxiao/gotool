package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

func GetReleaseNote(releaseInfo *TRelease) (err error) {
	sql := fmt.Sprintf(`
		SELECT
			lang,
			release_notes
		FROM
			us_release_notes
		WHERE
			release_id = %d
	`, releaseInfo.ReleaseId)

	var notes []*TReleaseNotes
	_, err = orm.NewOrm().Raw(sql).QueryRows(&notes)
	if err != nil {
		return err
	}

	releaseInfo.Notes = notes
	return nil
}
