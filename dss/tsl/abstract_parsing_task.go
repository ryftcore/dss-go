// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/AbstractParsingTask.java (DSS 6.5.RC1).
//
// It implements eu.europa.esig.dss.validation.job.parsing.ParsingTask (Supplier<ParsingResult>)
// through its concrete subclasses; Java's ParsingTask declares no member of its own, and Go
// satisfies interfaces structurally, so dss/validation/job is not imported here.
package tsl

import (
	"fmt"
	"io"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/trustedlist"
	"github.com/utain/esig/dss/trustedlist/jaxb"
	"github.com/utain/esig/dss/utils"
	xadestsl "github.com/utain/esig/dss/xades/tsl"
)

// trustedListFacade is the Go form of the eu.europa.esig.trustedlist.TrustedListFacade return
// type of createTrustedListFacade(). Upstream, MRAFacade extends TrustedListFacade so that
// LOTLParsingTask can substitute one for the other; the Go port of specs-trusted-list makes
// trustedlist.TrustedListFacade and trustedlist.MRAFacade two independent structs (Go has no
// implementation inheritance), so the substitutability is expressed by this interface, which both
// satisfy.
//
// Unmarshal takes the document bytes rather than an InputStream because that is the shape the
// ported facades expose; Java's second `validate` argument is always false here ("lax processing,
// validate XSD after") and the Go facades perform no runtime XSD validation at all, so it has no
// counterpart.
type trustedListFacade interface {
	Unmarshal(data []byte) (*jaxb.TrustStatusListType, error)
}

// AbstractParsingTaskOverrides captures the member AbstractParsingTaskBase reaches through
// virtual dispatch: the facade used to unmarshal the document (LOTLParsingTask swaps in the MRA
// one).
type AbstractParsingTaskOverrides interface {
	// CreateTrustedListFacade loads a TrustedListFacade. Port of the protected
	// createTrustedListFacade().
	CreateTrustedListFacade() trustedListFacade
}

// abstractParsingTaskStructureValidationTarget is the Go form of the AbstractParsingResult
// parameter of verifyTLVersionConformity: the only member the method touches is the structure
// validation message setter, so the parameter carries just that, keeping this file independent of
// how the dss-validation-job chunk spells the abstract result type.
type abstractParsingTaskStructureValidationTarget interface {
	SetStructureValidationMessages(structureValidationMessages []string)
}

// AbstractParsingTaskBase carries the concrete behaviour of the Java abstract class
// AbstractParsingTask. Concrete tasks (TLParsingTask, LOTLParsingTask) embed it and register
// themselves with InitAbstractParsingTask.
type AbstractParsingTaskBase struct {
	// overrides points back at the concrete task; see InitAbstractParsingTask.
	overrides AbstractParsingTaskOverrides

	// document is the document to parse.
	document model.DSSDocument
}

// NewAbstractParsingTaskBase instantiates the base state of a parsing task. Port of the protected
// AbstractParsingTask(DSSDocument) constructor; the concrete task must still call
// InitAbstractParsingTask.
//
// Panics with the Java message when document is nil (Objects.requireNonNull).
func NewAbstractParsingTaskBase(document model.DSSDocument) AbstractParsingTaskBase {
	if document == nil {
		panic("The document is null")
	}
	return AbstractParsingTaskBase{document: document}
}

// InitAbstractParsingTask registers the concrete task with its base so that the base can dispatch
// CreateTrustedListFacade the way Java reaches an overridden method through virtual dispatch. It
// must be called exactly once, by the concrete task's constructor, before any other method.
func (t *AbstractParsingTaskBase) InitAbstractParsingTask(overrides AbstractParsingTaskOverrides) {
	t.overrides = overrides
}

// abstractParsingTaskOverrides returns the registered overrides, panicking when the concrete task
// forgot to call InitAbstractParsingTask.
func (t *AbstractParsingTaskBase) abstractParsingTaskOverrides() AbstractParsingTaskOverrides {
	if t.overrides == nil {
		panic("AbstractParsingTask was not initialised: the concrete task must call InitAbstractParsingTask in its constructor")
	}
	return t.overrides
}

// Document returns the document to parse. It has no Java counterpart (upstream keeps the field
// private and reads it from getJAXBObject/verifyTLVersionConformity, both of which live on this
// same class); it exists because the concrete tasks in this package are separate Go types rather
// than subclasses sharing the field.
func (t *AbstractParsingTaskBase) Document() model.DSSDocument {
	return t.document
}

// JAXBObject gets the TrustStatusListType. Port of the protected getJAXBObject(); Java's
// DSSException becomes a returned error, per PORTING.md.
//
// NOTE: Java's `e.getMessage() == null && e.getCause() != null` branch, which substitutes the
// cause's message for a message-less exception, is unreachable in Go - an error always renders a
// message - so the two branches collapse into one.
func (t *AbstractParsingTaskBase) JAXBObject() (*jaxb.TrustStatusListType, error) {
	data, err := t.readDocument()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to parse binaries. Reason : '%s'", err.Error()), err)
	}
	// lax processing, validate XSD after
	jaxbObject, err := t.abstractParsingTaskOverrides().CreateTrustedListFacade().Unmarshal(data)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to parse binaries. Reason : '%s'", err.Error()), err)
	}
	return jaxbObject, nil
}

// readDocument reads the whole document, standing in for Java's try-with-resources over
// document.openStream().
func (t *AbstractParsingTaskBase) readDocument() ([]byte, error) {
	stream, err := t.document.OpenStream()
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	return io.ReadAll(stream)
}

// CreateTrustedListFacade loads a TrustedListFacade. Port of the protected
// createTrustedListFacade(), whose body is `return TrustedListFacade.newFacade();`.
func (t *AbstractParsingTaskBase) CreateTrustedListFacade() trustedListFacade {
	return trustedlist.NewTrustedListFacade()
}

// CommonParseSchemeInformation extracts the common values. Port of the protected
// commonParseSchemeInformation(AbstractTLParsingResult, TSLSchemeInformationType).
func (t *AbstractParsingTaskBase) CommonParseSchemeInformation(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	if schemeInformation != nil {
		t.extractTSLType(result, schemeInformation)
		t.extractSequenceNumber(result, schemeInformation)
		t.extractTerritory(result, schemeInformation)
		t.extractVersion(result, schemeInformation)
		t.extractIssueDate(result, schemeInformation)
		t.extractNextUpdateDate(result, schemeInformation)
		t.extractDistributionPoints(result, schemeInformation)
	}
}

func (t *AbstractParsingTaskBase) extractTSLType(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	tslType := schemeInformation.TSLType
	if utils.IsStringNotEmpty(tslType) {
		result.SetTSLType(enumerations.TSLTypeFromURI(tslType))
	}
}

func (t *AbstractParsingTaskBase) extractSequenceNumber(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	tslSequenceNumber := schemeInformation.TSLSequenceNumber
	if tslSequenceNumber != nil {
		// BigInteger#intValue() keeps the low-order 32 bits; every real TL sequence number
		// fits an int, so the plain conversion below reproduces it.
		value := int(tslSequenceNumber.Int64())
		result.SetSequenceNumber(&value)
	}
}

// extractTerritory ports the private extractTerritory, i.e.
// `result.setTerritory(schemeInformation.getSchemeTerritory())`.
//
// The model binds the optional tl:SchemeTerritory element to *string; an absent element becomes
// the empty string here, which is this codebase's uniform stand-in for a null Java String (the
// same value TrustServiceProviderConverter would then stamp onto every TrustServiceProvider).
func (t *AbstractParsingTaskBase) extractTerritory(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	territory := ""
	if schemeInformation.SchemeTerritory != nil {
		territory = *schemeInformation.SchemeTerritory
	}
	result.SetTerritory(territory)
}

func (t *AbstractParsingTaskBase) extractVersion(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	tslVersionIdentifier := schemeInformation.TSLVersionIdentifier
	if tslVersionIdentifier != nil {
		value := int(tslVersionIdentifier.Int64())
		result.SetVersion(&value)
	}
}

func (t *AbstractParsingTaskBase) extractIssueDate(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	result.SetIssueDate(abstractParsingTaskConvertToDate(schemeInformation.ListIssueDateTime))
}

func (t *AbstractParsingTaskBase) extractNextUpdateDate(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	nextUpdate := schemeInformation.NextUpdate
	if nextUpdate != nil {
		result.SetNextUpdateDate(abstractParsingTaskConvertToDate(nextUpdate.DateTime))
	}
}

// abstractParsingTaskConvertToDate ports the private convertToDate(XMLGregorianCalendar), i.e.
// `gregorianCalendar.toGregorianCalendar().getTime()`.
//
// The Go trusted-list model binds every xs:dateTime property to its raw lexical form (see
// dss/trustedlist/jaxb's header), so this is where the instant is recovered. A lexical form
// carrying no timezone is resolved in the local zone, exactly as
// XMLGregorianCalendar#toGregorianCalendar() does for an undefined timezone field
// (GregorianCalendar's default TimeZone). An unparseable form answers the zero time.Time, standing
// in for the null Java returns when the calendar itself is null.
func abstractParsingTaskConvertToDate(lexical *string) time.Time {
	if lexical == nil {
		return time.Time{}
	}
	value := *lexical
	if parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", value); err == nil {
		return parsed
	}
	if parsed, err := time.ParseInLocation("2006-01-02T15:04:05", value, time.Local); err == nil {
		return parsed
	}
	return time.Time{}
}

func (t *AbstractParsingTaskBase) extractDistributionPoints(result *AbstractTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	distributionPoints := schemeInformation.DistributionPoints
	if distributionPoints != nil && utils.IsCollectionNotEmpty(distributionPoints.URI) {
		result.SetDistributionPoints(distributionPoints.URI)
	} else {
		result.SetDistributionPoints([]string{})
	}
}

// VerifyTLVersionConformity verifies the structure conformity of the Trusted List against its
// schema version. Port of the protected verifyTLVersionConformity(AbstractParsingResult, Integer,
// List<Integer>); the DSSException TLStructureVerifier#validate raises when the document cannot be
// read as XML becomes a returned error, per PORTING.md.
func (t *AbstractParsingTaskBase) VerifyTLVersionConformity(result abstractParsingTaskStructureValidationTarget,
	tlVersion *int, tlVersions []int) error {
	if utils.IsCollectionNotEmpty(tlVersions) {
		structureValidationMessagesResult, err := xadestsl.NewTLStructureVerifier().
			SetAcceptedTLVersions(tlVersions).Validate(t.document, tlVersion)
		if err != nil {
			return err
		}
		if utils.IsCollectionNotEmpty(structureValidationMessagesResult) {
			result.SetStructureValidationMessages(structureValidationMessagesResult)
		}
	}
	return nil
}
