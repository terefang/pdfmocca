package pdf

// *************************************************************************
// interface for content streams
// *************************************************************************
type PdfStream interface {
	Length() int64
	UnCompressedLength() int64
	Filter() PdfValue
	StreamOut(stream *IoHelper) error
}

// *************************************************************************
// simple string content stream
// *************************************************************************
type PdfStringStream struct {
	length  int64
	content string
}

func (s PdfStringStream) Length() int64 { return s.length }

func (s PdfStringStream) UnCompressedLength() int64 { return s.length }

func (s PdfStringStream) Filter() PdfValue { return nil }

func (s PdfStringStream) StreamOut(stream *IoHelper) error {
	stream.PrintStreamStart()
	stream.PrintString(s.content)
	stream.Flush()
	return stream.PrintStreamEnd()
}

func NewPdfStringStream(content string) *PdfStringStream {
	return &PdfStringStream{content: content, length: int64(len(content))}
}

// *************************************************************************
// simple bytes content stream with optional filter
// *************************************************************************
type PdfBytesStream struct {
	content []byte
	length  int64
	ulength int64
	filter  PdfValue
}

func (s PdfBytesStream) Length() int64             { return s.length }
func (s PdfBytesStream) UnCompressedLength() int64 { return s.ulength }

func (s PdfBytesStream) Filter() PdfValue { return s.filter }

func (s PdfBytesStream) StreamOut(stream *IoHelper) error {
	stream.PrintStreamStart()
	stream.PrintBytes(s.content)
	stream.Flush()
	return stream.PrintStreamEnd()
}

func NewPdfBytesStreamWithFilter(content []byte, ulen int64, flen int64, filter string) *PdfBytesStream {
	return &PdfBytesStream{content: content, ulength: ulen, length: flen, filter: NewPdfLiteralValue(filter)}
}

func NewPdfBytesStream(content []byte, ulen int64) *PdfBytesStream {
	return &PdfBytesStream{content: content, ulength: ulen, length: ulen, filter: nil}
}

// *************************************************************************
// maker of flated bytestreams
// *************************************************************************
func NewPdfFlatedStream(content []byte) *PdfBytesStream {
	ulen := int64(len(content))
	_flated := flateCompress(content)
	olen := int64(len(_flated))
	return NewPdfBytesStreamWithFilter(_flated, ulen, olen, "/FlateDecode")
	//return NewPdfBytesStream(content, ulen)
}
