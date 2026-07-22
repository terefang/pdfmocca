package pdf

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/terefang/gocustomtokener/tokener"
)

type SVGFontRoot struct {
	XMLName   xml.Name    `xml:"svg"`
	Text      string      `xml:",chardata"`
	Xmlns     string      `xml:"xmlns,attr"`
	Xlink     string      `xml:"xlink,attr"`
	Version   string      `xml:"version,attr"`
	Metadata  string      `xml:"metadata"`
	Defs      SVGFontDefs `xml:"defs"`
	GlyphMap  map[string]int
	GlyphCode map[string]string
	UniMap    map[int]int
}

type SVGFontDefs struct {
	Text string  `xml:",chardata"`
	Font SVGFont `xml:"font"`
}

type SVGFont struct {
	Text         string              `xml:",chardata"`
	ID           string              `xml:"id,attr"`
	HorizAdvX    string              `xml:"horiz-adv-x,attr"`
	FontFace     SVGFontFace         `xml:"font-face"`
	MissingGlyph SVGFontMissingGlyph `xml:"missing-glyph"`
	Glyph        []SVGFontGlyph      `xml:"glyph"`
	Hkern        []SVGFontHkern      `xml:"hkern"`
}

type SVGFontFace struct {
	Text               string `xml:",chardata"`
	FontFamily         string `xml:"font-family,attr"`
	FontWeight         string `xml:"font-weight,attr"`
	FontStretch        string `xml:"font-stretch,attr"`
	UnitsPerEm         string `xml:"units-per-em,attr"`
	Panose1            string `xml:"panose-1,attr"`
	Ascent             string `xml:"ascent,attr"`
	Descent            string `xml:"descent,attr"`
	XHeight            string `xml:"x-height,attr"`
	CapHeight          string `xml:"cap-height,attr"`
	Bbox               string `xml:"bbox,attr"`
	UnderlineThickness string `xml:"underline-thickness,attr"`
	UnderlinePosition  string `xml:"underline-position,attr"`
	Stemh              string `xml:"stemh,attr"`
	Stemv              string `xml:"stemv,attr"`
	UnicodeRange       string `xml:"unicode-range,attr"`
}
type SVGFontMissingGlyph struct {
	Text      string `xml:",chardata"`
	HorizAdvX string `xml:"horiz-adv-x,attr"`
}

type SVGFontGlyph struct {
	Text      string `xml:",chardata"`
	GlyphName string `xml:"glyph-name,attr"`
	Unicode   string `xml:"unicode,attr"`
	HorizAdvX string `xml:"horiz-adv-x,attr"`
	D         string `xml:"d,attr"`
}

type SVGFontHkern struct {
	Text string `xml:",chardata"`
	U1   string `xml:"u1,attr"`
	U2   string `xml:"u2,attr"`
	K    string `xml:"k,attr"`
}

var sfntWeightClasses map[int]string = map[int]string{
	100: "Thin",
	200: "ExtraLight",
	300: "Light",
	400: "Regular",
	500: "Medium",
	600: "SemiBold",
	700: "Bold",
	800: "ExtraBold",
	900: "Black",
	950: "ExtraBlack",
}

var sfntWidthClasses map[int]string = map[int]string{
	500:  "UltraCondensed",
	625:  "ExtraCondensed",
	750:  "Condensed",
	875:  "SemiCondensed",
	1000: "Normal",
	1125: "SemiExpanded",
	1250: "Expanded",
	1500: "ExtraExpanded",
	2000: "UltraExpanded",
}

func (sfr *SVGFontRoot) MakeGlyphMapping() {
	for _i, _g := range sfr.Defs.Font.Glyph {
		sfr.GlyphMap[_g.GlyphName] = _i
	}
}

func (sfr *SVGFontRoot) MakeUniMapping() {
	for _i, _g := range sfr.Defs.Font.Glyph {
		_r := []rune(_g.Unicode)
		if len(_r) == 1 {
			sfr.UniMap[int(_r[0])] = _i
		} else if len(_r) == 0 {
			sfr.UniMap[0] = _i
		}
	}
}

func (sfr *SVGFontRoot) MakeGlyphCode() {
	_upemi, _ := strconv.ParseUint(sfr.Defs.Font.FontFace.UnitsPerEm, 10, 32)
	_upem := 1000. / float64(_upemi)
	for _, _g := range sfr.Defs.Font.Glyph {
		_h := NewPdfType3PathHandler(_upem)
		ParseDAndEmit(_g.D, _h)
		sfr.GlyphCode[_g.GlyphName] = _h.ToString()
	}
}

func ReadSvgGzFile(filename string) (*SVGFontRoot, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	var dec *xml.Decoder
	if strings.HasSuffix(filename, ".gz") {
		fz, err := gzip.NewReader(fi)
		if err != nil {
			return nil, err
		}
		defer fz.Close()
		dec = xml.NewDecoder(fz)
	} else {
		dec = xml.NewDecoder(fi)
	}

	doc := &SVGFontRoot{
		GlyphMap:  make(map[string]int),
		GlyphCode: make(map[string]string),
		UniMap:    make(map[int]int),
	}
	if err := dec.Decode(doc); err != nil {
		return nil, err
	}
	doc.MakeGlyphMapping()
	doc.MakeUniMapping()
	doc.MakeGlyphCode()
	return doc, nil
}

func ParseDAndEmit(path string, h PathHandler) {
	_tok := tokener.NewTokener(strings.NewReader(path))
	_tok.ResetSyntax()
	_tok.SetParseNumbers()
	_tok.SetWsChars(0, 32)
	_tok.ResetChar('-')
	_tok.ResetChar('+')
	_tok.SetCustomChar(42, 'z')
	_tok.SetCustomChar(42, 'Z')
	_tok.SetCustomChar(42, 'm')
	_tok.SetCustomChar(42, 'M')
	_tok.SetCustomChar(42, 'l')
	_tok.SetCustomChar(42, 'L')
	_tok.SetCustomChar(42, 'h')
	_tok.SetCustomChar(42, 'H')
	_tok.SetCustomChar(42, 'v')
	_tok.SetCustomChar(42, 'V')
	_tok.SetCustomChar(42, 'c')
	_tok.SetCustomChar(42, 'C')
	_tok.SetCustomChar(42, 'q')
	_tok.SetCustomChar(42, 'Q')
	_tok.SetCustomChar(42, 's')
	_tok.SetCustomChar(42, 'S')
	_tok.SetCustomChar(42, 't')
	_tok.SetCustomChar(42, 'T')
	_tok.SetCustomChar(42, 'a')
	_tok.SetCustomChar(42, 'A')

	for ParseAndEmitCommand(_tok, h) {
	}
	h.ClosePath()
}

func ParseAndEmitCommand(_tok *tokener.CustomTokener, h PathHandler) bool {
	_t, _ := _tok.NextToken()
	if _t == 42 {
		switch _tok.CharacterValue {
		case 'z':
			h.ClosePath()
		case 'Z':
			h.ClosePath()
		case 'm':
			h.MoveToRel(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'M':
			h.MoveToAbs(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'l':
			h.LineToRel(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'L':
			h.LineToAbs(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'h':
			h.LineToHorzizontalRel(ParseNextFloat(_tok))
		case 'H':
			h.LineToHorzizontalAbs(ParseNextFloat(_tok))
		case 'v':
			h.LineToVerticalRel(ParseNextFloat(_tok))
		case 'V':
			h.LineToVerticalAbs(ParseNextFloat(_tok))
		case 'c':
			h.CurveToCubicRel(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'C':
			h.CurveToCubicAbs(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'q':
			h.CurveToQuadRel(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'Q':
			h.CurveToQuadAbs(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 's':
			h.CurveToCubicSmoothRel(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'S':
			h.CurveToCubicSmoothAbs(ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 't':
			h.CurveToQuadSmoothRel(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'T':
			h.CurveToQuadSmoothAbs(ParseNextFloat(_tok), ParseNextFloat(_tok))
		case 'a':
			return false
		case 'A':
			return false
		}
		return true
	}
	return false
}

func ParseNextFloat(_tok *tokener.CustomTokener) float64 {
	_t, _ := _tok.NextToken()
	isMinus := false
	if _t == tokener.TOKEN_TYPE_UNKNOWN && _tok.CharacterValue == '-' {
		isMinus = true
		_t, _ = _tok.NextToken()
	}

	if _t == tokener.TOKEN_TYPE_CARDINAL {
		if isMinus {
			return -float64(_tok.CardinalValue)
		} else {
			return float64(_tok.CardinalValue)
		}
	} else if _t == tokener.TOKEN_TYPE_NUMBER {
		if isMinus {
			return -_tok.NumericValue
		} else {
			return _tok.NumericValue
		}
	}
	return 0
}

func _NewPdfSvgFont(fn string, pdf *PdfDoc) PdfFont {
	_svg, _ := ReadSvgGzFile(fn)
	obj := NewPdfType3FontObject(_svg.Defs.Font.FontFace.FontFamily, _svg.Defs.Font.FontFace.FontWeight, _svg.Defs.Font.FontFace.FontStretch)
	pdf.AddObject(obj)
	charProcs := pdf.NewDictObject()
	obj.Dict().Set("CharProcs", charProcs.Ref())
	_upem, _ := strconv.ParseUint(_svg.Defs.Font.FontFace.UnitsPerEm, 10, 32)
	_wm, _ := strconv.ParseUint(_svg.Defs.Font.MissingGlyph.HorizAdvX, 10, 32)
	for _, _n := range PDFDOC_ENCODING_CODEPOINTS {
		_i, _ok := _svg.UniMap[_n]
		if _ok {
			_g := _svg.Defs.Font.Glyph[_i]
			charProc := pdf.NewDictStreamObject()
			charProcs.Dict().Set(_g.GlyphName, charProc.Ref())
			_w, _err := strconv.ParseUint(_g.HorizAdvX, 10, 32)
			if _err != nil {
				_w = _wm
			}
			_w = 1000 * _w / _upem
			_s := fmt.Sprintf("%d 0 -1000 2000 -1000 2000 d1\nq\n%s f\nQ", _w, _svg.GlyphCode[_g.GlyphName])
			charProc.SetStringStream(_s)
		}
	}
	obj.Dict().Set("FirstChar", NewPdfIntValue(0))
	obj.Dict().Set("LastChar", NewPdfIntValue(255))
	_sb := strings.Builder{}
	_sb.WriteString("[ ")
	for _, _n := range PDFDOC_ENCODING_CODEPOINTS {
		_g, _ok := _svg.UniMap[_n]
		if !_ok {
			_sb.WriteString(fmt.Sprintf("%d ", 1000*_wm/_upem))
		} else if _svg.Defs.Font.Glyph[_g].HorizAdvX == "" {
			_sb.WriteString(fmt.Sprintf("%d ", 1000*_wm/_upem))
		} else {
			_w, _ := strconv.ParseUint(_svg.Defs.Font.Glyph[_g].HorizAdvX, 10, 32)
			_sb.WriteString(fmt.Sprintf("%d ", 1000*_w/_upem))
		}
	}
	_sb.WriteString("]")
	obj.Dict().Set("Widths", NewPdfLiteralValue(_sb.String()))
	_sb.Reset()
	_sb.WriteString("<</Type/Encoding /BaseEncoding/WinAnsiEncoding /Differences[ 0 ")
	for _, _n := range PDFDOC_ENCODING_CODEPOINTS {
		_g, _ok := _svg.UniMap[_n]
		if !_ok {
			_sb.WriteString("/.notdef ")
		} else {
			_sb.WriteString(fmt.Sprintf("/%s ", _svg.Defs.Font.Glyph[_g].GlyphName))
		}
	}
	_sb.WriteString("]>>")
	obj.Dict().Set("Encoding", NewPdfLiteralValue(_sb.String()))
	obj.Dict().Set("Encoding", NewPdfLiteralValue(_sb.String()))
	obj.Dict().Set("FontBBox", NewPdfLiteralValue("[ -1000 -1000 2000 2000 ]"))
	obj.Dict().Set("FontMatrix", NewPdfLiteralValue("[ 0.001 0 0 0.001 0 0 ]"))
	obj.MakeMap(PDFDOC_ENCODING_CODEPOINTS)
	return obj
}
