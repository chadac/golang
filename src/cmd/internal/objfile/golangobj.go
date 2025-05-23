// Copyright 2013 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Parsing of Golang intermediate object files and archives.

package objfile

import (
	"cmd/internal/archive"
	"cmd/internal/golangobj"
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"debug/dwarf"
	"debug/golangsym"
	"errors"
	"fmt"
	"io"
	"os"
)

type golangobjFile struct {
	golangobj *archive.GolangObj
	r     *golangobj.Reader
	f     *os.File
	arch  *sys.Arch
}

func openGolangFile(f *os.File) (*File, error) {
	a, err := archive.Parse(f, false)
	if err != nil {
		return nil, err
	}
	entries := make([]*Entry, 0, len(a.Entries))
L:
	for _, e := range a.Entries {
		switch e.Type {
		case archive.EntryPkgDef, archive.EntrySentinelNonObj:
			continue
		case archive.EntryGolangObj:
			o := e.Obj
			b := make([]byte, o.Size)
			_, err := f.ReadAt(b, o.Offset)
			if err != nil {
				return nil, err
			}
			r := golangobj.NewReaderFromBytes(b, false)
			var arch *sys.Arch
			for _, a := range sys.Archs {
				if a.Name == e.Obj.Arch {
					arch = a
					break
				}
			}
			entries = append(entries, &Entry{
				name: e.Name,
				raw:  &golangobjFile{e.Obj, r, f, arch},
			})
			continue
		case archive.EntryNativeObj:
			nr := io.NewSectionReader(f, e.Offset, e.Size)
			for _, try := range openers {
				if raw, err := try(nr); err == nil {
					entries = append(entries, &Entry{
						name: e.Name,
						raw:  raw,
					})
					continue L
				}
			}
		}
		return nil, fmt.Errorf("open %s: unrecognized archive member %s", f.Name(), e.Name)
	}
	return &File{f, entries}, nil
}

func golangobjName(name string, ver int) string {
	if ver == 0 {
		return name
	}
	return fmt.Sprintf("%s<%d>", name, ver)
}

type golangobjReloc struct {
	Off  int32
	Size uint8
	Type objabi.RelocType
	Add  int64
	Sym  string
}

func (r golangobjReloc) String(insnOffset uint64) string {
	delta := int64(r.Off) - int64(insnOffset)
	s := fmt.Sprintf("[%d:%d]%s", delta, delta+int64(r.Size), r.Type)
	if r.Sym != "" {
		if r.Add != 0 {
			return fmt.Sprintf("%s:%s+%d", s, r.Sym, r.Add)
		}
		return fmt.Sprintf("%s:%s", s, r.Sym)
	}
	if r.Add != 0 {
		return fmt.Sprintf("%s:%d", s, r.Add)
	}
	return s
}

func (f *golangobjFile) symbols() ([]Sym, error) {
	r := f.r
	var syms []Sym

	// Name of referenced indexed symbols.
	nrefName := r.NRefName()
	refNames := make(map[golangobj.SymRef]string, nrefName)
	for i := 0; i < nrefName; i++ {
		rn := r.RefName(i)
		refNames[rn.Sym()] = rn.Name(r)
	}

	abiToVer := func(abi uint16) int {
		var ver int
		if abi == golangobj.SymABIstatic {
			// Static symbol
			ver = 1
		}
		return ver
	}

	resolveSymRef := func(s golangobj.SymRef) string {
		var i uint32
		switch p := s.PkgIdx; p {
		case golangobj.PkgIdxInvalid:
			if s.SymIdx != 0 {
				panic("bad sym ref")
			}
			return ""
		case golangobj.PkgIdxHashed64:
			i = s.SymIdx + uint32(r.NSym())
		case golangobj.PkgIdxHashed:
			i = s.SymIdx + uint32(r.NSym()+r.NHashed64def())
		case golangobj.PkgIdxNone:
			i = s.SymIdx + uint32(r.NSym()+r.NHashed64def()+r.NHasheddef())
		case golangobj.PkgIdxBuiltin:
			name, abi := golangobj.BuiltinName(int(s.SymIdx))
			return golangobjName(name, abi)
		case golangobj.PkgIdxSelf:
			i = s.SymIdx
		default:
			return refNames[s]
		}
		sym := r.Sym(i)
		return golangobjName(sym.Name(r), abiToVer(sym.ABI()))
	}

	// Defined symbols
	ndef := uint32(r.NSym() + r.NHashed64def() + r.NHasheddef() + r.NNonpkgdef())
	for i := uint32(0); i < ndef; i++ {
		osym := r.Sym(i)
		if osym.Name(r) == "" {
			continue // not a real symbol
		}
		name := osym.Name(r)
		ver := osym.ABI()
		name = golangobjName(name, abiToVer(ver))
		typ := objabi.SymKind(osym.Type())
		var code rune = '?'
		switch typ {
		case objabi.STEXT, objabi.STEXTFIPS:
			code = 'T'
		case objabi.SRODATA, objabi.SRODATAFIPS:
			code = 'R'
		case objabi.SNOPTRDATA, objabi.SNOPTRDATAFIPS,
			objabi.SDATA, objabi.SDATAFIPS:
			code = 'D'
		case objabi.SBSS, objabi.SNOPTRBSS, objabi.STLSBSS:
			code = 'B'
		}
		if ver >= golangobj.SymABIstatic {
			code += 'a' - 'A'
		}

		sym := Sym{
			Name: name,
			Addr: uint64(r.DataOff(i)),
			Size: int64(osym.Siz()),
			Code: code,
		}

		relocs := r.Relocs(i)
		sym.Relocs = make([]Reloc, len(relocs))
		for j := range relocs {
			rel := &relocs[j]
			sym.Relocs[j] = Reloc{
				Addr: uint64(r.DataOff(i)) + uint64(rel.Off()),
				Size: uint64(rel.Siz()),
				Stringer: golangobjReloc{
					Off:  rel.Off(),
					Size: rel.Siz(),
					Type: objabi.RelocType(rel.Type()),
					Add:  rel.Add(),
					Sym:  resolveSymRef(rel.Sym()),
				},
			}
		}

		syms = append(syms, sym)
	}

	// Referenced symbols
	n := ndef + uint32(r.NNonpkgref())
	for i := ndef; i < n; i++ {
		osym := r.Sym(i)
		sym := Sym{Name: osym.Name(r), Code: 'U'}
		syms = append(syms, sym)
	}
	for i := 0; i < nrefName; i++ {
		rn := r.RefName(i)
		sym := Sym{Name: rn.Name(r), Code: 'U'}
		syms = append(syms, sym)
	}

	return syms, nil
}

func (f *golangobjFile) pcln() (textStart uint64, symtab, pclntab []byte, err error) {
	// Should never be called. We implement Liner below, callers
	// should use that instead.
	return 0, nil, nil, fmt.Errorf("pcln not available in golang object file")
}

// Find returns the file name, line, and function data for the given pc.
// Returns "",0,nil if unknown.
// This function implements the Liner interface in preference to pcln() above.
func (f *golangobjFile) PCToLine(pc uint64) (string, int, *golangsym.Func) {
	r := f.r
	if f.arch == nil {
		return "", 0, nil
	}
	getSymData := func(s golangobj.SymRef) []byte {
		if s.PkgIdx != golangobj.PkgIdxHashed {
			// We don't need the data for non-hashed symbols, yet.
			panic("not supported")
		}
		i := uint32(s.SymIdx + uint32(r.NSym()+r.NHashed64def()))
		return r.BytesAt(r.DataOff(i), r.DataSize(i))
	}

	ndef := uint32(r.NSym() + r.NHashed64def() + r.NHasheddef() + r.NNonpkgdef())
	for i := uint32(0); i < ndef; i++ {
		osym := r.Sym(i)
		addr := uint64(r.DataOff(i))
		if pc < addr || pc >= addr+uint64(osym.Siz()) {
			continue
		}
		var pcfileSym, pclineSym golangobj.SymRef
		for _, a := range r.Auxs(i) {
			switch a.Type() {
			case golangobj.AuxPcfile:
				pcfileSym = a.Sym()
			case golangobj.AuxPcline:
				pclineSym = a.Sym()
			}
		}
		if pcfileSym.IsZero() || pclineSym.IsZero() {
			continue
		}
		pcline := getSymData(pclineSym)
		line := int(pcValue(pcline, pc-addr, f.arch))
		pcfile := getSymData(pcfileSym)
		fileID := pcValue(pcfile, pc-addr, f.arch)
		fileName := r.File(int(fileID))
		// Note: we provide only the name in the Func structure.
		// We could provide more if needed.
		return fileName, line, &golangsym.Func{Sym: &golangsym.Sym{Name: osym.Name(r)}}
	}
	return "", 0, nil
}

// pcValue looks up the given PC in a pc value table. target is the
// offset of the pc from the entry point.
func pcValue(tab []byte, target uint64, arch *sys.Arch) int32 {
	val := int32(-1)
	var pc uint64
	for step(&tab, &pc, &val, pc == 0, arch) {
		if target < pc {
			return val
		}
	}
	return -1
}

// step advances to the next pc, value pair in the encoded table.
func step(p *[]byte, pc *uint64, val *int32, first bool, arch *sys.Arch) bool {
	uvdelta := readvarint(p)
	if uvdelta == 0 && !first {
		return false
	}
	if uvdelta&1 != 0 {
		uvdelta = ^(uvdelta >> 1)
	} else {
		uvdelta >>= 1
	}
	vdelta := int32(uvdelta)
	pcdelta := readvarint(p) * uint32(arch.MinLC)
	*pc += uint64(pcdelta)
	*val += vdelta
	return true
}

// readvarint reads, removes, and returns a varint from *p.
func readvarint(p *[]byte) uint32 {
	var v, shift uint32
	s := *p
	for shift = 0; ; shift += 7 {
		b := s[0]
		s = s[1:]
		v |= (uint32(b) & 0x7F) << shift
		if b&0x80 == 0 {
			break
		}
	}
	*p = s
	return v
}

// We treat the whole object file as the text section.
func (f *golangobjFile) text() (textStart uint64, text []byte, err error) {
	text = make([]byte, f.golangobj.Size)
	_, err = f.f.ReadAt(text, int64(f.golangobj.Offset))
	return
}

func (f *golangobjFile) golangarch() string {
	return f.golangobj.Arch
}

func (f *golangobjFile) loadAddress() (uint64, error) {
	return 0, fmt.Errorf("unknown load address")
}

func (f *golangobjFile) dwarf() (*dwarf.Data, error) {
	return nil, errors.New("no DWARF data in golang object file")
}
