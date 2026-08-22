package asic

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TestASiCUtilsConstants pins every filename, folder and extension constant verbatim against
// ASiCUtils.java: these strings end up as ZIP entry names and manifest URIs, so a typo is a
// silently broken container.
func TestASiCUtilsConstants(t *testing.T) {
	for _, tc := range []struct{ got, want, name string }{
		{ASiCUtilsManifestFilename, "Manifest", "MANIFEST_FILENAME"},
		{ASiCUtilsASiCManifestFilename, "ASiCManifest", "ASIC_MANIFEST_FILENAME"},
		{ASiCUtilsASiCArchiveManifestFilename, "ASiCArchiveManifest", "ASIC_ARCHIVE_MANIFEST_FILENAME"},
		{ASiCUtilsASiCEvidenceRecordManifestFilename, "ASiCEvidenceRecordManifest", "ASIC_EVIDENCE_RECORD_MANIFEST_FILENAME"},
		{ASiCUtilsASiCXAdESManifestFilename, "manifest", "ASIC_XAdES_MANIFEST_FILENAME"},
		{ASiCUtilsMimeType, "mimetype", "MIME_TYPE"},
		{ASiCUtilsMimeTypeComment, "mimetype=", "MIME_TYPE_COMMENT"},
		{ASiCUtilsMetaInfFolder, "META-INF/", "META_INF_FOLDER"},
		{ASiCUtilsPackageZip, "package.zip", "PACKAGE_ZIP"},
		{ASiCUtilsSignatureFilename, "signature", "SIGNATURE_FILENAME"},
		{ASiCUtilsSignaturesFilename, "signatures", "SIGNATURES_FILENAME"},
		{ASiCUtilsTimestampFilename, "timestamp", "TIMESTAMP_FILENAME"},
		{ASiCUtilsEvidenceRecordFilename, "evidencerecord", "EVIDENCE_RECORD_FILENAME"},
		{ASiCUtilsCAdESSignatureExtension, ".p7s", "CADES_SIGNATURE_EXTENSION"},
		{ASiCUtilsTSTExtension, ".tst", "TST_EXTENSION"},
		{ASiCUtilsERASN1Extension, ".ers", "ER_ASN1_EXTENSION"},
		{ASiCUtilsXMLExtension, ".xml", "XML_EXTENSION"},
		{ASiCUtilsSignaturesXML, "META-INF/signatures.xml", "SIGNATURES_XML"},
		{ASiCUtilsOpenDocumentSignatures, "META-INF/documentsignatures.xml", "OPEN_DOCUMENT_SIGNATURES"},
		{ASiCUtilsASiCEMetaInfManifest, "META-INF/manifest.xml", "ASICE_METAINF_MANIFEST"},
		{ASiCUtilsASiCEMetaInfXAdESSignature, "META-INF/signatures001.xml", "ASICE_METAINF_XADES_SIGNATURE"},
		{ASiCUtilsASiCEMetaInfCAdESSignature, "META-INF/signature001.p7s", "ASICE_METAINF_CADES_SIGNATURE"},
		{ASiCUtilsASiCEMetaInfCAdESTimestamp, "META-INF/timestamp001.tst", "ASICE_METAINF_CADES_TIMESTAMP"},
		{ASiCUtilsASiCEMetaInfCAdESEvidenceRecordASN1, "META-INF/evidencerecord001.ers", "ASICE_METAINF_CADES_EVIDENCE_RECORD_ASN1"},
		{ASiCUtilsASiCEMetaInfCAdESEvidenceRecordXML, "META-INF/evidencerecord001.xml", "ASICE_METAINF_CADES_EVIDENCE_RECORD_XML"},
		{ASiCUtilsASiCEMetaInfCAdESManifest, "META-INF/ASiCManifest001.xml", "ASICE_METAINF_CADES_MANIFEST"},
		{ASiCUtilsASiCEMetaInfCAdESArchiveManifest, "META-INF/ASiCArchiveManifest001.xml", "ASICE_METAINF_CADES_ARCHIVE_MANIFEST"},
		{ASiCUtilsASiCEMetaInfEvidenceRecordManifest, "META-INF/ASiCEvidenceRecordManifest001.xml", "ASICE_METAINF_EVIDENCE_RECORD_MANIFEST"},
		{ASiCUtilsSignatureP7S, "META-INF/signature.p7s", "SIGNATURE_P7S"},
		{ASiCUtilsTimestampTST, "META-INF/timestamp.tst", "TIMESTAMP_TST"},
		{ASiCUtilsEvidenceRecordERS, "META-INF/evidencerecord.ers", "EVIDENCE_RECORD_ERS"},
		{ASiCUtilsEvidenceRecordXML, "META-INF/evidencerecord.xml", "EVIDENCE_RECORD_XML"},
		{SecureContainerHandlerMimetype, "mimetype", "SecureContainerHandler.MIMETYPE"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
	if asicUtilsMaxToRead != 0xFFFF+2+4 {
		t.Errorf("MAX_TO_READ = %d, want %d", asicUtilsMaxToRead, 0xFFFF+2+4)
	}
}

// TestASiCUtilsEntryNamePredicates exercises the entry-name classifiers, including the traps in
// upstream's definitions: isSignature also matches "signatures*" and rejects anything containing
// "Manifest", and isTimestamp/isEvidenceRecord key off both the folder and the extension.
func TestASiCUtilsEntryNamePredicates(t *testing.T) {
	for _, tc := range []struct {
		name                                                              string
		signature, timestamp, evidenceRecord, xmlER, asn1ER, xades, cades bool
		manifest, archiveManifest, evidenceRecordManifest, mimetype       bool
	}{
		{name: "META-INF/signature001.p7s", signature: true, cades: true},
		{name: "META-INF/signatures001.xml", signature: true, xades: true},
		{name: "META-INF/documentsignatures.xml", signature: true, xades: true},
		{name: "META-INF/ASiCManifest001.xml", manifest: true},
		// "ASiCArchiveManifest..." does NOT contain the substring "ASiCManifest", so
		// isManifest is false for it - the two predicates are disjoint upstream.
		{name: "META-INF/ASiCArchiveManifest001.xml", archiveManifest: true},
		{name: "META-INF/ASiCEvidenceRecordManifest001.xml", evidenceRecordManifest: true},
		{name: "META-INF/manifest.xml"},
		{name: "META-INF/timestamp001.tst", timestamp: true},
		{name: "META-INF/timestamp.tst", timestamp: true},
		{name: "META-INF/evidencerecord.ers", evidenceRecord: true, asn1ER: true},
		{name: "META-INF/evidencerecord001.xml", evidenceRecord: true, xmlER: true},
		{name: "signature001.p7s"}, // outside META-INF/
		{name: "mimetype", mimetype: true},
		{name: "test.txt"},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			checks := []struct {
				got, want bool
				label     string
			}{
				{UtilsIsSignature(tc.name), tc.signature, "isSignature"},
				{UtilsIsTimestamp(tc.name), tc.timestamp, "isTimestamp"},
				{UtilsIsEvidenceRecord(tc.name), tc.evidenceRecord, "isEvidenceRecord"},
				{UtilsIsXmlEvidenceRecord(tc.name), tc.xmlER, "isXmlEvidenceRecord"},
				{UtilsIsAsn1EvidenceRecord(tc.name), tc.asn1ER, "isAsn1EvidenceRecord"},
				{UtilsIsXAdES(tc.name), tc.xades, "isXAdES"},
				{UtilsIsCAdES(tc.name), tc.cades, "isCAdES"},
				{UtilsIsManifest(tc.name), tc.manifest, "isManifest"},
				{UtilsIsArchiveManifest(tc.name), tc.archiveManifest, "isArchiveManifest"},
				{UtilsIsEvidenceRecordManifest(tc.name), tc.evidenceRecordManifest, "isEvidenceRecordManifest"},
				{UtilsIsMimetype(tc.name), tc.mimetype, "isMimetype"},
			}
			for _, check := range checks {
				if check.got != check.want {
					t.Errorf("%s(%q) = %t, want %t", check.label, tc.name, check.got, check.want)
				}
			}
		})
	}
}

// TestASiCUtilsZipCommentBuilders pins the "mimetype=" prefix every ASiC zip comment carries.
func TestASiCUtilsZipCommentBuilders(t *testing.T) {
	if got := UtilsZipCommentFromMimeTypeString("application/vnd.etsi.asic-e+zip"); got != "mimetype=application/vnd.etsi.asic-e+zip" {
		t.Errorf("zip comment = %q", got)
	}
	if got := UtilsZipCommentFromMimeType(enumerations.MimeTypeEnumASiCS); got != "mimetype=application/vnd.etsi.asic-s+zip" {
		t.Errorf("zip comment from MimeType = %q", got)
	}

	parameters := NewASiCParameters()
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	if got := UtilsZipCommentFromParameters(parameters); got != "" {
		t.Errorf("zip comment = %q with zipComment off, want \"\"", got)
	}
	parameters.SetZipComment(true)
	if got := UtilsZipCommentFromParameters(parameters); got != "mimetype=application/vnd.etsi.asic-e+zip" {
		t.Errorf("zip comment = %q", got)
	}
	if got := UtilsMimeTypeString(parameters); got != "application/vnd.etsi.asic-e+zip" {
		t.Errorf("mimetype string = %q", got)
	}
}

// TestASiCUtilsASiCContainerType pins the MimeType -> ASiCContainerType mapping and its
// IllegalArgumentException path.
func TestASiCUtilsASiCContainerType(t *testing.T) {
	for _, tc := range []struct {
		mimeType enumerations.MimeType
		want     enumerations.ASiCContainerType
	}{
		{enumerations.MimeTypeEnumASiCS, enumerations.ASiCContainerTypeASiCS},
		{enumerations.MimeTypeEnumASiCE, enumerations.ASiCContainerTypeASiCE},
		{enumerations.MimeTypeEnumODT, enumerations.ASiCContainerTypeASiCE},
		{enumerations.MimeTypeEnumODS, enumerations.ASiCContainerTypeASiCE},
		{enumerations.MimeTypeEnumODG, enumerations.ASiCContainerTypeASiCE},
		{enumerations.MimeTypeEnumODP, enumerations.ASiCContainerTypeASiCE},
	} {
		got, err := UtilsASiCContainerType(tc.mimeType)
		if err != nil || got != tc.want {
			t.Errorf("containerType(%s) = %s (err %v), want %s", tc.mimeType.MimeTypeString(), got, err, tc.want)
		}
	}
	if _, err := UtilsASiCContainerType(enumerations.MimeTypeEnumPDF); err == nil {
		t.Error("expected an error for a non-ASiC mimetype")
	} else if err.Error() != "Not allowed mimetype 'application/pdf'" {
		t.Errorf("error = %q, want the Java message", err.Error())
	}
}

// TestASiCUtilsContainerInspectionMatchesDSS runs isZip / isASiC / isContainerOpenDocument /
// getContainerType / getZipComment over the whole fixture corpus against real DSS. getContainerType
// in particular walks the whole detection ladder (mimetype document -> zip comment -> root-level
// document count -> enforced container mimetype).
func TestASiCUtilsContainerInspectionMatchesDSS(t *testing.T) {
	oracle := loadDSSOracle(t)
	for _, container := range oracle.Containers {
		container := container
		t.Run(container.Path, func(t *testing.T) {
			doc := zipCoreFileDocument(t, container.Path)

			isZip, err := UtilsIsZip(doc)
			if err != nil {
				t.Fatalf("isZip: %v", err)
			}
			if isZip != container.IsZip {
				t.Errorf("isZip = %t, DSS = %t", isZip, container.IsZip)
			}

			isASiC, err := UtilsIsASiC(doc)
			if err != nil {
				t.Fatalf("isASiC: %v", err)
			}
			if isASiC != container.IsASiC {
				t.Errorf("isASiC = %t, DSS = %t", isASiC, container.IsASiC)
			}

			isOpenDocument, err := UtilsIsContainerOpenDocument(doc)
			if err != nil {
				t.Fatalf("isContainerOpenDocument: %v", err)
			}
			if isOpenDocument != container.IsContainerOpenDocument {
				t.Errorf("isContainerOpenDocument = %t, DSS = %t", isOpenDocument, container.IsContainerOpenDocument)
			}

			containerType, err := UtilsContainerType(doc)
			if err != nil {
				t.Fatalf("containerType: %v", err)
			}
			wantContainerType := ""
			if container.ContainerType != nil {
				wantContainerType = *container.ContainerType
			}
			if string(containerType) != wantContainerType {
				t.Errorf("containerType = %q, DSS = %q", containerType, wantContainerType)
			}

			zipComment, err := UtilsZipCommentFromArchiveContainer(doc)
			if err != nil {
				t.Fatalf("zipComment: %v", err)
			}
			wantZipComment := ""
			if container.ZipComment != nil {
				wantZipComment = *container.ZipComment
			}
			if zipComment != wantZipComment {
				t.Errorf("zipComment = %q, DSS = %q", zipComment, wantZipComment)
			}
		})
	}
}

// TestASiCUtilsIsZipRejectsDigestDocument pins the DigestDocument short-circuit: a document with no
// readable content can never be a container, and must not be opened to find that out.
func TestASiCUtilsIsZipRejectsDigestDocument(t *testing.T) {
	digestDocument := model.NewDigestDocumentFromBase64(enumerations.DigestAlgorithmSHA256,
		"GTZDVjc8fdBQlkeXsQnPHYFXKVi7B6Nkzo9YDLBODcs=")
	isZip, err := UtilsIsZip(digestDocument)
	if err != nil {
		t.Fatalf("isZip: %v", err)
	}
	if isZip {
		t.Error("isZip(DigestDocument) = true, want false")
	}
	if isZip, err := UtilsIsZip(nil); err != nil || isZip {
		t.Errorf("isZip(nil) = %t (err %v), want false", isZip, err)
	}
}

// TestASiCUtilsAddOrReplaceDocument pins the replace-by-name-else-append contract the extension and
// merge flows rely on.
func TestASiCUtilsAddOrReplaceDocument(t *testing.T) {
	first := model.NewInMemoryDocumentWithName([]byte("1"), "a.txt")
	second := model.NewInMemoryDocumentWithName([]byte("2"), "b.txt")
	replacement := model.NewInMemoryDocumentWithName([]byte("3"), "a.txt")

	documents := []model.DSSDocument{first, second}
	documents = UtilsAddOrReplaceDocument(documents, replacement)
	if len(documents) != 2 || documents[0] != model.DSSDocument(replacement) || documents[1] != model.DSSDocument(second) {
		t.Fatalf("replace-in-place failed: %v", documents)
	}

	third := model.NewInMemoryDocumentWithName([]byte("4"), "c.txt")
	documents = UtilsAddOrReplaceDocument(documents, third)
	if len(documents) != 3 || documents[2] != model.DSSDocument(third) {
		t.Fatalf("append failed: %v", documents)
	}
}

// TestASiCUtilsRootLevelDocuments pins the root-level filters, including their exclusion of the
// mimetype entry and of anything inside a folder.
func TestASiCUtilsRootLevelDocuments(t *testing.T) {
	documents := []model.DSSDocument{
		model.NewInMemoryDocumentWithName([]byte("a"), "a.txt"),
		model.NewInMemoryDocumentWithName([]byte("b"), "folder/b.txt"),
		model.NewInMemoryDocumentWithName([]byte("m"), "mimetype"),
		model.NewInMemoryDocumentWithName([]byte("c"), "c.txt"),
	}
	rootLevel := UtilsRootLevelDocuments(documents)
	if len(rootLevel) != 2 || rootLevel[0].Name() != "a.txt" || rootLevel[1].Name() != "c.txt" {
		t.Fatalf("rootLevelDocuments = %v", modelNames(rootLevel))
	}

	asicContent := NewASiCContent()
	asicContent.SetSignedDocuments(documents)
	if got := UtilsRootLevelSignedDocuments(asicContent); len(got) != 2 {
		t.Fatalf("rootLevelSignedDocuments = %v", modelNames(got))
	}

	// A single signed document is returned as-is, even when it lives in a folder.
	single := []model.DSSDocument{model.NewInMemoryDocumentWithName([]byte("b"), "folder/b.txt")}
	asicContent.SetSignedDocuments(single)
	if got := UtilsRootLevelSignedDocuments(asicContent); len(got) != 1 || got[0].Name() != "folder/b.txt" {
		t.Fatalf("a single signed document must be returned unfiltered, got %v", modelNames(got))
	}
}

func modelNames(documents []model.DSSDocument) []string {
	names := make([]string, 0, len(documents))
	for _, doc := range documents {
		names = append(names, doc.Name())
	}
	return names
}

// TestASiCUtilsEnsureMimeTypeAndZipComment pins the mimetype document ensureMimeTypeAndZipComment
// synthesizes: named "mimetype", STORED, carrying the mimetype string as its content.
func TestASiCUtilsEnsureMimeTypeAndZipComment(t *testing.T) {
	parameters := NewASiCParameters()
	parameters.SetContainerType(enumerations.ASiCContainerTypeASiCE)
	parameters.SetZipComment(true)

	asicContent := NewASiCContent()
	if _, err := UtilsEnsureMimeTypeAndZipComment(asicContent, parameters); err != nil {
		t.Fatalf("ensureMimeTypeAndZipComment: %v", err)
	}
	mimetypeDocument := asicContent.MimeTypeDocument()
	if mimetypeDocument == nil {
		t.Fatal("no mimetype document was created")
	}
	if mimetypeDocument.Name() != ASiCUtilsMimeType {
		t.Errorf("mimetype document name = %q", mimetypeDocument.Name())
	}
	zipEntryDocument, ok := mimetypeDocument.(DSSZipEntryDocument)
	if !ok {
		t.Fatalf("mimetype document is %T, want a DSSZipEntryDocument", mimetypeDocument)
	}
	if zipEntryDocument.ZipEntry().CompressionMethod() != 0 {
		t.Errorf("mimetype compression method = %d, want STORED (0)", zipEntryDocument.ZipEntry().CompressionMethod())
	}
	if asicContent.ZipComment() != "mimetype=application/vnd.etsi.asic-e+zip" {
		t.Errorf("zip comment = %q", asicContent.ZipComment())
	}

	// A second call is a no-op (both fields are already populated).
	before := asicContent.MimeTypeDocument()
	if _, err := UtilsEnsureMimeTypeAndZipComment(asicContent, parameters); err != nil {
		t.Fatalf("ensureMimeTypeAndZipComment (second call): %v", err)
	}
	if asicContent.MimeTypeDocument() != before {
		t.Error("the mimetype document was replaced on the second call")
	}
}

// TestASiCUtilsToSimpleManifestEntries pins the "simple" manifest entry projection.
func TestASiCUtilsToSimpleManifestEntries(t *testing.T) {
	documents := []model.DSSDocument{
		model.NewInMemoryDocumentWithName([]byte("a"), "a.xml"),
		model.NewInMemoryDocumentWithName([]byte("b"), "b.txt"),
	}
	entries := UtilsToSimpleManifestEntries(documents)
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}
	for i, entry := range entries {
		if entry.Uri() != documents[i].Name() {
			t.Errorf("entry[%d] uri = %q, want %q", i, entry.Uri(), documents[i].Name())
		}
		if entry.MimeType() != documents[i].MimeType() {
			t.Errorf("entry[%d] mimeType mismatch", i)
		}
		if !entry.IsFound() {
			t.Errorf("entry[%d] found = false, want true", i)
		}
		if entry.Document() != documents[i] {
			t.Errorf("entry[%d] document mismatch", i)
		}
	}
}

// TestASiCUtilsFilenameCollectionPredicates pins the list-level classifiers used by the format
// detectors.
func TestASiCUtilsFilenameCollectionPredicates(t *testing.T) {
	cades := []string{"mimetype", "test.txt", "META-INF/signature001.p7s"}
	xades := []string{"mimetype", "test.txt", "META-INF/signatures001.xml"}
	timestamped := []string{"mimetype", "test.txt", "META-INF/timestamp.tst"}
	evidenceRecord := []string{"mimetype", "test.txt", "META-INF/evidencerecord.ers"}
	plain := []string{"a.txt", "b.txt"}

	if !UtilsFilesContainMetaInfFolder(cades) || UtilsFilesContainMetaInfFolder(plain) {
		t.Error("filesContainMetaInfFolder")
	}
	if !UtilsFilesContainSignatures(cades) || UtilsFilesContainSignatures(plain) {
		t.Error("filesContainSignatures")
	}
	if !UtilsFilesContainTimestamps(timestamped) || UtilsFilesContainTimestamps(cades) {
		t.Error("filesContainTimestamps")
	}
	if !UtilsFilesContainEvidenceRecords(evidenceRecord) || UtilsFilesContainEvidenceRecords(cades) {
		t.Error("filesContainEvidenceRecords")
	}
	if !UtilsAreFilesContainMimetype(cades) || UtilsAreFilesContainMimetype(plain) {
		t.Error("areFilesContainMimetype")
	}
	if !UtilsIsASiCWithCAdES(cades) || UtilsIsASiCWithCAdES(xades) {
		t.Error("isASiCWithCAdES")
	}
	if !UtilsIsASiCWithXAdES(xades) || UtilsIsASiCWithXAdES(cades) {
		t.Error("isASiCWithXAdES")
	}
	// Evidence records are a shared format: neither format-specific predicate claims them.
	if UtilsIsASiCWithCAdES(evidenceRecord) || UtilsIsASiCWithXAdES(evidenceRecord) {
		t.Error("an evidence record must not be claimed by either format predicate")
	}
	for _, filenames := range [][]string{cades, xades, timestamped, evidenceRecord} {
		if !UtilsIsAsicFileContent(filenames) {
			t.Errorf("isAsicFileContent(%v) = false", filenames)
		}
	}
	if UtilsIsAsicFileContent(plain) {
		t.Error("isAsicFileContent(plain) = true")
	}
}
