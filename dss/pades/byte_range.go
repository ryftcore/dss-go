// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/ByteRange.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart). Java's int[] byteRangeArray becomes a Go
// []int; internal/pdf itself parses /ByteRange into []int64 and stops there (see internal/pdf/
// doc.go: "eu.europa.esig.dss.pades.validation.ByteRange owns validate() and getLength()"), so
// callers building a ByteRange from an internal/pdf.SignatureDictionary.ByteRange convert the
// four int64s down to int first.
package pades

import (
	"fmt"
	"math/big"

	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// ByteRange represents a ByteRange of a PDF Revision.
type ByteRange struct {
	// byteRangeArray represents a PDF signature byteRange.
	byteRangeArray []int

	// valid defines whether /ByteRange is valid (shall be defined by external process). nil
	// (Java's null Boolean) means validation has not been performed yet.
	valid *bool
}

// NewByteRange represents a ByteRange extracted from a Signature Dictionary of a
// signed/timestamped revision. Port of the constructor ByteRange(int[]).
func NewByteRange(byteRangeArray []int) *ByteRange {
	return &ByteRange{byteRangeArray: byteRangeArray}
}

// Length returns a total revision length. Port of getLength().
func (b *ByteRange) Length() int {
	// (before signature value) + (signature value) + (after signature value)
	return (b.byteRangeArray[1] - b.byteRangeArray[0]) + (b.byteRangeArray[2] - b.byteRangeArray[1]) + b.byteRangeArray[3]
}

// FirstPartStart returns the first byte number of the first part of the revision.
// Port of getFirstPartStart().
func (b *ByteRange) FirstPartStart() int {
	return b.byteRangeArray[0]
}

// FirstPartEnd returns the last byte number of the first part of the revision.
// Port of getFirstPartEnd().
func (b *ByteRange) FirstPartEnd() int {
	return b.byteRangeArray[1]
}

// SecondPartStart returns the first byte number of the second part of the revision.
// Port of getSecondPartStart().
func (b *ByteRange) SecondPartStart() int {
	return b.byteRangeArray[2]
}

// SecondPartEnd returns the last byte number of the second part of the revision.
// Port of getSecondPartEnd().
func (b *ByteRange) SecondPartEnd() int {
	return b.byteRangeArray[3]
}

// ToBigIntegerList transforms the ByteRange to a list of *big.Int. Port of toBigIntegerList().
func (b *ByteRange) ToBigIntegerList() []*big.Int {
	return spi.DSSUtilsToBigIntegerList(b.byteRangeArray)
}

// IsValid returns whether the /ByteRange is valid. Port of isValid().
//
// Panics with the Java message when validation has not been performed yet, matching Java's
// IllegalStateException.
func (b *ByteRange) IsValid() bool {
	if b.valid == nil {
		panic("ByteRange validation has not been performed! " +
			"Validate the ByteRange and use setValid(valid) method to provide the result.")
	}
	return *b.valid
}

// SetValid sets whether /ByteRange has passed the verification against the PDF document.
// Port of setValid(boolean).
func (b *ByteRange) SetValid(valid bool) {
	b.valid = &valid
}

// Validate checks a validity of the ByteRange according to PDF specifications. This method
// verifies the array of integers representing the ByteRange itself without taking into account
// the PDF document, nor the /Contents octets.
//
// NOTE: this method returns an *exception.IllegalInputException error when an error is
// encountered and does not update the state of the object (Java throws it, uncaught, for the
// same reason). Use SetValid(valid) to define the validity of the ByteRange.
// Port of validate().
func (b *ByteRange) Validate() error {
	if b.byteRangeArray == nil || len(b.byteRangeArray) != 4 {
		return exception.NewIllegalInputException("Incorrect ByteRange size")
	}

	rangeStart := b.byteRangeArray[0]
	firstPartLength := b.byteRangeArray[1]
	secondPartStart := b.byteRangeArray[2]
	secondPartLength := b.byteRangeArray[3]

	if rangeStart != 0 {
		return exception.NewIllegalInputException("The ByteRange must cover start of file")
	}
	if firstPartLength < 0 {
		return exception.NewIllegalInputException("The first hash part doesn't cover anything")
	}
	if secondPartStart < rangeStart+firstPartLength {
		return exception.NewIllegalInputException("The second hash part must start after the first hash part")
	}
	if secondPartLength < 0 {
		return exception.NewIllegalInputException("The second hash part doesn't cover anything")
	}
	return nil
}

// String ports toString().
func (b *ByteRange) String() string {
	return fmt.Sprint(b.byteRangeArray)
}

// Equals ports equals(Object).
func (b *ByteRange) Equals(other *ByteRange) bool {
	if b == other {
		return true
	}
	if other == nil {
		return false
	}
	if len(b.byteRangeArray) != len(other.byteRangeArray) {
		return false
	}
	for i, v := range b.byteRangeArray {
		if v != other.byteRangeArray[i] {
			return false
		}
	}
	return true
}
