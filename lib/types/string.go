package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	. "github.com/cosiner/gohper/lib/errors"
)

// IsLower check letter is lower case or not
func IsLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

// IsUpper check letter is upper case or not
func IsUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// IsLetter check character is a letter or not
func IsLetter(b byte) bool {
	return IsLower(b) || IsUpper(b)
}

// IsSpaceQuote return wehter a byte is space or quote characters
func IsSpaceQuote(b byte) bool {
	return IsSpace(b) || b == '"' || b == '\''
}

// IsSpace only call unicode.IsSpace
func IsSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

// LowerCase convert a byte to lower case
func LowerCase(b byte) byte {
	if IsUpper(b) {
		b = b - 'A' + 'a'
	}
	return b
}

// UpperCase convert a byte to upper
func UpperCase(b byte) byte {
	if IsLower(b) {
		b = b - 'a' + 'A'
	}
	return b
}

// TrimQuote trim quote for string, if quote don't match, return an error
func TrimQuote(str string) (s string, err error) {
	s = TrimSpace(str)
	if l := len(s); l > 0 {
		c := s[0]
		if c == '\'' || c == '"' || c == '`' {
			if s[l-1] == c {
				s = s[1 : l-1]
			} else {
				err = Errorf("Quote don't match:%s", s)
			}
		}
	}
	return
}

// TrimSpace is only call strings.TrimSpace
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// TrimUpper return the trim and upper format of a string
func TrimUpper(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

// TrimLower return the trim and lower format of a string
func TrimLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// BytesTrim2Str trim bytes return as string
func TrimBytes2Str(s []byte) string {
	return string(bytes.TrimSpace(s))
}

// TrimSplit split string and return trim space string
func TrimSplit(s, sep string) []string {
	sp := strings.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = strings.TrimSpace(sp[i])
	}
	return sp
}

// TrimBytesSplit split string and return trim space string
func TrimBytesSplit(s, sep []byte) [][]byte {
	sp := bytes.Split(s, sep)
	for i, n := 0, len(sp); i < n; i++ {
		sp[i] = bytes.TrimSpace(sp[i])
	}
	return sp
}

// StartWith check whether str is start with another string
// it's a wrapper of strings.HasPrefix
func StartWith(str, start string) bool {
	return strings.HasPrefix(str, start)
}

// EndWith check whether str is end with another string
// it's a wrapper of strings.HasSuffix
func EndWith(str, end string) bool {
	return strings.HasSuffix(str, end)
}

// RepeatJoin repeat s count times as a string slice, then join with sep
func RepeatJoin(s, sep string, count int) string {
	switch {
	case count <= 0:
		return ""
	case count == 1:
		return s
	case count == 2:
		return s + sep + s
	default:
		bs := make([]byte, 0, (len(s)+len(sep))*count-len(sep))
		buf := bytes.NewBuffer(bs)
		buf.WriteString(s)
		for i := 1; i < count; i++ {
			buf.WriteString(sep)
			buf.WriteString(s)
		}
		return buf.String()
	}
}

// SuffixJoin join string slice with suffix
func SuffixJoin(s []string, suffix, sep string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0] + suffix
	}
	n := len(sep) * (len(s) - 1)
	for i, sl := 0, len(suffix); i < len(s); i++ {
		n += len(s[i]) + sl
	}
	b := make([]byte, n)
	bp := copy(b, s[0])
	bp += copy(b[bp:], suffix)
	for _, s := range s[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
		bp += copy(b[bp:], suffix)
	}
	return string(b)
}

// JoinInt join int slice as string
func JoinInt(v []int, sep string) (str string) {
	if len(v) > 0 {
		buf := bytes.NewBuffer([]byte{})
		buf.WriteString(strconv.Itoa(v[0]))
		for _, s := range v[1:] {
			buf.WriteString(fmt.Sprintf("%s%d", sep, s))
		}
		str = buf.String()
	}
	return
}

// StringReader is a wrapper of strings.NewReader
func StringReader(s string) *strings.Reader {
	return strings.NewReader(s)
}

// TrimBefore trim string and remove the section before delimiter and delimiter itself
func TrimBefore(s string, delimiter string) string {
	if idx := strings.Index(s, delimiter); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

// TrimBytesBefore trim bytes and remove the section before delimiter and delimiter itself
func TrimBytesBefore(s []byte, delimiter []byte) []byte {
	if idx := bytes.Index(s, delimiter); idx >= 0 {
		s = s[:idx]
	}
	return bytes.TrimSpace(s)
}

// TrimAfter trim string and remove the section after delimiter and delimiter itself
func TrimAfter(s string, delimiter string) (ret string) {
	if idx := strings.Index(s, delimiter); idx >= 0 {
		ret = TrimSpace(s[idx:])
	}
	return
}

// TrimBytesAfter trim bytes and remove the section after delimiter and delimiter itself
func TrimBytesAfter(s []byte, delimiter []byte) (ret []byte) {
	if idx := bytes.Index(s, delimiter); idx >= 0 {
		ret = bytes.TrimSpace(s[idx:])
	}
	return
}

// StrIndexN find index of n-th sep string
func StrIndexN(str, sep string, n int) (index int) {
	index, idx, seplen := 0, -1, len(sep)
	for i := 0; i < n; i++ {
		if idx = strings.Index(str, sep); idx == -1 {
			break
		}
		str = str[idx+seplen:]
		index += idx
	}
	if idx == -1 {
		index = -1
	} else {
		index += (n - 1) * seplen
	}
	return
}

// StrLastIndexN find last index of n-th sep string
func StrLastIndexN(str, sep string, n int) (index int) {
	for i := 0; i < n; i++ {
		if index = strings.LastIndex(str, sep); index == -1 {
			break
		}
		str = str[:index]
	}
	return
}

// snake string, XxYy to xx_yy, X_Y to x_y
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	num := len(s)
	need := false // need determin if it's necessery to add a '_'
	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
			if need {
				data = append(data, '_')
			}
		} else {
			// if previous is '_' or ' ',
			// there is no need to add extra '_' before
			need = (c != '_' && c != ' ')
		}
		data = append(data, c)
	}
	return string(data)
}

// camel string, xx_yy to XxYy, xx__yy to Xx_Yy
// xx _yy to Xx Yy, the rule is that a lower case letter
// after '_' will combine to a upper case letter
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	num := len(s)
	need := true
	var prev byte = ' '

	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			if need {
				c = c - 'a' + 'A'
				need = false
			}
		} else {
			if prev == '_' {
				data = append(data, '_')
			}
			need = (c == '_' || c == ' ')
			if c == '_' {
				prev = '_'
				continue
			}
		}
		prev = c
		data = append(data, c)
	}
	return string(data)
}

// AbridgeString extract first letter and all upper case letter
// from string as it's abridge case
func AbridgeString(str string) (s string) {
	if l := len(str); l != 0 {
		ret := []byte{str[0]}
		for i := 1; i < l; i++ {
			b := str[i]
			if IsUpper(b) {
				ret = append(ret, b)
			}
		}
		s = string(ret)
	}
	return
}

// AbridgeStringToLower extract first letter and all upper case letter
// from string as it's abridge case, and convert it to lower case
func AbridgeStringToLower(str string) (s string) {
	if l := len(str); l != 0 {
		ret := []byte{LowerCase(str[0])}
		for i := 1; i < l; i++ {
			b := str[i]
			if IsUpper(b) {
				ret = append(ret, LowerCase(b))
			}
		}
		s = string(ret)
	}
	return
}

// Compare compare two string, if equal, 0 was returned, if s1 > s2, 1 was returned,
// otherwise -1 was returned
func Compare(s1, s2 string) int {
	l1, l2 := len(s1), len(s2)
	for i := 0; i < l1 && i < l2; i++ {
		if s1[i] < s2[i] {
			return -1
		} else if s1[i] > s2[i] {
			return 1
		}
	}
	switch {
	case l1 < l2:
		return -1
	case l1 == l2:
		return 0
	default:
		return 1
	}
}

// AllCharIn check whether all chars of string is in encoding string
func AllCharsIn(s, encoding string) bool {
	for i := 0; i < len(s); i++ {
		if CharIn(s[i], encoding) < 0 {
			return false
		}
	}
	return true
}
