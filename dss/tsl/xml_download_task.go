// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/download/XmlDownloadTask.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY: implements eu.europa.esig.dss.validation.job.download.DownloadTask
// (job.DownloadTask, a Supplier<DownloadResult> - ported here as a Get() DownloadResult, error
// method, replacing Java's unchecked-exception-throwing get() per PORTING.md's throw->error
// rule) - see xml_download_result.go's header for the wider job.* cross-chunk convention.
package tsl

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/client/http"
	"github.com/ryftcore/dss-go/dss/validation/job"
	"github.com/ryftcore/dss-go/dss/xades"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// xmlDownloadTaskDefaultDigestAlgorithm is the default digest algorithm used for document
// integrity identification.
const xmlDownloadTaskDefaultDigestAlgorithm = enumerations.DigestAlgorithm_SHA256

// xmlDownloadTaskDefaultCanonicalizationMethod is the default canonicalization method to be used
// on a document's digest computation.
const xmlDownloadTaskDefaultCanonicalizationMethod = xmlutils.XMLCanonicalizerDefaultDSSC14NMethod

// XmlDownloadTask downloads the document and returns an XmlDownloadResult.
type XmlDownloadTask struct {
	// dssFileLoader is the file loader.
	dssFileLoader http.DSSFileLoader

	// url is the URL to download the document from.
	url string
}

var _ job.DownloadTask = (*XmlDownloadTask)(nil)

// NewXmlDownloadTask is the default constructor.
//
// Panics with the Java messages when dssFileLoader is nil or url is empty
// (Objects.requireNonNull).
func NewXmlDownloadTask(dssFileLoader http.DSSFileLoader, url string) *XmlDownloadTask {
	if dssFileLoader == nil {
		panic("The DSSFileLoader is null")
	}
	if url == "" {
		panic("The url is null")
	}
	return &XmlDownloadTask{dssFileLoader: dssFileLoader, url: url}
}

// Get ports get(). Java's unchecked DSSException becomes a returned error.
func (t *XmlDownloadTask) Get() (job.DownloadResult, error) {
	dssDocument, err := t.dssFileLoader.GetDocument(t.url)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the content for url '%s'. Reason : '%s'", t.url, err.Error())
	}
	if err := t.assertDocumentIsValidXML(dssDocument); err != nil {
		return nil, err
	}

	stream, err := dssDocument.OpenStream()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the content for url '%s'. Reason : '%s'", t.url, err.Error())
	}
	defer stream.Close()

	messageDigest, err := xades.DSSXMLUtilsGetDigestOnCanonicalizedInputStream(stream, xmlDownloadTaskDefaultDigestAlgorithm, xmlDownloadTaskDefaultCanonicalizationMethod)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the content for url '%s'. Reason : '%s'", t.url, err.Error())
	}
	return NewXmlDownloadResult(dssDocument, messageDigest.Digest), nil
}

// assertDocumentIsValidXML ports the private assertDocumentIsValidXML(DSSDocument).
func (t *XmlDownloadTask) assertDocumentIsValidXML(document model.DSSDocument) error {
	if document == nil {
		return fmt.Errorf("no document has been retrieved from URL '%s'!", t.url)
	}
	if !xmlutils.DomUtilsIsDOM(document) {
		return fmt.Errorf("the document obtained from URL '%s' is not a valid XML!", t.url)
	}
	return nil
}
