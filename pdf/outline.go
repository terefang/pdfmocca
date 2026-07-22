package pdf

import "fmt"

type PdfOutlinesObject struct {
	PdfBaseObject
	outline []PdfObjRefValue
	last    *PdfOutlineObject
}

func (p *PdfOutlinesObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(p.GetObjNum())
	ioh.PrintDictStart()
	ioh.PrintString("/type/Outlines ")
	_l := len(p.outline)
	if _l > 0 {
		ioh.PrintFmt("/First %s ", p.outline[0].AsString())
		ioh.PrintFmt("/Last %s ", p.outline[_l-1].AsString())
		ioh.PrintFmt("/Count %d ", _l)
	}
	ioh.PrintDictEnd()
	ioh.PrintObjectEnd()
	return nil
}

func (p *PdfOutlinesObject) Add(pdf *PdfDoc, text string, page PdfObjRefValue) *PdfOutlineObject {
	obj := NewPdfOutlineObject(text, p.Ref(), page)
	pdf.AddObject(obj)
	p.outline = append(p.outline, obj.Ref())
	if p.last != nil {
		obj.SetPrev(p.last.Ref())
		p.last.SetNext(obj.Ref())
	}
	p.last = obj
	return obj
}

func NewPdfOutlinesObject() *PdfOutlinesObject {
	obj := &PdfOutlinesObject{outline: make([]PdfObjRefValue, 0), PdfBaseObject: PdfBaseObject{}}
	return obj
}

type PdfOutlineObject struct {
	PdfDictObject
}

func NewPdfOutlineObject(text string, parent PdfObjRefValue, dest PdfObjRefValue) *PdfOutlineObject {
	obj := &PdfOutlineObject{PdfDictObject{object: NewPdfDictValue()}}
	//obj.Dict().Set("Type", NewPdfNameValue("Outline"))
	obj.SetParent(parent)
	obj.SetDest(dest)
	obj.SetTitle(text)
	return obj
}

func (poo *PdfOutlineObject) SetTitle(text string) {
	poo.Dict().Set("Title", NewPdfStringUtf8Value(text))
}

func (poo *PdfOutlineObject) SetDest(page PdfObjRefValue) {
	poo.Dict().Set("Dest", NewPdfLiteralValue(fmt.Sprintf("[%s /Fit]", page.AsString())))
}

func (poo *PdfOutlineObject) SetParent(p PdfObjRefValue) {
	poo.Dict().Set("Parent", p)
}

func (poo *PdfOutlineObject) SetLast(p PdfObjRefValue) {
	poo.Dict().Set("Last", p)
}

func (poo *PdfOutlineObject) SetNext(p PdfObjRefValue) {
	poo.Dict().Set("Next", p)
}

func (poo *PdfOutlineObject) SetPrev(ref PdfObjRefValue) {
	poo.Dict().Set("Prev", ref)
}
