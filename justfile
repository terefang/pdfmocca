#!/usr/bin/env just --justfile
# https://github.com/casey/just/

CURL := 'curl --progress-bar -L '
XDIR := justfile_directory()
XFNT := XDIR+"/pdf/fonts"
EXE := "pdfmocca"

set-drel: inc-level
    #!/bin/bash
    V=$(date '+%Y.%m.')
    V=$V$(cd {{XDIR}} && shtool version -n "{{EXE}}" -l short ./version.txt|cut -f 3 -d.)
    cd {{XDIR}} && just -f justfile set-version "$V"

inc-version:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i v -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

inc-major:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i r -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

inc-level:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i l -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

set-version _VERSION:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -s "{{_VERSION}}" -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

make-release: set-drel
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    MESSAGE="{{EXE}} automated release version $(shtool version -l text -d long ./version.txt)"
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt
    gh release create v$VERSION --notes "$MESSAGE"

#fetch-go:
#    mkdir -p
#    curl -L -o go.tar.gz https://go.dev/dl/go1.26.4.linux-amd64.tar.gz

fetch-fonts:
    #!/bin/sh -x
    mkdir -p {{XFNT}}
    curl -L -o {{XFNT}}/qdb.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitDingbats.otf
    curl -L -o {{XFNT}}/qsy.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSymbol.otf
    # TODO --- since the pdfium fonts lack the full range of pdfdoc encoding
    #    curl -L -o {{XFNT}}/qcrr.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitFixed.otf
    #    curl -L -o {{XFNT}}/qcrb.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitFixedBold.otf
    #    curl -L -o {{XFNT}}/qcrz.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitFixedBoldItalic.otf
    #    curl -L -o {{XFNT}}/qcri.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitFixedItalic.otf
    #    curl -L -o {{XFNT}}/qhvr.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSans.otf
    #    curl -L -o {{XFNT}}/qhvb.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSansBold.otf
    #    curl -L -o {{XFNT}}/qhvz.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSansBoldItalic.otf
    #    curl -L -o {{XFNT}}/qhvi.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSansItalic.otf
    #    curl -L -o {{XFNT}}/qtmr.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSerif.otf
    #    curl -L -o {{XFNT}}/qtmb.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSerifBold.otf
    #    curl -L -o {{XFNT}}/qtmz.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSerifBoldItalic.otf
    #    curl -L -o {{XFNT}}/qtmi.otf https://github.com/terefang/pdfium-fonts/raw/refs/heads/master/dist/otf/FoxitSerifItalic.otf
    # TODO --- we use the texgyre ones instead because they are anyhow derived from the urw substitutes
    curl -L -o {{XFNT}}/qcrb.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyrecursor-bold.otf
    curl -L -o {{XFNT}}/qcrz.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyrecursor-bolditalic.otf
    curl -L -o {{XFNT}}/qcri.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyrecursor-italic.otf
    curl -L -o {{XFNT}}/qcrr.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyrecursor-regular.otf
    curl -L -o {{XFNT}}/qhvb.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyreheros-bold.otf
    curl -L -o {{XFNT}}/qhvz.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyreheros-bolditalic.otf
    curl -L -o {{XFNT}}/qhvi.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyreheros-italic.otf
    curl -L -o {{XFNT}}/qhvr.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyreheros-regular.otf
    curl -L -o {{XFNT}}/qtmb.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyretermes-bold.otf
    curl -L -o {{XFNT}}/qtmz.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyretermes-bolditalic.otf
    curl -L -o {{XFNT}}/qtmi.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyretermes-italic.otf
    curl -L -o {{XFNT}}/qtmr.otf https://github.com/debian-tex/tex-gyre/raw/refs/heads/master/fonts/opentype/public/tex-gyre/texgyretermes-regular.otf

build: fetch-fonts
    #!/bin/sh -x
    #CGO_ENABLED=1 CC=musl-gcc go build -ldflags="-linkmode external -extldflags '-static'" -o test-me main

