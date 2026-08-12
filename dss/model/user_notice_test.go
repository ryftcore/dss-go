package model

import "testing"

func TestUserNoticeIsEmpty(t *testing.T) {
	u := NewUserNotice()
	if !u.IsEmpty() {
		t.Fatal("expected fresh UserNotice to be empty")
	}
	u.SetExplicitText("some text")
	if u.IsEmpty() {
		t.Fatal("expected UserNotice with explicit text to not be empty")
	}
}

func TestUserNoticeRoundTripAndEquals(t *testing.T) {
	a := NewUserNotice()
	a.SetOrganization("ACME")
	a.SetNoticeNumbers(1, 2, 3)
	a.SetExplicitText("notice")

	b := NewUserNotice()
	b.SetOrganization("ACME")
	b.SetNoticeNumbers(1, 2, 3)
	b.SetExplicitText("notice")

	if !a.Equals(b) {
		t.Fatalf("expected equal UserNotices to be Equals(): a=%s b=%s", a.String(), b.String())
	}

	b.SetNoticeNumbers(4, 5)
	if a.Equals(b) {
		t.Fatal("expected different notice numbers to not be Equals()")
	}
}
