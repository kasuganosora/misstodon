package mfm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// helpers for building expected nodes
func T(s string) MfmNode  { return MfmNode{Type: nodeTypeText, Props: map[string]any{"text": s}} }
func B(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeBold, Children: c} }
func I(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeItalic, Children: c} }
func S(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeStrike, Children: c} }
func SM(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeSmall, Children: c} }
func Q(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeQuote, Children: c} }
func CTR(c ...MfmNode) MfmNode { return MfmNode{Type: nodeTypeCenter, Children: c} }
func PL(s string) MfmNode { return MfmNode{Type: nodeTypePlain, Children: []MfmNode{T(s)}} }
func IC(code string) MfmNode { return MfmNode{Type: nodeTypeInlineCode, Props: map[string]any{"code": code}} }
func CB(code, lang string) MfmNode {
	p := map[string]any{"code": code}
	if lang != "" { p["lang"] = lang }
	return MfmNode{Type: nodeTypeBlockCode, Props: p}
}
func MI(f string) MfmNode { return MfmNode{Type: nodeTypeMathInline, Props: map[string]any{"formula": f}} }
func MB(f string) MfmNode { return MfmNode{Type: nodeTypeMathBlock, Props: map[string]any{"formula": f}} }
func EC(name string) MfmNode { return MfmNode{Type: nodeTypeEmojiCode, Props: map[string]any{"name": name}} }
func UE(emoji string) MfmNode { return MfmNode{Type: nodeTypeUnicodeEmoji, Props: map[string]any{"emoji": emoji}} }
func HT(tag string) MfmNode { return MfmNode{Type: nodeTypeHashtag, Props: map[string]any{"hashtag": tag}} }
func URL(u string) MfmNode { return MfmNode{Type: nodeTypeUrl, Props: map[string]any{"url": u}} }
func LNK(silent bool, url string, c ...MfmNode) MfmNode {
	return MfmNode{Type: nodeTypeLink, Props: map[string]any{"url": url, "silent": silent}, Children: c}
}
func MEN(user string, host any, acct string) MfmNode {
	return MfmNode{Type: nodeTypeMention, Props: map[string]any{"username": user, "host": host, "acct": acct}}
}
func FN(name string, args map[string]any, c ...MfmNode) MfmNode {
	return MfmNode{Type: nodeTypeFn, Props: map[string]any{"name": name, "args": args}, Children: c}
}
func SCH(query, content string) MfmNode {
	return MfmNode{Type: nodeTypeSearch, Props: map[string]any{"query": query, "content": content}}
}

func TestParserText(t *testing.T) {
	assert.Equal(t, []MfmNode{T("abc")}, parse("abc"))
	assert.Equal(t, []MfmNode{T("abc#abc")}, parse("abc#abc"), "ignore hashtag after alphanumeric")
}

func TestParserEmoji(t *testing.T) {
	assert.Equal(t, []MfmNode{EC("foo")}, parse(":foo:"))
	assert.Equal(t, []MfmNode{T("foo:bar:baz")}, parse("foo:bar:baz"), "between alphanumeric texts")
	assert.Equal(t, []MfmNode{T("12:34:56")}, parse("12:34:56"), "between numbers")
	assert.Equal(t, []MfmNode{T("あ"), EC("bar"), T("い")}, parse("あ:bar:い"), "between non-ascii")
}

func TestParserUnicodeEmoji(t *testing.T) {
	assert.Equal(t, []MfmNode{T("今起きた"), UE("😇")}, parse("今起きた😇"))
}

func TestParserBig(t *testing.T) {
	assert.Equal(t, []MfmNode{FN("tada", map[string]any{}, T("abc"))}, parse("***abc***"))
	assert.Equal(t, []MfmNode{FN("tada", map[string]any{}, T("123"), B(T("abc")), T("123"))}, parse("***123**abc**123***"), "inline syntax inside")
}

func TestParserBold(t *testing.T) {
	assert.Equal(t, []MfmNode{B(T("abc"))}, parse("**abc**"))
	assert.Equal(t, []MfmNode{B(T("123"), S(T("abc")), T("123"))}, parse("**123~~abc~~123**"), "inline syntax")
}

func TestParserBoldTag(t *testing.T) {
	assert.Equal(t, []MfmNode{B(T("abc"))}, parse("<b>abc</b>"))
	assert.Equal(t, []MfmNode{B(T("123"), S(T("abc")), T("123"))}, parse("<b>123~~abc~~123</b>"), "inline syntax")
}

func TestParserItalic(t *testing.T) {
	assert.Equal(t, []MfmNode{I(T("abc"))}, parse("<i>abc</i>"))
	assert.Equal(t, []MfmNode{I(T("abc"))}, parse("*abc*"))
	assert.Equal(t, []MfmNode{I(T("abc"))}, parse("_abc_"))
	assert.Equal(t, []MfmNode{T("before "), I(T("abc")), T(" after")}, parse("before *abc* after"))
	assert.Equal(t, []MfmNode{T("before*abc*after")}, parse("before*abc*after"), "ignore if before char is alphanumeric")
	assert.Equal(t, []MfmNode{T("before_abc_after")}, parse("before_abc_after"), "underscore same rule")
}

func TestParserSmall(t *testing.T) {
	assert.Equal(t, []MfmNode{SM(T("abc"))}, parse("<small>abc</small>"))
	assert.Equal(t, []MfmNode{SM(T("abc"), B(T("123")), T("abc"))}, parse("<small>abc**123**abc</small>"), "inline syntax")
}

func TestParserStrike(t *testing.T) {
	assert.Equal(t, []MfmNode{S(T("foo"))}, parse("~~foo~~"))
	assert.Equal(t, []MfmNode{S(T("foo"))}, parse("<s>foo</s>"))
}

func TestParserInlineCode(t *testing.T) {
	assert.Equal(t, []MfmNode{IC(`var x = "Strawberry Pasta";`)}, parse("`var x = \"Strawberry Pasta\";`"))
	assert.Equal(t, []MfmNode{T("`foo\nbar`")}, parse("`foo\nbar`"), "disallow line break")
}

func TestParserMathInline(t *testing.T) {
	assert.Equal(t, []MfmNode{MI(`x = {-b \pm \sqrt{b^2-4ac} \over 2a}`)}, parse(`\(x = {-b \pm \sqrt{b^2-4ac} \over 2a}\)`))
}

func TestParserMathBlock(t *testing.T) {
	assert.Equal(t, []MfmNode{MB("math1")}, parse(`\[math1\]`))
	assert.Equal(t, []MfmNode{MB("math1")}, parse("\\[\nmath1\n\\]"), "multiline")
}

func TestParserQuote(t *testing.T) {
	assert.Equal(t, []MfmNode{Q(T("abc"))}, parse("> abc"))
	assert.Equal(t, []MfmNode{Q(T("abc\n123"))}, parse("> abc\n> 123"), "multiline")
}

func TestParserSearch(t *testing.T) {
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 Search")}, parse("MFM 書き方 123 Search"))
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 [Search]")}, parse("MFM 書き方 123 [Search]"))
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 search")}, parse("MFM 書き方 123 search"))
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 [search]")}, parse("MFM 書き方 123 [search]"))
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 検索")}, parse("MFM 書き方 123 検索"))
	assert.Equal(t, []MfmNode{SCH("MFM 書き方 123", "MFM 書き方 123 [検索]")}, parse("MFM 書き方 123 [検索]"))
}

func TestParserCodeBlock(t *testing.T) {
	assert.Equal(t, []MfmNode{CB("abc", "")}, parse("```\nabc\n```"))
	assert.Equal(t, []MfmNode{CB("a\nb\nc", "")}, parse("```\na\nb\nc\n```"), "multiline")
	assert.Equal(t, []MfmNode{CB("const a = 1;", "js")}, parse("```js\nconst a = 1;\n```"), "with lang")
	assert.Equal(t, []MfmNode{CB("aaa```bbb", "")}, parse("```\naaa```bbb\n```"), "ignore internal marker")
}

func TestParserCenter(t *testing.T) {
	assert.Equal(t, []MfmNode{CTR(T("abc"))}, parse("<center>abc</center>"))
}

func TestParserEmojiCode(t *testing.T) {
	assert.Equal(t, []MfmNode{EC("abc")}, parse(":abc:"))
}

func TestParserMention(t *testing.T) {
	assert.Equal(t, []MfmNode{MEN("abc", nil, "@abc")}, parse("@abc"))
	assert.Equal(t, []MfmNode{T("before "), MEN("abc", nil, "@abc"), T(" after")}, parse("before @abc after"))
	assert.Equal(t, []MfmNode{MEN("abc", "misskey.io", "@abc@misskey.io")}, parse("@abc@misskey.io"), "remote")
	assert.Equal(t, []MfmNode{T("abc@example.com")}, parse("abc@example.com"), "ignore mail address")
	assert.Equal(t, []MfmNode{MEN("abc-d", nil, "@abc-d")}, parse("@abc-d"), "allow hyphen in username")
	assert.Equal(t, []MfmNode{T("@-abc")}, parse("@-abc"), "disallow hyphen at start")
	assert.Equal(t, []MfmNode{MEN("abc", nil, "@abc"), T("-")}, parse("@abc-"), "strip trailing hyphen")
}

func TestParserHashtag(t *testing.T) {
	assert.Equal(t, []MfmNode{HT("abc")}, parse("#abc"))
	assert.Equal(t, []MfmNode{T("before "), HT("abc"), T(" after")}, parse("before #abc after"))
	assert.Equal(t, []MfmNode{T("abc#abc")}, parse("abc#abc"), "ignore after alphanumeric")
	assert.Equal(t, []MfmNode{HT("foo123")}, parse("#foo123"), "allow numbers")
	assert.Equal(t, []MfmNode{T("#123")}, parse("#123"), "disallow number-only")
	assert.Equal(t, []MfmNode{T("("), HT("foo"), T(")")}, parse("(#foo)"), "brackets")
	assert.Equal(t, []MfmNode{T("「"), HT("foo"), T("」")}, parse("「#foo」"), "JP brackets")
}

func TestParserUrl(t *testing.T) {
	assert.Equal(t, []MfmNode{URL("https://misskey.io/@ai")}, parse("https://misskey.io/@ai"))
	assert.Equal(t, []MfmNode{URL("https://misskey.io/@ai"), T(".")}, parse("https://misskey.io/@ai."), "strip trailing period")
	assert.Equal(t, []MfmNode{URL("https://example.com/foo?bar=a,b")}, parse("https://example.com/foo?bar=a,b"), "with comma in url")
	assert.Equal(t, []MfmNode{URL("https://example.com/foo"), T(", bar")}, parse("https://example.com/foo, bar"), "strip trailing comma")
	assert.Equal(t, []MfmNode{URL("https://example.com/foo(bar)")}, parse("https://example.com/foo(bar)"), "with brackets")
	assert.Equal(t, []MfmNode{T("javascript:foo")}, parse("javascript:foo"), "prevent xss")
}

func TestParserUrlAlt(t *testing.T) {
	assert.Equal(t, []MfmNode{URL("https://misskey.io/@ai")}, parse("<https://misskey.io/@ai>"))
	assert.Equal(t, []MfmNode{URL("http://藍.jp/abc")}, parse("<http://藍.jp/abc>"), "non-ascii with brackets")
}

func TestParserLink(t *testing.T) {
	assert.Equal(t, []MfmNode{LNK(false, "https://misskey.io/@ai", T("official instance")), T(".")}, parse("[official instance](https://misskey.io/@ai)."))
	assert.Equal(t, []MfmNode{LNK(true, "https://misskey.io/@ai", T("official instance")), T(".")}, parse("?[official instance](https://misskey.io/@ai)."), "silent")
}

func TestParserFn(t *testing.T) {
	assert.Equal(t, []MfmNode{FN("tada", map[string]any{}, T("abc"))}, parse("$[tada abc]"))
	assert.Equal(t, []MfmNode{FN("spin", map[string]any{"speed": "1.1s"}, T("a"))}, parse("$[spin.speed=1.1s a]"), "string arg")
	assert.Equal(t, []MfmNode{FN("spin", map[string]any{"speed": "1.1s"}, FN("shake", map[string]any{}, T("a")))}, parse("$[spin.speed=1.1s $[shake a]]"), "nested")
}

func TestParserPlain(t *testing.T) {
	assert.Equal(t, []MfmNode{T("a\n"), PL("**Hello**\nworld"), T("\nb")}, parse("a\n<plain>\n**Hello**\nworld\n</plain>\nb"), "multiline")
	assert.Equal(t, []MfmNode{T("a\n"), PL("**Hello** world"), T("\nb")}, parse("a\n<plain>**Hello** world</plain>\nb"), "single line")
}
