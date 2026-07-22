package main

import (
	"pdfmocca/pdf"
)

func main() {
	_pdf := pdf.NewPdfDocument(pdf.PDF_VERSION_1_5)
	_pdf.Info().SetProducer("pdfmocca/go")
	_pdf.SetTitle("Lorem Ipsum!")

	_pdf.StreamBeginToFile("test.pdf")

	_hvb := _pdf.NewCoreFont(pdf.FONT_CORE_HELV_BOLD, 0, 255, pdf.PDFDOC_HELV_BOLD_W, pdf.PDFDOC_ENCODING)

	_f := _pdf.NewCoreFont(pdf.FONT_CORE_DINGBATS, 0, 255, pdf.PDF_DINGBATS_W, pdf.DINGBAT_ENCODING)

	_p := _pdf.NewPage()
	_c := _p.NewContent()
	_c.StartLayer("HV1")
	for _i := 0; _i < 16; _i++ {
		for _j := 0; _j < 16; _j++ {
			_c.PutText(50+(30*_i), 50+(30*_j), _f, 25, "dark-violet", string(rune(_i+(_j*16))))
		}
	}
	_c.EndLayer()

	_pdf.AddOutline("First Page", _p.Ref())

	var _fl []string = []string{
		pdf.FONT_CORE_SYMBOL, pdf.FONT_CORE_DINGBATS,
		"qhvr", "qhvi", "qhvb", "qhvz",
		"qcrr", "qcri", "qcrb", "qcrz",
		"qtmr", "qtmi", "qtmb", "qtmz",
	}

	for _, _ff := range _fl {
		_tt := _pdf.NewTtfFont(_ff, -1)
		_p = _pdf.NewPage()
		_pdf.AddOutline(_ff, _p.Ref())
		_c = _p.NewContent()
		_c.PutText(50, 830, _hvb, 10, "red", _ff)
		_c.StartLayer("TTF1")
		for _i := 0; _i < 16; _i++ {
			for _j := 0; _j < 30; _j++ {
				_u := _i + (_j * 16)
				_ct := string(rune(_u))
				_c.PutText(50+(30*_i), 50+(25*_j), _tt, 20, "dark-red", _ct)
			}
		}
		_c.EndLayer()

		for _ofs := 0; _ofs <= _tt.MaxCid(); _ofs += 16 * 30 {
			_p = _pdf.NewPage()
			_c = _p.NewContent()
			for _i := 0; _i < 16; _i++ {
				for _j := 0; _j < 30; _j++ {
					_u := _ofs + _i + (_j * 16)
					_c.PutCharCid(50+(30*_i), 50+(25*_j), _tt, 20, "black", _u)
				}
			}
		}
	}

	//	_pdf.NewEmbedFile("File1", "path1")
	//	_pdf.NewEmbedFile("File2", "path1")
	//	_pdf.NewEmbedFile("File3", "path1")
	_pdf.StreamEnd(true)
}
