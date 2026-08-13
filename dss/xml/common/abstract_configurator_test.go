package common

import (
	"errors"
	"testing"

	"github.com/utain/esig/dss/alert"
)

func TestAbstractConfigurator_SetSecurityFeatures_CollectsFailuresAndAlerts(t *testing.T) {
	setFeature := func(factory *int, feature string, value bool) error {
		if feature == "bad" {
			return NewSecurityConfigurationException(errors.New("boom"))
		}
		return nil
	}
	setAttribute := func(factory *int, attribute string, value any) error {
		return nil
	}

	c := NewAbstractConfigurator[*int](setFeature, setAttribute)
	c.EnableFeature("good")
	c.EnableFeature("bad")

	factory := new(int)
	err := c.SetSecurityFeatures(factory)
	if err == nil {
		t.Fatalf("SetSecurityFeatures() error = nil, want an alert error for the failing feature")
	}
}

func TestAbstractConfigurator_SetSecurityFeatures_NoFailuresSucceeds(t *testing.T) {
	c := NewAbstractConfigurator[*int](
		func(factory *int, feature string, value bool) error { return nil },
		func(factory *int, attribute string, value any) error { return nil },
	)
	c.EnableFeature("good")
	if err := c.SetSecurityFeatures(new(int)); err != nil {
		t.Fatalf("SetSecurityFeatures() error = %v, want nil", err)
	}
}

func TestAbstractConfigurator_SilentOnStatusAlertSuppressesError(t *testing.T) {
	c := NewAbstractConfigurator[*int](
		func(factory *int, feature string, value bool) error {
			return NewSecurityConfigurationException(errors.New("boom"))
		},
		func(factory *int, attribute string, value any) error { return nil },
	)
	c.SetSecurityExceptionAlert(alert.NewSilentOnStatusAlert())
	c.EnableFeature("bad")
	if err := c.SetSecurityFeatures(new(int)); err != nil {
		t.Fatalf("SetSecurityFeatures() error = %v, want nil under a silent alert", err)
	}
}

func TestAbstractConfigurator_SetSecurityExceptionAlertPanicsOnNil(t *testing.T) {
	c := NewAbstractConfigurator[*int](nil, nil)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when securityExceptionAlert is nil")
		}
	}()
	c.SetSecurityExceptionAlert(nil)
}

func TestAbstractConfigurator_RemoveAttribute(t *testing.T) {
	c := NewAbstractConfigurator[*int](
		func(factory *int, feature string, value bool) error { return nil },
		func(factory *int, attribute string, value any) error { return nil },
	)
	c.SetAttribute("a", 1)
	c.RemoveAttribute("a")
	if _, exists := c.attributes["a"]; exists {
		t.Errorf("attribute 'a' should have been removed")
	}
	if len(c.attributeNames) != 0 {
		t.Errorf("attributeNames should be empty after removal, got %v", c.attributeNames)
	}
}
