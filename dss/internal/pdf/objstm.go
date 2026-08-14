// Object streams. See DESIGN.md §2.4.
//
// Provenance: org.apache.pdfbox.pdfparser.PDFObjectStreamParser of pdfbox 3.0.7.

package pdf

import (
	"errors"
	"fmt"
	"strconv"
)

// objStm is a decoded /Type /ObjStm container.
type objStm struct {
	nums    []int64
	offsets []int64
	objects []Object // parallel to nums; nil when that entry failed to parse
	broken  bool
}

// objStmNumbers reads the (objnum, offset) pair list at the front of a container
// without parsing its contents.
func (d *Document) objStmNumbers(key ObjectKey, st *Stream) ([]int64, error) {
	data, err := d.StreamData(st)
	if err != nil {
		return nil, err
	}
	nums, _, err := parseObjStmHeader(data, st.Dict, &d.warnings, key)
	return nums, err
}

func parseObjStmHeader(data []byte, dict *Dict, warn *[]Warning, key ObjectKey) ([]int64, []int64, error) {
	n := dictInt(dict, "N", 0)
	first := dictInt(dict, "First", 0)
	if n < 0 || first < 0 || first > int64(len(data)) {
		return nil, nil, fmt.Errorf("pdf: object stream %s has a bad /N or /First", key)
	}
	l := newLexer(data[:first], warn)
	nums := make([]int64, 0, n)
	offs := make([]int64, 0, n)
	for i := int64(0); i < n; i++ {
		t1 := l.next()
		t2 := l.next()
		if t1.kind != tokInteger || t2.kind != tokInteger {
			addWarning(warn, WarnObjStmBroken, -1,
				"object stream "+key.String()+" header truncated after "+strconv.Itoa(len(nums))+" entries")
			break
		}
		nums = append(nums, t1.i)
		offs = append(offs, first+t2.i)
	}
	return nums, offs, nil
}

// loadObjStm decodes and parses a container, caching the result.
func (d *Document) loadObjStm(containerNum int64) (*objStm, error) {
	if os, ok := d.objStms[containerNum]; ok {
		return os, nil
	}
	key, ok := d.byNum[containerNum]
	if !ok {
		key = ObjectKey{Num: containerNum}
	}
	if e, ok := d.xref[key]; ok && e.typ == 2 {
		// An object stream may not live inside an object stream.
		return nil, fmt.Errorf("pdf: object stream %d is itself compressed", containerNum)
	}
	os := &objStm{}
	d.objStms[containerNum] = os // cache first: a cycle must not recurse forever

	obj, ok := d.object(key)
	if !ok {
		os.broken = true
		d.addWarning(WarnObjStmBroken, -1, "object stream "+key.String()+" is missing")
		return os, nil
	}
	st, isStream := obj.(*Stream)
	if !isStream {
		os.broken = true
		d.addWarning(WarnObjStmBroken, -1, "object "+key.String()+" is not a stream")
		return os, nil
	}
	data, err := d.StreamData(st)
	if err != nil {
		os.broken = true
		d.addWarning(WarnObjStmBroken, st.Offset,
			"object stream "+key.String()+" could not be decoded: "+err.Error())
		return os, nil
	}
	nums, offs, err := parseObjStmHeader(data, st.Dict, &d.warnings, key)
	if err != nil {
		os.broken = true
		d.addWarning(WarnObjStmBroken, st.Offset, err.Error())
		return os, nil
	}
	os.nums, os.offsets = nums, offs
	os.objects = make([]Object, len(nums))

	p := newParser(data, &d.warnings, d.opts.MaxDepth, d.opts.MaxStreamSize)
	p.resolveLength = d.lengthResolver()
	for i := range nums {
		if offs[i] < 0 || offs[i] > int64(len(data)) {
			d.addWarning(WarnObjStmBroken, -1,
				"object "+strconv.FormatInt(nums[i], 10)+" in stream "+key.String()+" is out of range")
			continue
		}
		p.lex.seek(offs[i])
		p.err = nil
		o := p.parseDirObject()
		if p.err != nil {
			d.addWarning(WarnObjStmBroken, -1, "object "+strconv.FormatInt(nums[i], 10)+
				" in stream "+key.String()+": "+p.err.Error())
			continue
		}
		// Streams cannot appear inside object streams: a `stream` keyword here is
		// a parse error for this object only.
		save := p.lex.pos
		p.lex.skipSpaces()
		if hasPrefixAt(p.lex.data, p.lex.pos, "stream") {
			d.addWarning(WarnObjStmBroken, -1,
				"object "+strconv.FormatInt(nums[i], 10)+" in stream "+key.String()+" contains a stream")
			p.lex.pos = save
			continue
		}
		p.lex.pos = save
		if o == nil {
			o = Null{}
		}
		os.objects[i] = o
	}
	return os, nil
}

// objectFromStm returns the object the type-2 entry names.
func (d *Document) objectFromStm(key ObjectKey, e xrefRec) (Object, bool) {
	os, err := d.loadObjStm(e.offset)
	if err != nil {
		d.addWarning(WarnObjStmBroken, -1, err.Error())
		return nil, false
	}
	if os.broken {
		return nil, false
	}
	idx := e.index
	if idx >= 0 && idx < len(os.nums) && os.nums[idx] == key.Num {
		if os.objects[idx] == nil {
			return nil, false
		}
		return os.objects[idx], true
	}
	// Trust the xref and search the container's number list (pdfbox
	// PDFObjectStreamParser behaviour), warning about the mismatch.
	d.addWarning(WarnObjStmBroken, -1, "object "+key.String()+
		" is not at index "+strconv.Itoa(idx)+" of its container; searching by number")
	for i, n := range os.nums {
		if n == key.Num {
			if os.objects[i] == nil {
				return nil, false
			}
			return os.objects[i], true
		}
	}
	return nil, false
}

var errObjStmCycle = errors.New("pdf: object stream cycle")
