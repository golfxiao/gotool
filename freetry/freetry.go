package freetry

import (
	"errors"
	"strings"
	"time"
)

var (
	defaultFreeConfig = make(map[string]FreeItem)
)

type FreeItem struct {
	Type int64 `json:"type" bson:"type"` // 试用方式，0：按次数，1：按期限
	Freq int64 `json:"freq" bson:"freq"` // 允许试用的次数
	Exp  int64 `json:"exp" bson:"exp"`   // 允许试用的期限，单位：天
}

type UserFreeItem struct {
	Type int64 `json:"type" bson:"type"` // 试用方式
	Freq int64 `json:"freq"`             // 允许试用的次数，
	Exp  int64 `json:"exp"`              // 允许试用的期限
	Used int64 `json:"used" bson:"used"` // 已试用次数，只适用于按次数
	St   int64 `json:"st" bson:"st"`     // 试用的开始时间，只适用于按期限
}

type UserFreeSetting struct {
	UserId     int64                    `json:"userId" bson:"userId"`
	FreeConfig map[string]*UserFreeItem `json:"freeconfig" bson:"freeconfig"`
}

func treeFreeUse(item *UserFreeItem) (result bool) {
	if item.Type == 1 { // 按期限
		result = time.Now().Unix() < item.St+item.Exp*3600*24
	} else { // 按次数
		result = item.Used < item.Freq
	}
	if result {
		item.Used++
	}
	return
}

func ApplyFreeUse(userConfig map[string]*UserFreeItem,
	feature string, defaultConfig map[string]FreeItem) error {
	item, ok := userConfig[feature]
	if !ok {
		for k, v := range defaultConfig {
			if strings.HasPrefix(feature, k) {
				item = &UserFreeItem{
					Type: v.Type,
					Freq: v.Freq,
					Used: 0,
					Exp:  v.Exp,
					St:   time.Now().Unix(),
				}
				userConfig[feature] = item
			}
		}
	}

	if item == nil {
		return errors.New("Not found this apply item: " + feature)
	}

	if !treeFreeUse(item) {
		return errors.New("Trial period expired or trial limit reached")
	}

	return nil
}
