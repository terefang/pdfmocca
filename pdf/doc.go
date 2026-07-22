package pdf

import (
	"fmt"
	"os"
	"time"
)

type PdfInfoObject struct {
	objnum   int
	subject  string
	creator  string
	author   string
	producer string
	title    string
	keywords string
}

func (p *PdfInfoObject) SetSubject(subject string) {
	p.subject = subject
}

func (p *PdfInfoObject) SetCreator(creator string) {
	p.creator = creator
}

func (p *PdfInfoObject) SetAuthor(author string) {
	p.author = author
}

func (p *PdfInfoObject) SetProducer(producer string) {
	p.producer = producer
}

func (p *PdfInfoObject) SetTitle(title string) {
	p.title = title
}

func (p *PdfInfoObject) SetKeywords(keywords string) {
	p.keywords = keywords
}

func NewPdfInfoObject() *PdfInfoObject {
	return &PdfInfoObject{}
}

func (p PdfInfoObject) Ref() PdfObjRefValue {
	return NewPdfObjRefValue(p.objnum)
}

func (p PdfInfoObject) GetObjNum() int {
	return p.objnum
}

func (p *PdfInfoObject) SetObjNum(n int, pdf *PdfDoc) {
	p.objnum = n
}

func (p PdfInfoObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(p.objnum)
	ioh.PrintDictStart()
	//ioh.PrintString("/Type /Info")
	ioh.PrintFmt("/Producer %s ", MakePdfStringUtf8WithBrackets(p.producer))
	ioh.PrintFmt("/Creator %s ", MakePdfStringUtf8WithBrackets(p.creator))
	ioh.PrintFmt("/Author %s ", MakePdfStringUtf8WithBrackets(p.author))
	ioh.PrintFmt("/Subject %s ", MakePdfStringUtf8WithBrackets(p.subject))
	ioh.PrintFmt("/Title %s ", MakePdfStringUtf8WithBrackets(p.title))
	ioh.PrintFmt("/Keywords %s ", MakePdfStringUtf8WithBrackets(p.keywords))
	_t := NewPdfTimeString(time.Now())
	ioh.PrintFmt("/ModDate %s ", _t.AsString())
	ioh.PrintFmt("/CreationDate %s ", _t.AsString())
	ioh.PrintDictEnd()
	return ioh.PrintObjectEnd()
}

type PdfDoc struct {
	version string

	objNum   int
	objects  []PdfObject
	offsets  []int64
	root     *PdfCatalogObject
	info     *PdfInfoObject
	pageTree *PdfPagesObject

	pageFront *PdfPagesObject
	pageCore  *PdfPagesObject
	pageBack  *PdfPagesObject
	layers    map[string]PdfResource
	/*

		layers map[string]PdfOptionalContentGroup

		outlines PdfOutlines
	*/
	ioh      *IoHelper
	ocgroups *PdfOptionalContentGroupObject
	outlines *PdfOutlinesObject
	embeds   []PdfResource
}

func NewPdfDocument(v string) *PdfDoc {
	pdf := &PdfDoc{
		version:  v,
		objNum:   0,
		objects:  make([]PdfObject, 0),
		offsets:  make([]int64, 0),
		layers:   make(map[string]PdfResource),
		embeds:   make([]PdfResource, 0),
		root:     NewPdfCatalogObject(),
		info:     NewPdfInfoObject(),
		pageTree: NewPdfPagesObjectNoParent(),
	}
	pdf.AddObject(pdf.root)
	pdf.AddObject(pdf.info)
	pdf.AddObject(pdf.pageTree)
	pdf.root.SetPages(pdf.pageTree.Ref())

	pdf.pageFront = NewPdfPagesObject(pdf.pageTree.Ref())
	pdf.AddObject(pdf.pageFront)
	pdf.pageTree.Add(pdf.pageFront)

	pdf.pageCore = NewPdfPagesObject(pdf.pageTree.Ref())
	pdf.AddObject(pdf.pageCore)
	pdf.pageTree.Add(pdf.pageCore)

	pdf.pageBack = NewPdfPagesObject(pdf.pageTree.Ref())
	pdf.AddObject(pdf.pageBack)
	pdf.pageTree.Add(pdf.pageBack)

	pdf.ocgroups = NewPdfOptionalContentGroupObject()
	pdf.AddObject(pdf.ocgroups)
	pdf.root.SetProps(pdf.ocgroups.Ref())
	//pdf.root.Dict().Set("Type", NewPdfNameValue("Catalog"))
	//	pdf.info = pdf.NewDictObject()
	//	pdf.info.Dict().Set("Type", NewPdfNameValue("Info"))
	//	pdf.pageTree = pdf.NewPages()
	//	pdf.root.Dict().Set("Pages", pdf.pageTree.obj.Ref())
	pdf.outlines = NewPdfOutlinesObject()
	pdf.AddObject(pdf.outlines)
	pdf.root.SetOutlines(pdf.outlines.Ref())
	return pdf
}

func (pdf *PdfDoc) Info() *PdfInfoObject {
	return pdf.info
}

func (pdf *PdfDoc) AddObject(obj PdfObject) {
	pdf.objNum++
	obj.SetObjNum(pdf.objNum, pdf)
	pdf.objects = append(pdf.objects, obj)
	pdf.offsets = append(pdf.offsets, -1)
}

func (pdf PdfDoc) NewObjNum() int {
	pdf.objNum = pdf.objNum + 1
	return pdf.objNum
}
func (pdf *PdfDoc) NewCoreFont(fn string, fc int, lc int, wl [256]int, el [256]string) PdfFont {
	obj := NewPdfCoreFontObject(pdf, fn, fc, lc, wl, el)
	return obj
}
func (pdf *PdfDoc) _NewSvgFont(f string) PdfFont {
	return _NewPdfSvgFont(f, pdf)
}

func (pdf *PdfDoc) NewTtfFont(s string, mapMode int) PdfFont {
	return NewPdfTTFont(s, pdf, mapMode)
}

func (pdf *PdfDoc) NewEmbedFile(name string, path string) {

	emfs := pdf.NewDictStreamObject()
	//emfs.Dict().Set("Params", NewPdfLiteralValue("<< /CheckSum <%s> /Size %d >>"))
	emfs.SetStringStream(path)

	emds := pdf.NewEmbedFileObject(name)
	emds.Dict().Set("EF", NewPdfLiteralValue(fmt.Sprintf("<< /F %d 0 R >>", emfs.GetObjNum())))

	pdf.root.AddEmbeds(emds)
}

func (pdf *PdfDoc) NewEmbedFileObject(name string) *PdfEmbedFileObject {
	emds := NewEmbedFileObject(pdf, name)
	emds.Dict().Set("Type", NewPdfNameValue("Filespec"))
	emds.Dict().Set("F", NewPdfStringValue(name))
	emds.Dict().Set("UF", NewPdfStringUcsValue(name))
	emds.Dict().Set("Desc", NewPdfStringUtf8Value(name))
	return emds
}

type PdfEmbedFileObject struct {
	PdfDictObject
	name string
}

func (p PdfEmbedFileObject) ResName() string {
	return p.name
}

func NewEmbedFileObject(pdf *PdfDoc, name string) *PdfEmbedFileObject {
	efo := &PdfEmbedFileObject{
		PdfDictObject: PdfDictObject{object: NewPdfDictValue()},
		name:          name,
	}
	pdf.AddObject(efo)
	return efo
}

func (pdf *PdfDoc) NewDictObject() *PdfDictObject {
	obj := NewPdfDictObject()
	pdf.AddObject(obj)
	return obj
}

func (pdf *PdfDoc) NewDictStreamObject() *PdfDictStreamObject {
	obj := NewPdfDictStreamObject()
	pdf.AddObject(obj)
	return obj
}

func (pdf *PdfDoc) AddOutline(text string, dest PdfObjRefValue) {
	pdf.outlines.Add(pdf, text, dest) // already AddObject
}

func (pdf *PdfDoc) GetObject(i int) PdfObject {
	return pdf.objects[i]
}

func (pdf *PdfDoc) NewPage() *PdfPageObject {
	obj := NewPdfPageObject(pdf.pageCore.Ref())
	pdf.AddObject(obj)
	pdf.pageCore.Add(obj)
	return obj
}

func (pdf *PdfDoc) NewEGState() *PdfEGStateObject {
	obj := NewPdfEGStateObject()
	pdf.AddObject(obj)
	return obj
}

func (pdf *PdfDoc) StreamBeginToFile(path string) error {
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	return pdf.StreamBeginTo(fh)
}

func (pdf *PdfDoc) StreamBeginTo(fh *os.File) error {
	pdf.ioh = NewIoHelper(fh)
	_err := pdf.ioh.PrintFmt("%%PDF-%s\n", pdf.version)
	if _err != nil {
		return _err
	}
	_err = pdf.ioh.PrintFmt("%c%c%c%c\n", 0xC5, 0xD6, 0xE7, 0xF8)
	if _err != nil {
		return _err
	}
	_err = pdf.ioh.PrintString("%% --- START ---\n")
	if _err != nil {
		return _err
	}
	return nil
}

func (pdf *PdfDoc) StreamOut(j int) error {
	if pdf.ioh == nil {
		return nil
	}
	_sz := len(pdf.objects)
	pdf.ioh.SyncOffset()
	if j <= _sz {
		_ref := pdf.ioh.Offset()
		if pdf.offsets[j-1] == -1 {
			pdf.offsets[j-1] = _ref
			_err := pdf.objects[j-1].StreamOut(pdf.ioh)
			if _err != nil {
				return _err
			}
		}
	}
	return nil
}

func (pdf *PdfDoc) StreamEnd(tc bool) error {
	// collect unstreamed objects
	_sz := len(pdf.objects)
	pdf.ioh.SyncOffset()
	for j := 0; j < _sz; j++ {
		_ref := pdf.ioh.Offset()
		if pdf.offsets[j] == -1 {
			pdf.offsets[j] = _ref
			_err := pdf.objects[j].StreamOut(pdf.ioh)
			if _err != nil {
				return _err
			}
		}
	}

	_xref := pdf.ioh.Offset()
	_md := pdf.ioh.Finalize()
	pdf.ioh.Flush()

	if true {
		_err := pdf.ioh.PrintObjectStartDict(_sz + 1)
		if _err != nil {
			return _err
		}
		pdf.ioh.PrintFmt("/Type /XRef\n/Size %d ", _sz+1)
		pdf.ioh.PrintFmt("/ID[<%s><%s>]", _md, _md)
		pdf.ioh.PrintFmt("/Root %s ", pdf.root.Ref().AsString())
		pdf.ioh.PrintFmt("/Info %s ", pdf.info.Ref().AsString())
		if _xref < 0xffff {
			pdf.ioh.PrintString("/W [ 1 2 1 ]")
			pdf.ioh.PrintFmt("/Length %d ", (_sz+1)*4)
		} else if _xref < 0xffffffff {
			pdf.ioh.PrintString("/W [ 1 4 1 ]")
			pdf.ioh.PrintFmt("/Length %d ", (_sz+1)*6)
		} else {
			pdf.ioh.PrintString("/W [ 1 8 1 ]")
			pdf.ioh.PrintFmt("/Length %d ", (_sz+1)*10)
		}
		pdf.ioh.PrintDictEnd()
		pdf.ioh.PrintStreamStart()
		if _xref < 0xffff {
			pdf.ioh.Print121(0, 0, -1)
			for j := 0; j < _sz; j++ {
				pdf.ioh.Print121(1, pdf.offsets[j], 0)
			}
		} else if _xref < 0xffffffff {
			pdf.ioh.Print141(0, 0, -1)
			for j := 0; j < _sz; j++ {
				pdf.ioh.Print141(1, pdf.offsets[j], 0)
			}
		} else {
			pdf.ioh.Print141(0, 0, -1)
			for j := 0; j < _sz; j++ {
				pdf.ioh.Print141(1, pdf.offsets[j], 0)
			}
		}
		pdf.ioh.PrintStreamEnd()
		pdf.ioh.PrintObjectEnd()
	}

	_md = pdf.ioh.Finalize()
	pdf.ioh.PrintFmt("%x-pdfmocca:version=%s;hash=%s\n")

	pdf.ioh.PrintFmt("startxref\n%d\n%%%%EOF\n", _xref)

	pdf.ioh.Flush()

	if tc {
		pdf.ioh.Close()
	}
	return nil
}

func (p *PdfDoc) SetSubject(subject string) {
	p.Info().SetSubject(subject)
}

func (p *PdfDoc) SetCreator(creator string) {
	p.Info().SetCreator(creator)
}

func (p *PdfDoc) SetAuthor(author string) {
	p.Info().SetAuthor(author)
}

func (p *PdfDoc) SetProducer(producer string) {
	p.Info().SetProducer(producer)
}

func (p *PdfDoc) SetTitle(title string) {
	p.Info().SetTitle(title)
}

func (p *PdfDoc) SetKeywords(keywords string) {
	p.Info().SetKeywords(keywords)
}

func (pdf *PdfDoc) AddLayer(n string) PdfResource {
	if pdf.layers[n] == nil {
		_l := NewPdfOptionalContentObject(n)
		pdf.AddObject(_l)
		pdf.OCGroups().Add(_l.Ref())
		pdf.layers[n] = _l
	}
	return pdf.layers[n]
}

func (pdf *PdfDoc) OCGroups() *PdfOptionalContentGroupObject {
	return pdf.ocgroups
}
