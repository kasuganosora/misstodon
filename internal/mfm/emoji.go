package mfm

import "unicode/utf8"

func isEmojiRune(r rune) bool {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F: // Emoticons
		return true
	case r >= 0x1F300 && r <= 0x1F5FF: // Misc Symbols & Pictographs
		return true
	case r >= 0x1F680 && r <= 0x1F6FF: // Transport & Map
		return true
	case r >= 0x1F700 && r <= 0x1F77F: // Alchemical Symbols
		return true
	case r >= 0x1F780 && r <= 0x1F7FF: // Geometric Shapes Extended
		return true
	case r >= 0x1F800 && r <= 0x1F8FF: // Supplemental Arrows-C
		return true
	case r >= 0x1F900 && r <= 0x1F9FF: // Supplemental Symbols & Pictographs
		return true
	case r >= 0x1FA00 && r <= 0x1FA6F: // Chess Symbols
		return true
	case r >= 0x1FA70 && r <= 0x1FAFF: // Symbols & Pictographs Extended-A
		return true
	case r >= 0x2600 && r <= 0x26FF: // Misc symbols
		return true
	case r >= 0x2700 && r <= 0x27BF: // Dingbats
		return true
	case r >= 0x2300 && r <= 0x23FF: // Misc Technical
		return true
	case r >= 0x2B50 && r <= 0x2B55: // Stars, circles
		return true
	case r >= 0x1F1E0 && r <= 0x1F1FF: // Regional indicators (flags)
		return true
	case r == 0x200D: // ZWJ
		return true
	case r >= 0xFE00 && r <= 0xFE0F: // Variation selectors
		return true
	case r >= 0x1F3FB && r <= 0x1F3FF: // Skin tone modifiers
		return true
	case r == 0x20E3: // Combining Enclosing Keycap
		return true
	case r >= 0xE0020 && r <= 0xE007F: // Tags (flag subdivision)
		return true
	case r == 0x00A9 || r == 0x00AE: // (c), (r)
		return true
	case r == 0x203C || r == 0x2049: // !!, !?
		return true
	case r >= 0x2100 && r <= 0x21FF: // Letterlike symbols & Arrows
		return true
	default:
		return false
	}
}

func isSkinToneModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

func (s *parserState) tryUnicodeEmoji() (MfmNode, bool) {
	start := s.pos
	r, size := utf8.DecodeRuneInString(s.input[s.pos:])
	if r == utf8.RuneError && size <= 1 {
		return MfmNode{}, false
	}

	// Keycap sequences: [0-9#*] + FE0F? + 20E3
	if r == '#' || r == '*' || (r >= '0' && r <= '9') {
		saved := s.pos
		s.pos += size
		// Optional FE0F
		if s.pos < len(s.input) {
			r2, sz := utf8.DecodeRuneInString(s.input[s.pos:])
			if r2 == 0xFE0F {
				s.pos += sz
			}
		}
		if s.pos < len(s.input) {
			r2, sz := utf8.DecodeRuneInString(s.input[s.pos:])
			if r2 == 0x20E3 {
				s.pos += sz
				return MfmNode{
					Type:  nodeTypeUnicodeEmoji,
					Props: map[string]any{"emoji": s.input[start:s.pos]},
				}, true
			}
		}
		s.pos = saved
		return MfmNode{}, false
	}

	if !isEmojiRune(r) {
		return MfmNode{}, false
	}

	// Single FE0F alone is not an emoji
	if r == 0xFE0F {
		return MfmNode{}, false
	}

	// Regional indicator: consume pairs for flags
	if isRegionalIndicator(r) {
		s.pos += size
		if s.pos < len(s.input) {
			r2, sz := utf8.DecodeRuneInString(s.input[s.pos:])
			if isRegionalIndicator(r2) {
				s.pos += sz
			}
		}
		// Consume trailing FE0F
		if s.pos < len(s.input) {
			r2, sz := utf8.DecodeRuneInString(s.input[s.pos:])
			if r2 == 0xFE0F {
				s.pos += sz
			}
		}
		return MfmNode{
			Type:  nodeTypeUnicodeEmoji,
			Props: map[string]any{"emoji": s.input[start:s.pos]},
		}, true
	}

	// Consume the initial emoji rune
	s.pos += size

	// Consume continuation: ZWJ sequences, variation selectors, skin tones, more emoji
	for s.pos < len(s.input) {
		r2, sz := utf8.DecodeRuneInString(s.input[s.pos:])
		if r2 == 0x200D { // ZWJ - consume and expect next emoji
			s.pos += sz
			if s.pos < len(s.input) {
				r3, sz3 := utf8.DecodeRuneInString(s.input[s.pos:])
				if isEmojiRune(r3) || r3 >= 0x2600 {
					s.pos += sz3
					continue
				}
			}
			continue
		}
		if r2 == 0xFE0F || isSkinToneModifier(r2) || r2 == 0x20E3 {
			s.pos += sz
			continue
		}
		// Tag sequences (for subdivision flags like 🏴󠁧󠁢󠁥󠁮󠁧󠁿)
		if r2 >= 0xE0020 && r2 <= 0xE007F {
			s.pos += sz
			continue
		}
		break
	}

	emoji := s.input[start:s.pos]
	return MfmNode{
		Type:  nodeTypeUnicodeEmoji,
		Props: map[string]any{"emoji": emoji},
	}, true
}
