package pdf

import _ "embed"

//go:embed fonts/qsy.otf
var qsyBytes []byte

//go:embed fonts/qdb.otf
var qdbBytes []byte

//go:embed fonts/qhvr.otf
var qhvrBytes []byte

//go:embed fonts/qhvb.otf
var qhvbBytes []byte

//go:embed fonts/qhvz.otf
var qhvzBytes []byte

//go:embed fonts/qhvi.otf
var qhviBytes []byte

//go:embed fonts/qcrr.otf
var qcrrBytes []byte

//go:embed fonts/qcrb.otf
var qcrbBytes []byte

//go:embed fonts/qcrz.otf
var qcrzBytes []byte

//go:embed fonts/qcri.otf
var qcriBytes []byte

//go:embed fonts/qtmr.otf
var qtmrBytes []byte

//go:embed fonts/qtmb.otf
var qtmbBytes []byte

//go:embed fonts/qtmz.otf
var qtmzBytes []byte

//go:embed fonts/qtmi.otf
var qtmiBytes []byte

type EncodedRune struct {
	Char  rune
	Gid   int
	Width int
}

type EncodedWord struct {
	Word  []EncodedRune
	Width int
}

type KerningRune struct {
	Char    rune
	Gid     int
	Width   int
	WAdjust float64
}

type PdfFontWithKerning interface {
	KerningText(s string) []KerningRune
	KerningRunes(s []rune) []KerningRune
	EncodeKerned(rs []KerningRune) string
}

type LayoutRune struct {
	char     rune
	gid      int
	width    int
	xadvance float64
	yadvance float64
	xoffset  float64
	yoffset  float64
}

type PdfFontWithLayout interface {
	LayoutText(s string) []LayoutRune
	LayoutRunes(s []rune) []LayoutRune
	EncodeLayout(rs []LayoutRune) string
}

type PdfFont interface {
	PdfResource
	EncodeText(s string) []EncodedRune
	EncodeWords(s string) []EncodedWord
	IsCid() bool
	MaxCid() int
	//	WidthText(s string) []int
	//	SizeText(s string) int
	//	WidthRunes(s []rune) []int
	//	SizeRunes(s []rune) int
	//	SizeRune(r rune) int
}
