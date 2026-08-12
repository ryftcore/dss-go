// Ported from dss-model/.../UserNotice.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"
)

// UserNotice may be used to define a User Notice signature policy
// qualifier.
type UserNotice struct {
	// organization is the name of the organization.
	organization string

	// noticeNumbers are numbers identifying a group of textual statements
	// prepared by the organization.
	noticeNumbers []int

	// explicitText is the text of the notice.
	explicitText string
}

// NewUserNotice instantiates the object with null values. Ports the empty
// constructor.
func NewUserNotice() *UserNotice {
	return &UserNotice{}
}

// Organization gets the organization name.
func (u *UserNotice) Organization() string { return u.organization }

// SetOrganization sets the organization name.
//
// NOTE: when the property is not empty, the NoticeNumbers also shall be
// set!
func (u *UserNotice) SetOrganization(organization string) { u.organization = organization }

// NoticeNumbers gets the notice numbers.
func (u *UserNotice) NoticeNumbers() []int { return u.noticeNumbers }

// SetNoticeNumbers sets the notice numbers identifying a group of textual
// statements prepared by the organization.
//
// NOTE: when the property is not empty, the Organization also shall be
// set!
func (u *UserNotice) SetNoticeNumbers(noticeNumbers ...int) { u.noticeNumbers = noticeNumbers }

// ExplicitText gets the notice text.
func (u *UserNotice) ExplicitText() string { return u.explicitText }

// SetExplicitText sets the text of the notice to be displayed.
func (u *UserNotice) SetExplicitText(explicitText string) { u.explicitText = explicitText }

// IsEmpty checks if the content of the UserNotice is empty or not.
func (u *UserNotice) IsEmpty() bool {
	if u.organization != "" {
		return false
	}
	if len(u.noticeNumbers) > 0 {
		return false
	}
	if u.explicitText != "" {
		return false
	}
	return true
}

// Equals ports UserNotice#equals.
func (u *UserNotice) Equals(other *UserNotice) bool {
	if u == other {
		return true
	}
	if other == nil {
		return false
	}
	return u.organization == other.organization &&
		reflect.DeepEqual(u.noticeNumbers, other.noticeNumbers) &&
		u.explicitText == other.explicitText
}

// String ports UserNotice#toString.
func (u *UserNotice) String() string {
	return fmt.Sprintf("UserNotice {organization='%s', noticeNumbers=%v, explicitText='%s'}",
		u.organization, u.noticeNumbers, u.explicitText)
}
