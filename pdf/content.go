package pdf

import (
	"fmt"
	"math"
	"strings"
)

type PdfContentObject struct {
	PdfDictStreamObject
	page          *PdfPageObject
	streamBuilder *strings.Builder
	font          PdfFont
	fontSize      float64
	px            float64
	py            float64
	pmx           float64
	pmy           float64
}

func (page *PdfPageObject) NewContent() *PdfContentObject {
	pco := &PdfContentObject{
		streamBuilder: new(strings.Builder),
		page:          page,
		px:            0, py: 0, pmx: 0, pmy: 0,
		PdfDictStreamObject: PdfDictStreamObject{
			PdfDictObject: PdfDictObject{
				PdfBaseObject: PdfBaseObject{
					objnum: -1,
					pdf:    nil,
				}, object: NewPdfDictValue()}}}
	page.pdf.AddObject(pco)
	page.Add(pco)
	return pco
}

func (pco *PdfContentObject) Font() PdfFont     { return pco.font }
func (pco *PdfContentObject) FontSize() float64 { return pco.fontSize }
func (pco *PdfContentObject) StreamOut(ioh *IoHelper) error {
	_b := []byte(pco.streamBuilder.String())
	//pco.SetBytesStream(_b, int64(len(_b)))
	pco.PdfDictStreamObject.SetFlateStream(_b)
	return pco.PdfDictStreamObject.StreamOut(ioh)
}

func (pco *PdfContentObject) UseFont(font PdfFont) {
	pco.font = font
	pco.page.UseFont(font)
}

func (pco *PdfContentObject) AddContent(s ...string) {
	for _, _s := range s {
		pco.streamBuilder.WriteString(_s)
	}
}

func (pco *PdfContentObject) AddContentFmt(s string, v ...any) {
	pco.streamBuilder.WriteString(fmt.Sprintf(s, v...))
}

func (pco *PdfContentObject) StateSave() { pco.AddContent("q ") }

func (pco *PdfContentObject) StateRestore() { pco.AddContent("Q ") }

func (pco *PdfContentObject) SetInt(n int64)     { pco.AddContent(fmt.Sprintf("%d ", n)) }
func (pco *PdfContentObject) SetFloat(n float64) { pco.AddContent(fmt.Sprintf("%f ", n)) }

func (pco *PdfContentObject) BeginText() { pco.AddContent("BT ") }
func (pco *PdfContentObject) EndText()   { pco.AddContent("ET ") }

func (pco *PdfContentObject) SetText(s string) {
	_en := pco.Font().EncodeText(s)
	if pco.Font().IsCid() {
		pco.AddContent("<", MakePdfRunesCidString(_en), ">")
	} else {
		pco.AddContent("(", MakePdfRunesString(_en), ")")
	}
}

func (pco *PdfContentObject) SetCharCid(r int) {
	if pco.Font().IsCid() {
		pco.AddContentFmt("<%04x>", r)
	} else {
		pco.AddContentFmt("(\\%03o)", r)
	}
}

func (pco *PdfContentObject) SetTextShow()    { pco.AddContent("Tj ") }
func (pco *PdfContentObject) StartTextShowN() { pco.AddContent("[ ") }
func (pco *PdfContentObject) EndTextShowN()   { pco.AddContent("] TJ ") }

func (pco *PdfContentObject) SetFillBlack()   { pco.AddContent("0 g ") }
func (pco *PdfContentObject) SetStrokeBlack() { pco.AddContent("0 G ") }

func (pco *PdfContentObject) SetFillWhite()   { pco.AddContent("1 g ") }
func (pco *PdfContentObject) SetStrokeWhite() { pco.AddContent("1 G ") }

func (pco *PdfContentObject) SetFillColorGrey(f float64)   { pco.AddContentFmt("%f g ", f) }
func (pco *PdfContentObject) SetStrokeColorGrey(f float64) { pco.AddContentFmt("%f G ", f) }

func (pco *PdfContentObject) SetFillColorRgb(r float64, g float64, b float64) {
	pco.AddContentFmt("%f %f %f rg ", r, g, b)
}
func (pco *PdfContentObject) SetStrokeColorRgb(r float64, g float64, b float64) {
	pco.AddContentFmt("%f %f %f RG ", r, g, b)
}

func (pco *PdfContentObject) SetFillColorCmyk(c float64, m float64, y float64, k float64) {
	pco.AddContentFmt("%f %f %f %f k ", c, m, y, k)
}
func (pco *PdfContentObject) SetStrokeColorCmyk(c float64, m float64, y float64, k float64) {
	pco.AddContentFmt("%f %f %f %f K ", c, m, y, k)
}

func (pco *PdfContentObject) SetTextMove(x int, y int) {
	pco.AddContentFmt("%d %d Td ", x, y)
}
func (pco *PdfContentObject) SetTextMoveF(x float64, y float64) {
	pco.AddContentFmt("%f %f Td ", x, y)
}

func (pco *PdfContentObject) SetTextFont(fn PdfFont, fs int) {
	pco.UseFont(fn)
	pco.fontSize = float64(fs)
	pco.AddContentFmt("/%s %d Tf ", fn.ResName(), fs)
}
func (pco *PdfContentObject) SetTextFontF(fn PdfFont, fs float64) {
	pco.UseFont(fn)
	pco.fontSize = fs
	pco.AddContentFmt("/%s %f Tf ", fn.ResName(), fs)
}

func (pco *PdfContentObject) SetTextNextline() {
	pco.AddContent("T* ")
}

func (pco *PdfContentObject) SetTextLeading(n int) {
	pco.AddContentFmt("%d TL ", n)
}
func (pco *PdfContentObject) SetTextLeadingF(n float64) {
	pco.AddContentFmt("%f TL ", n)
}

func (pco *PdfContentObject) SetTextMatrix(a int, b int, c int, d int, e int, f int) {
	pco.AddContentFmt("%d %d %d %d %d %d Tm ", a, b, c, d, e, f)
}
func (pco *PdfContentObject) SetTextMatrixF(a float64, b float64, c float64, d float64, e float64, f float64) {
	pco.AddContentFmt("%f %f %f %f %f %f Tm ", a, b, c, d, e, f)
}

func (pco *PdfContentObject) SetTextRendering(r int) {
	pco.AddContentFmt("%d Tr ", r)
}

func (pco *PdfContentObject) SetTextRise(n int) {
	pco.AddContentFmt("%d Ts ", n)
}
func (pco *PdfContentObject) SetTextRiseF(n float64) {
	pco.AddContentFmt("%f Ts ", n)
}

func (pco *PdfContentObject) SetTextHorizonScale(n int) {
	pco.AddContentFmt("%d Tz ", n)
}
func (pco *PdfContentObject) SetTextHorizonScaleF(n float64) {
	pco.AddContentFmt("%f Tz ", n)
}

func (pco *PdfContentObject) SetTextWordSpacing(n int) {
	pco.AddContentFmt("%d Tw ", n)
}
func (pco *PdfContentObject) SetTextWordSpacingF(n float64) {
	pco.AddContentFmt("%f Tw ", n)
}

func (pco *PdfContentObject) SetTextCharSpacing(n int) {
	pco.AddContentFmt("%d Tc ", n)
}
func (pco *PdfContentObject) SetTextCharSpacingF(n float64) {
	pco.AddContentFmt("%f Tc ", n)
}

func (pco *PdfContentObject) PutText(x int, y int, fn PdfFont, fs int, cl string, text string) {
	pco.StateSave()
	pco.SetStrokeWhite()
	pco.SetFillColor(cl)
	pco.BeginText()
	pco.SetTextFont(fn, fs)
	pco.SetTextMove(x, y)
	pco.StartTextShowN()
	pco.SetText(text)
	pco.EndTextShowN()
	pco.EndText()
	pco.StateRestore()
}

func (pco *PdfContentObject) PutCharCid(x int, y int, fn PdfFont, fs int, cl string, r int) {
	pco.StateSave()
	pco.SetStrokeWhite()
	pco.SetFillColor(cl)
	pco.BeginText()
	pco.SetTextFont(fn, fs)
	pco.SetTextMove(x, y)
	pco.StartTextShowN()
	pco.SetCharCid(r)
	pco.EndTextShowN()
	pco.EndText()
	pco.StateRestore()
}

func (pco *PdfContentObject) StartLayer(n string) {
	_l := pco.AddLayer(n)
	pco.AddContentFmt("q /OC /%s BDC q ", _l.ResName())
}
func (pco *PdfContentObject) EndLayer() {
	pco.AddContent("Q EMC Q ")
}
func (pco *PdfContentObject) AddLayer(n string) PdfResource {
	return pco.page.AddLayer(n)
}

func (pco *PdfContentObject) SetMatrix(gOrT bool, a int, b int, c int, d int, e int, f int) {
	if gOrT {
		pco.AddContentFmt("%d %d %d %d %d %d cm", a, b, c, d, e, f)
	} else {
		pco.AddContentFmt("%d %d %d %d %d %d Tm", a, b, c, d, e, f)
	}
}

func (pco *PdfContentObject) SetMatrixF(gOrT bool, a float64, b float64, c float64, d float64, e float64, f float64) {
	if gOrT {
		pco.AddContentFmt("%0.4f %0.4f %0.4f %0.4f %0.4f %0.4f cm ", a, b, c, d, e, f)
	} else {
		pco.AddContentFmt("%0.4f %0.4f %0.4f %0.4f %0.4f %0.4f Tm ", a, b, c, d, e, f)
	}
}

func (pco *PdfContentObject) SetLinecap(n int) {
	pco.AddContentFmt("%d J ", n)
}

func (pco *PdfContentObject) SetLinejoin(n int) {
	pco.AddContentFmt("%d j ", n)
}

func (pco *PdfContentObject) SetFlatness(n int) {
	pco.AddContentFmt("%d i ", n)
}

func (pco *PdfContentObject) SetEGState(egs PdfResource) {
	pco.page.UseEGState(egs)
	pco.AddContentFmt("/%s gs", egs.ResName())
}

func (pco *PdfContentObject) SetLinewidth(w int) {
	pco.AddContentFmt("%d w ", w)
}

func (pco *PdfContentObject) SetLinewidthF(w float64) {
	pco.AddContentFmt("%0.4f w ", w)
}

func (pco *PdfContentObject) SetMeterlimit(w int) {
	pco.AddContentFmt("%d M ", w)
}

func (pco *PdfContentObject) SetMeterlimitF(w float64) {
	pco.AddContentFmt("%0.4f m ", w)
}

func (pco *PdfContentObject) SetLinedashNone() {
	pco.AddContent("[ ] 0 d ")
}

func (pco *PdfContentObject) SetLinedashOnOff(n int) {
	pco.AddContentFmt("[ %d ] 0 d ", n)
}

func (pco *PdfContentObject) SetLinedashOnOffF(n float64) {
	pco.AddContentFmt("[ %0.4f ] 0 d ", n)
}

func (pco *PdfContentObject) SetLinedashWithOffset(ofs int, w ...int) {
	_sb := strings.Builder{}
	for _, x := range w {
		_sb.WriteString(fmt.Sprintf("%d ", x))
	}
	pco.AddContentFmt("[ %s ] %d d ", _sb.String(), ofs)
}

func (pco *PdfContentObject) SetLinedashWithOffsetF(ofs float64, w ...float64) {
	_sb := strings.Builder{}
	for _, x := range w {
		_sb.WriteString(fmt.Sprintf("%0.4f ", x))
	}
	pco.AddContentFmt("[ %s ] %0.4f d ", _sb.String(), ofs)
}

func (pco *PdfContentObject) MoveTo(x int, y int) {
	pco.AddContentFmt("%d %d m ", x, y)
	pco.pmx, pco.pmy = float64(x), float64(y)
}

func (pco *PdfContentObject) MoveToF(x float64, y float64) {
	pco.AddContentFmt("%0.4f %0.4f m ", x, y)
	pco.pmx, pco.pmy = x, y
}

func (pco *PdfContentObject) LineTo(x int, y int) {
	pco.AddContentFmt("%d %d l ", x, y)
	pco.px, pco.py = float64(x), float64(y)
}

func (pco *PdfContentObject) LineToF(x float64, y float64) {
	pco.AddContentFmt("%0.4f %0.4f l ", x, y)
	pco.px, pco.py = x, y
}

func (pco *PdfContentObject) CurveTo(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
	pco.AddContentFmt("%d %d %d %d %d %d c ", x1, y1, x2, y2, x3, y3)
	pco.px, pco.py = float64(x3), float64(y3)
}

func (pco *PdfContentObject) CurveToF(x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64) {
	pco.AddContentFmt("%0.4f %0.4f %0.4f %0.4f %0.4f %0.4f c ", x1, y1, x2, y2, x3, y3)
	pco.px, pco.py = x3, y3
}

func (pco *PdfContentObject) Rectangle(x1 int, y1 int, x2 int, y2 int) {
	pco.AddContentFmt("%d %d %d %d re ", x1, y1, x2-x1, y2-y1)
}

func (pco *PdfContentObject) RectangleF(x1 float64, y1 float64, x2 float64, y2 float64) {
	pco.AddContentFmt("%0.4f %0.4f %0.4f %0.4f re ", x1, y1, x2-x1, y2-y1)
}

func (pco *PdfContentObject) RectangleHw(x1 int, y1 int, w int, h int) {
	pco.AddContentFmt("%d %d %d %d re ", x1, y1, w, h)
}

func (pco *PdfContentObject) RectangleHwF(x1 float64, y1 float64, w float64, h float64) {
	pco.AddContentFmt("%0.4f %0.4f %0.4f %0.4f re ", x1, y1, w, h)
}

// TODO pie circle ellipse arc

func (pco *PdfContentObject) ClosePath() {
	pco.AddContent("h ")
	pco.px, pco.py = pco.pmx, pco.pmy
}

func (pco *PdfContentObject) EndPath() {
	pco.AddContent("n ")
}

func (pco *PdfContentObject) PolyLine(_p ...int) {
	pco.MoveTo(_p[0], _p[1])
	_l := len(_p)
	for _i := 3; _i < _l; _i += 2 {
		pco.LineTo(_p[_i-1], _p[_i])
	}
}

func (pco *PdfContentObject) PolyLineF(_p ...float64) {
	pco.MoveToF(_p[0], _p[1])
	_l := len(_p)
	for _i := 3; _i < _l; _i += 2 {
		pco.LineToF(_p[_i-1], _p[_i])
	}
}

func (pco *PdfContentObject) Spline(_p ...int) {
	_l := len(_p)
	for _i := 3; _i < _l; _i += 2 {
		_cx := float64(_p[_i-3])
		_cy := float64(_p[_i-2])
		_x := float64(_p[_i-1])
		_y := float64(_p[_i])
		_c1x := (2.*_cx + pco.px) / 3.
		_c1y := (2.*_cy + pco.py) / 3.
		_c2x := (2.*_cx + _x) / 3.
		_c2y := (2.*_cy + _y) / 3.
		pco.CurveToF(_c1x, _c1y, _c2x, _c2y, _x, _y)
	}
}

func (pco *PdfContentObject) SplineF(_p ...float64) {
	_l := len(_p)
	for _i := 3; _i < _l; _i += 2 {
		_cx := _p[_i-3]
		_cy := _p[_i-2]
		_x := _p[_i-1]
		_y := _p[_i]
		_c1x := (2.*_cx + pco.px) / 3.
		_c1y := (2.*_cy + pco.py) / 3.
		_c2x := (2.*_cx + _x) / 3.
		_c2y := (2.*_cy + _y) / 3.
		pco.CurveToF(_c1x, _c1y, _c2x, _c2y, _x, _y)
	}
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180.
}

func _arcToCurve(_a float64, _b float64, _alpha float64, _beta float64) (_ret []float64) {
	if math.Abs(_beta-_alpha) > 15. {
		_part1 := _arcToCurve(_a, _b, _alpha, (_beta+_alpha)/2.)
		_part2 := _arcToCurve(_a, _b, (_beta+_alpha)/2., _beta)
		_l1 := len(_part1)
		_l2 := len(_part2)
		_ret = make([]float64, _l1+_l2)
		copy(_ret, _part1)
		copy(_ret[_l1:], _part2)
	} else {
		_alpha = degToRad(_alpha)
		_beta = degToRad(_beta)

		_bcp := (4.0 / 3 * (1 - math.Cos((_beta-_alpha)/2)) / math.Sin((_beta-_alpha)/2))
		_sin_alpha := math.Sin(_alpha)
		_sin_beta := math.Sin(_beta)
		_cos_alpha := math.Cos(_alpha)
		_cos_beta := math.Cos(_beta)

		_p0_x := _a * _cos_alpha
		_p0_y := _b * _sin_alpha
		_p1_x := _a * (_cos_alpha - _bcp*_sin_alpha)
		_p1_y := _b * (_sin_alpha + _bcp*_cos_alpha)
		_p2_x := _a * (_cos_beta + _bcp*_sin_beta)
		_p2_y := _b * (_sin_beta - _bcp*_cos_beta)
		_p3_x := _a * _cos_beta
		_p3_y := _b * _sin_beta
		_ret = []float64{_p0_x, _p0_y, _p1_x, _p1_y, _p2_x, _p2_y, _p3_x, _p3_y}
	}
	return
}

func (pco *PdfContentObject) Arc(x int, y int, a int, b int, alpha float64, beta float64, move bool) {
	pco.ArcF(float64(x), float64(y), float64(a), float64(b), alpha, beta, move)
}
func (pco *PdfContentObject) ArcF(x float64, y float64, a float64, b float64, alpha float64, beta float64, move bool) {
	_pts := _arcToCurve(a, b, alpha, beta)
	if move {
		pco.MoveToF(x+_pts[0], y+_pts[1])
	}
	_lp := len(_pts)
	for _j := 7; _j < _lp; _j += 8 {
		pco.CurveToF(_pts[_j-5], _pts[_j-4], _pts[_j-3], _pts[_j-2], _pts[_j-1], _pts[_j])
	}
}

func (pco *PdfContentObject) Ellipse(x int, y int, a int, b int) {
	pco.Arc(x, y, a, b, 0, 360., true)
	pco.ClosePath()
}

func (pco *PdfContentObject) EllipseF(x float64, y float64, a float64, b float64) {
	pco.ArcF(x, y, a, b, 0, 360., true)
	pco.ClosePath()
}

func (pco *PdfContentObject) Circle(x int, y int, r int) {
	pco.Arc(x, y, r, r, 0, 360., true)
	pco.ClosePath()
}

func (pco *PdfContentObject) CircleF(x float64, y float64, r float64) {
	pco.ArcF(x, y, r, r, 0, 360., true)
	pco.ClosePath()
}

func (pco *PdfContentObject) Pie(x int, y int, a int, b int, alpha float64, beta float64) {
	pco.PieF(float64(x), float64(y), float64(a), float64(b), alpha, beta)
}
func (pco *PdfContentObject) PieF(x float64, y float64, a float64, b float64, alpha float64, beta float64) {
	_pts := _arcToCurve(a, b, alpha, beta)
	pco.MoveToF(x, y)
	pco.LineToF(x+_pts[0], y+_pts[1])
	pco.ArcF(x, y, a, b, alpha, beta, false)
	pco.ClosePath()
}

func (pco *PdfContentObject) Stroke() {
	pco.AddContent("S ")
}

func (pco *PdfContentObject) Fill() {
	pco.AddContent("f ")
}

func (pco *PdfContentObject) Clip() {
	pco.AddContent("W ")
}

func (pco *PdfContentObject) ClipEvenOdd(evenOdd bool) {
	if evenOdd {
		pco.AddContent("W* ")
	} else {
		pco.AddContent("W ")
	}
}

func (pco *PdfContentObject) FillStroke() {
	pco.AddContent("B ")
}

func (pco *PdfContentObject) FillStrokeEvenOdd(evenOdd bool) {
	if evenOdd {
		pco.AddContent("B* ")
	} else {
		pco.AddContent("B ")
	}
}

func (pco *PdfContentObject) FillEvenOdd(evenOdd bool) {
	if evenOdd {
		pco.AddContent("f* ")
	} else {
		pco.AddContent("f ")
	}
}

func (pco *PdfContentObject) SetFillColorSpace(cs PdfResource) {
	pco.page.UseColorSpace(cs)
	pco.AddContentFmt("/%s cs", cs.ResName())
}

func (pco *PdfContentObject) SetStrokeColorSpace(cs PdfResource) {
	pco.page.UseColorSpace(cs)
	pco.AddContentFmt("/%s CS", cs.ResName())
}

func (pco *PdfContentObject) SetFillColorSpaceColor(_nf bool, _v ...float64) {
	for _, v := range _v {
		pco.AddContentFmt("%0.4f ", v)
	}
	if _nf {
		pco.AddContent("scn ")
	} else {
		pco.AddContent("sc ")
	}
}

func (pco *PdfContentObject) SetStrokeColorSpaceColor(_nf bool, _v ...float64) {
	for _, v := range _v {
		pco.AddContentFmt("%0.4f ", v)
	}
	if _nf {
		pco.AddContent("SCN ")
	} else {
		pco.AddContent("SC ")
	}
}

func (pco *PdfContentObject) SetFillColorSpacedColor(_nf bool, cs PdfResource, _v ...float64) {
	pco.SetFillColorSpace(cs)
	pco.SetFillColorSpaceColor(_nf, _v...)
}

func (pco *PdfContentObject) SetStrokeColorSpacedColor(_nf bool, cs PdfResource, _v ...float64) {
	pco.SetStrokeColorSpace(cs)
	pco.SetStrokeColorSpaceColor(_nf, _v...)
}
