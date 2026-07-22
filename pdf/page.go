package pdf

import (
	"fmt"
	"slices"
)

// ******************************************************************

type PdfCatalogObject struct {
	objnum   int
	pages    *PdfObjRefValue
	props    *PdfObjRefValue
	intents  []PdfObjRefValue
	outlines *PdfObjRefValue
	embeds   []PdfResource
}

func (p *PdfCatalogObject) AddEmbeds(em PdfResource) {
	p.embeds = append(p.embeds, em)
}

func NewPdfCatalogObject() *PdfCatalogObject {
	return &PdfCatalogObject{
		intents: make([]PdfObjRefValue, 0),
		embeds:  make([]PdfResource, 0),
	}
}

func (p PdfCatalogObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(p.objnum)
	ioh.PrintDictStart()
	ioh.PrintString("/Type /Catalog")
	if p.pages != nil {
		ioh.PrintFmt("/Pages %s ", p.pages.AsString())
	}
	if p.outlines != nil {
		ioh.PrintFmt("/Outlines %s ", p.outlines.AsString())
	}
	if p.props != nil {
		ioh.PrintFmt("/OCProperties %s ", p.props.AsString())
	}
	ioh.PrintString("/OutputIntents [ ")
	for _, intent := range p.intents {
		ioh.PrintFmt(" %s ", intent.AsString())
	}
	ioh.PrintString("] ")
	ioh.PrintString("/EmbeddedFiles << /Names [ ")
	for _, _em := range p.embeds {
		ioh.PrintFmt("%s %d 0 R ", NewPdfStringValue(_em.ResName()).AsString(), _em.ObjId())
	}
	ioh.PrintString("] >> ")
	ioh.PrintDictEnd()
	return ioh.PrintObjectEnd()
}

func (p *PdfCatalogObject) SetPages(pages PdfObjRefValue) {
	p.pages = &pages
}

func (p *PdfCatalogObject) SetProps(props PdfObjRefValue) {
	p.props = &props
}

func (p *PdfCatalogObject) AddIntent(intent PdfObjRefValue) {
	p.intents = append(p.intents, intent)
}

func (p PdfCatalogObject) Ref() PdfObjRefValue {
	return NewPdfObjRefValue(p.objnum)
}

func (p PdfCatalogObject) GetObjNum() int {
	return p.objnum
}

func (p *PdfCatalogObject) SetObjNum(n int, pdf *PdfDoc) {
	p.objnum = n
}

func (p *PdfCatalogObject) SetOutlines(ref PdfObjRefValue) {
	p.outlines = &ref
}

// ******************************************************************

type PdfOptionalContentGroupObject struct {
	PdfBaseObject
	groups []PdfObjRefValue
}

func (p *PdfOptionalContentGroupObject) Add(ocg PdfObjRefValue) {
	p.groups = append(p.groups, ocg)
}

func (p *PdfOptionalContentGroupObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(p.objnum)
	ioh.PrintDictStart()
	ioh.PrintString("/OCGs [ ")
	for _, res := range p.groups {
		ioh.PrintFmt(" %s ", res.AsString())
	}
	ioh.PrintString("] ")
	ioh.PrintString("/D << >> ")
	ioh.PrintDictEnd()
	return ioh.PrintObjectEnd()
}

func NewPdfOptionalContentGroupObject() *PdfOptionalContentGroupObject {
	return &PdfOptionalContentGroupObject{
		PdfBaseObject: PdfBaseObject{},
		groups:        make([]PdfObjRefValue, 0),
	}
}

type PdfOptionalContentObject struct {
	PdfBaseObject
	name string
}

func (p PdfOptionalContentObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(p.objnum)
	ioh.PrintFmt("<</Type /OCG /Name %s>>", MakePdfStringUtf8WithBrackets(p.name))
	return ioh.PrintObjectEnd()
}

func (p PdfOptionalContentObject) ResName() string {
	return fmt.Sprintf("OC%d", p.objnum)
}

func NewPdfOptionalContentObject(n string) *PdfOptionalContentObject {
	return &PdfOptionalContentObject{name: n}
}

// ******************************************************************

type PdfPageBase interface {
	Ref() PdfObjRefValue
	GetCount() int
}

type PdfPages interface {
	Ref() PdfObjRefValue
	Add(p PdfPageBase)
	Prepend(p PdfPageBase)
	InsertAt(idx int, p PdfPageBase)
	GetCount() int
}

type PdfPageObject struct {
	PdfBaseObject
	parent      PdfObjRefValue
	fonts       map[string]PdfResource
	xobjects    map[string]PdfResource
	extgstates  map[string]PdfResource
	colorspaces map[string]PdfResource
	properties  map[string]PdfResource
	contents    []PdfObjRefValue
	mediabox    string
	artbox      string
	cropbox     string
	trimbox     string
	bleedbox    string
}

func NewPdfPageObject(parent PdfObjRefValue) *PdfPageObject {
	return &PdfPageObject{
		mediabox:      PAGE_SIZE_A4_LITERAL,
		PdfBaseObject: PdfBaseObject{objnum: -1},
		parent:        parent,
		fonts:         make(map[string]PdfResource),
		xobjects:      make(map[string]PdfResource),
		extgstates:    make(map[string]PdfResource),
		colorspaces:   make(map[string]PdfResource),
		properties:    make(map[string]PdfResource),
		contents:      make([]PdfObjRefValue, 0),
	}
}

func (ppo *PdfPageObject) SetMediaBox(x int, y int, w int, h int) {
	ppo.mediabox = fmt.Sprintf("[%d %d %d %d]", x, y, w, h)
}
func (ppo *PdfPageObject) SetTrimBox(x int, y int, w int, h int) {
	ppo.trimbox = fmt.Sprintf("[%d %d %d %d]", x, y, w, h)
}
func (ppo *PdfPageObject) SetCropBox(x int, y int, w int, h int) {
	ppo.cropbox = fmt.Sprintf("[%d %d %d %d]", x, y, w, h)
}
func (ppo *PdfPageObject) SetBleedBox(x int, y int, w int, h int) {
	ppo.bleedbox = fmt.Sprintf("[%d %d %d %d]", x, y, w, h)
}
func (ppo *PdfPageObject) SetArtBox(x int, y int, w int, h int) {
	ppo.artbox = fmt.Sprintf("[%d %d %d %d]", x, y, w, h)
}

func (ppo *PdfPageObject) Add(pco *PdfContentObject) {
	pco.pdf = ppo.pdf
	ppo.contents = append(ppo.contents, pco.Ref())
}

func (ppo *PdfPageObject) Prepend(pco *PdfContentObject) {
	pco.pdf = ppo.pdf
	ppo.contents = slices.Insert(ppo.contents, 0, pco.Ref())
}

func (ppo *PdfPageObject) InsertAt(idx int, pco *PdfContentObject) {
	pco.pdf = ppo.pdf
	ppo.contents = slices.Insert(ppo.contents, idx, pco.Ref())
}

func (ppo PdfPageObject) StreamOut(ioh *IoHelper) error {
	ioh.PrintObjectStart(ppo.objnum)
	ioh.PrintDictStart()
	ioh.PrintString("/Type /Page")
	ioh.PrintFmt("/Parent %s ", ppo.parent.AsString())
	ioh.PrintString("/Resources << ")
	if len(ppo.fonts) > 0 {
		ioh.PrintString("/Font << ")
		for _, fnt := range ppo.fonts {
			ioh.PrintFmt("/%s %s ", fnt.ResName(), fnt.Ref().AsString())
		}
		ioh.PrintString(">> ")
	}
	if len(ppo.xobjects) > 0 {
		ioh.PrintString("/XObject << ")
		for _, res := range ppo.xobjects {
			ioh.PrintFmt("/%s %s ", res.ResName(), res.Ref().AsString())
		}
		ioh.PrintString(">> ")
	}
	if len(ppo.extgstates) > 0 {
		ioh.PrintString("/ExtGstate << ")
		for _, res := range ppo.extgstates {
			ioh.PrintFmt("/%s %s ", res.ResName(), res.Ref().AsString())
		}
		ioh.PrintString(">> ")
	}
	if len(ppo.colorspaces) > 0 {
		ioh.PrintString("/ColorSpace << ")
		for _, res := range ppo.colorspaces {
			ioh.PrintFmt("/%s %s ", res.ResName(), res.Ref().AsString())
		}
		ioh.PrintString(">> ")
	}
	if len(ppo.properties) > 0 {
		ioh.PrintString("/Properties << ")
		for _, res := range ppo.properties {
			ioh.PrintFmt("/%s %s ", res.ResName(), res.Ref().AsString())
		}
		ioh.PrintString(">> ")
	}
	ioh.PrintString("/ProcSet [ /PDF /Text /ImageC /ImageB /ImageI ]")
	ioh.PrintString(">> ")
	ioh.PrintFmt("/MediaBox %s ", ppo.mediabox)
	if ppo.artbox != "" {
		ioh.PrintFmt("/ArtBox %s ", ppo.artbox)
	}
	if ppo.cropbox != "" {
		ioh.PrintFmt("/CropBox %s ", ppo.cropbox)
	}
	if ppo.bleedbox != "" {
		ioh.PrintFmt("/BleedBox %s ", ppo.bleedbox)
	}
	if ppo.trimbox != "" {
		ioh.PrintFmt("/TrimBox %s ", ppo.trimbox)
	}
	ioh.PrintString("/Contents [ ")
	for _, cnt := range ppo.contents {
		ioh.PrintFmt(" %s ", cnt.AsString())
	}
	ioh.PrintString("] ")
	ioh.PrintDictEnd()
	return ioh.PrintObjectEnd()
}

func (p PdfPageObject) GetCount() int {
	return 1
}

func (ppo *PdfPageObject) UseFont(res PdfResource) {
	ppo.fonts[res.ResName()] = res
}

func (ppo *PdfPageObject) UseColorSpace(res PdfResource) {
	ppo.colorspaces[res.ResName()] = res
}

func (ppo *PdfPageObject) UseEGState(egs PdfResource) {
	ppo.extgstates[egs.ResName()] = egs
}

func (ppo *PdfPageObject) SetMediabox(mediabox string) {
	ppo.mediabox = mediabox
}

func (ppo *PdfPageObject) AddLayer(n string) PdfResource {
	if ppo.properties[n] == nil {
		ppo.properties[n] = ppo.pdf.AddLayer(n)
	}
	return ppo.properties[n]
}

type PdfPagesObject struct {
	objnum int
	pages  []PdfPageBase
	parent *PdfObjRefValue
}

func (ppo *PdfPagesObject) Add(p PdfPageBase) {
	ppo.pages = append(ppo.pages, p)
}

func (ppo *PdfPagesObject) Prepend(p PdfPageBase) {
	ppo.pages = slices.Insert(ppo.pages, 0, p)
}

func (ppo *PdfPagesObject) InsertAt(idx int, p PdfPageBase) {
	ppo.pages = slices.Insert(ppo.pages, idx, p)
}

func NewPdfPagesObject(parent PdfObjRefValue) *PdfPagesObject {
	return &PdfPagesObject{pages: make([]PdfPageBase, 0), parent: &parent}
}

func NewPdfPagesObjectNoParent() *PdfPagesObject {
	return &PdfPagesObject{pages: make([]PdfPageBase, 0)}
}

func (p PdfPagesObject) Ref() PdfObjRefValue {
	return NewPdfObjRefValue(p.objnum)
}

func (p PdfPagesObject) GetObjNum() int {
	return p.objnum
}

func (p *PdfPagesObject) SetObjNum(n int, pdf *PdfDoc) {
	p.objnum = n
}

func (p PdfPagesObject) GetCount() int {
	_sz := 0
	for _, page := range p.pages {
		_sz += page.GetCount()
	}
	return _sz
}

func (p PdfPagesObject) StreamOut(ioh *IoHelper) error {
	_sz := p.GetCount()
	ioh.PrintObjectStart(p.objnum)
	ioh.PrintDictStart()
	ioh.PrintString("/Type /Pages")
	ioh.PrintFmt("/Count %d ", _sz)
	if p.parent != nil {
		ioh.PrintFmt("/Parent %s ", p.parent.AsString())
	}
	ioh.PrintString("/Kids [ ")
	for _, page := range p.pages {
		ioh.PrintFmt(" %s ", page.Ref().AsString())
	}
	ioh.PrintString("] ")
	ioh.PrintDictEnd()
	return ioh.PrintObjectEnd()
}

const (
	PAGE_SIZE_A1_LITERAL      = "[0 0 1684 2384]"
	PAGE_SIZE_A2_LITERAL      = "[0 0 1191 1684]"
	PAGE_SIZE_A3_LITERAL      = "[0 0 842 1191]"
	PAGE_SIZE_A4_LITERAL      = "[0 0 595 842]"
	PAGE_SIZE_A5_LITERAL      = "[0 0 421 595]"
	PAGE_SIZE_A6_LITERAL      = "[0 0 298 421]"
	PAGE_SIZE_A7_LITERAL      = "[0 0 210 298]"
	PAGE_SIZE_LETTER_LITERAL  = "[0 0 612 792]"
	PAGE_SIZE_LEGAL_LITERAL   = "[0 0 612 1008]"
	PAGE_SIZE_TABLOID_LITERAL = "[0 0 792 1224]"
)
