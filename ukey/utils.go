package ukey

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var encodeTable = map[int]string{0: "0", 1: "a", 2: "2", 3: "3", 4: "4", 5: "b", 6: "6", 7: "7", 8: "8", 9: "9", 10: "1", 11: "5", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: "A", 37: "B", 38: "C", 39: "D", 40: "E", 41: "F", 42: "G", 43: "H", 44: "I", 45: "J", 46: "K", 47: "L", 48: "M", 49: "N", 50: "O", 51: "P", 52: "Q", 53: "R", 54: "S", 55: "T", 56: "U", 57: "V", 58: "W", 59: "X", 60: "Y", 61: "Z", 62: ":", 63: ";", 64: "<", 65: "=", 66: ">", 67: "?", 68: "@", 69: "[", 70: "]", 71: "^", 72: "_", 73: "{", 74: "|", 75: "}"}

// 10进制转62进制
func from10To62(num, n int) string {
	if n > 76 {
		n = 76
	}
	new_num_str := ""
	var remainder int
	var remainder_string string
	for num != 0 {
		remainder = num % n
		remainder_string = encodeTable[remainder]
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	return new_num_str
}

func from62To10(num string, n int) int {
	var new_num float64
	new_num = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(findkey(value))
		if tmp != -1 {
			new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(new_num)
}

func strrev(src string) string {
	s := []byte(src)
	l := len(s)
	d := make([]byte, l)
	for i := 0; i < l; i++ {
		d[l-i-1] = s[i]
	}
	return string(d)
}

func findkey(in string) int {
	result := -1
	for k, v := range encodeTable {
		if in == v {
			result = k
		}
	}
	return result
}

// 生成随机字符串
func getRandomString(lang int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lang; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func parseUrl(strUrl string) (path string, host string, port int, err error) {
	if strUrl == "" {
		err = fmt.Errorf("Invalid params of strUrl empty")
		return
	}
	u, err := url.Parse(strUrl)
	if err != nil {
		return
	}

	path = u.Path
	h := strings.Split(u.Host, ":")
	host = h[0]
	if len(h) >= 2 {
		port, _ = strconv.Atoi(h[1])
	}
	return
}
