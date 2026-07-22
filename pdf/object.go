package pdf

type PdfObject interface {
	Ref() PdfObjRefValue
	GetObjNum() int
	SetObjNum(n int, pdf *PdfDoc)
	StreamOut(ioh *IoHelper) error
}

type PdfBaseObject struct {
	offset int64
	objnum int
	pdf    *PdfDoc
}

func (obj PdfBaseObject) Ref() PdfObjRefValue { return NewPdfObjRefValue(obj.objnum) }

func (obj PdfBaseObject) GetObjNum() int { return obj.objnum }

func (obj *PdfBaseObject) SetObjNum(n int, pdf *PdfDoc) {
	obj.objnum = n
	obj.pdf = pdf
}

func (obj *PdfBaseObject) ObjId() int {
	return obj.GetObjNum()
}

type PdfDictObject struct {
	PdfBaseObject
	object PdfDictValue
}

func NewPdfDictObject() *PdfDictObject {
	return &PdfDictObject{PdfBaseObject: PdfBaseObject{objnum: -1, pdf: nil}, object: NewPdfDictValue()}
}

func (obj PdfDictObject) StreamOut(ioh *IoHelper) error {
	obj.offset = ioh.offset
	ioh.PrintObjectStart(obj.objnum)
	obj.object.StreamOut(ioh)
	ioh.PrintObjectEnd()
	return nil
}

func (obj PdfDictObject) Dict() *PdfDictValue { return &(obj.object) }

type PdfDictStreamObject struct {
	PdfDictObject
	stream PdfStream
}

func NewPdfDictStreamObject() *PdfDictStreamObject {
	return &PdfDictStreamObject{PdfDictObject: PdfDictObject{PdfBaseObject: PdfBaseObject{objnum: -1, pdf: nil}, object: NewPdfDictValue()}}
}

func (obj PdfDictStreamObject) StreamOut(ioh *IoHelper) error {
	obj.offset = ioh.offset
	ioh.PrintObjectStart(obj.objnum)
	if obj.stream != nil {
		if obj.stream.Filter() != nil {
			obj.Dict().Set("Filter", obj.stream.Filter())
		}
		obj.Dict().Set("Length", NewPdfIntValue(obj.stream.Length()))
	}
	obj.object.StreamOut(ioh)
	if obj.stream != nil {
		obj.stream.StreamOut(ioh)
	}
	return ioh.PrintObjectEnd()
}

func (obj PdfDictStreamObject) Stream() PdfStream { return obj.stream }

func (obj *PdfDictStreamObject) SetStream(stream PdfStream) {
	obj.stream = stream
}

func (obj *PdfDictStreamObject) SetStringStream(s string) {
	obj.stream = NewPdfStringStream(s)
}

func (obj *PdfDictStreamObject) SetBytesStreamWithFilter(s []byte, l int64, lf int64, f string) {
	obj.stream = NewPdfBytesStreamWithFilter(s, l, lf, f)
}

func (obj *PdfDictStreamObject) SetBytesStream(s []byte, l int64) {
	obj.stream = NewPdfBytesStream(s, l)
}

func (obj *PdfDictStreamObject) SetFlateStream(s []byte) {
	obj.stream = NewPdfFlatedStream(s)
	obj.Dict().Set("Length", NewPdfIntValue(obj.stream.Length()))
}
