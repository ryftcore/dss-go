package eaa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
)

// TestInitChain_NoFormatCheckingResult covers T32-SEC-001. In a build without the `eaa` build
// tag blocks.BasicBuildingBlocks has no EAA format-checking block, so the basic building blocks
// of an EAA token carry no FC while a later step (here the cryptographic verification) does. As
// upstream's `item = firstItem` (null) would, the process cannot chain that step onto nothing:
// it must fail closed with an explicit, descriptive panic (which the facade turns into an
// error) rather than a bare nil-pointer dereference - and never build a result that skipped the
// format-checking gate.
func TestInitChain_NoFormatCheckingResult(t *testing.T) {
	id := diagnosticjaxb.CollapsedString("EAA-1")
	eaaWrapper := diagnostic.NewEAAWrapper(&diagnosticjaxb.XmlEAA{XmlAbstractTokenAttrs: diagnosticjaxb.XmlAbstractTokenAttrs{Id: &id}})
	bbbs := map[string]*jaxb.XmlBasicBuildingBlocks{
		"EAA-1": {CV: &jaxb.XmlCV{}},
	}
	process := NewValidationProcess(i18n.NewProvider(), eaaWrapper, nil, bbbs, nil)

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected InitChain to panic when the EAA has no Format Checking result")
		}
		message := fmt.Sprint(recovered)
		if _, isRuntimeError := recovered.(error); isRuntimeError {
			t.Fatalf("got a bare runtime error instead of an explicit message: %v", recovered)
		}
		for _, want := range []string{"EAA-1", "Format Checking", "'eaa' build tag"} {
			if !strings.Contains(message, want) {
				t.Errorf("panic message %q does not mention %q", message, want)
			}
		}
	}()
	process.InitChain()
}

// TestInitChain_NothingToChainWithoutFormatChecking pins the corner upstream also tolerates: with
// no FC block AND nothing else to verify, there is no step that dereferences the null head, so
// the chain is simply empty.
func TestInitChain_NothingToChainWithoutFormatChecking(t *testing.T) {
	id := diagnosticjaxb.CollapsedString("EAA-2")
	eaaWrapper := diagnostic.NewEAAWrapper(&diagnosticjaxb.XmlEAA{XmlAbstractTokenAttrs: diagnosticjaxb.XmlAbstractTokenAttrs{Id: &id}})
	bbbs := map[string]*jaxb.XmlBasicBuildingBlocks{"EAA-2": {}}
	process := NewValidationProcess(i18n.NewProvider(), eaaWrapper, nil, bbbs, nil)

	process.InitChain() // must not panic
}
