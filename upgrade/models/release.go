package models

import (
	"fmt"
	"gotool/upgrade/util"
	"strings"

	"github.com/astaxie/beego/orm"
)

func SaveRelease(releaseInfo *TRelease) error {
	success := false
	o := orm.NewOrm()
	defer func() {
		if success {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()

	releaseId, err := NewReleaseId()
	if err != nil {
		return err
	}
	o.Begin()

	//insert release notes
	sqlVals := releaseInfo.GetReleaseNotes(releaseId)
	if len(sqlVals) > 0 {
		sql := fmt.Sprintf(`
		INSERT INTO
  			us_release_notes(release_id, lang, release_notes)
		VALUES
  			%s
		`, strings.Join(sqlVals, ","))
		if _, err := o.Raw(sql).Exec(); err != nil {
			return fmt.Errorf("save release notes error: %s", err.Error())
		}
	}

	//insert release file element
	sqlVals = releaseInfo.GetReleaseFileInfo(releaseId)
	if len(sqlVals) > 0 {
		sql := fmt.Sprintf(`
		INSERT INTO
  			us_file_element(
				path,
				version,
				storage_path,
				release_date,
				checksum,
				size,
				release_id
  			)
		VALUES 
			%s
		`, strings.Join(sqlVals, ","))
		if _, err := o.Raw(sql).Exec(); err != nil {
			return fmt.Errorf("save release file error: %s", err.Error())
		}
	}

	//insert release site
	sql := fmt.Sprintf(`
		INSERT INTO
			us_site_release (
				release_id,
				release_version,
				site_id,
				application_id,
				client_type,
				status,
				extend_attr,
				create_time,
				update_time
			)
		VALUES
			(% d, '%s', '%s', % d, % d, %d, '%s', '%s', '%s')
		`,
		releaseId, releaseInfo.Version, releaseInfo.Site, releaseInfo.Application,
		releaseInfo.ClientType, 1, releaseInfo.SetWSMiddleAttr(),
		releaseInfo.Time, releaseInfo.Time)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		return err
	}

	success = true
	return nil
}

func NewReleaseId() (int64, error) {
	success := false

	o := orm.NewOrm()
	defer func() {
		if success {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()

	o.Begin()
	sql := fmt.Sprintf("UPDATE us_ticket set max_id=max_id+1 WHERE `key_name`='release'")
	_, err := o.Raw(sql).Exec()
	if err != nil {
		return 0, err
	}
	var releaseId int64
	sql = "SELECT max_id FROM us_ticket WHERE `key_name`='release'"
	err = o.Raw(sql).QueryRow(&releaseId)
	if err != nil {
		return 0, err
	}

	success = true
	return releaseId, nil
}

func GetDefaultRelease(appId, siteId int64) (releaseInfo *TRelease, err error) {
	sql := fmt.Sprintf(`
		SELECT
			site_id,
			release_version,
			release_id,
			extend_attr,
			create_time
		FROM
			gnet_tang_us_site_release
		WHERE
			application_id = %d
		AND 
			status = 20
		AND 
			site_id = '%d'
		ORDER BY
			id DESC
		LIMIT
			1
	`, appId, siteId)

	var row []orm.Params
	num, err := orm.NewOrm().Raw(sql).Values(&row)
	if err != nil {
		return nil, fmt.Errorf("get site release info error, error msg:%s", err.Error())
	}
	if num < 1 {
		return nil, nil
	}
	return getReleaseRow(row[0]), nil
}

func getReleaseRow(row orm.Params) (r *TRelease) {
	r = &TRelease{}
	r.Site = util.ToString(row["site_id"])
	r.Version = util.ToString(row["release_version"])
	r.ReleaseId = util.ToInt64(row["release_id"], 0)
	r.Time = util.ToString(row["create_time"])
	return r
}
