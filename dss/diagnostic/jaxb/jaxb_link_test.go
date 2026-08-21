package jaxb

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// collectTokens indexes every object carrying an xs:ID by that ID, walking the
// owned tree independently of Link.
func collectTokens(v reflect.Value, seen map[uintptr]bool, out map[string]reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		if v.Elem().Kind() == reflect.Struct {
			if seen[v.Pointer()] {
				return
			}
			seen[v.Pointer()] = true
			if id, ok := tokenIDOf(v); ok && id != "" {
				if _, dup := out[id]; !dup {
					out[id] = v
				}
			}
		}
		collectTokens(v.Elem(), seen, out)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			collectTokens(v.Index(i), seen, out)
		}
	case reflect.Interface:
		if !v.IsNil() {
			collectTokens(v.Elem(), seen, out)
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}
			if isAttrField(f) {
				continue // an IDREF is a back-edge, not ownership
			}
			collectTokens(v.Field(i), seen, out)
		}
	}
}

func tokenIDOf(pv reflect.Value) (string, bool) {
	if tok, ok := pv.Interface().(XmlToken); ok {
		return tok.TokenID(), true
	}
	return "", false
}

// TestLinkResolvesIDREFs checks the object graph Unmarshal builds: every IDREF
// attribute whose target is present in the document must point at that very
// object, so the wrappers can navigate the graph the way the Java model does.
func TestLinkResolvesIDREFs(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(oracleDir(t), "*.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatalf("no oracle dumps under %s", oracleDir(t))
	}
	total := 0
	for _, file := range files {
		if strings.HasPrefix(filepath.Base(file), "model-") {
			continue // the synthetic dumps reference tokens that are not in them
		}
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			dd, err := Unmarshal(data)
			if err != nil {
				t.Fatal(err)
			}
			root := reflect.ValueOf(dd)
			tokens := map[string]reflect.Value{}
			collectTokens(root, map[uintptr]bool{}, tokens)
			if len(tokens) == 0 {
				t.Fatal("document declares no xs:ID at all")
			}
			resolved, dangling := 0, 0
			walk(root, map[uintptr]bool{}, func(sv reflect.Value) {
				st := sv.Type()
				for i := 0; i < st.NumField(); i++ {
					f := st.Field(i)
					if !isAttrField(f) {
						continue
					}
					fv := sv.Field(i)
					if fv.Kind() != reflect.Pointer || fv.IsNil() {
						continue
					}
					var id string
					if f.Type == tokenRefType {
						ref := fv.Interface().(*XmlTokenRef)
						id = ref.ID
						target, ok := tokens[id]
						if !ok {
							dangling++
							continue
						}
						if ref.Token == nil {
							t.Errorf("%s.%s: reference %q left unresolved", st.Name(), f.Name, id)
							continue
						}
						if reflect.ValueOf(ref.Token).Pointer() != target.Pointer() {
							t.Errorf("%s.%s: reference %q resolved to another object", st.Name(), f.Name, id)
							continue
						}
						resolved++
						continue
					}
					if !f.Type.Implements(identifiableTyp) {
						continue
					}
					id = fv.Interface().(identifiable).TokenID()
					target, ok := tokens[id]
					if !ok {
						dangling++
						continue
					}
					if target.Pointer() != fv.Pointer() {
						t.Errorf("%s.%s: reference %q is a stub, not the referenced %s",
							st.Name(), f.Name, id, target.Type())
						continue
					}
					resolved++
				}
			})
			if resolved == 0 {
				t.Skipf("no resolvable IDREF in this document (%d dangling)", dangling)
			}
			total += resolved
		})
	}
	if total < 100 {
		t.Fatalf("only %d IDREF attributes resolved across the corpus", total)
	}
}
