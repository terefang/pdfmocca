package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"seehuhn.de/go/sfnt"
	"seehuhn.de/go/sfnt/cmap"
	"seehuhn.de/go/sfnt/glyph"
	"seehuhn.de/go/sfnt/parser"
)

type PdfTtFont struct {
	PdfFont0Object
	theFont *sfnt.Font
}

func (p PdfTtFont) MaxCid() int {
	return (p.theFont.NumGlyphs() - 1)
}

func InitPdfTTFont(_bytes []byte, pdf *PdfDoc, obj *PdfTtFont, des *PdfDictObject, mapMode int) {
	obj.Dict().Set("DescendantFonts", NewPdfLiteralValue("["+des.Ref().AsString()+"]"))
	des.Dict().Set("Type", NewPdfNameValue("Font"))
	des.Dict().Set("CIDSystemInfo", NewPdfLiteralValue("<</Registry(Adobe)/Ordering(Identity)/Supplement 0>>"))
	des.Dict().Set("DW", NewPdfIntValue(0))
	//des.Dict().Set("BaseFont", NewPdfNameValue("TTFSSF+"+_sfnt.FullName()))
	if obj.theFont.IsCFF() {
		des.Dict().Set("Subtype", NewPdfNameValue("CIDFontType0"))
	} else {
		des.Dict().Set("Subtype", NewPdfNameValue("CIDFontType2"))
	}
	descr := pdf.NewDictObject()
	des.Dict().Set("FontDescriptor", descr.Ref())
	descr.Dict().Set("Type", NewPdfNameValue("FontDescriptor"))
	descr.Dict().Set("FontName", NewPdfNameValue(obj.theFont.FullName()))
	descr.Dict().Set("FontFamily", NewPdfStringValue(obj.theFont.FamilyName))
	descr.Dict().Set("FontBBox", NewPdfLiteralValue(fmt.Sprintf("[%d %d %d %d]", int(obj.theFont.FontBBoxPDF().LLx), int(obj.theFont.FontBBoxPDF().LLy), int(obj.theFont.FontBBoxPDF().Dx()), int(obj.theFont.FontBBoxPDF().Dy()))))
	ff2 := pdf.NewDictStreamObject()
	if obj.theFont.IsCFF() {
		descr.Dict().Set("FontFile3", ff2.Ref())
	} else {
		descr.Dict().Set("FontFile2", ff2.Ref())
	}
	ff2.SetFlateStream(_bytes)
	if obj.theFont.IsCFF() {
		ff2.Dict().Set("Subtype", NewPdfNameValue("OpenType"))
	}
	ff2.Dict().Set("Length1", NewPdfIntValue(ff2.Stream().UnCompressedLength()))

	_gnum := obj.theFont.NumGlyphs()
	_w := make([]int, _gnum)
	_wl := obj.theFont.WidthsPDF()
	_sb := strings.Builder{}
	_sb.WriteString("[")
	for _i := 0; _i < _gnum; _i++ {
		if (_i % 256) == 0 {
			if _i > 0 {
				_sb.WriteString("]")
			}
			_sb.WriteString(fmt.Sprintf(" %d [", _i))
		}
		_sb.WriteString(fmt.Sprintf(" %d ", int(_wl[_i]*1000)))
		_w[_i] = int(_wl[_i] * 1000)
	}
	_sb.WriteString("]]")
	des.Dict().Set("W", NewPdfLiteralValue(_sb.String()))
	obj.SetWidthMap(_w)
	_cmap, _err := obj.theFont.CMapTable.GetBest()
	_def := _cmap.Lookup(rune(0xFFFD))
	_uc := make([]int, _gnum)
	if _err == nil {
		if mapMode == UNICODE_MAP_PDFDOC {
			MakeBaseMapping(obj, 0, 0x500, _def, _cmap, &_uc, &PDFDOC_ENCODING_CODEPOINTS)
		} else if mapMode == UNICODE_MAP_SYMBOL {
			MakeBaseMapping(obj, 0, 0x3000, _def, _cmap, &_uc, &SYMBOL_ENCODING_CODEPOINTS)
		} else if mapMode == UNICODE_MAP_DINGBAT {
			MakeBaseMapping(obj, 0, 0x3000, _def, _cmap, &_uc, &DINGBAT_ENCODING_CODEPOINTS)
		} else if mapMode == UNICODE_MAP_FORCE {
			MakeBaseMapping(obj, 0, 0x300000, _def, _cmap, &_uc, nil)
		} else if mapMode == UNICODE_MAP_FULL {
			MakeBaseMapping(obj, 0, 0x20000, _def, _cmap, &_uc, nil)
		} else if mapMode == UNICODE_MAP_MS_DINGBAT {
			MakeBaseMapping(obj, 0, 0x100, _def, _cmap, &_uc, nil)
			MakeBaseMapping(obj, 0xf020, 0xf100, _def, _cmap, &_uc, nil)
		} else if mapMode == UNICODE_MAP_MS_UNICODE {
			MakeBaseMapping(obj, 0, 0x10000, _def, _cmap, &_uc, nil)
		} else {
			// just scan sample set
			MakeBaseMapping(obj, 0, 0x3000, _def, _cmap, &_uc, nil)

			//_map = []int{0, 0x2fff}
			//for _j := 0; _j < len(_map); _j += 2 {
			//	MakeBaseMapping(obj, _map[_j], _map[_j+1]+1, _def, _cmap, &_uc, nil)
			//}
		}
	}

	_sb = strings.Builder{}

	_sb.WriteString(`/CIDInit /ProcSet findresource begin
12 dict begin
begincmap
/CIDSystemInfo << /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def
/CMapName /Adobe-Identity-UCS def
/CMapType 2 def
1 begincodespacerange
<0000> <FFFF>
endcodespacerange
`)

	for start := 0; start < _gnum; start += 100 {
		end := start + 100
		if end > _gnum {
			end = _gnum
		}
		_sb.WriteString(fmt.Sprintf("%d beginbfchar\n", end-start))
		for _i := start; _i < end; _i++ {
			_sb.WriteString(fmt.Sprintf("<%04X> <%04X>\n", _i, _uc[_i]))
		}
		_sb.WriteString("endbfchar\n")
	}
	_sb.WriteString(`endcmap
CMapName currentdict /CMap defineresource pop
end
end
`)
	tucmap := pdf.NewDictStreamObject()
	tucmap.Dict().Set("Type", NewPdfNameValue("CMap"))
	tucmap.Dict().Set("CMapName", NewPdfNameValue(fmt.Sprintf("F%d-Adobe-Identity-UCS", obj.GetObjNum())))
	obj.Dict().Set("ToUnicode", tucmap.Ref())
	des.Dict().Set("ToUnicode", tucmap.Ref())
	ucmap := _sb.String()
	tucmap.SetFlateStream([]byte(ucmap))
	pdf.StreamOut(obj.ObjId())
	pdf.StreamOut(des.ObjId())
	pdf.StreamOut(descr.ObjId())
	pdf.StreamOut(tucmap.ObjId())
	pdf.StreamOut(ff2.ObjId())

}

func MakeBaseMapping(obj *PdfTtFont, firstChar int, lastChar int, _def glyph.ID, _cmap cmap.Subtable, _uc *[]int, codepoints *[256]int) {
	if codepoints != nil {
		for _u := 0; _u < 256; _u++ {
			_r := rune((*codepoints)[_u])
			_g := _cmap.Lookup(_r)
			if _g < 1 {
				//obj.SetUnicodeMapping(_r, int(_def))
			} else {
				obj.SetUnicodeMapping(rune(_u), int(_g))
				obj.SetUnicodeMapping(_r, int(_g))
				if (*_uc)[_g] == 0 {
					(*_uc)[_g] = _u
				}
			}
		}
	}

	for _u := firstChar; _u < lastChar; _u++ {
		_r := rune(_u)
		_g := _cmap.Lookup(_r)
		if _g < 1 {
			//obj.SetUnicodeMapping(_r, int(_def))
		} else {
			obj.SetUnicodeMapping(_r, int(_g))
			if (*_uc)[_g] == 0 {
				(*_uc)[_g] = _u
			}
		}
	}
}

func NewPdfTTFont(fn string, pdf *PdfDoc, mapMode int) PdfFont {

	var _bytes []byte

	switch fn {
	case "qdb", "qzdb", FONT_CORE_DINGBATS:
		_bytes = qdbBytes
		mapMode = UNICODE_MAP_DINGBAT
	case "qsy", "qsym", FONT_CORE_SYMBOL:
		_bytes = qsyBytes
		mapMode = UNICODE_MAP_SYMBOL

	case "qtmr", FONT_CORE_TIMES_REGULAR:
		_bytes = qtmrBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qtmi", FONT_CORE_TIMES_ITALIC:
		_bytes = qtmiBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qtmb", FONT_CORE_TIMES_BOLD:
		_bytes = qtmbBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qtmz", FONT_CORE_TIMES_BOLD_ITALIC:
		_bytes = qtmzBytes
		mapMode = UNICODE_MAP_PDFDOC

	case "qcrr", FONT_CORE_COURIER_REGULAR:
		_bytes = qcrrBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qcri", FONT_CORE_COURIER_ITALIC:
		_bytes = qcriBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qcrb", FONT_CORE_COURIER_BOLD:
		_bytes = qcrbBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qcrz", FONT_CORE_COURIER_BOLD_ITALIC:
		_bytes = qcrzBytes
		mapMode = UNICODE_MAP_PDFDOC

	case "qhvr", FONT_CORE_HELV_REGULAR:
		_bytes = qhvrBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qhvi", FONT_CORE_HELV_ITALIC:
		_bytes = qhviBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qhvb", FONT_CORE_HELV_BOLD:
		_bytes = qhvbBytes
		mapMode = UNICODE_MAP_PDFDOC
	case "qhvz", FONT_CORE_HELV_BOLD_ITALIC:
		_bytes = qhvzBytes
		mapMode = UNICODE_MAP_PDFDOC

	default:
		_bytes, _ = os.ReadFile(fn)
	}

	_rdr := bytes.NewReader(_bytes)
	_sfnt, _ := sfnt.Read(_rdr, parser.NewBudget(int64(len(_bytes))))
	//_sfnt, _ := sfnt.ReadFile(fn)

	obj := &PdfTtFont{PdfFont0Object: *NewPdfType0FontObject(fmt.Sprintf("%s-W%d-%s", _sfnt.PostScriptName(), int(_sfnt.Weight), _sfnt.Width.String()))}
	obj.theFont = _sfnt
	pdf.AddObject(obj)
	des := pdf.NewDictObject()
	InitPdfTTFont(_bytes, pdf, obj, des, mapMode)
	return obj
}

const UNICODE_MAP_FORCE = 0
const UNICODE_MAP_FULL = 1
const UNICODE_MAP_PDFDOC = 2
const UNICODE_MAP_SYMBOL = 3
const UNICODE_MAP_DINGBAT = 4
const UNICODE_MAP_MS_DINGBAT = 5
const UNICODE_MAP_MS_UNICODE = 6

/* https://learn.microsoft.com/en-us/typography/opentype/spec/os2#ur */
var unicodePageRanges map[int][]int = map[int][]int{
	0:/*Basic Latin*/ []int{0x0000, 0x007F},        //
	1:/*Latin-1 Supplement*/ []int{0x0080, 0x00FF}, //
	2:/*Latin Extended-A*/ []int{0x0100, 0x017F},   //
	3:/*Latin Extended-B*/ []int{0x0180, 0x024F},   // }, //
	4: /*IPA Extensions*/ []int{0x0250, 0x02AF,
		/*Phonetic Extensions*/ 0x1D00, 0x1D7F, //Added in OpenType 1.5 for OS/2 version 4.
		/*Phonetic Extensions Supplement*/ 0x1D80, 0x1DBF}, //Added in OpenType 1.5 for OS/2 version 4.
	5: /*Spacing Modifier Letters*/ []int{0x02B0, 0x02FF, //
		/*Modifier Tone Letters*/ 0xA700, 0xA71F}, //Added in OpenType 1.5 for OS/2 version 4.
	6: /*Combining Diacritical Marks*/ []int{0x0300, 0x036F,
		/*Combining Diacritical Marks Supplement*/ 0x1DC0, 0x1DFF}, //	Added in OpenType 1.5 for OS/2 version 4.
	7:/*Greek and Coptic*/ []int{0x0370, 0x03FF}, //
	8:/*Coptic*/ []int{0x2C80, 0x2CFF},           //	Added in OpenType 1.5 for OS/2 version 4. See below for other version differences.
	9: /*Cyrillic*/ []int{0x0400, 0x04FF, //
		/*Cyrillic Supplement*/ 0x0500, 0x052F, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Cyrillic Extended-A*/ 0x2DE0, 0x2DFF, //	Added in OpenType 1.5 for OS/2 version 4.
		/*Cyrillic Extended-B*/ 0xA640, 0xA69F}, //	Added in OpenType 1.5 for OS/2 version 4.
	10:/*Armenian*/ []int{0x0530, 0x058F}, //
	11:/*Hebrew*/ []int{0x0590, 0x05FF},   //
	12:/*Vai*/ []int{0xA500, 0xA63F},      //	Added in OpenType 1.5 for OS/2 version 4. See below for other version differences.
	13: /*Arabic*/ []int{0x0600, 0x06FF,
		/*Arabic Supplement*/ 0x0750, 0x077F}, //	Added in OpenType 1.5 for OS/2 version 4.
	14:/*NKo*/ []int{0x07C0, 0x07FF},        //	Added in OpenType 1.5 for OS/2 version 4. See below for other version differences.
	15:/*Devanagari*/ []int{0x0900, 0x097F}, //
	16:/*Bangla*/ []int{0x0980, 0x09FF},     //
	17:/*Gurmukhi*/ []int{0x0A00, 0x0A7F},   //
	18:/*Gujarati*/ []int{0x0A80, 0x0AFF},   //
	19:/*Odia*/ []int{0x0B00, 0x0B7F},       //
	20:/*Tamil*/ []int{0x0B80, 0x0BFF},      //
	21:/*Telugu*/ []int{0x0C00, 0x0C7F},     //
	22:/*Kannada*/ []int{0x0C80, 0x0CFF},    //
	23:/*Malayalam*/ []int{0x0D00, 0x0D7F},  //
	24:/*Thai*/ []int{0x0E00, 0x0E7F},       //
	25:/*Lao*/ []int{0x0E80, 0x0EFF},        //
	26: /*Georgian*/ []int{0x10A0, 0x10FF, //
		/*Georgian Supplement*/ 0x2D00, 0x2D2F}, //	Added in OpenType 1.5 for OS/2 version 4.
	27:/*Balinese*/ []int{0x1B00, 0x1B7F},    //	Added in OpenType 1.5 for OS/2 version 4. See below for other version differences.
	28:/*Hangul Jamo*/ []int{0x1100, 0x11FF}, //
	29: /*Latin Extended Additional*/ []int{0x1E00, 0x1EFF, //
		/*Latin Extended-C*/ 0x2C60, 0x2C7F, //	Added in OpenType 1.5 for OS/2 version 4.
		/*Latin Extended-D*/ 0xA720, 0xA7FF}, //	Added in OpenType 1.5 for OS/2 version 4.
	30:/*Greek Extended*/ []int{0x1F00, 0x1FFF}, //
	31: /*General Punctuation*/ []int{0x2000, 0x206F, //
		/*Supplemental Punctuation*/ 0x2E00, 0x2E7F}, //	Added in OpenType 1.5 for OS/2 version 4.
	32:/*Superscripts And Subscripts*/ []int{0x2070, 0x209F},             //
	33:/*Currency Symbols*/ []int{0x20A0, 0x20CF},                        //
	34:/*Combining Diacritical Marks For Symbols*/ []int{0x20D0, 0x20FF}, //
	35:/*Letterlike Symbols*/ []int{0x2100, 0x214F},                      //
	36:/*Number Forms*/ []int{0x2150, 0x218F},                            //
	37: /*Arrows*/ []int{0x2190, 0x21FF, //
		/*Supplemental Arrows-A*/ 0x27F0, 0x27FF, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Supplemental Arrows-B*/ 0x2900, 0x297F, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Miscellaneous Symbols and Arrows*/ 0x2B00, 0x2BFF}, //	Added in OpenType 1.5 for OS/2 version 4.
	38: /*Mathematical Operators*/ []int{0x2200, 0x22FF, //
		/*Supplemental Mathematical Operators*/ 0x2A00, 0x2AFF, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Miscellaneous Mathematical Symbols-A*/ 0x27C0, 0x27EF, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Miscellaneous Mathematical Symbols-B*/ 0x2980, 0x29FF}, //	Added in OpenType 1.4 for OS/2 version 3.
	39:/*Miscellaneous Technical*/ []int{0x2300, 0x23FF},       //
	40:/*Control Pictures*/ []int{0x2400, 0x243F},              //
	41:/*Optical Character Recognition*/ []int{0x2440, 0x245F}, //
	42:/*Enclosed Alphanumerics*/ []int{0x2460, 0x24FF},        //
	43:/*Box Drawing*/ []int{0x2500, 0x257F},                   //
	44:/*Block Elements*/ []int{0x2580, 0x259F},                //
	45:/*Geometric Shapes*/ []int{0x25A0, 0x25FF},              //
	46:/*Miscellaneous Symbols*/ []int{0x2600, 0x26FF},         //
	47:/*Dingbats*/ []int{0x2700, 0x27BF},                      //
	48:/*CJK Symbols And Punctuation*/ []int{0x3000, 0x303F},   //
	49:/*Hiragana*/ []int{0x3040, 0x309F},                      //
	50: /*Katakana*/ []int{0x30A0, 0x30FF, //
		/*Katakana Phonetic Extensions*/ 0x31F0, 0x31FF}, //	Added in OpenType 1.4 for OS/2 version 3.
	51: /*Bopomofo*/ []int{0x3100, 0x312F, //
		/*Bopomofo Extended*/ 0x31A0, 0x31BF}, //	Added in OpenType 1.3, extending OS/2 version 2.
	52:/*Hangul Compatibility Jamo*/ []int{0x3130, 0x318F},       //
	53:/*Phags-pa*/ []int{0xA840, 0xA87F},                        //	Added in OpenType 1.5 for OS/2 version 4. See below for other version differences.
	54:/*Enclosed CJK Letters And Months*/ []int{0x3200, 0x32FF}, //
	55:/*CJK Compatibility*/ []int{0x3300, 0x33FF},               //
	56:/*Hangul Syllables*/ []int{0xAC00, 0xD7AF},                //
	57:/*Non-Plane 0*/ []int{0x10000, 0x10FFFF},                  //	Setting this bit implies there is at least one character beyond the Basic Multilingual Plane supported by this font. First assigned in OpenType 1.3 for OS/2 version 2.
	58:/*Phoenician*/ []int{0x10900, 0x1091F},                    //	First assigned in OpenType 1.5 for OS/2 version 4.
	59: /*CJK Unified Ideographs*/ []int{0x4E00, 0x9FFF, //
		/*CJK Radicals Supplement*/ 0x2E80, 0x2EFF, //	Added in OpenType 1.3 for OS/2 version 2.
		/*Kangxi Radicals*/ 0x2F00, 0x2FDF, //	Added in OpenType 1.3 for OS/2 version 2.
		/*Ideographic Description Characters*/ 0x2FF0, 0x2FFF, //	Added in OpenType 1.3 for OS/2 version 2.
		/*CJK Unified Ideographs Extension A*/ 0x3400, 0x4DBF, //	Added in OpenType 1.3 for OS/2 version 2.
		/*CJK Unified Ideographs Extension B*/ 0x20000, 0x2A6DF, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Kanbun*/ 0x3190, 0x319F}, //	Added in OpenType 1.4 for OS/2 version 3.
	60:/*Private Use Area (plane 0)*/ []int{0xE000, 0xF8FF}, //
	61: /*CJK Strokes*/ []int{0x31C0, 0x31EF, //	Added in OpenType 1.5 for OS/2 version 4.
		/*CJK Compatibility Ideographs*/ 0xF900, 0xFAFF, //
		/*CJK Compatibility Ideographs Supplement*/ 0x2F800, 0x2FA1F}, //	Added in OpenType 1.4 for OS/2 version 3.
	62:/*Alphabetic Presentation Forms*/ []int{0xFB00, 0xFB4F}, //
	63:/*Arabic Presentation Forms-A*/ []int{0xFB50, 0xFDFF},   //
	64:/*Combining Half Marks*/ []int{0xFE20, 0xFE2F},          //
	65: /*Vertical Forms*/ []int{0xFE10, 0xFE1F, //	Added in OpenType 1.5 for OS/2 version 4.
		/*CJK Compatibility Forms*/ 0xFE30, 0xFE4F}, //
	66:/*Small Form Variants*/ []int{0xFE50, 0xFE6F},           //
	67:/*Arabic Presentation Forms-B*/ []int{0xFE70, 0xFEFF},   //
	68:/*Halfwidth And Fullwidth Forms*/ []int{0xFF00, 0xFFEF}, //
	69:/*Specials*/ []int{0xFFF0, 0xFFFF},                      //
	70:/*Tibetan*/ []int{0x0F00, 0x0FFF},                       //	First assigned in OpenType 1.3, extending OS/2 version 2.
	71:/*Syriac*/ []int{0x0700, 0x074F},                        //	First assigned in OpenType 1.3, extending OS/2 version 2.
	72:/*Thaana*/ []int{0x0780, 0x07BF},                        //	First assigned in OpenType 1.3, extending OS/2 version 2.
	73:/*Sinhala*/ []int{0x0D80, 0x0DFF},                       //	First assigned in OpenType 1.3, extending OS/2 version 2.
	74:/*Myanmar*/ []int{0x1000, 0x109F},                       //	First assigned in OpenType 1.3, extending OS/2 version 2.
	75: /*Ethiopic*/ []int{0x1200, 0x137F, //	First assigned in OpenType 1.3, extending OS/2 version 2.
		/*Ethiopic Supplement*/ 0x1380, 0x139F, //	Added in OpenType 1.5 for OS/2 version 4.
		/*Ethiopic Extended*/ 0x2D80, 0x2DDF}, //	Added in OpenType 1.5 for OS/2 version 4.
	76:/*Cherokee*/ []int{0x13A0, 0x13FF},                              //	First assigned in OpenType 1.3, extending OS/2 version 2.
	77:/*Unified Canadian Aboriginal Syllabics*/ []int{0x1400, 0x167F}, //	First assigned in OpenType 1.3, extending OS/2 version 2.
	78:/*Ogham*/ []int{0x1680, 0x169F},                                 //	First assigned in OpenType 1.3, extending OS/2 version 2.
	79:/*Runic*/ []int{0x16A0, 0x16FF},                                 //	First assigned in OpenType 1.3, extending OS/2 version 2.
	80: /*Khmer*/ []int{0x1780, 0x17FF, //	First assigned in OpenType 1.3, extending OS/2 version 2.
		/*Khmer Symbols*/ 0x19E0, 0x19FF}, //	Added in OpenType 1.5 for OS/2 version 4.
	81:/*Mongolian*/ []int{0x1800, 0x18AF},        //	First assigned in OpenType 1.3, extending OS/2 version 2.
	82:/*Braille Patterns*/ []int{0x2800, 0x28FF}, //	First assigned in OpenType 1.3, extending OS/2 version 2.
	83: /*Yi Syllables*/ []int{0xA000, 0xA48F, //	First assigned in OpenType 1.3, extending OS/2 version 2.
		/*Yi Radicals*/ 0xA490, 0xA4CF}, //	Added in OpenType 1.3, extending OS/2 version 2.
	84: /*Tagalog*/ []int{0x1700, 0x171F, //	First assigned in OpenType 1.4 for OS/2 version 3.
		/*Hanunoo*/ 0x1720, 0x173F, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Buhid*/ 0x1740, 0x175F, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Tagbanwa*/ 0x1760, 0x177F}, //	Added in OpenType 1.4 for OS/2 version 3.
	85:/*Old Italic*/ []int{0x10300, 0x1032F}, //	First assigned in OpenType 1.4 for OS/2 version 3.
	86:/*Gothic*/ []int{0x10330, 0x1034F},     //	First assigned in OpenType 1.4 for OS/2 version 3.
	87:/*Deseret*/ []int{0x10400, 0x1044F},    //	First assigned in OpenType 1.4 for OS/2 version 3.
	88: /*Byzantine Musical Symbols*/ []int{0x1D000, 0x1D0FF, //	First assigned in OpenType 1.4 for OS/2 version 3.
		/*Musical Symbols*/ 0x1D100, 0x1D1FF, //	Added in OpenType 1.4 for OS/2 version 3.
		/*Ancient Greek Musical Notation*/ 0x1D200, 0x1D24F}, //	Added in OpenType 1.5 for OS/2 version 4.
	89:/*Mathematical Alphanumeric Symbols*/ []int{0x1D400, 0x1D7FF}, //	First assigned in OpenType 1.4 for OS/2 version 3.
	90: /*Private Use (plane 15)*/ []int{0xF0000, 0xFFFFD, //	First assigned in OpenType 1.4 for OS/2 version 3.
		/*Private Use (plane 16)*/ 0x100000, 0x10FFFD}, //	Added in OpenType 1.4 for OS/2 version 3.
	91: /*Variation Selectors*/ []int{0xFE00, 0xFE0F, //	First assigned in OpenType 1.4 for OS/2 version 3.
		/*Variation Selectors Supplement*/ 0xE0100, 0xE01EF}, //	Added in OpenType 1.4 for OS/2 version 3.
	92:/*Tags*/ []int{0xE0000, 0xE007F},                  //	First assigned in OpenType 1.4 for OS/2 version 3.
	93:/*Limbu*/ []int{0x1900, 0x194F},                   //	First assigned in OpenType 1.5 for OS/2 version 4.
	94:/*Tai Le*/ []int{0x1950, 0x197F},                  //	First assigned in OpenType 1.5 for OS/2 version 4.
	95:/*New Tai Lue*/ []int{0x1980, 0x19DF},             //	First assigned in OpenType 1.5 for OS/2 version 4.
	96:/*Buginese*/ []int{0x1A00, 0x1A1F},                //	First assigned in OpenType 1.5 for OS/2 version 4.
	97:/*Glagolitic*/ []int{0x2C00, 0x2C5F},              //	First assigned in OpenType 1.5 for OS/2 version 4.
	98:/*Tifinagh*/ []int{0x2D30, 0x2D7F},                //	First assigned in OpenType 1.5 for OS/2 version 4.
	99:/*Yijing Hexagram Symbols*/ []int{0x4DC0, 0x4DFF}, //	First assigned in OpenType 1.5 for OS/2 version 4.
	100:/*Syloti Nagri*/ []int{0xA800, 0xA82F},           //	First assigned in OpenType 1.5 for OS/2 version 4.
	101: /*Linear B Syllabary*/ []int{0x10000, 0x1007F, //	First assigned in OpenType 1.5 for OS/2 version 4.
		/*Linear B Ideograms*/ 0x10080, 0x100FF, //	Added in OpenType 1.5 for OS/2 version 4.
		/*Aegean Numbers*/ 0x10100, 0x1013F}, //	Added in OpenType 1.5 for OS/2 version 4.
	102:/*Ancient Greek Numbers*/ []int{0x10140, 0x1018F}, //	First assigned in OpenType 1.5 for OS/2 version 4.
	103:/*Ugaritic*/ []int{0x10380, 0x1039F},              //	First assigned in OpenType 1.5 for OS/2 version 4.
	104:/*Old Persian*/ []int{0x103A0, 0x103DF},           //	First assigned in OpenType 1.5 for OS/2 version 4.
	105:/*Shavian*/ []int{0x10450, 0x1047F},               //	First assigned in OpenType 1.5 for OS/2 version 4.
	106:/*Osmanya*/ []int{0x10480, 0x104AF},               //	First assigned in OpenType 1.5 for OS/2 version 4.
	107:/*Cypriot Syllabary*/ []int{0x10800, 0x1083F},     //	First assigned in OpenType 1.5 for OS/2 version 4.
	108:/*Kharoshthi*/ []int{0x10A00, 0x10A5F},            //	First assigned in OpenType 1.5 for OS/2 version 4.
	109:/*Tai Xuan Jing Symbols*/ []int{0x1D300, 0x1D35F}, //	First assigned in OpenType 1.5 for OS/2 version 4.
	110: /*Cuneiform*/ []int{0x12000, 0x123FF, //	First assigned in OpenType 1.5 for OS/2 version 4.
		/*Cuneiform Numbers and Punctuation*/ 0x12400, 0x1247F}, //	Added in OpenType 1.5 for OS/2 version 4.
	111:/*Counting Rod Numerals*/ []int{0x1D360, 0x1D37F}, //	First assigned in OpenType 1.5 for OS/2 version 4.
	112:/*Sundanese*/ []int{0x1B80, 0x1BBF},               //	First assigned in OpenType 1.5 for OS/2 version 4.
	113:/*Lepcha*/ []int{0x1C00, 0x1C4F},                  //	First assigned in OpenType 1.5 for OS/2 version 4.
	114:/*Ol Chiki*/ []int{0x1C50, 0x1C7F},                //	First assigned in OpenType 1.5 for OS/2 version 4.
	115:/*Saurashtra*/ []int{0xA880, 0xA8DF},              //	First assigned in OpenType 1.5 for OS/2 version 4.
	116:/*Kayah Li*/ []int{0xA900, 0xA92F},                //	First assigned in OpenType 1.5 for OS/2 version 4.
	117:/*Rejang*/ []int{0xA930, 0xA95F},                  //	First assigned in OpenType 1.5 for OS/2 version 4.
	118:/*Cham*/ []int{0xAA00, 0xAA5F},                    //	First assigned in OpenType 1.5 for OS/2 version 4.
	119:/*Ancient Symbols*/ []int{0x10190, 0x101CF},       //	First assigned in OpenType 1.5 for OS/2 version 4.
	120:/*Phaistos Disc*/ []int{0x101D0, 0x101FF},         //	First assigned in OpenType 1.5 for OS/2 version 4.
	121: /*Carian*/ []int{0x102A0, 0x102DF, //	First assigned in OpenType 1.5 for OS/2 version 4.
		/*Lycian*/ 0x10280, 0x1029F, //	Added in OpenType 1.5 for OS/2 version 4.
		/*Lydian*/ 0x10920, 0x1093F}, //	Added in OpenType 1.5 for OS/2 version 4.
	122: /*Domino Tiles*/ []int{0x1F030, 0x1F09F, //	First assigned in OpenType 1.5 for OS/2 version 4.
		/*Mahjong Tiles*/ 0x1F000, 0x1F02F}, //	First assigned in OpenType 1.5 for OS/2 version 4.
}
