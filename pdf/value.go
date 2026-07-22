package pdf

import (
	"fmt"
	"strings"
	"time"
)

// *************************************************************************
// the value interface
// *************************************************************************
type PdfValue interface {
	AsString() string
	StreamOut(stream *IoHelper) error
}

// *************************************************************************
// value reference – refers to document level objects
// *************************************************************************
type PdfObjRefValue struct {
	objnum int
}

func (s PdfObjRefValue) ObjId() int {
	return s.objnum
}

func NewPdfObjRefValue(n int) PdfObjRefValue {
	return PdfObjRefValue{n}
}

func (s PdfObjRefValue) AsString() string {
	return fmt.Sprintf("%d 0 R ", s.objnum)
}

func (s PdfObjRefValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintFmt("%d 0 R ", s.objnum)
}

// *************************************************************************
// literal value will be written literally to the pdf stream
// *************************************************************************
type PdfLiteralValue struct {
	literal string
}

func NewPdfLiteralValue(s string) PdfLiteralValue {
	return PdfLiteralValue{literal: s}
}

func (s PdfLiteralValue) AsString() string {
	return s.literal
}

func (s PdfLiteralValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintString(s.literal)
}

// *************************************************************************
// pdf integer
// *************************************************************************
type PdfIntValue struct {
	value int64
}

func NewPdfIntValue(n int64) PdfIntValue {
	return PdfIntValue{value: n}
}

func (s PdfIntValue) AsString() string {
	return fmt.Sprintf("%d", s.value)
}

func (s PdfIntValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintString(s.AsString())
}

// *************************************************************************
// pdf integer
// *************************************************************************
type PdfFloatValue struct {
	value float64
}

func NewPdfFloatValue(n float64) PdfFloatValue {
	return PdfFloatValue{value: n}
}

func (s PdfFloatValue) AsString() string {
	return fmt.Sprintf("%f", s.value)
}

func (s PdfFloatValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintString(s.AsString())
}

// *************************************************************************
// postscript/pdf name value
// *************************************************************************
type PdfNameValue struct {
	str string
}

func NewPdfNameValue(s string) *PdfNameValue { return &PdfNameValue{str: s} }

func (d PdfNameValue) AsString() string {
	return MakePdfNameString(d.str)
}

func (d PdfNameValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

// *************************************************************************
// simple string value in pdfdoc encoding
// *************************************************************************
type PdfStringValue struct {
	str string
}

func NewPdfStringValue(s string) *PdfStringValue { return &PdfStringValue{str: s} }

func (d PdfStringValue) AsString() string {
	return MakePdfStringWithBrackets(d.str)
}

func (d PdfStringValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

func NewPdfTimeString(t time.Time) *PdfStringValue {
	return NewPdfStringValue(fmt.Sprintf("D:%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

// *************************************************************************
// utf8 encoded string value
// *************************************************************************
type PdfStringUtf8Value struct {
	str string
}

func NewPdfStringUtf8Value(s string) *PdfStringUtf8Value { return &PdfStringUtf8Value{str: s} }

func (d PdfStringUtf8Value) AsString() string {
	return MakePdfStringUtf8WithBrackets(d.str)
}

func (d PdfStringUtf8Value) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

// *************************************************************************
// ucs2 encoded string value
// *************************************************************************
type PdfStringUcsValue struct {
	str string
}

func NewPdfStringUcsValue(s string) *PdfStringUcsValue { return &PdfStringUcsValue{str: s} }

func (d PdfStringUcsValue) AsString() string {
	return MakePdfStringUcsWithBrackets(d.str)
}

func (d PdfStringUcsValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

// *************************************************************************
// dictionary value
// *************************************************************************
type PdfDictValue struct {
	dict map[string]PdfValue
}

func NewPdfDictValue() PdfDictValue { return PdfDictValue{dict: make(map[string]PdfValue)} }

func (d PdfDictValue) Set(name string, value PdfValue) {
	d.dict[name] = value
}

func (d PdfDictValue) AsString() string {
	_sb := strings.Builder{}
	_sb.WriteString("<< ")
	for k, v := range d.dict {
		_sb.WriteString(fmt.Sprintf("%s %s ", MakePdfNameString(k), v.AsString()))
	}
	_sb.WriteString(">> ")
	return _sb.String()
}

func (d PdfDictValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

// *************************************************************************
// array/list value
// *************************************************************************
type PdfArrayValue struct {
	array []PdfValue
}

func NewPdfArrayValue() *PdfArrayValue { return &PdfArrayValue{array: make([]PdfValue, 0)} }

func (d PdfArrayValue) Add(value PdfValue) {
	d.array = append(d.array, value)
}

func (d PdfArrayValue) AsString() string {
	_sb := strings.Builder{}
	_sb.WriteString("[ ")
	_i := 0
	for _, v := range d.array {
		_i++
		_sb.WriteString(fmt.Sprintf("%s ", v.AsString()))
		if (_i % 20) == 0 {
			_sb.WriteString("\n")
		}
	}
	_sb.WriteString(" ] ")
	return _sb.String()
}

func (d PdfArrayValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

// *************************************************************************
// int array/list value
// *************************************************************************
type PdfIntArrayValue struct {
	array []int
}

func NewPdfIntArrayValue() *PdfIntArrayValue { return &PdfIntArrayValue{array: make([]int, 0)} }

func (d PdfIntArrayValue) Add(value int) {
	d.array = append(d.array, value)
}

func (d PdfIntArrayValue) AsString() string {
	_sb := strings.Builder{}
	_sb.WriteString("[ ")
	_i := 0
	for _, v := range d.array {
		_i++
		_sb.WriteString(fmt.Sprintf("%d ", v))
		if (_i % 20) == 0 {
			_sb.WriteString("\n")
		}
	}
	_sb.WriteString(" ] ")
	return _sb.String()
}

func (d PdfIntArrayValue) StreamOut(ioh *IoHelper) error {
	return ioh.PrintStringLn(d.AsString())
}

func NewPdfIntBBoxValue(a int, b int, c int, d int) *PdfIntArrayValue {
	arr := NewPdfIntArrayValue()
	arr.Add(a)
	arr.Add(b)
	arr.Add(c)
	arr.Add(d)
	return arr
}

func NewPdfBBoxValue(a int, b int, c int, d int) *PdfLiteralValue {
	obj := NewPdfLiteralValue(fmt.Sprintf("[%d %d %d %d]", a, b, c, d))
	return &obj
}

func NewPdfBBoxValueF(a float64, b float64, c float64, d float64) *PdfLiteralValue {
	obj := NewPdfLiteralValue(fmt.Sprintf("[%f %f %f %f]", a, b, c, d))
	return &obj
}

func NewPdfMatrixValue(a int, b int, c int, d int, e int, f int) *PdfLiteralValue {
	obj := NewPdfLiteralValue(fmt.Sprintf("[%d %d %d %d %d %d]", a, b, c, d, e, f))
	return &obj
}

func NewPdfMatrixValueF(a float64, b float64, c float64, d float64, e float64, f float64) *PdfLiteralValue {
	obj := NewPdfLiteralValue(fmt.Sprintf("[%f %f %f %f %f %f]", a, b, c, d, e, f))
	return &obj
}
