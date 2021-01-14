package wordscount

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"mvdan.cc/xurls/v2"
)

type Counter struct {
	Total     int // 总字数 = Words + Puncts
	Words     int // 只包含字符数
	Puncts    int // 标点数
	Links     int // 链接数
	Pics      int // 图片数
	CodeLines int // 代码行数
}

func (wc *Counter) Stat(str string) {
	wc.Links = len(rxStrict.FindAllString(str, -1))
	wc.Pics = len(imgReg.FindAllString(str, -1))

	// 剔除 HTML
	str = StripHTML(str)

	str = AutoSpace(str)

	// 普通的链接去除（非 HTML 标签链接）
	str = rxStrict.ReplaceAllString(str, " ")
	plainWords := strings.Fields(str)

	for _, plainWord := range plainWords {
		words := strings.FieldsFunc(plainWord, func(r rune) bool {
			if unicode.IsPunct(r) {
				wc.Puncts++
				return true
			}
			return false
		})

		for _, word := range words {
			runeCount := utf8.RuneCountInString(word)
			if len(word) == runeCount {
				wc.Words++
			} else {
				wc.Words += runeCount
			}
		}
	}

	wc.Total = wc.Words + wc.Puncts
}

// AutoSpace 自动给中英文之间加上空格
func AutoSpace(str string) string {
	out := ""

	for _, r := range str {
		out = addSpaceAtBoundary(out, r)
	}

	return out
}

func addSpaceAtBoundary(prefix string, nextChar rune) string {
	if len(prefix) == 0 {
		return string(nextChar)
	}

	r, size := utf8.DecodeLastRuneInString(prefix)
	if isLatin(size) != isLatin(utf8.RuneLen(nextChar)) &&
		isAllowSpace(nextChar) && isAllowSpace(r) {
		return prefix + " " + string(nextChar)
	}

	return prefix + string(nextChar)
}

var (
	rxStrict          = xurls.Strict()
	imgReg            = regexp.MustCompile(`<img [^>]*>`)
	stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")
)

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {
	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)

	// Walk through the string removing all tags
	b := GetBuffer()
	defer PutBuffer(b)
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}

		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}

		wasSpace = isSpace

	}
	return b.String()
}

func isLatin(size int) bool {
	return size == 1
}

func isAllowSpace(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsPunct(r)
}
