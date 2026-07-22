package pdf

import "fmt"

type PdfEGStateObject struct {
	PdfDictObject
}

func (ppo *PdfEGStateObject) ResName() string {
	return fmt.Sprintf("E%d", ppo.GetObjNum())
}

func NewPdfEGStateObject() *PdfEGStateObject {
	obj := PdfEGStateObject{PdfDictObject{object: NewPdfDictValue()}}
	obj.Dict().Set("Type", NewPdfNameValue("ExtGState"))
	return &obj
}

func (ppo *PdfEGStateObject) SetFillAlpha(v float64) {
	ppo.Dict().Set("ca", NewPdfLiteralValue(fmt.Sprintf("%0.4f", v)))
}

func (ppo *PdfEGStateObject) SetStrokeAlpha(v float64) {
	ppo.Dict().Set("CA", NewPdfLiteralValue(fmt.Sprintf("%0.4f", v)))
}
