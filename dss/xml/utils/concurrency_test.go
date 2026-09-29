package utils

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// TestConcurrentRegistrationAndLookup drives the process-global namespace registry and the
// shared executor loader the way parallel validations do: xades.NewSignature re-registers the
// XAdES namespaces on every construction while other goroutines compile XPath queries. Java's
// HashMap merely loses an update in that case; a Go map aborts the process ("fatal error:
// concurrent map writes"), so the registry must be synchronized. Meaningful under go test -race.
func TestConcurrentRegistrationAndLookup(t *testing.T) {
	doc, err := xmldom.Parse([]byte(`<a xmlns:p="urn:p"><p:b Id="x"/></a>`), nil)
	if err != nil {
		t.Fatal(err)
	}
	query := common.XPathQueryBuilderAllFromCurrentPosition().Build()

	const workers, rounds = 8, 200
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < rounds; i++ {
				XPathUtilsRegisterNamespace(common.NewDSSNamespace(fmt.Sprintf("urn:test:%d", w), fmt.Sprintf("t%d", w)))
				_ = XPathUtilsGetNamespaceContextMap().PrefixMap()
				_ = XPathUtilsGetNamespaceContextMap().NamespaceURI("t0")
				_, _ = XPathUtilsGetNamespaceContextMap().Prefix("urn:test:0")
				_ = XPathUtilsGetNamespaceContextMap().Prefixes("urn:test:0")
				if _, err := XPathUtilsGetNodeList(doc, query); err != nil {
					t.Errorf("XPathUtilsGetNodeList: %v", err)
					return
				}
			}
		}(w)
	}
	wg.Wait()
}

// TestConcurrentLoaderFirstUse is the lazy initialisation of a fresh loader raced from many
// goroutines.
func TestConcurrentLoaderFirstUse(t *testing.T) {
	loader := NewXPathQueryExecutorLoader()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if loader.GetXPathQueryExecutor() == nil || loader.GetXPathStringExecutor() == nil {
				t.Error("loader returned a nil executor")
			}
		}()
	}
	wg.Wait()
}
