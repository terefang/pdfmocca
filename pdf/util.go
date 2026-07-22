package pdf

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

type IoHelper struct {
	fh     *os.File
	md     hash.Hash
	offset int64
}

func (ioh *IoHelper) Write(p []byte) (n int, err error) {
	ioh.md.Write(p)
	_count, _err := ioh.fh.Write(p)
	ioh.offset += int64(_count)
	return _count, _err
}

func (ioh *IoHelper) Print(p []byte) (err error) {
	_, _err := ioh.Write(p)
	return _err
}

func NewIoHelper(fh *os.File) *IoHelper {
	return &IoHelper{fh: fh, offset: 0, md: md5.New()}
}

func (ioh *IoHelper) Flush() {
	ioh.fh.Sync()
}

func (ioh *IoHelper) Close() {
	ioh.fh.Close()
}

func (ioh *IoHelper) PrintBytes(p []byte) (err error) {
	return ioh.Print(p)
}

func (ioh *IoHelper) PrintString(p string) (err error) {
	return ioh.Print([]byte(p))
}

func (ioh *IoHelper) PrintStringLn(p string) (err error) {
	ioh.Print([]byte(p))
	return ioh.Print([]byte{'\n'})
}

func (ioh *IoHelper) PrintFmt(_fmt string, v ...any) (err error) {
	return ioh.PrintString(fmt.Sprintf(_fmt, v...))
}

func (ioh *IoHelper) PrintFmtLn(_fmt string, v ...any) (err error) {
	return ioh.PrintStringLn(fmt.Sprintf(_fmt, v...))
}

func (ioh *IoHelper) PrintObjectStartDict(n int) error {
	return ioh.PrintFmtLn("%d 0 obj <<", n)
}

func (ioh *IoHelper) PrintObjectStart(n int) error {
	return ioh.PrintFmt("%d 0 obj\n", n)
}

func (ioh *IoHelper) PrintDictStart() error {
	return ioh.PrintString(" << ")
}

func (ioh *IoHelper) PrintDictEnd() error {
	return ioh.PrintString(" >> ")
}

func (ioh *IoHelper) PrintObjectEnd() error {
	return ioh.PrintStringLn("endobj")
}

func (ioh *IoHelper) PrintStreamStart() error {
	return ioh.PrintStringLn("stream")
}

func (ioh *IoHelper) PrintStreamEnd() error {
	return ioh.PrintStringLn("\nendstream")
}

func (ioh *IoHelper) PrintInt8BE(n int) error {
	_err := binary.Write(ioh, binary.BigEndian, int8(n))
	return _err
}

func (ioh *IoHelper) Print121(a int, b int64, c int) error {
	_err := binary.Write(ioh, binary.BigEndian, int8(a))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int16(b))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int8(c))
	return _err
}

func (ioh *IoHelper) Print141(a int, b int64, c int) error {
	_err := binary.Write(ioh, binary.BigEndian, int8(a))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int32(b))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int8(c))
	return _err
}

func (ioh *IoHelper) Print181(a int, b int64, c int) error {
	_err := binary.Write(ioh, binary.BigEndian, int8(a))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int64(b))
	if _err != nil {
		return _err
	}
	_err = binary.Write(ioh, binary.BigEndian, int8(c))
	return _err
}

func (ioh *IoHelper) SyncOffset() error {
	offset, err := ioh.fh.Seek(0, io.SeekCurrent)
	ioh.offset = offset
	return err
}

func (ioh *IoHelper) Offset() int64 {
	return ioh.offset
}

func (ioh *IoHelper) Finalize() string {
	_md := ioh.md.Sum(nil)
	return hex.EncodeToString(_md)
}

func MakePdfNameString(name string) string {
	_sb := strings.Builder{}
	_sb.WriteString("/")
	runes := []rune(name)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '#' {
			_sb.WriteString("#23")
		} else if runes[i] >= '!' && runes[i] <= '~' {
			_sb.WriteRune(runes[i])
		} else if runes[i] >= 256 {
			_sb.WriteString("#20")
		} else {
			_sb.WriteString(fmt.Sprintf("#%02X", int(runes[i])))
		}
	}
	return _sb.String()
}

func MakePdfRunesStringWithBrackets(str []EncodedRune) string {
	_sb := strings.Builder{}
	_sb.WriteString("(")
	_sb.WriteString(MakePdfRunesString(str))
	_sb.WriteString(")")
	return _sb.String()
}

func MakePdfRuneString(r int, _sb *strings.Builder) {
	switch r {
	case '\\':
		_sb.WriteString("\\\\")
	case '\n':
		_sb.WriteString("\\n")
	case '\r':
		_sb.WriteString("\\r")
	case '\t':
		_sb.WriteString("\\t")
	case '\f':
		_sb.WriteString("\\f")
	case '(':
		_sb.WriteString("\\(")
	case ')':
		_sb.WriteString("\\)")
	default:
		{
			if r < 32 {
				_sb.WriteString(fmt.Sprintf("\\%03o", r))
			} else if r > 127 && r < 256 {
				_sb.WriteString(fmt.Sprintf("\\%03o", r))
			} else if r > 256 {
				_sb.WriteString(fmt.Sprintf("[u%04x]", r))
			} else {
				_sb.WriteRune(rune(r))
			}
		}
	}

}

func MakePdfRuneCidString(r int, _sb *strings.Builder) {
	if r > 0xffff {
		_sb.WriteString("0000")
	} else {
		_sb.WriteString(fmt.Sprintf("%04x", r))
	}
}

func MakePdfRunesString(str []EncodedRune) string {
	_sb := strings.Builder{}
	_sz := len(str)
	for i := 0; i < _sz; i++ {
		MakePdfRuneString(str[i].Gid, &_sb)
	}
	return _sb.String()
}

func MakePdfRunesCidString(str []EncodedRune) string {
	_sb := strings.Builder{}
	_sz := len(str)
	for i := 0; i < _sz; i++ {
		MakePdfRuneCidString(str[i].Gid, &_sb)
	}
	return _sb.String()
}

func MakePdfStringWithBrackets(str string) string {
	_sb := strings.Builder{}
	_sb.WriteString("(")
	_sb.WriteString(MakePdfStringNoBrackets(str))
	_sb.WriteString(")")
	return _sb.String()
}

func MakePdfStringNoBrackets(str string) string {
	_sb := strings.Builder{}
	_runes := []rune(str)
	_sz := len(_runes)
	for i := 0; i < _sz; i++ {
		MakePdfRuneString(int(_runes[i]), &_sb)
	}
	return _sb.String()
}

func MakePdfStringUtf8WithBrackets(str string) string {
	_sb := strings.Builder{}
	_sb.WriteString("<efbbbf")
	_sb.WriteString(MakePdfStringUtf8NoBrackets(str))
	_sb.WriteString(">")
	return _sb.String()
}

func MakePdfStringUtf8NoBrackets(str string) string {
	_sb := strings.Builder{}
	_sz := len(str)
	for i := 0; i < _sz; i++ {
		_sb.WriteString(fmt.Sprintf("%02x", str[i]))
	}
	return _sb.String()
}

func MakePdfStringUcsWithBrackets(str string) string {
	_sb := strings.Builder{}
	_sb.WriteString("<feff")
	_sb.WriteString(MakePdfStringUcsNoBrackets(str))
	_sb.WriteString(">")
	return _sb.String()
}

func MakePdfStringUcsNoBrackets(str string) string {
	_sb := strings.Builder{}
	_runes := []rune(str)
	_sz := len(_runes)
	for i := 0; i < _sz; i++ {
		if _runes[i] > 0xffff {
			_sb.WriteString("003f")
		} else {
			_sb.WriteString(fmt.Sprintf("%04x", int(_runes[i])&0xffff))
		}
	}
	return _sb.String()
}

func flateCompress(plain []byte) []byte {
	var _buf bytes.Buffer
	_bufwr := bufio.NewWriter(&_buf)
	_flate := zlib.NewWriter(_bufwr)
	_flate.Write(plain)
	_flate.Flush()
	_flate.Close()
	_bufwr.Flush()
	io.CopyN()
	return _buf.Bytes()
}
