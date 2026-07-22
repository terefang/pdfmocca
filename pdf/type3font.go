package pdf

import (
	"fmt"
	"strings"
)

func NewPdfType3FontObject(ffn string, fwn string, fsn string) *PdfFontObject {
	pfo := &PdfFontObject{
		unimap: make(map[rune]int),
		PdfDictObject: PdfDictObject{
			PdfBaseObject: PdfBaseObject{
				objnum: -1,
				pdf:    nil},
			object: NewPdfDictValue(),
		},
	}
	pfo.Dict().Set("Type", NewPdfNameValue("Font"))
	pfo.Dict().Set("Subtype", NewPdfNameValue("Type3"))
	pfo.Dict().Set("Name", NewPdfNameValue(fmt.Sprintf("%s-%s-%s", ffn, fwn, fsn)))
	return pfo
}

type PathHandler interface {
	StartPath()
	EndPath()
	MoveToRel(x float64, y float64)
	MoveToAbs(x float64, y float64)
	ClosePath()
	LineToRel(x float64, y float64)
	LineToAbs(x float64, y float64)
	LineToHorzizontalRel(x float64)
	LineToHorzizontalAbs(x float64)
	LineToVerticalRel(y float64)
	LineToVerticalAbs(y float64)
	CurveToCubicRel(x1 float64, y1 float64, x2 float64, y2 float64, x float64, y float64)
	CurveToCubicAbs(x1 float64, y1 float64, x2 float64, y2 float64, x float64, y float64)
	CurveToCubicSmoothRel(x2 float64, y2 float64, x float64, y float64)
	CurveToCubicSmoothAbs(x2 float64, y2 float64, x float64, y float64)
	CurveToQuadRel(x1 float64, y1 float64, x float64, y float64)
	CurveToQuadAbs(x1 float64, y1 float64, x float64, y float64)
	CurveToQuadSmoothRel(x float64, y float64)
	CurveToQuadSmoothAbs(x float64, y float64)
	ArcRel(rx float64, ry float64, xAxisRot float64, largeArcFlag bool, sweepFlag bool, x float64, y float64)
	ArcAbs(rx float64, ry float64, xAxisRot float64, largeArcFlag bool, sweepFlag bool, x float64, y float64)
}

type PdfType3PathHandler struct {
	stream strings.Builder
	upem   float64
	lx     float64
	ly     float64
	cx     float64
	cy     float64
}

func NewPdfType3PathHandler(upem float64) *PdfType3PathHandler {
	return &PdfType3PathHandler{upem: upem, stream: strings.Builder{}}
}

func (p *PdfType3PathHandler) Reset() {
	p.stream.Reset()
}

func (p *PdfType3PathHandler) ToString() string {
	return p.stream.String()
}

func (p *PdfType3PathHandler) StartPath() {
	p.lx = 0
	p.ly = 0
	p.cx = 0
	p.cy = 0
	p.stream.WriteString("n ")
}

func (p *PdfType3PathHandler) EndPath() {
	p.stream.WriteString("f*")
}

func (p *PdfType3PathHandler) MoveToRel(x float64, y float64) {
	p.MoveToAbs(p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) MoveToAbs(x float64, y float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d m ", int(x*p.upem), int(y*p.upem)))
	p.cx = x
	p.lx = x
	p.cy = y
	p.ly = y
}

func (p *PdfType3PathHandler) ClosePath() {
	p.stream.WriteString("h ")
}

func (p *PdfType3PathHandler) LineToRel(x float64, y float64) {
	p.LineToAbs(p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) LineToAbs(x float64, y float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d l ", int(x*p.upem), int(y*p.upem)))
	p.cx = x
	p.lx = x
	p.cy = y
	p.ly = y
}

func (p *PdfType3PathHandler) LineToHorzizontalRel(x float64) {
	p.LineToHorzizontalAbs(p.lx + x)
}

func (p *PdfType3PathHandler) LineToHorzizontalAbs(x float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d l ", int(x*p.upem), int(p.ly*p.upem)))
	p.cx = x
	p.lx = x
}

func (p *PdfType3PathHandler) LineToVerticalRel(y float64) {
	p.LineToVerticalAbs(p.ly + y)
}

func (p *PdfType3PathHandler) LineToVerticalAbs(y float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d l ", int(p.lx*p.upem), int(y*p.upem)))
	p.cy = y
	p.ly = y
}

func (p *PdfType3PathHandler) CurveToCubicRel(x1 float64, y1 float64, x2 float64, y2 float64, x float64, y float64) {
	p.CurveToCubicAbs(p.lx+x1, p.ly+x2, p.lx+x, p.ly+y, p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) CurveToCubicAbs(x1 float64, y1 float64, x2 float64, y2 float64, x float64, y float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d %d %d %d %d c ", int(x1*p.upem), int(y1*p.upem), int(x2*p.upem), int(y2*p.upem), int(x*p.upem), int(y*p.upem)))
	p.lx = x
	p.ly = y
	p.cx = x2
	p.cy = y2
}

func (p *PdfType3PathHandler) CurveToCubicSmoothRel(x2 float64, y2 float64, x float64, y float64) {
	p.CurveToCubicSmoothAbs(p.lx+x2, p.ly+y2, p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) CurveToCubicSmoothAbs(x2 float64, y2 float64, x float64, y float64) {
	p.CurveToCubicAbs(p.lx*2-p.cx, p.ly*2-p.cy, x2, y2, x, y)
}

func (p *PdfType3PathHandler) CurveToQuadRel(x1 float64, y1 float64, x float64, y float64) {
	p.CurveToQuadAbs(p.lx+x1, p.ly+y1, p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) CurveToQuadAbs(x1 float64, y1 float64, x float64, y float64) {
	p.stream.WriteString(fmt.Sprintf("%d %d %d %d v ", int(x1*p.upem), int(y1*p.upem), int(x*p.upem), int(y*p.upem)))
	p.lx = x
	p.ly = y
	p.cx = x1
	p.cy = y1
}

func (p *PdfType3PathHandler) CurveToQuadSmoothRel(x float64, y float64) {
	p.CurveToQuadSmoothAbs(p.lx+x, p.ly+y)
}

func (p *PdfType3PathHandler) CurveToQuadSmoothAbs(x float64, y float64) {
	p.CurveToQuadAbs(p.lx*2-p.cx, p.ly*2-p.cy, x, y)
}

func (p *PdfType3PathHandler) ArcRel(rx float64, ry float64, xAxisRot float64, largeArcFlag bool, sweepFlag bool, x float64, y float64) {
	// IGNORE
}

func (p *PdfType3PathHandler) ArcAbs(rx float64, ry float64, xAxisRot float64, largeArcFlag bool, sweepFlag bool, x float64, y float64) {
	// IGNORE
}
