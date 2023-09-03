package ukey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUkey(t *testing.T) {
	// 不重复
	assert.NotEqual(t, BuildUkey(45454545, 0, 0, 0), BuildUkey(45454545, 1, 0, 0))
	assert.NotEqual(t, BuildUkey(45454545, 333, 0, 0), BuildUkey(45454545, 332, 0, 0))
	assert.NotEqual(t, BuildUkey(45454545, 333, 2, 0), BuildUkey(45454545, 333, 1, 0))
	assert.NotEqual(t, BuildUkey(45454545, 333, 444, 1), BuildUkey(45454545, 333, 444, 2))

	// 两次生成的是否相同
	assert.Equal(t, BuildUkey(45454545, 0, 0, 0), BuildUkey(45454545, 0, 0, 0))
	assert.Equal(t, BuildUkey(45454545, 5, 0, 0), BuildUkey(45454545, 5, 0, 0))
	assert.Equal(t, BuildUkey(45454545, 5, 100, 0), BuildUkey(45454545, 5, 100, 0))

	// 非法数据
	assert.Equal(t, "", BuildUkey(0, 0, 0, 0))
	assert.Equal(t, "", BuildUkey(-1, -3223222, -3292323832823823212, -10))
	assert.Empty(t, BuildUkey(-1, -10, -12, 10))
	assert.NotEmpty(t, BuildUkey(-1, 10, -3292323832823823212, -10)) // TODO
	assert.NotEmpty(t, BuildUkey(1, 0, 3939, -10))
	assert.NotEmpty(t, BuildUkey(-1, -10, 12, -10))

	t.Log("buildUkey: ", BuildUkey(45454545, 0, 0, 0), BuildUkey(45454545, 333, 2, 0), BuildUkey(45454545, 333, 444, 1))
}

func TestParseUkey(t *testing.T) {
	ukey := BuildUkey(45454545, 0, 0, 0)
	data1, data2, data3, randomLen := ParseUkey(ukey)
	assert.Equal(t, 45454545, data1)
	assert.Equal(t, 0, data2)
	assert.Equal(t, 0, data3)
	assert.Equal(t, 0, randomLen)

	ukey = BuildUkey(1234567890, 45678, 0, 0)
	data1, data2, data3, randomLen = ParseUkey(ukey)
	assert.Equal(t, 1234567890, data1)
	assert.Equal(t, 45678, data2)
	assert.Equal(t, 0, data3)
	assert.Equal(t, 0, randomLen)

	ukey = BuildUkey(9876543210, 45678, 1649769576000, 0)
	data1, data2, data3, randomLen = ParseUkey(ukey)
	assert.Equal(t, 9876543210, data1)
	assert.Equal(t, 45678, data2)
	assert.Equal(t, 1649769576000, data3)
	assert.Equal(t, 0, randomLen)

	ukey = BuildUkey(38546, 2, 1649769576, 3)
	data1, data2, data3, randomLen = ParseUkey(ukey)
	assert.Equal(t, 38546, data1)
	assert.Equal(t, 2, data2)
	assert.Equal(t, 1649769576, data3)
	assert.Equal(t, 3, randomLen) // TODO
}

func TestExceptionCase(t *testing.T) {
	// 异常Case测试
	data1, data2, data3, randomLen := ParseUkey("gdwPAGI")
	assert.Equal(t, 0, data1)
	assert.Equal(t, 0, data2)
	assert.Equal(t, 0, data3)
	assert.Equal(t, 0, randomLen)
}
