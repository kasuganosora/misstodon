package mfm

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// parse is the entry point - parses MFM text into a tree of MfmNodes.
func parse(input string) []MfmNode {
	s := &parserState{input: input, nestLimit: 20}
	nodes := s.parseFull()
	return mergeText(nodes)
}

// --- Domain: Parser State (Value Object) ---

type parserState struct {
	input     string
	pos       int
	depth     int
	nestLimit int
	linkLabel bool // inside [label](url), disables mention/hashtag/url
}

func (s *parserState) remaining() string { return s.input[s.pos:] }
func (s *parserState) eof() bool         { return s.pos >= len(s.input) }
func (s *parserState) char() byte        { return s.input[s.pos] }

func (s *parserState) consume(prefix string) bool {
	if strings.HasPrefix(s.input[s.pos:], prefix) {
		s.pos += len(prefix)
		return true
	}
	return false
}

func (s *parserState) atLineBegin() bool {
	return s.pos == 0 || (s.pos > 0 && s.input[s.pos-1] == '\n')
}

// --- Service: Parse Dispatchers ---

// parseFull tries all syntax elements (block + inline).
func (s *parserState) parseFull() []MfmNode {
	var nodes []MfmNode
	for !s.eof() {
		node, ok := s.tryBlock()
		if !ok {
			node, ok = s.tryInline()
		}
		if !ok {
			// fallback: consume one rune as text
			_, sz := utf8.DecodeRuneInString(s.remaining())
			nodes = append(nodes, textNode(s.input[s.pos:s.pos+sz]))
			s.pos += sz
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// parseInline tries only inline syntax elements.
func (s *parserState) parseInline() []MfmNode {
	var nodes []MfmNode
	for !s.eof() {
		node, ok := s.tryInline()
		if !ok {
			_, sz := utf8.DecodeRuneInString(s.remaining())
			nodes = append(nodes, textNode(s.input[s.pos:s.pos+sz]))
			s.pos += sz
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// parseInlineUntil parses inline content until `end` delimiter is found.
// Returns (children, true) if end found, or (nil, false) if not.
func (s *parserState) parseInlineUntil(end string) ([]MfmNode, bool) {
	if s.depth >= s.nestLimit {
		return nil, false
	}
	s.depth++
	defer func() { s.depth-- }()

	var nodes []MfmNode
	for !s.eof() {
		if strings.HasPrefix(s.remaining(), end) {
			s.pos += len(end)
			return mergeText(nodes), true
		}
		node, ok := s.tryInline()
		if !ok {
			_, sz := utf8.DecodeRuneInString(s.remaining())
			nodes = append(nodes, textNode(s.input[s.pos:s.pos+sz]))
			s.pos += sz
			continue
		}
		nodes = append(nodes, node)
	}
	return nil, false
}

// --- Service: Block-level Parsers ---

func (s *parserState) tryBlock() (MfmNode, bool) {
	if !s.atLineBegin() {
		return MfmNode{}, false
	}
	for _, fn := range []func() (MfmNode, bool){
		s.tryQuote,
		s.tryCodeBlock,
		s.tryMathBlock,
		s.tryCenterTag,
		s.trySearch,
	} {
		if node, ok := fn(); ok {
			return node, true
		}
	}
	return MfmNode{}, false
}

// --- Service: Inline-level Parsers ---

func (s *parserState) tryInline() (MfmNode, bool) {
	for _, fn := range []func() (MfmNode, bool){
		s.tryUnicodeEmoji,
		s.tryUrlAlt,
		s.trySmallTag,
		s.tryPlainTag,
		s.tryBoldTag,
		s.tryItalicTag,
		s.tryStrikeTag,
		s.tryBig,
		s.tryBoldAsta,
		s.tryItalicAsta,
		s.tryBoldUnder,
		s.tryItalicUnder,
		s.tryInlineCode,
		s.tryMathInline,
		s.tryStrikeWave,
		s.tryFn,
		s.tryMention,
		s.tryHashtag,
		s.tryEmojiCode,
		s.tryLink,
		s.tryUrl,
	} {
		if node, ok := fn(); ok {
			return node, true
		}
	}
	return MfmNode{}, false
}

// --- Leaf Parsers ---

func (s *parserState) tryInlineCode() (MfmNode, bool) {
	if !s.consume("`") {
		return MfmNode{}, false
	}
	start := s.pos
	for !s.eof() {
		ch := s.char()
		if ch == '`' {
			code := s.input[start:s.pos]
			s.pos++ // consume closing `
			if code == "" {
				break
			}
			return MfmNode{Type: nodeTypeInlineCode, Props: map[string]any{"code": code}}, true
		}
		if ch == '\n' {
			break
		}
		s.pos++
	}
	s.pos = start - 1 // restore
	return MfmNode{}, false
}

func (s *parserState) tryMathInline() (MfmNode, bool) {
	saved := s.pos
	if !s.consume(`\(`) {
		return MfmNode{}, false
	}
	start := s.pos
	idx := strings.Index(s.input[start:], `\)`)
	if idx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	formula := s.input[start : start+idx]
	s.pos = start + idx + 2
	return MfmNode{Type: nodeTypeMathInline, Props: map[string]any{"formula": formula}}, true
}

func (s *parserState) tryEmojiCode() (MfmNode, bool) {
	if s.pos > 0 {
		prev, _ := utf8.DecodeLastRuneInString(s.input[:s.pos])
		if isAlphanumeric(prev) {
			return MfmNode{}, false
		}
	}
	if !s.consume(":") {
		return MfmNode{}, false
	}
	start := s.pos
	for !s.eof() && s.char() != ':' && s.char() != '\n' && s.char() != ' ' {
		ch := s.char()
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' || ch == '+') {
			break
		}
		s.pos++
	}
	if s.eof() || s.char() != ':' || s.pos == start {
		s.pos = start - 1
		return MfmNode{}, false
	}
	name := s.input[start:s.pos]
	s.pos++ // consume closing :
	// Check that next char is not alphanumeric
	if !s.eof() {
		next, _ := utf8.DecodeRuneInString(s.remaining())
		if isAlphanumeric(next) {
			s.pos = start - 1
			return MfmNode{}, false
		}
	}
	return MfmNode{Type: nodeTypeEmojiCode, Props: map[string]any{"name": name}}, true
}

// --- Block Parsers ---

func (s *parserState) tryCodeBlock() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("```") {
		return MfmNode{}, false
	}
	// optional language tag
	langEnd := strings.IndexByte(s.input[s.pos:], '\n')
	if langEnd < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	lang := strings.TrimSpace(s.input[s.pos : s.pos+langEnd])
	s.pos += langEnd + 1 // skip newline

	// find closing ``` (exactly 3 backticks at line start, not more)
	searchFrom := s.pos
	closeIdx := -1
	for {
		idx := strings.Index(s.input[searchFrom:], "\n```")
		if idx < 0 {
			break
		}
		afterClose := searchFrom + idx + 4
		// Must NOT be followed by another backtick
		if afterClose < len(s.input) && s.input[afterClose] == '`' {
			searchFrom = afterClose
			continue
		}
		closeIdx = searchFrom + idx - s.pos
		break
	}
	if closeIdx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	code := s.input[s.pos : s.pos+closeIdx]
	s.pos += closeIdx + 4 // skip \n```
	// consume optional newline after closing
	if !s.eof() && s.char() == '\n' {
		s.pos++
	}
	props := map[string]any{"code": code}
	if lang != "" {
		props["lang"] = lang
	}
	return MfmNode{Type: nodeTypeBlockCode, Props: props}, true
}

func (s *parserState) tryMathBlock() (MfmNode, bool) {
	saved := s.pos
	if !s.consume(`\[`) {
		return MfmNode{}, false
	}
	// skip optional newline
	if !s.eof() && s.char() == '\n' {
		s.pos++
	}
	start := s.pos
	// find \]
	idx := strings.Index(s.input[start:], `\]`)
	if idx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	formula := s.input[start : start+idx]
	// strip trailing newline from formula
	formula = strings.TrimRight(formula, "\n")
	s.pos = start + idx + 2
	return MfmNode{Type: nodeTypeMathBlock, Props: map[string]any{"formula": formula}}, true
}

func (s *parserState) tryQuote() (MfmNode, bool) {
	if !strings.HasPrefix(s.remaining(), "> ") && !strings.HasPrefix(s.remaining(), ">") {
		return MfmNode{}, false
	}
	saved := s.pos
	var lines []string
	for !s.eof() && strings.HasPrefix(s.remaining(), ">") {
		s.pos++ // skip >
		if !s.eof() && s.char() == ' ' {
			s.pos++ // skip space after >
		}
		lineEnd := strings.IndexByte(s.input[s.pos:], '\n')
		if lineEnd < 0 {
			lines = append(lines, s.input[s.pos:])
			s.pos = len(s.input)
		} else {
			lines = append(lines, s.input[s.pos:s.pos+lineEnd])
			s.pos += lineEnd + 1
		}
	}
	if len(lines) == 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	inner := strings.Join(lines, "\n")
	// Recursively parse the inner content
	innerParser := &parserState{input: inner, nestLimit: s.nestLimit, depth: s.depth}
	children := mergeText(innerParser.parseFull())
	return MfmNode{Type: nodeTypeQuote, Children: children}, true
}

func (s *parserState) tryCenterTag() (MfmNode, bool) {
	return s.tryHtmlBlock("center", nodeTypeCenter)
}

func (s *parserState) tryHtmlBlock(tag string, nodeType mfmNodeType) (MfmNode, bool) {
	saved := s.pos
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	if !s.consume(open) {
		return MfmNode{}, false
	}
	// skip optional newline after open
	if !s.eof() && s.char() == '\n' {
		s.pos++
	}
	start := s.pos
	idx := strings.Index(s.input[start:], close)
	if idx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	inner := s.input[start : start+idx]
	// strip trailing newline from inner
	inner = strings.TrimRight(inner, "\n")
	s.pos = start + idx + len(close)
	// parse inner content inline
	innerParser := &parserState{input: inner, nestLimit: s.nestLimit, depth: s.depth + 1}
	children := mergeText(innerParser.parseInline())
	return MfmNode{Type: nodeType, Children: children}, true
}

func (s *parserState) trySearch() (MfmNode, bool) {
	if !s.atLineBegin() {
		return MfmNode{}, false
	}
	saved := s.pos
	// find end of line
	lineEnd := strings.IndexByte(s.input[s.pos:], '\n')
	var line string
	if lineEnd < 0 {
		line = s.input[s.pos:]
	} else {
		line = s.input[s.pos : s.pos+lineEnd]
	}
	// Check for search suffixes
	suffixes := []string{" Search", " search", " 検索", " [Search]", " [search]", " [検索]"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(line, suffix) {
			query := strings.TrimSuffix(line, suffix)
			if query == "" {
				continue
			}
			content := query + suffix
			if lineEnd < 0 {
				s.pos = len(s.input)
			} else {
				s.pos += lineEnd + 1
			}
			return MfmNode{Type: nodeTypeSearch, Props: map[string]any{
				"query":   query,
				"content": content,
			}}, true
		}
	}
	s.pos = saved
	return MfmNode{}, false
}

// --- HTML Tag Inline Parsers ---

func (s *parserState) trySmallTag() (MfmNode, bool) {
	return s.tryHtmlInline("small", nodeTypeSmall)
}

func (s *parserState) tryPlainTag() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("<plain>") {
		return MfmNode{}, false
	}
	idx := strings.Index(s.input[s.pos:], "</plain>")
	if idx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	inner := s.input[s.pos : s.pos+idx]
	s.pos += idx + 8
	// Strip surrounding newlines
	inner = strings.TrimPrefix(inner, "\n")
	inner = strings.TrimSuffix(inner, "\n")
	// plain content is literal text, no further parsing
	children := []MfmNode{textNode(inner)}
	return MfmNode{Type: nodeTypePlain, Children: children}, true
}

func (s *parserState) tryBoldTag() (MfmNode, bool) {
	return s.tryHtmlInline("b", nodeTypeBold)
}

func (s *parserState) tryItalicTag() (MfmNode, bool) {
	return s.tryHtmlInline("i", nodeTypeItalic)
}

func (s *parserState) tryStrikeTag() (MfmNode, bool) {
	return s.tryHtmlInline("s", nodeTypeStrike)
}

func (s *parserState) tryHtmlInline(tag string, nodeType mfmNodeType) (MfmNode, bool) {
	saved := s.pos
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	if !s.consume(open) {
		return MfmNode{}, false
	}
	children, ok := s.parseInlineUntil(close)
	if !ok {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{Type: nodeType, Children: children}, true
}

// --- Markdown-style Inline Parsers ---

func (s *parserState) tryBig() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("***") {
		return MfmNode{}, false
	}
	children, ok := s.parseInlineUntil("***")
	if !ok || len(children) == 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{
		Type:     nodeTypeFn,
		Props:    map[string]any{"name": "tada", "args": map[string]any{}},
		Children: children,
	}, true
}

func (s *parserState) tryBoldAsta() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("**") {
		return MfmNode{}, false
	}
	children, ok := s.parseInlineUntil("**")
	if !ok || len(children) == 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{Type: nodeTypeBold, Children: children}, true
}

func (s *parserState) tryItalicAsta() (MfmNode, bool) {
	// *text* - content must be [a-z0-9 \n], must not be preceded by alphanumeric
	if s.pos > 0 && isAlphanumeric(rune(s.input[s.pos-1])) {
		return MfmNode{}, false
	}
	saved := s.pos
	if !s.consume("*") {
		return MfmNode{}, false
	}
	start := s.pos
	for !s.eof() {
		ch := s.char()
		if ch == '*' {
			if s.pos == start {
				break
			}
			content := s.input[start:s.pos]
			s.pos++ // consume *
			children := []MfmNode{textNode(content)}
			return MfmNode{Type: nodeTypeItalic, Children: children}, true
		}
		if !isItalicContent(ch) {
			break
		}
		s.pos++
	}
	s.pos = saved
	return MfmNode{}, false
}

func (s *parserState) tryBoldUnder() (MfmNode, bool) {
	if s.pos > 0 && isAlphanumeric(rune(s.input[s.pos-1])) {
		return MfmNode{}, false
	}
	saved := s.pos
	if !s.consume("__") {
		return MfmNode{}, false
	}
	start := s.pos
	for !s.eof() {
		if strings.HasPrefix(s.remaining(), "__") {
			if s.pos == start {
				break
			}
			content := s.input[start:s.pos]
			s.pos += 2
			children := []MfmNode{textNode(content)}
			return MfmNode{Type: nodeTypeBold, Children: children}, true
		}
		ch := s.char()
		if !isItalicContent(ch) {
			break
		}
		s.pos++
	}
	s.pos = saved
	return MfmNode{}, false
}

func (s *parserState) tryItalicUnder() (MfmNode, bool) {
	if s.pos > 0 && isAlphanumeric(rune(s.input[s.pos-1])) {
		return MfmNode{}, false
	}
	saved := s.pos
	if !s.consume("_") {
		return MfmNode{}, false
	}
	start := s.pos
	for !s.eof() {
		ch := s.char()
		if ch == '_' {
			if s.pos == start {
				break
			}
			content := s.input[start:s.pos]
			s.pos++
			children := []MfmNode{textNode(content)}
			return MfmNode{Type: nodeTypeItalic, Children: children}, true
		}
		if !isItalicContent(ch) {
			break
		}
		s.pos++
	}
	s.pos = saved
	return MfmNode{}, false
}

func (s *parserState) tryStrikeWave() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("~~") {
		return MfmNode{}, false
	}
	children, ok := s.parseInlineUntil("~~")
	if !ok || len(children) == 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{Type: nodeTypeStrike, Children: children}, true
}

// --- Complex Parsers ---

func (s *parserState) tryFn() (MfmNode, bool) {
	saved := s.pos
	if !s.consume("$[") {
		return MfmNode{}, false
	}
	// Parse function name
	nameStart := s.pos
	for !s.eof() && s.char() != ' ' && s.char() != '.' && s.char() != '\n' && s.char() != ']' {
		s.pos++
	}
	if s.pos == nameStart {
		s.pos = saved
		return MfmNode{}, false
	}
	name := s.input[nameStart:s.pos]

	// Parse optional args (after '.')
	args := map[string]any{}
	if !s.eof() && s.char() == '.' {
		s.pos++ // skip .
		argsStart := s.pos
		for !s.eof() && s.char() != ' ' && s.char() != '\n' && s.char() != ']' {
			s.pos++
		}
		argsStr := s.input[argsStart:s.pos]
		for _, part := range strings.Split(argsStr, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if eqIdx := strings.IndexByte(part, '='); eqIdx >= 0 {
				args[part[:eqIdx]] = part[eqIdx+1:]
			} else {
				args[part] = true
			}
		}
	}

	// Expect space before content
	if s.eof() || s.char() != ' ' {
		s.pos = saved
		return MfmNode{}, false
	}
	s.pos++ // skip space

	children, ok := s.parseInlineUntil("]")
	if !ok {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{
		Type:     nodeTypeFn,
		Props:    map[string]any{"name": name, "args": args},
		Children: children,
	}, true
}

func (s *parserState) tryUrl() (MfmNode, bool) {
	if s.linkLabel {
		return MfmNode{}, false
	}
	rem := s.remaining()
	if !strings.HasPrefix(rem, "https://") && !strings.HasPrefix(rem, "http://") {
		return MfmNode{}, false
	}
	start := s.pos
	// Find scheme end
	schemeEnd := strings.Index(rem, "://") + 3
	s.pos += schemeEnd

	// Consume URL characters with balanced parentheses
	parenDepth := 0
	for !s.eof() {
		ch := s.char()
		if ch == '(' {
			parenDepth++
			s.pos++
		} else if ch == ')' {
			if parenDepth > 0 {
				parenDepth--
				s.pos++
			} else {
				break
			}
		} else if ch <= ' ' || ch == '"' || ch == '<' || ch == '>' || ch == '[' || ch == ']' {
			break
		} else {
			s.pos++
		}
	}
	// Strip trailing punctuation
	for s.pos > start && (s.input[s.pos-1] == '.' || s.input[s.pos-1] == ',') {
		s.pos--
	}
	url := s.input[start:s.pos]
	if len(url) <= schemeEnd {
		s.pos = start
		return MfmNode{}, false
	}
	return MfmNode{Type: nodeTypeUrl, Props: map[string]any{"url": url}}, true
}

func (s *parserState) tryUrlAlt() (MfmNode, bool) {
	if s.linkLabel {
		return MfmNode{}, false
	}
	saved := s.pos
	if !s.consume("<") {
		return MfmNode{}, false
	}
	rem := s.remaining()
	if !strings.HasPrefix(rem, "https://") && !strings.HasPrefix(rem, "http://") {
		s.pos = saved
		return MfmNode{}, false
	}
	idx := strings.IndexByte(rem, '>')
	if idx < 0 {
		s.pos = saved
		return MfmNode{}, false
	}
	url := rem[:idx]
	s.pos += idx + 1
	return MfmNode{Type: nodeTypeUrl, Props: map[string]any{"url": url}}, true
}

func (s *parserState) tryMention() (MfmNode, bool) {
	if s.linkLabel {
		return MfmNode{}, false
	}
	if s.pos > 0 {
		prev := s.input[s.pos-1]
		if isAlphanumeric(rune(prev)) {
			return MfmNode{}, false
		}
	}
	saved := s.pos
	if !s.consume("@") {
		return MfmNode{}, false
	}
	// Parse username: [a-zA-Z0-9_-], must start with [a-zA-Z0-9_]
	start := s.pos
	if s.eof() {
		s.pos = saved
		return MfmNode{}, false
	}
	// First char must not be - or .
	ch0 := s.char()
	if ch0 == '-' || ch0 == '.' {
		s.pos = saved
		return MfmNode{}, false
	}
	for !s.eof() {
		ch := s.char()
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
			s.pos++
		} else {
			break
		}
	}
	if s.pos == start {
		s.pos = saved
		return MfmNode{}, false
	}
	username := s.input[start:s.pos]
	// Strip trailing hyphens/dots from username
	for len(username) > 0 && (username[len(username)-1] == '-' || username[len(username)-1] == '.') {
		username = username[:len(username)-1]
		s.pos--
	}

	// Optional @host
	host := ""
	if !s.eof() && s.char() == '@' {
		s.pos++ // skip @
		hostStart := s.pos
		for !s.eof() {
			ch := s.char()
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' || ch == '_' {
				s.pos++
			} else {
				break
			}
		}
		if s.pos == hostStart {
			s.pos-- // put back @
		} else {
			host = s.input[hostStart:s.pos]
			// Strip trailing dots
			for len(host) > 0 && host[len(host)-1] == '.' {
				host = host[:len(host)-1]
				s.pos--
			}
		}
	}

	acct := "@" + username
	if host != "" {
		acct += "@" + host
	}
	return MfmNode{Type: nodeTypeMention, Props: map[string]any{
		"username": username,
		"host":     nilIfEmpty(host),
		"acct":     acct,
	}}, true
}

func (s *parserState) tryHashtag() (MfmNode, bool) {
	if s.linkLabel {
		return MfmNode{}, false
	}
	// Must not be preceded by alphanumeric
	if s.pos > 0 && isAlphanumeric(rune(s.input[s.pos-1])) {
		return MfmNode{}, false
	}
	saved := s.pos
	if !s.consume("#") {
		return MfmNode{}, false
	}
	start := s.pos
	// Balanced bracket pairs
	type bracketPair struct{ open, close rune }
	brackets := []bracketPair{
		{'(', ')'}, {'[', ']'}, {'「', '」'}, {'（', '）'},
	}
	bracketStack := []rune{}

	for !s.eof() {
		r, sz := utf8.DecodeRuneInString(s.remaining())
		if r == ' ' || r == '\n' || r == '\t' || r == '.' || r == ',' || r == '!' || r == '?' || r == '\'' || r == '"' || r == '#' {
			break
		}
		// Check brackets
		pushed := false
		for _, bp := range brackets {
			if r == bp.open {
				bracketStack = append(bracketStack, bp.close)
				pushed = true
				break
			}
			if r == bp.close {
				if len(bracketStack) > 0 && bracketStack[len(bracketStack)-1] == r {
					bracketStack = bracketStack[:len(bracketStack)-1]
					pushed = true
				} else {
					goto done
				}
				break
			}
		}
		_ = pushed
		s.pos += sz
	}
done:
	if s.pos == start {
		s.pos = saved
		return MfmNode{}, false
	}
	tag := s.input[start:s.pos]
	// Reject pure-numeric hashtags
	if isAllDigits(tag) {
		s.pos = saved
		return MfmNode{}, false
	}
	return MfmNode{Type: nodeTypeHashtag, Props: map[string]any{"hashtag": tag}}, true
}

func (s *parserState) tryLink() (MfmNode, bool) {
	saved := s.pos
	silent := false
	if s.consume("?") {
		silent = true
	}
	if !s.consume("[") {
		s.pos = saved
		return MfmNode{}, false
	}
	// Parse label (with linkLabel=true to disable mention/hashtag/url inside)
	oldLinkLabel := s.linkLabel
	s.linkLabel = true
	children, ok := s.parseInlineUntil("]")
	s.linkLabel = oldLinkLabel
	if !ok {
		s.pos = saved
		return MfmNode{}, false
	}
	if !s.consume("(") {
		s.pos = saved
		return MfmNode{}, false
	}
	// Parse URL
	urlStart := s.pos
	for !s.eof() && s.char() != ')' && s.char() != ' ' && s.char() != '\n' {
		s.pos++
	}
	if s.eof() || s.char() != ')' {
		s.pos = saved
		return MfmNode{}, false
	}
	url := s.input[urlStart:s.pos]
	s.pos++ // consume )
	return MfmNode{
		Type:     nodeTypeLink,
		Props:    map[string]any{"url": url, "silent": silent},
		Children: children,
	}, true
}

// --- Utility Functions ---

func textNode(text string) MfmNode {
	return MfmNode{Type: nodeTypeText, Props: map[string]any{"text": text}}
}

func mergeText(nodes []MfmNode) []MfmNode {
	if len(nodes) == 0 {
		return nodes
	}
	var result []MfmNode
	for _, n := range nodes {
		if n.Type == nodeTypeText && len(result) > 0 && result[len(result)-1].Type == nodeTypeText {
			result[len(result)-1].Props["text"] = result[len(result)-1].Props["text"].(string) + n.Props["text"].(string)
		} else {
			result = append(result, n)
		}
	}
	return result
}

func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func isItalicContent(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == ' ' || ch == '\n'
}

func isAllDigits(s string) bool {
	for _, ch := range s {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
