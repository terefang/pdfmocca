package pdf

import (
	"math"
	"strconv"
	"strings"
)

var rgbColors map[string]*[]float64 = map[string]*[]float64{
	//! Red colors
	"lightsalmon": &[]float64{float64(0xFF) / 255., float64(0xA0) / 255., float64(0x7A) / 255.},
	"salmon":      &[]float64{float64(0xFA) / 255., float64(0x80) / 255., float64(0x72) / 255.},
	"darksalmon":  &[]float64{float64(0xE9) / 255., float64(0x96) / 255., float64(0x7A) / 255.},
	"lightcoral":  &[]float64{float64(0xF0) / 255., float64(0x80) / 255., float64(0x80) / 255.},
	"indianred":   &[]float64{float64(0xCD) / 255., float64(0x5C) / 255., float64(0x5C) / 255.},
	"crimson":     &[]float64{float64(0xDC) / 255., float64(0x14) / 255., float64(0x3C) / 255.},
	"firebrick":   &[]float64{float64(0xB2) / 255., float64(0x22) / 255., float64(0x22) / 255.},
	"red":         &[]float64{float64(0xFF) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	"darkred":     &[]float64{float64(0x8B) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	//! Orange colors
	"coral":      &[]float64{float64(0xFF) / 255., float64(0x7F) / 255., float64(0x50) / 255.},
	"tomato":     &[]float64{float64(0xFF) / 255., float64(0x63) / 255., float64(0x47) / 255.},
	"orangered":  &[]float64{float64(0xFF) / 255., float64(0x45) / 255., float64(0x00) / 255.},
	"gold":       &[]float64{float64(0xFF) / 255., float64(0xD7) / 255., float64(0x00) / 255.},
	"orange":     &[]float64{float64(0xFF) / 255., float64(0xA5) / 255., float64(0x00) / 255.},
	"darkorange": &[]float64{float64(0xFF) / 255., float64(0x8C) / 255., float64(0x00) / 255.},
	//! Yellow colors
	"lightyellow":          &[]float64{float64(0xFF) / 255., float64(0xFF) / 255., float64(0xE0) / 255.},
	"lemonchiffon":         &[]float64{float64(0xFF) / 255., float64(0xFA) / 255., float64(0xCD) / 255.},
	"lightgoldenrodyellow": &[]float64{float64(0xFA) / 255., float64(0xFA) / 255., float64(0xD2) / 255.},
	"papayawhip":           &[]float64{float64(0xFF) / 255., float64(0xEF) / 255., float64(0xD5) / 255.},
	"moccasin":             &[]float64{float64(0xFF) / 255., float64(0xE4) / 255., float64(0xB5) / 255.},
	"peachpuff":            &[]float64{float64(0xFF) / 255., float64(0xDA) / 255., float64(0xB9) / 255.},
	"palegoldenrod":        &[]float64{float64(0xEE) / 255., float64(0xE8) / 255., float64(0xAA) / 255.},
	"khaki":                &[]float64{float64(0xF0) / 255., float64(0xE6) / 255., float64(0x8C) / 255.},
	"darkkhaki":            &[]float64{float64(0xBD) / 255., float64(0xB7) / 255., float64(0x6B) / 255.},
	"yellow":               &[]float64{float64(0xFF) / 255., float64(0xFF) / 255., float64(0x00) / 255.},
	//! Green colors
	"lawngreen":         &[]float64{float64(0x7C) / 255., float64(0xFC) / 255., float64(0x00) / 255.},
	"chartreuse":        &[]float64{float64(0x7F) / 255., float64(0xFF) / 255., float64(0x00) / 255.},
	"limegreen":         &[]float64{float64(0x32) / 255., float64(0xCD) / 255., float64(0x32) / 255.},
	"lime":              &[]float64{float64(0x00) / 255., float64(0xFF) / 255., float64(0x00) / 255.},
	"forestgreen":       &[]float64{float64(0x22) / 255., float64(0x8B) / 255., float64(0x22) / 255.},
	"green":             &[]float64{float64(0x00) / 255., float64(0x80) / 255., float64(0x00) / 255.},
	"darkgreen":         &[]float64{float64(0x00) / 255., float64(0x64) / 255., float64(0x00) / 255.},
	"greenyellow":       &[]float64{float64(0xAD) / 255., float64(0xFF) / 255., float64(0x2F) / 255.},
	"yellowgreen":       &[]float64{float64(0x9A) / 255., float64(0xCD) / 255., float64(0x32) / 255.},
	"springgreen":       &[]float64{float64(0x00) / 255., float64(0xFF) / 255., float64(0x7F) / 255.},
	"mediumspringgreen": &[]float64{float64(0x00) / 255., float64(0xFA) / 255., float64(0x9A) / 255.},
	"lightgreen":        &[]float64{float64(0x90) / 255., float64(0xEE) / 255., float64(0x90) / 255.},
	"palegreen":         &[]float64{float64(0x98) / 255., float64(0xFB) / 255., float64(0x98) / 255.},
	"darkseagreen":      &[]float64{float64(0x8F) / 255., float64(0xBC) / 255., float64(0x8F) / 255.},
	"mediumseagreen":    &[]float64{float64(0x3C) / 255., float64(0xB3) / 255., float64(0x71) / 255.},
	"seagreen":          &[]float64{float64(0x2E) / 255., float64(0x8B) / 255., float64(0x57) / 255.},
	"olive":             &[]float64{float64(0x80) / 255., float64(0x80) / 255., float64(0x00) / 255.},
	"darkolivegreen":    &[]float64{float64(0x55) / 255., float64(0x6B) / 255., float64(0x2F) / 255.},
	"olivedrab":         &[]float64{float64(0x6B) / 255., float64(0x8E) / 255., float64(0x23) / 255.},
	//! Cyan colors
	"lightcyan":        &[]float64{float64(0xE0) / 255., float64(0xFF) / 255., float64(0xFF) / 255.},
	"cyan":             &[]float64{float64(0x00) / 255., float64(0xFF) / 255., float64(0xFF) / 255.},
	"aqua":             &[]float64{float64(0x00) / 255., float64(0xFF) / 255., float64(0xFF) / 255.},
	"aquamarine":       &[]float64{float64(0x7F) / 255., float64(0xFF) / 255., float64(0xD4) / 255.},
	"mediumaquamarine": &[]float64{float64(0x66) / 255., float64(0xCD) / 255., float64(0xAA) / 255.},
	"paleturquoise":    &[]float64{float64(0xAF) / 255., float64(0xEE) / 255., float64(0xEE) / 255.},
	"turquoise":        &[]float64{float64(0x40) / 255., float64(0xE0) / 255., float64(0xD0) / 255.},
	"mediumturquoise":  &[]float64{float64(0x48) / 255., float64(0xD1) / 255., float64(0xCC) / 255.},
	"darkturquoise":    &[]float64{float64(0x00) / 255., float64(0xCE) / 255., float64(0xD1) / 255.},
	"lightseagreen":    &[]float64{float64(0x20) / 255., float64(0xB2) / 255., float64(0xAA) / 255.},
	"cadetblue":        &[]float64{float64(0x5F) / 255., float64(0x9E) / 255., float64(0xA0) / 255.},
	"darkcyan":         &[]float64{float64(0x00) / 255., float64(0x8B) / 255., float64(0x8B) / 255.},
	"teal":             &[]float64{float64(0x00) / 255., float64(0x80) / 255., float64(0x80) / 255.},
	//! Blue colors
	"powderblue":      &[]float64{float64(0xB0) / 255., float64(0xE0) / 255., float64(0xE6) / 255.},
	"lightblue":       &[]float64{float64(0xAD) / 255., float64(0xD8) / 255., float64(0xE6) / 255.},
	"lightskyblue":    &[]float64{float64(0x87) / 255., float64(0xCE) / 255., float64(0xFA) / 255.},
	"skyblue":         &[]float64{float64(0x87) / 255., float64(0xCE) / 255., float64(0xEB) / 255.},
	"deepskyblue":     &[]float64{float64(0x00) / 255., float64(0xBF) / 255., float64(0xFF) / 255.},
	"lightsteelblue":  &[]float64{float64(0xB0) / 255., float64(0xC4) / 255., float64(0xDE) / 255.},
	"dodgerblue":      &[]float64{float64(0x1E) / 255., float64(0x90) / 255., float64(0xFF) / 255.},
	"cornflowerblue":  &[]float64{float64(0x64) / 255., float64(0x95) / 255., float64(0xED) / 255.},
	"steelblue":       &[]float64{float64(0x46) / 255., float64(0x82) / 255., float64(0xB4) / 255.},
	"royalblue":       &[]float64{float64(0x41) / 255., float64(0x69) / 255., float64(0xE1) / 255.},
	"blue":            &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0xFF) / 255.},
	"mediumblue":      &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0xCD) / 255.},
	"darkblue":        &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0x8B) / 255.},
	"navy":            &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0x80) / 255.},
	"midnightblue":    &[]float64{float64(0x19) / 255., float64(0x19) / 255., float64(0x70) / 255.},
	"mediumslateblue": &[]float64{float64(0x7B) / 255., float64(0x68) / 255., float64(0xEE) / 255.},
	"slateblue":       &[]float64{float64(0x6A) / 255., float64(0x5A) / 255., float64(0xCD) / 255.},
	"darkslateblue":   &[]float64{float64(0x48) / 255., float64(0x3D) / 255., float64(0x8B) / 255.},
	//! Purple colors
	"lavender":     &[]float64{float64(0xE6) / 255., float64(0xE6) / 255., float64(0xFA) / 255.},
	"thistle":      &[]float64{float64(0xD8) / 255., float64(0xBF) / 255., float64(0xD8) / 255.},
	"plum":         &[]float64{float64(0xDD) / 255., float64(0xA0) / 255., float64(0xDD) / 255.},
	"violet":       &[]float64{float64(0xEE) / 255., float64(0x82) / 255., float64(0xEE) / 255.},
	"orchid":       &[]float64{float64(0xDA) / 255., float64(0x70) / 255., float64(0xD6) / 255.},
	"fuchsia":      &[]float64{float64(0xFF) / 255., float64(0x00) / 255., float64(0xFF) / 255.},
	"magenta":      &[]float64{float64(0xFF) / 255., float64(0x00) / 255., float64(0xFF) / 255.},
	"mediumorchid": &[]float64{float64(0xBA) / 255., float64(0x55) / 255., float64(0xD3) / 255.},
	"mediumpurple": &[]float64{float64(0x93) / 255., float64(0x70) / 255., float64(0xDB) / 255.},
	"blueviolet":   &[]float64{float64(0x8A) / 255., float64(0x2B) / 255., float64(0xE2) / 255.},
	"darkviolet":   &[]float64{float64(0x94) / 255., float64(0x00) / 255., float64(0xD3) / 255.},
	"darkorchid":   &[]float64{float64(0x99) / 255., float64(0x32) / 255., float64(0xCC) / 255.},
	"darkmagenta":  &[]float64{float64(0x8B) / 255., float64(0x00) / 255., float64(0x8B) / 255.},
	"purple":       &[]float64{float64(0x80) / 255., float64(0x00) / 255., float64(0x80) / 255.},
	"indigo":       &[]float64{float64(0x4B) / 255., float64(0x00) / 255., float64(0x82) / 255.},
	//! Pink colors
	"pink":            &[]float64{float64(0xFF) / 255., float64(0xC0) / 255., float64(0xCB) / 255.},
	"lightpink":       &[]float64{float64(0xFF) / 255., float64(0xB6) / 255., float64(0xC1) / 255.},
	"hotpink":         &[]float64{float64(0xFF) / 255., float64(0x69) / 255., float64(0xB4) / 255.},
	"deeppink":        &[]float64{float64(0xFF) / 255., float64(0x14) / 255., float64(0x93) / 255.},
	"palevioletred":   &[]float64{float64(0xDB) / 255., float64(0x70) / 255., float64(0x93) / 255.},
	"mediumvioletred": &[]float64{float64(0xC7) / 255., float64(0x15) / 255., float64(0x85) / 255.},
	//! White colors
	"white":         &[]float64{float64(0xFF) / 255., float64(0xFF) / 255., float64(0xFF) / 255.},
	"snow":          &[]float64{float64(0xFF) / 255., float64(0xFA) / 255., float64(0xFA) / 255.},
	"honeydew":      &[]float64{float64(0xF0) / 255., float64(0xFF) / 255., float64(0xF0) / 255.},
	"mintcream":     &[]float64{float64(0xF5) / 255., float64(0xFF) / 255., float64(0xFA) / 255.},
	"azure":         &[]float64{float64(0xF0) / 255., float64(0xFF) / 255., float64(0xFF) / 255.},
	"aliceblue":     &[]float64{float64(0xF0) / 255., float64(0xF8) / 255., float64(0xFF) / 255.},
	"ghostwhite":    &[]float64{float64(0xF8) / 255., float64(0xF8) / 255., float64(0xFF) / 255.},
	"whitesmoke":    &[]float64{float64(0xF5) / 255., float64(0xF5) / 255., float64(0xF5) / 255.},
	"seashell":      &[]float64{float64(0xFF) / 255., float64(0xF5) / 255., float64(0xEE) / 255.},
	"beige":         &[]float64{float64(0xF5) / 255., float64(0xF5) / 255., float64(0xDC) / 255.},
	"oldlace":       &[]float64{float64(0xFD) / 255., float64(0xF5) / 255., float64(0xE6) / 255.},
	"floralwhite":   &[]float64{float64(0xFF) / 255., float64(0xFA) / 255., float64(0xF0) / 255.},
	"ivory":         &[]float64{float64(0xFF) / 255., float64(0xFF) / 255., float64(0xF0) / 255.},
	"antiquewhite":  &[]float64{float64(0xFA) / 255., float64(0xEB) / 255., float64(0xD7) / 255.},
	"linen":         &[]float64{float64(0xFA) / 255., float64(0xF0) / 255., float64(0xE6) / 255.},
	"lavenderblush": &[]float64{float64(0xFF) / 255., float64(0xF0) / 255., float64(0xF5) / 255.},
	"mistyrose":     &[]float64{float64(0xFF) / 255., float64(0xE4) / 255., float64(0xE1) / 255.},
	//! Gray colors
	"gainsboro":      &[]float64{float64(0xDC) / 255., float64(0xDC) / 255., float64(0xDC) / 255.},
	"lightgray":      &[]float64{float64(0xD3) / 255., float64(0xD3) / 255., float64(0xD3) / 255.},
	"silver":         &[]float64{float64(0xC0) / 255., float64(0xC0) / 255., float64(0xC0) / 255.},
	"darkgray":       &[]float64{float64(0xA9) / 255., float64(0xA9) / 255., float64(0xA9) / 255.},
	"gray":           &[]float64{float64(0x80) / 255., float64(0x80) / 255., float64(0x80) / 255.},
	"dimgray":        &[]float64{float64(0x69) / 255., float64(0x69) / 255., float64(0x69) / 255.},
	"lightslategray": &[]float64{float64(0x77) / 255., float64(0x88) / 255., float64(0x99) / 255.},
	"slategray":      &[]float64{float64(0x70) / 255., float64(0x80) / 255., float64(0x90) / 255.},
	"darkslategray":  &[]float64{float64(0x2F) / 255., float64(0x4F) / 255., float64(0x4F) / 255.},
	"black":          &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	//! Brown colors
	"cornsilk":       &[]float64{float64(0xFF) / 255., float64(0xF8) / 255., float64(0xDC) / 255.},
	"blanchedalmond": &[]float64{float64(0xFF) / 255., float64(0xEB) / 255., float64(0xCD) / 255.},
	"bisque":         &[]float64{float64(0xFF) / 255., float64(0xE4) / 255., float64(0xC4) / 255.},
	"navajowhite":    &[]float64{float64(0xFF) / 255., float64(0xDE) / 255., float64(0xAD) / 255.},
	"wheat":          &[]float64{float64(0xF5) / 255., float64(0xDE) / 255., float64(0xB3) / 255.},
	"burlywood":      &[]float64{float64(0xDE) / 255., float64(0xB8) / 255., float64(0x87) / 255.},
	"tan":            &[]float64{float64(0xD2) / 255., float64(0xB4) / 255., float64(0x8C) / 255.},
	"rosybrown":      &[]float64{float64(0xBC) / 255., float64(0x8F) / 255., float64(0x8F) / 255.},
	"sandybrown":     &[]float64{float64(0xF4) / 255., float64(0xA4) / 255., float64(0x60) / 255.},
	"goldenrod":      &[]float64{float64(0xDA) / 255., float64(0xA5) / 255., float64(0x20) / 255.},
	"peru":           &[]float64{float64(0xCD) / 255., float64(0x85) / 255., float64(0x3F) / 255.},
	"chocolate":      &[]float64{float64(0xD2) / 255., float64(0x69) / 255., float64(0x1E) / 255.},
	"saddlebrown":    &[]float64{float64(0x8B) / 255., float64(0x45) / 255., float64(0x13) / 255.},
	"sienna":         &[]float64{float64(0xA0) / 255., float64(0x52) / 255., float64(0x2D) / 255.},
	"brown":          &[]float64{float64(0xA5) / 255., float64(0x2A) / 255., float64(0x2A) / 255.},
	"maroon":         &[]float64{float64(0x80) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	//! X11 COLORS
	"alice-blue":    &[]float64{float64(0xf0) / 255., float64(0xf8) / 255., float64(0xff) / 255.},
	"antique-white": &[]float64{float64(0xfa) / 255., float64(0xeb) / 255., float64(0xd7) / 255.},
	//"aqua": &[]float64{float64(0x00)/255.,float64(0xff)/255.,float64(0xff)/255.},
	//"aquamarine": &[]float64{float64(0x7f)/255.,float64(0xff)/255.,float64(0xd4)/255.},
	//"azure": &[]float64{float64(0xf0)/255.,float64(0xff)/255.,float64(0xff)/255.},
	//"beige": &[]float64{float64(0xf5)/255.,float64(0xf5)/255.,float64(0xdc)/255.},
	//"bisque": &[]float64{float64(0xff)/255.,float64(0xe4)/255.,float64(0xc4)/255.},
	//"black": &[]float64{float64(0x00)/255.,float64(0x00)/255.,float64(0x00)/255.},
	"blanched-almond": &[]float64{float64(0xff) / 255., float64(0xeb) / 255., float64(0xcd) / 255.},
	//"blue": &[]float64{float64(0x00)/255.,float64(0x00)/255.,float64(0xff)/255.},
	"blue-violet": &[]float64{float64(0x8a) / 255., float64(0x2b) / 255., float64(0xe2) / 255.},
	//"brown": &[]float64{float64(0xa5)/255.,float64(0x2a)/255.,float64(0x2a)/255.},
	//"burlywood": &[]float64{float64(0xde)/255.,float64(0xb8)/255.,float64(0x87)/255.},
	"cadet-blue": &[]float64{float64(0x5f) / 255., float64(0x9e) / 255., float64(0xa0) / 255.},
	//"chartreuse": &[]float64{float64(0x7f)/255.,float64(0xff)/255.,float64(0x00)/255.},
	//"chocolate": &[]float64{float64(0xd2)/255.,float64(0x69)/255.,float64(0x1e)/255.},
	//"coral": &[]float64{float64(0xff)/255.,float64(0x7f)/255.,float64(0x50)/255.},
	"cornflower-blue": &[]float64{float64(0x64) / 255., float64(0x95) / 255., float64(0xed) / 255.},
	//"cornsilk": &[]float64{float64(0xff)/255.,float64(0xf8)/255.,float64(0xdc)/255.},
	//"crimson": &[]float64{float64(0xdc)/255.,float64(0x14)/255.,float64(0x3c)/255.},
	//"cyan": &[]float64{float64(0x00)/255.,float64(0xff)/255.,float64(0xff)/255.},
	"dark-blue":        &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0x8b) / 255.},
	"dark-cyan":        &[]float64{float64(0x00) / 255., float64(0x8b) / 255., float64(0x8b) / 255.},
	"dark-goldenrod":   &[]float64{float64(0xb8) / 255., float64(0x86) / 255., float64(0x0b) / 255.},
	"dark-gray":        &[]float64{float64(0xa9) / 255., float64(0xa9) / 255., float64(0xa9) / 255.},
	"dark-green":       &[]float64{float64(0x00) / 255., float64(0x64) / 255., float64(0x00) / 255.},
	"dark-khaki":       &[]float64{float64(0xbd) / 255., float64(0xb7) / 255., float64(0x6b) / 255.},
	"dark-magenta":     &[]float64{float64(0x8b) / 255., float64(0x00) / 255., float64(0x8b) / 255.},
	"dark-olive-green": &[]float64{float64(0x55) / 255., float64(0x6b) / 255., float64(0x2f) / 255.},
	"dark-orange":      &[]float64{float64(0xff) / 255., float64(0x8c) / 255., float64(0x00) / 255.},
	"dark-orchid":      &[]float64{float64(0x99) / 255., float64(0x32) / 255., float64(0xcc) / 255.},
	"dark-red":         &[]float64{float64(0x8b) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	"dark-salmon":      &[]float64{float64(0xe9) / 255., float64(0x96) / 255., float64(0x7a) / 255.},
	"dark-sea-green":   &[]float64{float64(0x8f) / 255., float64(0xbc) / 255., float64(0x8f) / 255.},
	"dark-slate-blue":  &[]float64{float64(0x48) / 255., float64(0x3d) / 255., float64(0x8b) / 255.},
	"dark-slate-gray":  &[]float64{float64(0x2f) / 255., float64(0x4f) / 255., float64(0x4f) / 255.},
	"dark-turquoise":   &[]float64{float64(0x00) / 255., float64(0xce) / 255., float64(0xd1) / 255.},
	"dark-violet":      &[]float64{float64(0x94) / 255., float64(0x00) / 255., float64(0xd3) / 255.},
	"deep-pink":        &[]float64{float64(0xff) / 255., float64(0x14) / 255., float64(0x93) / 255.},
	"deep-sky-blue":    &[]float64{float64(0x00) / 255., float64(0xbf) / 255., float64(0xff) / 255.},
	"dim-gray":         &[]float64{float64(0x69) / 255., float64(0x69) / 255., float64(0x69) / 255.},
	"dodger-blue":      &[]float64{float64(0x1e) / 255., float64(0x90) / 255., float64(0xff) / 255.},
	//"firebrick": &[]float64{float64(0xb2)/255.,float64(0x22)/255.,float64(0x22)/255.},
	"floral-white": &[]float64{float64(0xff) / 255., float64(0xfa) / 255., float64(0xf0) / 255.},
	"forest-green": &[]float64{float64(0x22) / 255., float64(0x8b) / 255., float64(0x22) / 255.},
	//"fuchsia": &[]float64{float64(0xff)/255.,float64(0x00)/255.,float64(0xff)/255.},
	//"gainsboro": &[]float64{float64(0xdc)/255.,float64(0xdc)/255.,float64(0xdc)/255.},
	"ghost-white": &[]float64{float64(0xf8) / 255., float64(0xf8) / 255., float64(0xff) / 255.},
	//"gold": &[]float64{float64(0xff)/255.,float64(0xd7)/255.,float64(0x00)/255.},
	//"goldenrod": &[]float64{float64(0xda)/255.,float64(0xa5)/255.,float64(0x20)/255.},
	//"gray": &[]float64{float64(0xbe)/255.,float64(0xbe)/255.,float64(0xbe)/255.},
	"web-gray": &[]float64{float64(0x80) / 255., float64(0x80) / 255., float64(0x80) / 255.},
	//"green": &[]float64{float64(0x00)/255.,float64(0xff)/255.,float64(0x00)/255.},
	"web-green":    &[]float64{float64(0x00) / 255., float64(0x80) / 255., float64(0x00) / 255.},
	"green-yellow": &[]float64{float64(0xad) / 255., float64(0xff) / 255., float64(0x2f) / 255.},
	//"honeydew": &[]float64{float64(0xf0)/255.,float64(0xff)/255.,float64(0xf0)/255.},
	"hot-pink":   &[]float64{float64(0xff) / 255., float64(0x69) / 255., float64(0xb4) / 255.},
	"indian-red": &[]float64{float64(0xcd) / 255., float64(0x5c) / 255., float64(0x5c) / 255.},
	//"indigo": &[]float64{float64(0x4b)/255.,float64(0x00)/255.,float64(0x82)/255.},
	//"ivory": &[]float64{float64(0xff)/255.,float64(0xff)/255.,float64(0xf0)/255.},
	//"khaki": &[]float64{float64(0xf0)/255.,float64(0xe6)/255.,float64(0x8c)/255.},
	//"lavender": &[]float64{float64(0xe6)/255.,float64(0xe6)/255.,float64(0xfa)/255.},
	"lavender-blush":   &[]float64{float64(0xff) / 255., float64(0xf0) / 255., float64(0xf5) / 255.},
	"lawn-green":       &[]float64{float64(0x7c) / 255., float64(0xfc) / 255., float64(0x00) / 255.},
	"lemon-chiffon":    &[]float64{float64(0xff) / 255., float64(0xfa) / 255., float64(0xcd) / 255.},
	"light-blue":       &[]float64{float64(0xad) / 255., float64(0xd8) / 255., float64(0xe6) / 255.},
	"light-coral":      &[]float64{float64(0xf0) / 255., float64(0x80) / 255., float64(0x80) / 255.},
	"light-cyan":       &[]float64{float64(0xe0) / 255., float64(0xff) / 255., float64(0xff) / 255.},
	"light-goldenrod":  &[]float64{float64(0xfa) / 255., float64(0xfa) / 255., float64(0xd2) / 255.},
	"light-gray":       &[]float64{float64(0xd3) / 255., float64(0xd3) / 255., float64(0xd3) / 255.},
	"light-green":      &[]float64{float64(0x90) / 255., float64(0xee) / 255., float64(0x90) / 255.},
	"light-pink":       &[]float64{float64(0xff) / 255., float64(0xb6) / 255., float64(0xc1) / 255.},
	"light-salmon":     &[]float64{float64(0xff) / 255., float64(0xa0) / 255., float64(0x7a) / 255.},
	"light-sea-green":  &[]float64{float64(0x20) / 255., float64(0xb2) / 255., float64(0xaa) / 255.},
	"light-sky-blue":   &[]float64{float64(0x87) / 255., float64(0xce) / 255., float64(0xfa) / 255.},
	"light-slate-gray": &[]float64{float64(0x77) / 255., float64(0x88) / 255., float64(0x99) / 255.},
	"light-steel-blue": &[]float64{float64(0xb0) / 255., float64(0xc4) / 255., float64(0xde) / 255.},
	"light-yellow":     &[]float64{float64(0xff) / 255., float64(0xff) / 255., float64(0xe0) / 255.},
	//"lime": &[]float64{float64(0x00)/255.,float64(0xff)/255.,float64(0x00)/255.},
	"lime-green": &[]float64{float64(0x32) / 255., float64(0xcd) / 255., float64(0x32) / 255.},
	//"linen": &[]float64{float64(0xfa)/255.,float64(0xf0)/255.,float64(0xe6)/255.},
	//"magenta": &[]float64{float64(0xff)/255.,float64(0x00)/255.,float64(0xff)/255.},
	//"maroon": &[]float64{float64(0xb0)/255.,float64(0x30)/255.,float64(0x60)/255.},
	"web-maroon":          &[]float64{float64(0x80) / 255., float64(0x00) / 255., float64(0x00) / 255.},
	"medium-aquamarine":   &[]float64{float64(0x66) / 255., float64(0xcd) / 255., float64(0xaa) / 255.},
	"medium-blue":         &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0xcd) / 255.},
	"medium-orchid":       &[]float64{float64(0xba) / 255., float64(0x55) / 255., float64(0xd3) / 255.},
	"medium-purple":       &[]float64{float64(0x93) / 255., float64(0x70) / 255., float64(0xdb) / 255.},
	"medium-sea-green":    &[]float64{float64(0x3c) / 255., float64(0xb3) / 255., float64(0x71) / 255.},
	"medium-slate-blue":   &[]float64{float64(0x7b) / 255., float64(0x68) / 255., float64(0xee) / 255.},
	"medium-spring-green": &[]float64{float64(0x00) / 255., float64(0xfa) / 255., float64(0x9a) / 255.},
	"medium-turquoise":    &[]float64{float64(0x48) / 255., float64(0xd1) / 255., float64(0xcc) / 255.},
	"medium-violet-red":   &[]float64{float64(0xc7) / 255., float64(0x15) / 255., float64(0x85) / 255.},
	"midnight-blue":       &[]float64{float64(0x19) / 255., float64(0x19) / 255., float64(0x70) / 255.},
	"mint-cream":          &[]float64{float64(0xf5) / 255., float64(0xff) / 255., float64(0xfa) / 255.},
	"misty-rose":          &[]float64{float64(0xff) / 255., float64(0xe4) / 255., float64(0xe1) / 255.},
	//"moccasin": &[]float64{float64(0xff)/255.,float64(0xe4)/255.,float64(0xb5)/255.},
	"navajo-white": &[]float64{float64(0xff) / 255., float64(0xde) / 255., float64(0xad) / 255.},
	"navy-blue":    &[]float64{float64(0x00) / 255., float64(0x00) / 255., float64(0x80) / 255.},
	"old-lace":     &[]float64{float64(0xfd) / 255., float64(0xf5) / 255., float64(0xe6) / 255.},
	//"olive": &[]float64{float64(0x80)/255.,float64(0x80)/255.,float64(0x00)/255.},
	"olive-drab": &[]float64{float64(0x6b) / 255., float64(0x8e) / 255., float64(0x23) / 255.},
	//"orange": &[]float64{float64(0xff)/255.,float64(0xa5)/255.,float64(0x00)/255.},
	"orange-red": &[]float64{float64(0xff) / 255., float64(0x45) / 255., float64(0x00) / 255.},
	//"orchid": &[]float64{float64(0xda)/255.,float64(0x70)/255.,float64(0xd6)/255.},
	"pale-goldenrod":  &[]float64{float64(0xee) / 255., float64(0xe8) / 255., float64(0xaa) / 255.},
	"pale-green":      &[]float64{float64(0x98) / 255., float64(0xfb) / 255., float64(0x98) / 255.},
	"pale-turquoise":  &[]float64{float64(0xaf) / 255., float64(0xee) / 255., float64(0xee) / 255.},
	"pale-violet-red": &[]float64{float64(0xdb) / 255., float64(0x70) / 255., float64(0x93) / 255.},
	"papaya-whip":     &[]float64{float64(0xff) / 255., float64(0xef) / 255., float64(0xd5) / 255.},
	"peach-puff":      &[]float64{float64(0xff) / 255., float64(0xda) / 255., float64(0xb9) / 255.},
	//"peru": &[]float64{float64(0xcd)/255.,float64(0x85)/255.,float64(0x3f)/255.},
	//"pink": &[]float64{float64(0xff)/255.,float64(0xc0)/255.,float64(0xcb)/255.},
	//"plum": &[]float64{float64(0xdd)/255.,float64(0xa0)/255.,float64(0xdd)/255.},
	"powder-blue": &[]float64{float64(0xb0) / 255., float64(0xe0) / 255., float64(0xe6) / 255.},
	//"purple": &[]float64{float64(0xa0)/255.,float64(0x20)/255.,float64(0xf0)/255.},
	"web-purple":     &[]float64{float64(0x80) / 255., float64(0x00) / 255., float64(0x80) / 255.},
	"rebecca-purple": &[]float64{float64(0x66) / 255., float64(0x33) / 255., float64(0x99) / 255.},
	//"red": &[]float64{float64(0xff)/255.,float64(0x00)/255.,float64(0x00)/255.},
	"rosy-brown":   &[]float64{float64(0xbc) / 255., float64(0x8f) / 255., float64(0x8f) / 255.},
	"royal-blue":   &[]float64{float64(0x41) / 255., float64(0x69) / 255., float64(0xe1) / 255.},
	"saddle-brown": &[]float64{float64(0x8b) / 255., float64(0x45) / 255., float64(0x13) / 255.},
	//"salmon": &[]float64{float64(0xfa)/255.,float64(0x80)/255.,float64(0x72)/255.},
	"sandy-brown": &[]float64{float64(0xf4) / 255., float64(0xa4) / 255., float64(0x60) / 255.},
	"sea-green":   &[]float64{float64(0x2e) / 255., float64(0x8b) / 255., float64(0x57) / 255.},
	//"seashell": &[]float64{float64(0xff)/255.,float64(0xf5)/255.,float64(0xee)/255.},
	//"sienna": &[]float64{float64(0xa0)/255.,float64(0x52)/255.,float64(0x2d)/255.},
	//"silver": &[]float64{float64(0xc0)/255.,float64(0xc0)/255.,float64(0xc0)/255.},
	"sky-blue":   &[]float64{float64(0x87) / 255., float64(0xce) / 255., float64(0xeb) / 255.},
	"slate-blue": &[]float64{float64(0x6a) / 255., float64(0x5a) / 255., float64(0xcd) / 255.},
	"slate-gray": &[]float64{float64(0x70) / 255., float64(0x80) / 255., float64(0x90) / 255.},
	//"snow": &[]float64{float64(0xff)/255.,float64(0xfa)/255.,float64(0xfa)/255.},
	"spring-green": &[]float64{float64(0x00) / 255., float64(0xff) / 255., float64(0x7f) / 255.},
	"steel-blue":   &[]float64{float64(0x46) / 255., float64(0x82) / 255., float64(0xb4) / 255.},
	//"tan": &[]float64{float64(0xd2)/255.,float64(0xb4)/255.,float64(0x8c)/255.},
	//"teal": &[]float64{float64(0x00)/255.,float64(0x80)/255.,float64(0x80)/255.},
	//"thistle": &[]float64{float64(0xd8)/255.,float64(0xbf)/255.,float64(0xd8)/255.},
	//"tomato": &[]float64{float64(0xff)/255.,float64(0x63)/255.,float64(0x47)/255.},
	//"turquoise": &[]float64{float64(0x40)/255.,float64(0xe0)/255.,float64(0xd0)/255.},
	//"violet": &[]float64{float64(0xee)/255.,float64(0x82)/255.,float64(0xee)/255.},
	//"wheat": &[]float64{float64(0xf5)/255.,float64(0xde)/255.,float64(0xb3)/255.},
	//"white": &[]float64{float64(0xff)/255.,float64(0xff)/255.,float64(0xff)/255.},
	"white-smoke": &[]float64{float64(0xf5) / 255., float64(0xf5) / 255., float64(0xf5) / 255.},
	//"yellow": &[]float64{float64(0xff)/255.,float64(0xff)/255.,float64(0x00)/255.},
	"yellow-green": &[]float64{float64(0x9a) / 255., float64(0xcd) / 255., float64(0x32) / 255.},
	//! PWG 5101.1 Media Color Names
	"dark-brown":      &[]float64{float64(0x5c) / 255., float64(0x40) / 255., float64(0x33) / 255.},
	"light-brown":     &[]float64{float64(0x99) / 255., float64(0x66) / 255., float64(0xff) / 255.},
	"dark-buff":       &[]float64{float64(0x97) / 255., float64(0x66) / 255., float64(0x38) / 255.},
	"light-buff":      &[]float64{float64(0xec) / 255., float64(0xd9) / 255., float64(0xb0) / 255.},
	"dark-gold":       &[]float64{float64(0xee) / 255., float64(0xbc) / 255., float64(0x1d) / 255.},
	"light-gold":      &[]float64{float64(0xf1) / 255., float64(0xe5) / 255., float64(0xac) / 255.},
	"dark-ivory":      &[]float64{float64(0xf2) / 255., float64(0xe5) / 255., float64(0x8f) / 255.},
	"light-ivory":     &[]float64{float64(0xff) / 255., float64(0xf8) / 255., float64(0xc9) / 255.},
	"light-magenta":   &[]float64{float64(0xff) / 255., float64(0x77) / 255., float64(0xff) / 255.},
	"mustard":         &[]float64{float64(0xff) / 255., float64(0xdb) / 255., float64(0x58) / 255.},
	"dark-mustard":    &[]float64{float64(0x7c) / 255., float64(0x7c) / 255., float64(0x40) / 255.},
	"light-mustard":   &[]float64{float64(0xee) / 255., float64(0xdd) / 255., float64(0x62) / 255.},
	"light-orange":    &[]float64{float64(0xd9) / 255., float64(0xa4) / 255., float64(0x65) / 255.},
	"dark-pink":       &[]float64{float64(0xe7) / 255., float64(0x54) / 255., float64(0x80) / 255.},
	"light-red":       &[]float64{float64(0xff) / 255., float64(0x33) / 255., float64(0x33) / 255.},
	"dark-silver":     &[]float64{float64(0xaf) / 255., float64(0xaf) / 255., float64(0xaf) / 255.},
	"light-silver":    &[]float64{float64(0xe1) / 255., float64(0xe1) / 255., float64(0xe1) / 255.},
	"light-turquoise": &[]float64{float64(0xaf) / 255., float64(0xe4) / 255., float64(0xde) / 255.},
	"light-violet":    &[]float64{float64(0x7a) / 255., float64(0x52) / 255., float64(0x99) / 255.},
	"dark-yellow":     &[]float64{float64(0xff) / 255., float64(0xcc) / 255., float64(0x00) / 255.},
	//! material design flat ui colors
	//! https://flatuicolors.com/
	"md-flat-ui-turquoise":    &[]float64{float64(0x1A) / 255., float64(0xBC) / 255., float64(0x9C) / 255.},
	"md-flat-ui-emerland":     &[]float64{float64(0x2E) / 255., float64(0xCC) / 255., float64(0x71) / 255.},
	"md-flat-ui-peterriver":   &[]float64{float64(0x34) / 255., float64(0x98) / 255., float64(0xdb) / 255.},
	"md-flat-ui-amethyst":     &[]float64{float64(0x9b) / 255., float64(0x59) / 255., float64(0xb6) / 255.},
	"md-flat-ui-wetasphalt":   &[]float64{float64(0x34) / 255., float64(0x49) / 255., float64(0x5e) / 255.},
	"md-flat-ui-greensea":     &[]float64{float64(0x16) / 255., float64(0xa0) / 255., float64(0x85) / 255.},
	"md-flat-ui-nephritis":    &[]float64{float64(0x27) / 255., float64(0xae) / 255., float64(0x60) / 255.},
	"md-flat-ui-belizehole":   &[]float64{float64(0x29) / 255., float64(0x80) / 255., float64(0xb9) / 255.},
	"md-flat-ui-wisteria":     &[]float64{float64(0x8e) / 255., float64(0x44) / 255., float64(0xad) / 255.},
	"md-flat-ui-midnightblue": &[]float64{float64(0x2c) / 255., float64(0x3e) / 255., float64(0x50) / 255.},
	"md-flat-ui-sunflower":    &[]float64{float64(0xf1) / 255., float64(0xc4) / 255., float64(0x0f) / 255.},
	"md-flat-ui-carrot":       &[]float64{float64(0xe6) / 255., float64(0x7e) / 255., float64(0x22) / 255.},
	"md-flat-ui-alizarin":     &[]float64{float64(0xe7) / 255., float64(0x4c) / 255., float64(0x3c) / 255.},
	"md-flat-ui-clouds":       &[]float64{float64(0xec) / 255., float64(0xf0) / 255., float64(0xf1) / 255.},
	"md-flat-ui-concrete":     &[]float64{float64(0x95) / 255., float64(0xa5) / 255., float64(0xa6) / 255.},
	"md-flat-ui-orange":       &[]float64{float64(0xf3) / 255., float64(0x9c) / 255., float64(0x12) / 255.},
	"md-flat-ui-pumpkin":      &[]float64{float64(0xd3) / 255., float64(0x54) / 255., float64(0x00) / 255.},
	"md-flat-ui-pomegranate":  &[]float64{float64(0xc0) / 255., float64(0x39) / 255., float64(0x2b) / 255.},
	"md-flat-ui-silver":       &[]float64{float64(0xbd) / 255., float64(0xc3) / 255., float64(0xc7) / 255.},
	"md-flat-ui-asbestos":     &[]float64{float64(0x7f) / 255., float64(0x8c) / 255., float64(0x8d) / 255.},
}

func ConvertColorToFloats(c string) []float64 {
	c = strings.ToLower(c)
	//lookup colors first
	if rgbColors[c] != nil {
		return *rgbColors[c]
	} else if c[0] == '#' {
		if len(c) == 4 {
			_ri, _ := strconv.ParseUint(c[1:2], 16, 32)
			_gi, _ := strconv.ParseUint(c[2:3], 16, 32)
			_bi, _ := strconv.ParseUint(c[3:4], 16, 32)

			_rf := float64(_ri) / 15.0
			_gf := float64(_gi) / 15.0
			_bf := float64(_bi) / 15.0
			return []float64{_rf, _gf, _bf}
		} else if len(c) == 7 {
			_ri, _ := strconv.ParseUint(c[1:3], 16, 32)
			_gi, _ := strconv.ParseUint(c[3:5], 16, 32)
			_bi, _ := strconv.ParseUint(c[5:7], 16, 32)

			_rf := float64(_ri) / 255.0
			_gf := float64(_gi) / 255.0
			_bf := float64(_bi) / 255.0
			return []float64{_rf, _gf, _bf}
		} else if len(c) == 10 {
			_ri, _ := strconv.ParseUint(c[1:4], 16, 32)
			_gi, _ := strconv.ParseUint(c[4:7], 16, 32)
			_bi, _ := strconv.ParseUint(c[7:10], 16, 32)

			_rf := float64(_ri) / 4095.0
			_gf := float64(_gi) / 4095.0
			_bf := float64(_bi) / 4095.0
			return []float64{_rf, _gf, _bf}
		} else if len(c) == 9 {
			_ci, _ := strconv.ParseUint(c[1:3], 16, 32)
			_mi, _ := strconv.ParseUint(c[3:5], 16, 32)
			_yi, _ := strconv.ParseUint(c[5:7], 16, 32)
			_ki, _ := strconv.ParseUint(c[7:9], 16, 32)

			_cf := float64(_ci) / 255.0
			_mf := float64(_mi) / 255.0
			_yf := float64(_yi) / 255.0
			_kf := float64(_ki) / 255.0
			return []float64{_cf, _mf, _yf, _kf}
		} else if len(c) == 5 {
			_ci, _ := strconv.ParseUint(c[1:2], 16, 32)
			_mi, _ := strconv.ParseUint(c[2:3], 16, 32)
			_yi, _ := strconv.ParseUint(c[3:4], 16, 32)
			_ki, _ := strconv.ParseUint(c[4:5], 16, 32)

			_cf := float64(_ci) / 15.0
			_mf := float64(_mi) / 15.0
			_yf := float64(_yi) / 15.0
			_kf := float64(_ki) / 15.0
			return []float64{_cf, _mf, _yf, _kf}
		} else if len(c) == 2 {
			_i, _ := strconv.ParseUint(c[1:2], 16, 32)

			_f := float64(_i) / 16.0
			return []float64{_f}
		} else if len(c) == 3 {
			_i, _ := strconv.ParseUint(c[1:3], 16, 32)

			_f := float64(_i) / 255.0
			return []float64{_f}
		}
	}
	return []float64{0.}
}

func (pco *PdfContentObject) SetFillColor(c string) {
	_col := ConvertColorToFloats(c)
	if len(_col) == 1 {
		pco.SetFillColorGrey(_col[0])
	} else if len(_col) == 3 {
		pco.SetFillColorRgb(_col[0], _col[1], _col[2])
	} else if len(_col) == 4 {
		pco.SetFillColorCmyk(_col[0], _col[1], _col[2], _col[3])
	}
}
func (pco *PdfContentObject) SetStrokeColor(c string) {
	_col := ConvertColorToFloats(c)
	if len(_col) == 1 {
		pco.SetStrokeColorGrey(_col[0])
	} else if len(_col) == 3 {
		pco.SetStrokeColorRgb(_col[0], _col[1], _col[2])
	} else if len(_col) == 4 {
		pco.SetStrokeColorCmyk(_col[0], _col[1], _col[2], _col[3])
	}
}

func ColorHslToRgb(h float64, s float64, l float64) (r float64, g float64, b float64) {
	h = math.Mod(h, 360.)
	h = h / 360.
	s = s / 100.
	l = l / 100.
	var q float64 = 0
	if l < .5 {
		q = l * (1. + s)
	} else {
		q = (l + s) - (s * l)
	}
	p := 2*l - q
	r = max(0, HueToRgb(p, q, h+(1./3.)))
	g = max(0, HueToRgb(p, q, h))
	b = max(0, HueToRgb(p, q, h-(1./3.)))
	return
}

func ColorHsvToRgb(h float64, s float64, v float64) (r float64, g float64, b float64) {
	h = math.Mod(h, 360.)
	h = h / 60.
	s = s / 100.
	v = v / 100.

	var hi int = int(h) % 6

	f := h - math.Floor(h)
	p := v * (1. - s)
	q := v * (1. - (s * f))
	t := v * (1. - (s * (1. - f)))

	switch hi {
	case 1:
		{
			r = q
			g = v
			b = p
		}
	case 2:
		{
			r = p
			g = v
			b = t
		}
	case 3:
		{
			r = p
			g = q
			b = v
		}
	case 4:
		{
			r = t
			g = p
			b = v
		}
	case 5:
		{
			r = v
			g = p
			b = q
		}
	default:
		{
			r = v
			g = t
			b = p
		}
	}
	return
}

func HueToRgb(p float64, q float64, h float64) float64 {
	if h < 0 {
		h += 1.
	}

	if h > 1. {
		h -= 1.
	}

	if 6*h < 1 {
		return p + ((q - p) * 6 * h)
	}
	if 2*h < 1 {
		return q
	}
	if 3*h < 2 {
		return p + ((q - p) * 6 * ((2. / 3.) + h))
	}
	return p
}

// http://dev.w3.org/csswg/css-color/#hwb-to-rgb
func ColorHwbToRgb(h float64, wh float64, bi float64) (r float64, g float64, b float64) {
	h = math.Mod(h, 360.)
	h = h / 360.
	wh = wh / 100.
	bi = bi / 100.

	var ratio float64 = wh + bi

	if ratio > 1. {
		wh /= ratio
		bi /= ratio
	}

	i := int(math.Floor(h * 6))
	v := 1. - bi
	f := (6 * h) - float64(i)

	if (i & 1) != 0 {
		f = 1. - f
	}

	n := wh + f*(v-wh)

	switch i {
	case 1:
		{
			r = n
			g = v
			b = wh
		}
	case 2:
		{
			r = wh
			g = v
			b = n
		}
	case 3:
		{
			r = wh
			g = n
			b = v
		}
	case 4:
		{
			r = n
			g = wh
			b = v
		}
	case 5:
		{
			r = v
			g = wh
			b = n
		}
	default:
		{
			r = v
			g = n
			b = wh
		}
	}
	return
}

// TODO HCG XYZ Lab Lch

func (pco *PdfContentObject) SetFillColorHsl(h float64, s float64, l float64) {
	r, g, b := ColorHslToRgb(h, s, l)
	pco.SetFillColorRgb(r, g, b)
}
func (pco *PdfContentObject) SetStrokeColorHsl(h float64, s float64, l float64) {
	r, g, b := ColorHslToRgb(h, s, l)
	pco.SetStrokeColorRgb(r, g, b)
}

func (pco *PdfContentObject) SetFillColorHsv(h float64, s float64, v float64) {
	r, g, b := ColorHsvToRgb(h, s, v)
	pco.SetFillColorRgb(r, g, b)
}
func (pco *PdfContentObject) SetStrokeColorHsv(h float64, s float64, v float64) {
	r, g, b := ColorHsvToRgb(h, s, v)
	pco.SetStrokeColorRgb(r, g, b)
}

func (pco *PdfContentObject) SetFillColorHwb(h float64, w float64, b float64) {
	r, g, b := ColorHsvToRgb(h, w, b)
	pco.SetFillColorRgb(r, g, b)
}
func (pco *PdfContentObject) SetStrokeColorHwb(h float64, w float64, b float64) {
	r, g, b := ColorHsvToRgb(h, w, b)
	pco.SetStrokeColorRgb(r, g, b)
}
