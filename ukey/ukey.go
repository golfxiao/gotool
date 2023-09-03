package ukey

import (
	"strings"
)

const (
	MinVersion = 0x00
	MetaLength = 3
)

// 短链ukey的构造-解析器，工作流程如下：
// 1. 生成ukey时，read input -> buildMiddle -> buildMeta -> buildUkey -> return output.ukey
// 2. 解析ukey时，read output -> parseMeta -> parseMiddle -> parseInput -> return input
type ukeyBuilder struct {
	// 元信息，包括版本、各数据信息的长度
	meta struct {
		version int // ukey的版本号
		len1    int // segment1的62进制字符长度
		len2    int // segment2的62进制字符长度
		len3    int // segment3的62进制字符长度
	}

	// 业务输入信息
	input struct {
		segment1  int // 数据信息-1
		segment2  int // 数据信息-2
		segment3  int // 数据信息-3
		randomLen int // 随机数长度
	}

	// 中间过程中生成的信息
	middle struct {
		segment1Chars string // 数据信息-1对应的62进制字符
		segment2Chars string // 数据信息-2对应的62进制字符
		segment3Chars string // 数据信息-3对应的62进制字符
		randomChars   string // 随机字符
		metaChars     string // 元信息字符
	}

	// 最终输出的ukey信息
	output struct {
		ukey string // ukey串
	}
}

// 构建ukey，支持传多个正整数形式的数据信息，数据信息不要大于毫秒时间戳，请确保至少有一个数据信息有效
// @param data1 数据信息1，小于等于0时会忽略此数据项
// @param data2 数据信息2，小于等于0时会忽略此数据项
// @param data3 数据信息3，小于等于0时会忽略此数据项
// @param randomLen 生成随机字符的位数，小于等于0时不生成
// @return string 生成的ukey串
func BuildUkey(data1, data2, data3, randomLen int) string {
	if data1 <= 0 && data2 <= 0 && data3 <= 0 {
		return ""
	}
	builder := &ukeyBuilder{}
	builder.meta.version = MinVersion
	builder.input.segment1 = data1
	builder.input.segment2 = data2
	builder.input.segment3 = data3
	builder.input.randomLen = randomLen

	return builder.buildMiddle().buildMeta().buildUkey().output.ukey
}

// 解析ukey获取业务信息
// @param ukey 短链接上的ukey串
// @return data1 数据信息1，未找到时返回0
// @return data2 数据信息1，未找到时返回0
// @return data3 数据信息1，未找到时返回0
// @return randomLen 随机字符长度，未找到时返回0
// @note 兼容落地页老版本的ukey，解析到的数据会在data1中返回
func ParseUkey(ukey string) (data1, data2, data3, randomLen int) {
	if len(ukey) <= MetaLength {
		return 0, 0, 0, 0
	}

	parser := &ukeyBuilder{}
	parser.output.ukey = ukey

	// 先解析meta，判断如果是老版本ukey，则走老的解析方法
	if parser.parseMeta(); parser.isOldVersion() {
		// 如果有多个版本，这里可以分版本处理兼容
	}

	// 非法ukey检测（主要判断长度）
	// 背景：有QA手输链接，输错一个字符，非法ukey导致数组越界
	if parser.isInvalidUkey() {
		return 0, 0, 0, 0
	}

	// 新的ukey走新的解析方法
	parser.parseMiddle().parseInput()
	return parser.input.segment1, parser.input.segment2, parser.input.segment3, parser.input.randomLen
}

func (this *ukeyBuilder) buildMiddle() *ukeyBuilder {
	if this.input.segment1 > 0 {
		this.middle.segment1Chars = from10To62(this.input.segment1, 62)
	}
	if this.input.segment2 > 0 {
		this.middle.segment2Chars = from10To62(this.input.segment2, 62)
	}
	if this.input.segment3 > 0 {
		this.middle.segment3Chars = from10To62(this.input.segment3, 62)
	}
	if this.input.randomLen > 0 {
		this.middle.randomChars = getRandomString(this.input.randomLen)
	}
	return this
}

func (this *ukeyBuilder) buildMeta() *ukeyBuilder {
	this.meta.len1 = len(this.middle.segment1Chars)
	this.meta.len2 = len(this.middle.segment2Chars)
	this.meta.len3 = len(this.middle.segment3Chars)

	metaInt := this.meta.version<<12 | this.meta.len1<<8 | this.meta.len2<<4 | this.meta.len3
	metaChars := from10To62(metaInt, 62)
	if len(metaChars) < MetaLength {
		metaChars = strings.Repeat("0", MetaLength-len(metaChars)) + metaChars
	}
	this.middle.metaChars = strrev(metaChars)
	return this
}

func (this *ukeyBuilder) buildUkey() *ukeyBuilder {
	this.output.ukey = this.middle.segment1Chars +
		this.middle.segment2Chars +
		this.middle.segment3Chars +
		this.middle.randomChars +
		this.middle.metaChars
	return this
}

func (this *ukeyBuilder) parseMeta() *ukeyBuilder {
	this.middle.metaChars = this.output.ukey[len(this.output.ukey)-MetaLength:]
	metaInt := from62To10(strrev(this.middle.metaChars), 62)
	this.meta.version = (metaInt >> 12) & 0x1F
	this.meta.len1 = (metaInt >> 8) & 0x0F
	this.meta.len2 = (metaInt >> 4) & 0x0F
	this.meta.len3 = metaInt & 0x0F
	return this
}

func (this *ukeyBuilder) parseMiddle() *ukeyBuilder {
	this.middle.segment1Chars = this.output.ukey[:this.meta.len1]
	this.middle.segment2Chars = this.output.ukey[this.meta.len1 : this.meta.len1+this.meta.len2]
	this.middle.segment3Chars = this.output.ukey[this.meta.len1+this.meta.len2 : this.meta.len1+this.meta.len2+this.meta.len3]
	this.middle.randomChars = this.output.ukey[this.meta.len1+this.meta.len2+this.meta.len3 : len(this.output.ukey)-MetaLength]
	return this
}

func (this *ukeyBuilder) parseInput() *ukeyBuilder {
	this.input.segment1 = from62To10(this.middle.segment1Chars, 62)
	this.input.segment2 = from62To10(this.middle.segment2Chars, 62)
	this.input.segment3 = from62To10(this.middle.segment3Chars, 62)
	this.input.randomLen = len(this.middle.randomChars)
	return this
}

func (this *ukeyBuilder) isOldVersion() bool {
	return false
}

// 用于解析场景的保护：如果前三段加meta的长度超过ukey的总长度，则认为非法ukey
func (this *ukeyBuilder) isInvalidUkey() bool {
	return this.meta.len1+this.meta.len2+this.meta.len3+MetaLength > len(this.output.ukey)
}

func GetUKeyByLink(link string) string {
	path, _, _, _ := parseUrl(link)
	pathArr := strings.Split(path, "/")
	uKey := pathArr[len(pathArr)-1]
	return uKey
}
