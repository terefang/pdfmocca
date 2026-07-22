package pdf

import (
	"fmt"
)

type PdfFont0Object struct {
	PdfDictObject
	unimap   map[rune]int
	widthMap []int
}

func (pfo *PdfFont0Object) MaxCid() int {
	return 0
}

func (pfo *PdfFont0Object) IsCid() bool {
	return true
}

func (pfo *PdfFont0Object) WidthRunes(s []rune) []int {
	_l := len(s)
	_ret := make([]int, _l)
	for _i := 0; _i < _l; _i++ {
		_ret[_i] = pfo.widthMap[pfo.mapRune(s[_i])]
	}
	return _ret
}

func (pfo *PdfFont0Object) SizeRunes(s []rune) int {
	_l := len(s)
	_w := 0
	for _i := 0; _i < _l; _i++ {
		_w += pfo.widthMap[pfo.mapRune(s[_i])]
	}
	return _w
}

func (pfo *PdfFont0Object) SizeRune(r rune) int {
	return pfo.widthMap[pfo.mapRune(r)]
}

func (pfo *PdfFont0Object) mapRune(_r rune) int {
	_l := len(pfo.widthMap)
	if len(pfo.unimap) != 0 {
		_e, ok := pfo.unimap[_r]
		if ok {
			return _e
		} else {
			return '?'
		}
	} else {
		if int(_r) >= _l {
			return '?'
		} else {
			return int(_r)
		}
	}
}

func (pfo *PdfFont0Object) WidthText(s string) []int {
	return pfo.WidthRunes([]rune(s))
}

func (pfo *PdfFont0Object) SizeText(s string) int {
	return pfo.SizeRunes([]rune(s))
}

func (p PdfFont0Object) ResName() string {
	return fmt.Sprintf("FT%d", p.GetObjNum())
}

func (p *PdfFont0Object) SetUnicodeMapping(r rune, n int) {
	p.unimap[r] = n
}

func (p PdfFont0Object) EncodeText(s string) []EncodedRune {
	runes := []rune(s)
	_l := len(runes)
	_ret := make([]EncodedRune, _l)
	for _i, _rune := range runes {
		_e, _ok := p.unimap[_rune]
		if _ok && (_e < 0x10000) {
			_ret[_i] = EncodedRune{Char: _rune, Gid: _e, Width: p.widthMap[_e]}
		} else {
			_ret[_i] = EncodedRune{Char: _rune, Gid: 0, Width: p.widthMap[0]}
		}
	}
	return _ret
}

func (pfo *PdfFont0Object) EncodeWords(s string) []EncodedWord {
	return DoEncodeWords(pfo, s)
}

func (p *PdfFont0Object) SetWidthMap(w []int) {
	p.widthMap = w
}

func (p *PdfFont0Object) GetCidWidth(i int) int {
	return p.widthMap[i]
}

func (p *PdfFont0Object) GetUniWidth(u int) int {
	_i, _ok := p.unimap[rune(u)]
	if _ok {
		return p.widthMap[_i]
	} else {
		return p.widthMap[0]
	}
}

func NewPdfType0FontObject(ffn string) *PdfFont0Object {
	pfo := &PdfFont0Object{
		unimap: make(map[rune]int),
		PdfDictObject: PdfDictObject{
			PdfBaseObject: PdfBaseObject{
				objnum: -1,
				pdf:    nil},
			object: NewPdfDictValue(),
		},
	}
	pfo.Dict().Set("Type", NewPdfNameValue("Font"))
	pfo.Dict().Set("Subtype", NewPdfNameValue("Type0"))
	pfo.Dict().Set("Encoding", NewPdfNameValue("Identity-H"))
	pfo.Dict().Set("CIDToGIDMap", NewPdfNameValue("Identity"))
	pfo.Dict().Set("Name", NewPdfNameValue(ffn))
	return pfo
}
