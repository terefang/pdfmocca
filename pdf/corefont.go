package pdf

import (
	"fmt"
	"strings"
)

type PdfResource interface {
	Ref() PdfObjRefValue
	ResName() string
	ObjId() int
}

type PdfFontObject struct {
	PdfDictObject
	unimap map[rune]int
	wmap   [256]int
}

func (pfo *PdfFontObject) EncodeWords(s string) []EncodedWord {
	return DoEncodeWords(pfo, s)
}

func DoEncodeWords(pf PdfFont, s string) []EncodedWord {
	wds := strings.Split(s, " ")
	_ret := make([]EncodedWord, 0)
	for _, wd := range wds {
		er := pf.EncodeText(wd)
		_l := 0
		for _, e := range er {
			_l += e.Width
		}
		ew := EncodedWord{Word: er, Width: _l}
		_ret = append(_ret, ew)
	}
	return _ret
}

func (pfo *PdfFontObject) MaxCid() int {
	return 255
}

func (pfo *PdfFontObject) IsCid() bool {
	return false
}

func (pfo *PdfFontObject) WidthRunes(s []rune) []int {
	_l := len(s)
	_ret := make([]int, _l)
	for _i := 0; _i < _l; _i++ {
		_ret[_i] = pfo.wmap[pfo.mapRune(s[_i])]
	}
	return _ret
}

func (pfo *PdfFontObject) SizeRunes(s []rune) int {
	_l := len(s)
	_w := 0
	for _i := 0; _i < _l; _i++ {
		_w += pfo.wmap[pfo.mapRune(s[_i])]
	}
	return _w
}

func (pfo *PdfFontObject) SizeRune(r rune) int {
	return pfo.wmap[pfo.mapRune(r)]
}

func (pfo *PdfFontObject) mapRune(_r rune) int {
	if len(pfo.unimap) != 0 {
		_e, ok := pfo.unimap[_r]
		if ok {
			return _e
		} else {
			return '?'
		}
	} else {
		if _r > 255 {
			return '?'
		} else {
			return int(_r)
		}
	}
}

func (pfo *PdfFontObject) WidthText(s string) []int {
	return pfo.WidthRunes([]rune(s))
}

func (pfo *PdfFontObject) SizeText(s string) int {
	return pfo.SizeRunes([]rune(s))
}

func (pfo *PdfFontObject) ResName() string {
	return fmt.Sprintf("F%d", pfo.GetObjNum())
}

func (pfo *PdfFontObject) ObjId() int {
	return pfo.GetObjNum()
}

func (pfo *PdfFontObject) EncodeText(s string) []EncodedRune {
	if len(pfo.unimap) == 0 {
		return MakeEncodeString(s, nil, &pfo.wmap)
	}
	return MakeEncodeString(s, &pfo.unimap, &pfo.wmap)
}

func MakeEncodeString(s string, unimap *map[rune]int, wmap *[256]int) []EncodedRune {
	return MakeEncodeRunes([]rune(s), unimap, wmap)
}

func MakeEncodeRunes(runes []rune, unimap *map[rune]int, wmap *[256]int) []EncodedRune {
	_l := len(runes)
	_ret := make([]EncodedRune, _l)
	for _i := 0; _i < _l; _i++ {
		if unimap == nil {
			_ret[_i] = EncodedRune{Char: runes[_i], Gid: int(runes[_i] & 0xff), Width: wmap[int(runes[_i]&0xff)]}
			continue
		}
		_e, _ok := (*unimap)[runes[_i]]
		if _ok {
			_ret[_i] = EncodedRune{Char: runes[_i], Gid: _e, Width: wmap[_e]}
		} else {
			_ret[_i] = EncodedRune{Char: runes[_i], Gid: '?', Width: wmap['?']}
		}
	}
	return _ret
}

func (pfo *PdfFontObject) MakeMap(cmap [256]int) {
	for _i, _r := range cmap {
		pfo.unimap[rune(_r)] = (_i & 0xff)
	}
}

func NewBaseCrrFont(pdf *PdfDoc, e bool) PdfFont {
	if e {
		return NewPdfTTFont(FONT_CORE_COURIER_REGULAR, pdf, UNICODE_MAP_PDFDOC)
	} else {
		return NewPdfCoreFontObject(pdf, FONT_CORE_COURIER_REGULAR, 0, 255, PDFDOC_COURIER_W, PDFDOC_ENCODING)
	}
}

func NewBaseCriFont(pdf *PdfDoc, e bool) PdfFont {
	if e {
		return NewPdfTTFont(FONT_CORE_COURIER_ITALIC, pdf, UNICODE_MAP_PDFDOC)
	} else {
		return NewPdfCoreFontObject(pdf, FONT_CORE_COURIER_ITALIC, 0, 255, PDFDOC_COURIER_W, PDFDOC_ENCODING)
	}
}

func NewBaseCrzFont(pdf *PdfDoc, e bool) PdfFont {
	if e {
		return NewPdfTTFont(FONT_CORE_COURIER_BOLD_ITALIC, pdf, UNICODE_MAP_PDFDOC)
	} else {
		return NewPdfCoreFontObject(pdf, FONT_CORE_COURIER_BOLD_ITALIC, 0, 255, PDFDOC_COURIER_W, PDFDOC_ENCODING)
	}
}

func NewBaseCrbFont(pdf *PdfDoc, e bool) PdfFont {
	if e {
		return NewPdfTTFont(FONT_CORE_COURIER_BOLD, pdf, UNICODE_MAP_PDFDOC)
	} else {
		return NewPdfCoreFontObject(pdf, FONT_CORE_COURIER_BOLD, 0, 255, PDFDOC_COURIER_W, PDFDOC_ENCODING)
	}
}

func NewPdfCoreFontObject(pdf *PdfDoc, fn string, fc int, lc int, wl [256]int, el [256]string) *PdfFontObject {
	pfo := &PdfFontObject{
		unimap: make(map[rune]int),
		wmap:   wl,
		PdfDictObject: PdfDictObject{
			PdfBaseObject: PdfBaseObject{
				objnum: -1,
				pdf:    nil},
			object: NewPdfDictValue(),
		},
	}
	pfo.Dict().Set("Type", NewPdfNameValue("Font"))
	pfo.Dict().Set("Subtype", NewPdfNameValue("Type1"))
	pfo.Dict().Set("BaseFont", NewPdfNameValue(fn))
	pfo.Dict().Set("FirstChar", NewPdfIntValue(int64(fc)))
	pfo.Dict().Set("LastChar", NewPdfIntValue(int64(lc)))
	_sb := strings.Builder{}
	_sb.WriteString("[ ")
	for _i := fc; _i <= lc && _i < 256; _i++ {
		_sb.WriteString(fmt.Sprintf("%d", wl[_i]))
		if (_i % 16) == 0 {
			_sb.WriteString("\n")
		} else {
			_sb.WriteString(" ")
		}
	}
	_sb.WriteString("]")
	pfo.Dict().Set("Widths", NewPdfLiteralValue(_sb.String()))
	_sb.Reset()
	_sb.WriteString(fmt.Sprintf("<</Type/Encoding /BaseEncoding/WinAnsiEncoding /Differences[ %d ", fc))
	for _i := fc; _i <= lc; _i++ {
		_sb.WriteString("/" + el[_i])
		if (_i % 16) == 0 {
			_sb.WriteString("\n")
		} else {
			_sb.WriteString(" ")
		}
	}
	_sb.WriteString("]>>")
	pfo.Dict().Set("Encoding", NewPdfLiteralValue(_sb.String()))
	_sb.Reset()
	pdf.AddObject(pfo)
	pdf.StreamOut(pfo.ObjId())
	return pfo
}

const (
	FONT_CORE_TIMES_REGULAR     = "Times-Roman"
	FONT_CORE_TIMES_BOLD        = "Times-Bold"
	FONT_CORE_TIMES_ITALIC      = "Times-Italic"
	FONT_CORE_TIMES_BOLD_ITALIC = "Times-BoldItalic"

	FONT_CORE_HELV_REGULAR     = "Helvetica"
	FONT_CORE_HELV_BOLD        = "Helvetica-Bold"
	FONT_CORE_HELV_ITALIC      = "Helvetica-Oblique"
	FONT_CORE_HELV_BOLD_ITALIC = "Helvetica-BoldOblique"

	FONT_CORE_COURIER_REGULAR     = "Courier"
	FONT_CORE_COURIER_BOLD        = "Courier-Bold"
	FONT_CORE_COURIER_ITALIC      = "Courier-Oblique"
	FONT_CORE_COURIER_BOLD_ITALIC = "Courier-BoldOblique"

	FONT_CORE_SYMBOL = "Symbol"

	FONT_CORE_DINGBATS = "ZapfDingbats"
)

var PDFDOC_ENCODING_CODEPOINTS = [256]int{
	0x000, 0x001, 0x002, 0x003, 0x004, 0x005, 0x006, 0x007, 0x008, 0x009, 0x00a, 0x00b, 0x00c, 0x00d, 0x00e, 0x00f,
	0x010, 0x011, 0x012, 0x013, 0x014, 0x015, 0x016, 0x017, 0x2D8, 0x2C7, 0x2C6, 0x2D9, 0x2DD, 0x2DB, 0x2DA, 0x2DC,
	0x020, 0x021, 0x022, 0x023, 0x024, 0x025, 0x026, 0x027, 0x028, 0x029, 0x02a, 0x02b, 0x02c, 0x02d, 0x02e, 0x02f,
	0x030, 0x031, 0x032, 0x033, 0x034, 0x035, 0x036, 0x037, 0x038, 0x039, 0x03a, 0x03b, 0x03c, 0x03d, 0x03e, 0x03f,
	0x040, 0x041, 0x042, 0x043, 0x044, 0x045, 0x046, 0x047, 0x048, 0x049, 0x04a, 0x04b, 0x04c, 0x04d, 0x04e, 0x04f,
	0x050, 0x051, 0x052, 0x053, 0x054, 0x055, 0x056, 0x057, 0x058, 0x059, 0x05a, 0x05b, 0x05c, 0x05d, 0x05e, 0x05f,
	0x060, 0x061, 0x062, 0x063, 0x064, 0x065, 0x066, 0x067, 0x068, 0x069, 0x06a, 0x06b, 0x06c, 0x06d, 0x06e, 0x06f,
	0x070, 0x071, 0x072, 0x073, 0x074, 0x075, 0x076, 0x077, 0x078, 0x079, 0x07a, 0x07b, 0x07c, 0x07d, 0x07e, 0xfffd,
	0x2022, 0x2020, 0x2021, 0x2026, 0x2014, 0x2013, 0x192, 0x2044,
	0x2039, 0x203A, 0x2212, 0x2030, 0x201E, 0x201C, 0x201D, 0x2018,
	0x2019, 0x201A, 0x2122, 0xFB01, 0xFB02, 0x141, 0x152, 0x160,
	0x178, 0x17D, 0x131, 0x142, 0x153, 0x161, 0x17E, 0xfffd,
	0x20AC, 0x0a1, 0x0a2, 0x0a3, 0x0a4, 0x0a5, 0x0a6, 0x0a7, 0x0a8, 0x0a9, 0x0aa, 0x0ab, 0x0ac, 0x0ad, 0x0ae, 0x0af,
	0x0b0, 0x0b1, 0x0b2, 0x0b3, 0x0b4, 0x0b5, 0x0b6, 0x0b7, 0x0b8, 0x0b9, 0x0ba, 0x0bb, 0x0bc, 0x0bd, 0x0be, 0x0bf,
	0x0c0, 0x0c1, 0x0c2, 0x0c3, 0x0c4, 0x0c5, 0x0c6, 0x0c7, 0x0c8, 0x0c9, 0x0ca, 0x0cb, 0x0cc, 0x0cd, 0x0ce, 0x0cf,
	0x0d0, 0x0d1, 0x0d2, 0x0d3, 0x0d4, 0x0d5, 0x0d6, 0x0d7, 0x0d8, 0x0d9, 0x0da, 0x0db, 0x0dc, 0x0dd, 0x0de, 0x0df,
	0x0e0, 0x0e1, 0x0e2, 0x0e3, 0x0e4, 0x0e5, 0x0e6, 0x0e7, 0x0e8, 0x0e9, 0x0ea, 0x0eb, 0x0ec, 0x0ed, 0x0ee, 0x0ef,
	0x0f0, 0x0f1, 0x0f2, 0x0f3, 0x0f4, 0x0f5, 0x0f6, 0x0f7, 0x0f8, 0x0f9, 0x0fa, 0x0fb, 0x0fc, 0x0fd, 0x0fe, 0x0ff,
}

var PDFDOC_ENCODING = [256]string{
	"uni0000", "controlSTX", "controlSOT", "controlETX", "controlEOT", "controlENQ", "controlACK", "controlBEL",
	"controlBS", "controlHT", "controlLF", "controlVT", "controlFF", "controlCR", "controlSO", "controlSI",
	"controlDLE", "controlDC1", "controlDC2", "controlDC3", "controlDC4", "controlNAK", "controlSYN", "controlETB",
	"breve", "caron", "circumflex", "dotaccent", "hungarumlaut", "ogonek", "ring", "tilde",
	"space", "exclam", "quotedbl", "numbersign", "dollar", "percent", "ampersand", "quotesingle",
	"parenleft", "parenright", "asterisk", "plus", "comma", "hyphen", "period", "slash",
	"zero", "one", "two", "three", "four", "five", "six", "seven",
	"eight", "nine", "colon", "semicolon", "less", "equal", "greater", "question",
	"at", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "bracketleft", "backslash", "bracketright", "asciicircum", "underscore",
	"grave", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o",
	"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "braceleft", "bar", "braceright", "asciitilde",
	"uniFFFD", "bullet", "dagger", "daggerdbl", "ellipsis", "emdash", "endash", "florin",
	"fraction", "guilsinglleft", "guilsinglright", "minus", "perthousand", "quotedblbase", "quotedblleft", "quotedblright",
	"quoteleft", "quoteright", "quotesinglbase", "trademark", "uniFB01", "uniFB02", "Lslash", "OE",
	"Scaron", "Ydieresis", "Zcaron", "dotlessi", "lslash", "oe", "scaron", "zcaron", "uniFFFD",
	"Euro", "exclamdown", "cent", "sterling", "currency", "yen", "brokenbar", "section",
	"dieresis", "copyright", "ordfeminine", "guillemotleft", "logicalnot", "uni00AD", "registered", "macron",
	"degree", "plusminus", "uni00B2", "uni00B3", "acute", "mu", "paragraph", "periodcentered",
	"cedilla", "uni00B9", "ordmasculine", "guillemotright", "onequarter", "onehalf", "threequarters", "questiondown",
	"Agrave", "Aacute", "Acircumflex", "Atilde", "Adieresis", "Aring", "AE", "Ccedilla",
	"Egrave", "Eacute", "Ecircumflex", "Edieresis", "Igrave", "Iacute", "Icircumflex", "Idieresis",
	"Eth", "Ntilde", "Ograve", "Oacute", "Ocircumflex", "Otilde", "Odieresis", "multiply",
	"Oslash", "Ugrave", "Uacute", "Ucircumflex", "Udieresis", "Yacute", "Thorn", "germandbls",
	"agrave", "aacute", "acircumflex", "atilde", "adieresis", "aring", "ae", "ccedilla",
	"egrave", "eacute", "ecircumflex", "edieresis", "igrave", "iacute", "icircumflex", "idieresis",
	"eth", "ntilde", "ograve", "oacute", "ocircumflex", "otilde", "odieresis", "divide",
	"oslash", "ugrave", "uacute", "ucircumflex", "udieresis", "yacute", "thorn", "ydieresis"}

var SYMBOL_ENCODING = [256]string{
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", "space", "exclam", "universal", "numbersign",
	"existential", "percent", "ampersand", "suchthat", "parenleft", "parenright", "asteriskmath",
	"plus", "comma", "minus", "period", "slash", "zero", "one", "two", "three", "four", "five", "six", "seven",
	"eight", "nine", "colon", "semicolon", "less", "equal", "greater", "question", "congruent", "Alpha",
	"Beta", "Chi", "Delta", "Epsilon", "Phi", "Gamma", "Eta", "Iota", "theta1", "Kappa", "Lambda", "Mu", "Nu",
	"Omicron", "Pi", "Theta", "Rho", "Sigma", "Tau", "Upsilon", "sigma1", "Omega", "Xi", "Psi", "Zeta", "bracketleft",
	"therefore", "bracketright", "perpendicular", "underscore", "radicalex", "alpha", "beta", "chi",
	"delta", "epsilon", "phi", "gamma", "eta", "iota", "phi1", "kappa", "lambda", "mu", "nu", "omicron", "pi",
	"theta", "rho", "sigma", "tau", "upsilon", "omega1", "omega", "xi", "psi", "zeta", "braceleft", "bar", "braceright",
	"similar", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", "Euro", "Upsilon1", "minute",
	"lessequal", "fraction", "infinity", "florin", "club", "diamond", "heart", "spade", "arrowboth", "arrowleft",
	"arrowup", "arrowright", "arrowdown", "degree", "plusminus", "second", "greaterequal", "multiply",
	"proportional", "partialdiff", "bullet", "divide", "notequal", "equivalence", "approxequal", "ellipsis",
	"arrowvertex", "arrowhorizex", "carriagereturn", "aleph", "Ifraktur", "Rfraktur", "weierstrass",
	"circlemultiply", "circleplus", "emptyset", "intersection", "union", "propersuperset", "reflexsuperset",
	"notsubset", "propersubset", "reflexsubset", "element", "notelement", "angle", "gradient", "registerserif",
	"copyrightserif", "trademarkserif", "product", "radical", "dotmath", "logicalnot", "logicaland",
	"logicalor", "arrowdblboth", "arrowdblleft", "arrowdblup", "arrowdblright", "arrowdbldown",
	"lozenge", "angleleft", "registersans", "copyrightsans", "trademarksans", "summation", "parenlefttp",
	"parenleftex", "parenleftbt", "bracketlefttp", "bracketleftex", "bracketleftbt", "bracelefttp",
	"braceleftmid", "braceleftbt", "braceex", ".notdef", "angleright", "integral", "integraltp", "integralex",
	"integralbt", "parenrighttp", "parenrightex", "parenrightbt", "bracketrighttp", "bracketrightex",
	"bracketrightbt", "bracerighttp", "bracerightmid", "bracerightbt", ".notdef",
}

var DINGBAT_ENCODING = [256]string{
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", "space", "a1", "a2", "a202", "a3", "a4", "a5", "a119",
	"a118", "a117", "a11", "a12", "a13", "a14", "a15", "a16", "a105", "a17", "a18", "a19", "a20", "a21", "a22", "a23",
	"a24", "a25", "a26", "a27", "a28", "a6", "a7", "a8", "a9", "a10", "a29", "a30", "a31", "a32", "a33", "a34", "a35",
	"a36", "a37", "a38", "a39", "a40", "a41", "a42", "a43", "a44", "a45", "a46", "a47", "a48", "a49", "a50", "a51",
	"a52", "a53", "a54", "a55", "a56", "a57", "a58", "a59", "a60", "a61", "a62", "a63", "a64", "a65", "a66", "a67",
	"a68", "a69", "a70", "a71", "a72", "a73", "a74", "a203", "a75", "a204", "a76", "a77", "a78", "a79", "a81", "a82",
	"a83", "a84", "a97", "a98", "a99", "a100", ".notdef", "a89", "a90", "a93", "a94", "a91", "a92", "a205", "a85",
	"a206", "a86", "a87", "a88", "a95", "a96", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef", ".notdef",
	".notdef", ".notdef", ".notdef", ".notdef", "a101", "a102", "a103", "a104", "a106", "a107", "a108", "a112",
	"a111", "a110", "a109", "a120", "a121", "a122", "a123", "a124", "a125", "a126", "a127", "a128", "a129", "a130",
	"a131", "a132", "a133", "a134", "a135", "a136", "a137", "a138", "a139", "a140", "a141", "a142", "a143", "a144",
	"a145", "a146", "a147", "a148", "a149", "a150", "a151", "a152", "a153", "a154", "a155", "a156", "a157", "a158",
	"a159", "a160", "a161", "a163", "a164", "a196", "a165", "a192", "a166", "a167", "a168", "a169", "a170", "a171",
	"a172", "a173", "a162", "a174", "a175", "a176", "a177", "a178", "a179", "a193", "a180", "a199", "a181", "a200",
	"a182", ".notdef", "a201", "a183", "a184", "a197", "a185", "a194", "a198", "a186", "a195", "a187", "a188",
	"a189", "a190", "a191", ".notdef",
}

var DINGBAT_ENCODING_CODEPOINTS = [256]int{
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000,
	0x2701, //0x21 0x2701
	0x2702, //0x22 0x2702
	0x2703, //0x23 0x2703
	0x2704, //0x24 0x2704
	0x260E, //0x25 0x260E
	0x2706, //0x26 0x2706
	0x2707, //0x27 0x2707
	0x2708, //0x28 0x2708
	0x2709, //0x29 0x2709
	0x261B, //0x2A 0x261B
	0x261E, //0x2B 0x261E
	0x270C, //0x2C 0x270C
	0x270D, //0x2D 0x270D
	0x270E, //0x2E 0x270E
	0x270F, //0x2F 0x270F
	0x2710, //0x30 0x2710
	0x2711, //0x31 0x2711
	0x2712, //0x32 0x2712
	0x2713, //0x33 0x2713
	0x2714, //0x34 0x2714
	0x2715, //0x35 0x2715
	0x2716, //0x36 0x2716
	0x2717, //0x37 0x2717
	0x2718, //0x38 0x2718
	0x2719, //0x39 0x2719
	0x271A, //0x3A 0x271A
	0x271B, //0x3B 0x271B
	0x271C, //0x3C 0x271C
	0x271D, //0x3D 0x271D
	0x271E, //0x3E 0x271E
	0x271F, //0x3F 0x271F
	0x2720, //0x40 0x2720
	0x2721, //0x41 0x2721
	0x2722, //0x42 0x2722
	0x2723, //0x43 0x2723
	0x2724, //0x44 0x2724
	0x2725, //0x45 0x2725
	0x2726, //0x46 0x2726
	0x2727, //0x47 0x2727
	0x2605, //0x48 0x2605
	0x2729, //0x49 0x2729
	0x272A, //0x4A 0x272A
	0x272B, //0x4B 0x272B
	0x272C, //0x4C 0x272C
	0x272D, //0x4D 0x272D
	0x272E, //0x4E 0x272E
	0x272F, //0x4F 0x272F
	0x2730, //0x50 0x2730
	0x2731, //0x51 0x2731
	0x2732, //0x52 0x2732
	0x2733, //0x53 0x2733
	0x2734, //0x54 0x2734
	0x2735, //0x55 0x2735
	0x2736, //0x56 0x2736
	0x2737, //0x57 0x2737
	0x2738, //0x58 0x2738
	0x2739, //0x59 0x2739
	0x273A, //0x5A 0x273A
	0x273B, //0x5B 0x273B
	0x273C, //0x5C 0x273C
	0x273D, //0x5D 0x273D
	0x273E, //0x5E 0x273E
	0x273F, //0x5F 0x273F
	0x2740, //0x60 0x2740
	0x2741, //0x61 0x2741
	0x2742, //0x62 0x2742
	0x2743, //0x63 0x2743
	0x2744, //0x64 0x2744
	0x2745, //0x65 0x2745
	0x2746, //0x66 0x2746
	0x2747, //0x67 0x2747
	0x2748, //0x68 0x2748
	0x2749, //0x69 0x2749
	0x274A, //0x6A 0x274A
	0x274B, //0x6B 0x274B
	0x274D, //0x6D 0x274D
	0x25A0, //0x6E 0x25A0
	0x274F, //0x6F 0x274F
	0x2750, //0x70 0x2750
	0x2751, //0x71 0x2751
	0x2752, //0x72 0x2752
	0x25B2, //0x73 0x25B2
	0x25BC, //0x74 0x25BC
	0x25C6, //0x75 0x25C6
	0x2756, //0x76 0x2756
	0x000,
	0x2758, //0x78 0x2758
	0x2759, //0x79 0x2759
	0x275A, //0x7A 0x275A
	0x275B, //0x7B 0x275B
	0x275C, //0x7C 0x275C
	0x275D, //0x7D 0x275D
	0x275E, //0x7E 0x275E
	0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000, 0x000,
	0x000,
	0x2761, //0xA1 0x2761
	0x2762, //0xA2 0x2762
	0x2763, //0xA3 0x2763
	0x2764, //0xA4 0x2764
	0x2765, //0xA5 0x2765
	0x2766, //0xA6 0x2766
	0x2767, //0xA7 0x2767
	0x2663, //0xA8 0x2663
	0x2666, //0xA9 0x2666
	0x2665, //0xAA 0x2665
	0x2660, //0xAB 0x2660
	0x2460, //0xAC 0x2460
	0x2461, //0xAD 0x2461
	0x2462, //0xAE 0x2462
	0x2463, //0xAF 0x2463
	0x2464, //0xB0 0x2464
	0x2465, //0xB1 0x2465
	0x2466, //0xB2 0x2466
	0x2467, //0xB3 0x2467
	0x2468, //0xB4 0x2468
	0x2469, //0xB5 0x2469
	0x2776, //0xB6 0x2776
	0x2777, //0xB7 0x2777
	0x2778, //0xB8 0x2778
	0x2779, //0xB9 0x2779
	0x277A, //0xBA 0x277A
	0x277B, //0xBB 0x277B
	0x277C, //0xBC 0x277C
	0x277D, //0xBD 0x277D
	0x277E, //0xBE 0x277E
	0x277F, //0xBF 0x277F
	0x2780, //0xC0 0x2780
	0x2781, //0xC1 0x2781
	0x2782, //0xC2 0x2782
	0x2783, //0xC3 0x2783
	0x2784, //0xC4 0x2784
	0x2785, //0xC5 0x2785
	0x2786, //0xC6 0x2786
	0x2787, //0xC7 0x2787
	0x2788, //0xC8 0x2788
	0x2789, //0xC9 0x2789
	0x278A, //0xCA 0x278A
	0x278B, //0xCB 0x278B
	0x278C, //0xCC 0x278C
	0x278D, //0xCD 0x278D
	0x278E, //0xCE 0x278E
	0x278F, //0xCF 0x278F
	0x2790, //0xD0 0x2790
	0x2791, //0xD1 0x2791
	0x2792, //0xD2 0x2792
	0x2793, //0xD3 0x2793
	0x2794, //0xD4 0x2794
	0x2192, //0xD5 0x2192
	0x2194, //0xD6 0x2194
	0x2195, //0xD7 0x2195
	0x2798, //0xD8 0x2798
	0x2799, //0xD9 0x2799
	0x279A, //0xDA 0x279A
	0x279B, //0xDB 0x279B
	0x279C, //0xDC 0x279C
	0x279D, //0xDD 0x279D
	0x279E, //0xDE 0x279E
	0x279F, //0xDF 0x279F
	0x27A0, //0xE0 0x27A0
	0x27A1, //0xE1 0x27A1
	0x27A2, //0xE2 0x27A2
	0x27A3, //0xE3 0x27A3
	0x27A4, //0xE4 0x27A4
	0x27A5, //0xE5 0x27A5
	0x27A6, //0xE6 0x27A6
	0x27A7, //0xE7 0x27A7
	0x27A8, //0xE8 0x27A8
	0x27A9, //0xE9 0x27A9
	0x27AA, //0xEA 0x27AA
	0x27AB, //0xEB 0x27AB
	0x27AC, //0xEC 0x27AC
	0x27AD, //0xED 0x27AD
	0x27AE, //0xEE 0x27AE
	0x27AF, //0xEF 0x27AF
	0x000,
	0x27B1, //0xF1 0x27B1
	0x27B2, //0xF2 0x27B2
	0x27B3, //0xF3 0x27B3
	0x27B4, //0xF4 0x27B4
	0x27B5, //0xF5 0x27B5
	0x27B6, //0xF6 0x27B6
	0x27B7, //0xF7 0x27B7
	0x27B8, //0xF8 0x27B8
	0x27B9, //0xF9 0x27B9
	0x27BA, //0xFA 0x27BA
	0x27BB, //0xFB 0x27BB
	0x27BC, //0xFC 0x27BC
	0x27BD, //0xFD 0x27BD
	0x27BE, //0xFE 0x27BE
	0x000,
}

var PDF_SYMBOL_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 250, 333, 713, 500, 549, 833, 778, 439,
	333, 333, 500, 549, 250, 549, 250, 278, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 278, 278,
	549, 549, 549, 444, 549, 722, 667, 722, 612, 611, 763, 603, 722, 333, 631, 722, 686, 889, 722, 722,
	768, 741, 556, 592, 611, 690, 439, 768, 645, 795, 611, 333, 863, 333, 658, 500, 500, 631, 549, 549,
	494, 439, 521, 411, 603, 329, 603, 549, 549, 576, 521, 549, 549, 521, 549, 603, 439, 576, 713, 686,
	493, 686, 494, 480, 200, 480, 549, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	750, 620, 247, 549, 167, 713, 500, 753, 753, 753, 753, 1042, 987, 603, 987, 603, 400, 549, 411, 549,
	549, 713, 494, 460, 549, 549, 549, 549, 1000, 603, 1000, 658, 823, 686, 795, 987, 768, 768, 823, 768,
	768, 713, 713, 713, 713, 713, 713, 713, 768, 713, 790, 790, 890, 823, 549, 250, 713, 603, 603, 1042,
	987, 603, 987, 603, 494, 329, 790, 790, 786, 713, 384, 384, 384, 384, 384, 384, 494, 494, 494, 494,
	300, 329, 274, 686, 686, 686, 384, 384, 384, 384, 384, 384, 494, 494, 494, 300,
}

var PDF_DINGBATS_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 278, 974, 961, 974, 980, 719, 789, 790,
	791, 690, 960, 939, 549, 855, 911, 933, 911, 945, 974, 755, 846, 762, 761, 571, 677, 763, 760, 759,
	754, 494, 552, 537, 577, 692, 786, 788, 788, 790, 793, 794, 816, 823, 789, 841, 823, 833, 816, 831,
	923, 744, 723, 749, 790, 792, 695, 776, 768, 792, 759, 707, 708, 682, 701, 826, 815, 789, 789, 707,
	687, 696, 689, 786, 787, 713, 791, 785, 791, 873, 761, 762, 762, 759, 759, 892, 892, 788, 784, 438,
	138, 277, 415, 392, 392, 668, 668, 300, 390, 390, 317, 317, 276, 276, 509, 509, 410, 410, 234, 234,
	334, 334, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 732, 544, 544, 910, 667, 760, 760, 776, 595, 694, 626, 788, 788, 788, 788, 788, 788, 788, 788,
	788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788,
	788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 788, 894, 838, 1016, 458, 748, 924, 748, 918,
	927, 928, 928, 834, 873, 828, 924, 924, 917, 930, 931, 463, 883, 836, 836, 867, 867, 696, 696, 874,
	300, 874, 760, 946, 771, 865, 771, 888, 967, 888, 831, 873, 927, 970, 918, 300,
}

var PDFDOC_TIMES_REG_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 250, 333, 408, 500, 500, 833, 778, 180, 333, 333, 500, 564,
	250, 333, 250, 278, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 278, 278, 564, 564, 564, 444, 921, 722,
	667, 667, 722, 611, 556, 722, 722, 333, 389, 722, 611, 889, 722, 722, 556, 722, 667, 556, 611, 722, 722, 944,
	722, 722, 611, 333, 278, 333, 469, 500, 333, 444, 500, 444, 500, 444, 333, 500, 500, 278, 278, 500, 278, 778,
	500, 500, 500, 500, 333, 389, 278, 500, 500, 722, 500, 500, 444, 480, 200, 480, 541, 300, 350, 500, 500, 1000,
	1000, 500, 500, 167, 333, 333, 564, 1000, 444, 444, 444, 333, 333, 333, 980, 300, 300, 611, 889, 556, 722, 611,
	278, 278, 722, 389, 444, 300, 500, 333, 500, 500, 500, 500, 200, 500, 333, 760, 276, 500, 564, 300, 760, 333,
	400, 564, 300, 300, 333, 500, 453, 250, 333, 300, 310, 500, 750, 750, 750, 444, 722, 722, 722, 722, 722, 722,
	889, 667, 611, 611, 611, 611, 333, 333, 333, 333, 722, 722, 722, 722, 722, 722, 722, 564, 722, 722, 722, 722,
	722, 722, 556, 500, 444, 444, 444, 444, 444, 444, 667, 444, 444, 444, 444, 444, 278, 278, 278, 278, 500, 500,
	500, 500, 500, 500, 500, 564, 500, 500, 500, 500, 500, 500, 500, 500}

var PDFDOC_TIMES_ITALIC_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 333, 333, 333,
	333, 333, 333, 333, 333, 250, 333, 420, 500, 500, 833, 778, 214, 333, 333, 500, 675, 250, 333, 250, 278, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 333, 333, 675, 675, 675, 500, 920, 611, 611, 667, 722, 611, 611, 722, 722, 333, 444, 667, 556, 833, 667, 722, 611,
	722, 611, 500, 556, 722, 611, 833, 611, 556, 556, 389, 278, 389, 422, 500, 333, 500, 500, 444, 500, 444, 278, 500, 500, 278, 278, 444,
	278, 722, 500, 500, 500, 500, 389, 389, 278, 500, 444, 667, 444, 444, 389, 400, 275, 400, 541, 300, 350, 500, 500, 889, 889, 500, 500,
	167, 333, 333, 675, 1000, 556, 556, 556, 333, 333, 333, 980, 300, 300, 556, 944, 500, 556, 556, 278, 278, 667, 389, 389, 300, 500, 389,
	500, 500, 500, 500, 275, 500, 333, 760, 276, 500, 675, 300, 760, 333, 400, 675, 300, 300, 333, 500, 523, 250, 333, 300, 310, 500, 750,
	750, 750, 500, 611, 611, 611, 611, 611, 611, 889, 667, 611, 611, 611, 611, 333, 333, 333, 333, 722, 667, 722, 722, 722, 722, 722, 675,
	722, 722, 722, 722, 722, 556, 611, 500, 500, 500, 500, 500, 500, 500, 667, 444, 444, 444, 444, 444, 278, 278, 278, 278, 500, 500, 500,
	500, 500, 500, 500, 675, 500, 500, 500, 500, 500, 444, 500, 444}

var PDFDOC_TIMES_BOLD_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 250, 333, 555, 500, 500, 1000, 833, 278, 333, 333, 500, 570,
	250, 333, 250, 278, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500, 333, 333, 570, 570, 570, 500, 930, 722,
	667, 722, 722, 667, 611, 778, 778, 389, 500, 778, 667, 944, 722, 778, 611, 778, 722, 556, 667, 722, 722, 1000,
	722, 722, 667, 333, 278, 333, 581, 500, 333, 500, 556, 444, 556, 444, 333, 500, 556, 278, 333, 556, 278, 833,
	556, 500, 556, 556, 444, 389, 333, 556, 500, 722, 500, 500, 444, 394, 220, 394, 520, 300, 350, 500, 500, 1000,
	1000, 500, 500, 167, 333, 333, 570, 1000, 500, 500, 500, 333, 333, 333, 1000, 300, 300, 667, 1000, 556, 722,
	667, 278, 278, 722, 389, 444, 300, 500, 333, 500, 500, 500, 500, 220, 500, 333, 747, 300, 500, 570, 300, 747,
	333, 400, 570, 300, 300, 333, 556, 540, 250, 333, 300, 330, 500, 750, 750, 750, 500, 722, 722, 722, 722, 722,
	722, 1000, 722, 667, 667, 667, 667, 389, 389, 389, 389, 722, 722, 778, 778, 778, 778, 778, 570, 778, 722, 722,
	722, 722, 722, 611, 556, 500, 500, 500, 500, 500, 500, 722, 444, 444, 444, 444, 444, 278, 278, 278, 278, 500,
	556, 500, 500, 500, 500, 500, 570, 500, 556, 556, 556, 556, 500, 556, 500}

var PDFDOC_TIMES_BOLD_ITALIC_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 333, 333, 333,
	333, 333, 333, 333, 333, 250, 389, 555, 500, 500, 833, 778, 278, 333, 333, 500, 570, 250, 333, 250, 278, 500, 500, 500, 500, 500, 500,
	500, 500, 500, 500, 333, 333, 570, 570, 570, 500, 832, 667, 667, 667, 722, 667, 667, 722, 778, 389, 500, 667, 611, 889, 722, 722, 611,
	722, 667, 556, 611, 722, 667, 889, 667, 611, 611, 333, 278, 333, 570, 500, 333, 500, 500, 444, 500, 444, 333, 500, 556, 278, 278, 500,
	278, 778, 556, 500, 500, 500, 389, 389, 278, 556, 444, 667, 500, 444, 389, 348, 220, 348, 570, 300, 350, 500, 500, 1000, 1000, 500,
	500, 167, 333, 333, 606, 1000, 500, 500, 500, 333, 333, 333, 1000, 300, 300, 611, 944, 556, 611, 611, 278, 278, 722, 389, 389, 300,
	500, 389, 500, 500, 500, 500, 220, 500, 333, 747, 266, 500, 606, 300, 747, 333, 400, 570, 300, 300, 333, 576, 500, 250, 333, 300, 300,
	500, 750, 750, 750, 500, 667, 667, 667, 667, 667, 667, 944, 667, 667, 667, 667, 667, 389, 389, 389, 389, 722, 722, 722, 722, 722, 722,
	722, 570, 722, 722, 722, 722, 722, 611, 611, 500, 500, 500, 500, 500, 500, 500, 722, 444, 444, 444, 444, 444, 278, 278, 278, 278, 500,
	556, 500, 500, 500, 500, 500, 570, 500, 556, 556, 556, 556, 444, 500, 444}

var PDFDOC_HELV_ITALIC_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 278, 278, 355, 556, 556, 889, 667, 191,
	333, 333, 389, 584, 278, 333, 278, 278, 556, 556, 556, 556, 556, 556, 556, 556, 556, 556, 278, 278,
	584, 584, 584, 556, 1015, 667, 667, 722, 722, 667, 611, 778, 722, 278, 500, 667, 556, 833, 722, 778,
	667, 778, 722, 667, 611, 722, 667, 944, 667, 667, 611, 278, 278, 278, 469, 556, 333, 556, 556, 500,
	556, 556, 278, 556, 556, 222, 222, 500, 222, 833, 556, 556, 556, 556, 333, 500, 278, 556, 500, 722,
	500, 500, 500, 334, 260, 334, 584, 300, 350, 556, 556, 1000, 1000, 556, 556, 167, 333, 333, 584, 1000,
	333, 333, 333, 222, 222, 222, 1000, 300, 300, 556, 1000, 667, 667, 611, 278, 222, 944, 500, 500, 300,
	556, 333, 556, 556, 556, 556, 260, 556, 333, 737, 370, 556, 584, 300, 737, 333, 400, 584, 300, 300,
	333, 556, 537, 278, 333, 300, 365, 556, 834, 834, 834, 611, 667, 667, 667, 667, 667, 667, 1000, 722,
	667, 667, 667, 667, 278, 278, 278, 278, 722, 722, 778, 778, 778, 778, 778, 584, 778, 722, 722, 722,
	722, 667, 667, 611, 556, 556, 556, 556, 556, 556, 889, 500, 556, 556, 556, 556, 278, 278, 278, 278,
	556, 556, 556, 556, 556, 556, 556, 584, 611, 556, 556, 556, 556, 500, 556, 500}

var PDFDOC_HELV_BOLD_ITALIC_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 278, 333, 474, 556, 556, 889, 722, 238,
	333, 333, 389, 584, 278, 333, 278, 278, 556, 556, 556, 556, 556, 556, 556, 556, 556, 556, 333, 333,
	584, 584, 584, 611, 975, 722, 722, 722, 722, 667, 611, 778, 722, 278, 556, 722, 611, 833, 722, 778,
	667, 778, 722, 667, 611, 722, 667, 944, 667, 667, 611, 333, 278, 333, 584, 556, 333, 556, 611, 556,
	611, 556, 333, 611, 611, 278, 278, 556, 278, 889, 611, 611, 611, 611, 389, 556, 333, 611, 556, 778,
	556, 556, 500, 389, 280, 389, 584, 300, 350, 556, 556, 1000, 1000, 556, 556, 167, 333, 333, 584, 1000,
	500, 500, 500, 278, 278, 278, 1000, 300, 300, 611, 1000, 667, 667, 611, 278, 278, 944, 556, 500, 300,
	556, 333, 556, 556, 556, 556, 280, 556, 333, 737, 370, 556, 584, 300, 737, 333, 400, 584, 300, 300,
	333, 611, 556, 278, 333, 300, 365, 556, 834, 834, 834, 611, 722, 722, 722, 722, 722, 722, 1000, 722,
	667, 667, 667, 667, 278, 278, 278, 278, 722, 722, 778, 778, 778, 778, 778, 584, 778, 722, 722, 722,
	722, 667, 667, 611, 556, 556, 556, 556, 556, 556, 889, 556, 556, 556, 556, 556, 278, 278, 278, 278,
	611, 611, 611, 611, 611, 611, 611, 584, 611, 611, 611, 611, 611, 556, 611, 556}

var PDFDOC_HELV_REG_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 278, 278, 355, 556, 556, 889, 667, 191,
	333, 333, 389, 584, 278, 333, 278, 278, 556, 556, 556, 556, 556, 556, 556, 556, 556, 556, 278, 278,
	584, 584, 584, 556, 1015, 667, 667, 722, 722, 667, 611, 778, 722, 278, 500, 667, 556, 833, 722, 778,
	667, 778, 722, 667, 611, 722, 667, 944, 667, 667, 611, 278, 278, 278, 469, 556, 333, 556, 556, 500,
	556, 556, 278, 556, 556, 222, 222, 500, 222, 833, 556, 556, 556, 556, 333, 500, 278, 556, 500, 722,
	500, 500, 500, 334, 260, 334, 584, 300, 350, 556, 556, 1000, 1000, 556, 556, 167, 333, 333, 584, 1000,
	333, 333, 333, 222, 222, 222, 1000, 300, 300, 556, 1000, 667, 667, 611, 278, 222, 944, 500, 500, 300,
	556, 333, 556, 556, 556, 556, 260, 556, 333, 737, 370, 556, 584, 300, 737, 333, 400, 584, 300, 300,
	333, 556, 537, 278, 333, 300, 365, 556, 834, 834, 834, 611, 667, 667, 667, 667, 667, 667, 1000, 722,
	667, 667, 667, 667, 278, 278, 278, 278, 722, 722, 778, 778, 778, 778, 778, 584, 778, 722, 722, 722,
	722, 667, 667, 611, 556, 556, 556, 556, 556, 556, 889, 500, 556, 556, 556, 556, 278, 278, 278, 278,
	556, 556, 556, 556, 556, 556, 556, 584, 611, 556, 556, 556, 556, 500, 556, 500}

var PDFDOC_HELV_BOLD_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 333, 333, 333, 333, 333, 333, 333, 333, 278, 333, 474, 556, 556, 889, 722, 238,
	333, 333, 389, 584, 278, 333, 278, 278, 556, 556, 556, 556, 556, 556, 556, 556, 556, 556, 333, 333,
	584, 584, 584, 611, 975, 722, 722, 722, 722, 667, 611, 778, 722, 278, 556, 722, 611, 833, 722, 778,
	667, 778, 722, 667, 611, 722, 667, 944, 667, 667, 611, 333, 278, 333, 584, 556, 333, 556, 611, 556,
	611, 556, 333, 611, 611, 278, 278, 556, 278, 889, 611, 611, 611, 611, 389, 556, 333, 611, 556, 778,
	556, 556, 500, 389, 280, 389, 584, 300, 350, 556, 556, 1000, 1000, 556, 556, 167, 333, 333, 584, 1000,
	500, 500, 500, 278, 278, 278, 1000, 300, 300, 611, 1000, 667, 667, 611, 278, 278, 944, 556, 500, 300,
	556, 333, 556, 556, 556, 556, 280, 556, 333, 737, 370, 556, 584, 300, 737, 333, 400, 584, 300, 300,
	333, 611, 556, 278, 333, 300, 365, 556, 834, 834, 834, 611, 722, 722, 722, 722, 722, 722, 1000, 722,
	667, 667, 667, 667, 278, 278, 278, 278, 722, 722, 778, 778, 778, 778, 778, 584, 778, 722, 722, 722,
	722, 667, 667, 611, 556, 556, 556, 556, 556, 556, 889, 556, 556, 556, 556, 556, 278, 278, 278, 278,
	611, 611, 611, 611, 611, 611, 611, 584, 611, 611, 611, 611, 611, 556, 611, 556}

var PDFDOC_COURIER_W = [256]int{
	300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
	300, 300, 300, 300, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 300, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 300, 300, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 300,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 300, 600, 600, 600, 600, 300, 300,
	600, 600, 600, 600, 600, 300, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
	600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600, 600,
}

var SYMBOL_ENCODING_CODEPOINTS = [256]int{
	0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007, 0x0008, 0x0009, 0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F,
	0x0010, 0x0011, 0x0012, 0x0013, 0x0014, 0x0015, 0x0016, 0x0017, 0x0018, 0x0019, 0x001A, 0x001B, 0x001C, 0x001D, 0x001E, 0x001F,
	0x0020, 0x0021, 0x2200, 0x0023, 0x2203, 0x0025, 0x0026, 0x220B, 0x0028, 0x0029, 0x2217, 0x002B, 0x002C, 0x2212, 0x002E, 0x002F,
	0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037, 0x0038, 0x0039, 0x003A, 0x003B, 0x003C, 0x003D, 0x003E, 0x003F,
	0x2245, 0x0391, 0x0392, 0x03A7, 0x0394, 0x0395, 0x03A6, 0x0393, 0x0397, 0x0399, 0x03D1, 0x039A, 0x039B, 0x039C, 0x039D, 0x039F,
	0x03A0, 0x0398, 0x03A1, 0x03A3, 0x03A4, 0x03A5, 0x03C2, 0x03A9, 0x039E, 0x03A8, 0x0396, 0x005B, 0x2234, 0x005D, 0x22A5, 0x005F,
	0x203E, 0x03B1, 0x03B2, 0x03C7, 0x03B4, 0x03B5, 0x03C6, 0x03B3, 0x03B7, 0x03B9, 0x03D5, 0x03BA, 0x03BB, 0x03BC, 0x03BD, 0x03BF,
	0x03C0, 0x03B8, 0x03C1, 0x03C3, 0x03C4, 0x03C5, 0x03D6, 0x03C9, 0x03BE, 0x03C8, 0x03B6, 0x007B, 0x007C, 0x007D, 0x223C, 0x007F,
	0x0080, 0x0081, 0x0082, 0x0083, 0x0084, 0x0085, 0x0086, 0x0087, 0x0088, 0x0089, 0x008A, 0x008B, 0x008C, 0x008D, 0x008E, 0x008F,
	0x0090, 0x0091, 0x0092, 0x0093, 0x0094, 0x0095, 0x0096, 0x0097, 0x0098, 0x0099, 0x009A, 0x009B, 0x009C, 0x009D, 0x009E, 0x009F,
	0x00A0, 0x03D2, 0x2032, 0x2264, 0x2044, 0x221E, 0x0192, 0x2663, 0x2666, 0x2665, 0x2660, 0x2194, 0x2190, 0x2191, 0x2192, 0x2193,
	0x00B0, 0x00B1, 0x2033, 0x2265, 0x00D7, 0x221D, 0x2202, 0x2022, 0x00F7, 0x2260, 0x2261, 0x2248, 0x2026, 0x00bd, 0x00be, 0x21B5,
	0x2135, 0x2111, 0x211C, 0x2118, 0x2297, 0x2295, 0x2205, 0x2229, 0x222A, 0x2283, 0x2287, 0x2284, 0x2282, 0x2286, 0x2208, 0x2209,
	0x2220, 0x2207, 0x00AE, 0x00A9, 0x2122, 0x220F, 0x221A, 0x22C5, 0x00AC, 0x2227, 0x2228, 0x21D4, 0x21D0, 0x21D1, 0x21D2, 0x21D3,
	0x25CA, 0x2329, 0x00AE, 0x00A9, 0x2122, 0x2211, 0x00E6, 0x00E7, 0x00E8, 0x00E9, 0x00Ea, 0x00Eb, 0x00Ec, 0x00Ed, 0x00Ee, 0x00Ef,
	0x00f0, 0x232A, 0x222B, 0x2320, 0x00f4, 0x2321, 0x00f6, 0x00f7, 0x00f8, 0x00f9, 0x00fa, 0x00fb, 0x00fc, 0x00fd, 0x00fe, 0x00ff,
}
