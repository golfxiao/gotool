package freetry

import (
	"testing"
	"time"
)

func TestTreeFreeUse(t *testing.T) {
	// Testing for item type 1 and current time is before the expiration date
	item1 := &UserFreeItem{
		Type: 1,
		St:   time.Now().Unix() - 3600*24,
		Exp:  7,
	}
	expected1 := true
	if result1 := treeFreeUse(item1); result1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, result1)
	}

	// Testing for item type 1 and current time is after the expiration date
	item2 := &UserFreeItem{
		Type: 1,
		St:   time.Now().Unix() - 3600*24*8,
		Exp:  7,
	}
	expected2 := false
	if result2 := treeFreeUse(item2); result2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, result2)
	}

	// Testing for item type 2 and used count is less than frequency
	item3 := &UserFreeItem{
		Type: 2,
		Used: 2,
		Freq: 5,
	}
	expected3 := true
	if result3 := treeFreeUse(item3); result3 != expected3 {
		t.Errorf("Expected %v but got %v", expected3, result3)
	}

	// Testing for item type 2 and used count is equal to frequency
	item4 := &UserFreeItem{
		Type: 2,
		Used: 5,
		Freq: 5,
	}
	expected4 := false
	if result4 := treeFreeUse(item4); result4 != expected4 {
		t.Errorf("Expected %v but got %v", expected4, result4)
	}

	// Testing for item type 2 and used count is greater than frequency
	item5 := &UserFreeItem{
		Type: 2,
		Used: 6,
		Freq: 5,
	}
	expected5 := false
	if result5 := treeFreeUse(item5); result5 != expected5 {
		t.Errorf("Expected %v but got %v", expected5, result5)
	}
}

func TestApplyFreeUse(t *testing.T) {
	defaultConfig := map[string]FreeItem{
		"feature1": {
			Type: 0,
			Freq: 10,
			Exp:  time.Now().Add(time.Hour).Unix(),
		},
		"feature2": {
			Type: 1,
			Freq: 0,
			Exp:  time.Now().Add(2 * time.Hour).Unix(),
		},
		"feature3_": {
			Type: 0,
			Freq: 5,
			Exp:  time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	userConfig := map[string]*UserFreeItem{
		"feature1": {
			Type: 0,
			Freq: 10,
			Used: 9,
			Exp:  time.Now().Add(time.Hour).Unix(),
		},
	}

	// Test for existing userConfig item
	err := ApplyFreeUse(userConfig, "feature1", defaultConfig)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	// Test for existing userConfig item with used count greater than frequency
	err = ApplyFreeUse(userConfig, "feature1", defaultConfig)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	// Test for non-existing userConfig item with default config
	err = ApplyFreeUse(userConfig, "feature2", defaultConfig)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	// Test for non-existing userConfig item with prefix match
	err = ApplyFreeUse(userConfig, "feature3_123456", defaultConfig)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}

	// Test for nil userConfig item
	err = ApplyFreeUse(userConfig, "feature4", defaultConfig)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}
