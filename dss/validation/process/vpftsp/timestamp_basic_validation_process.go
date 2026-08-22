// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftsp/TimestampBasicValidationProcess.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.vpftsp.checks (BasicTimestampValidationCheck,
// BasicTimestampValidationWithIdCheck) is ported into package vpfbs, not this
// package, to break a Go import cycle - see vpfbs/basic_timestamp_validation_check.go.
package vpftsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfbs"
)

// TimestampBasicValidationProcess performs Time-stamp validation building
// block as per clause 5.4.
type TimestampBasicValidationProcess struct {
	*vpfbs.AbstractBasicValidationProcess[*jaxb.XmlValidationProcessBasicTimestamp]

	// timestamp is the timestamp being validated.
	timestamp *diagnostic.TimestampWrapper
}

// NewTimestampBasicValidationProcess is the default constructor. Port of
// TimestampBasicValidationProcess(I18nProvider, DiagnosticData, TimestampWrapper, Map).
func NewTimestampBasicValidationProcess(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	timestamp *diagnostic.TimestampWrapper, bbbs map[string]*jaxb.XmlBasicBuildingBlocks) *TimestampBasicValidationProcess {
	xmlResult := &jaxb.XmlValidationProcessBasicTimestamp{}
	result := process.NewResult(xmlResult, &xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)
	c := &TimestampBasicValidationProcess{
		AbstractBasicValidationProcess: vpfbs.NewAbstractBasicValidationProcess(i18nProvider, result, diagnosticData, timestamp, bbbs),
		timestamp:                      timestamp,
	}
	c.InitAbstractBasicValidationProcess(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *TimestampBasicValidationProcess) Title() i18n.MessageTag {
	return i18n.MessageTag_VPFTSP
}

// AddAdditionalInfo adds additional info to the chain. Port of
// addAdditionalInfo().
func (c *TimestampBasicValidationProcess) AddAdditionalInfo() {
	c.Result.Value.Type = string(c.timestamp.Type())
	if productionTime := c.timestamp.ProductionTime(); productionTime != nil {
		c.Result.Value.ProductionTime = jaxb.XSDateTime(*productionTime)
	}
}
