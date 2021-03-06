package openapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecSecuritySchemeGetTerraformConfigurationName(t *testing.T) {
	Convey("Given a SpecSecurityScheme with a terraform compliant name", t, func() {
		expectedName := "some_compliant_name"
		specSecurityScheme := SpecSecurityScheme{Name: expectedName}
		Convey("When newAPIKeySecurityDefinition method is called", func() {
			secSchemeTerraformName := specSecurityScheme.getTerraformConfigurationName()
			Convey("Then the secSchemeTerraformName name should match", func() {
				So(secSchemeTerraformName, ShouldEqual, expectedName)
			})
		})
	})
	Convey("Given a SpecSecurityScheme with a Non terraform compliant name", t, func() {
		specSecurityScheme := SpecSecurityScheme{Name: "nonCompliantName"}
		Convey("When newAPIKeySecurityDefinition method is called", func() {
			secSchemeTerraformName := specSecurityScheme.getTerraformConfigurationName()
			Convey("Then the secSchemeTerraformName name should be compliant", func() {
				So(secSchemeTerraformName, ShouldEqual, "non_compliant_name")
			})
		})
	})
}

func TestSecuritySchemeExists(t *testing.T) {
	Convey("Given a list of specSecuritySchemes", t, func() {
		securitySchemes := []map[string][]string{
			{
				"secDef1": {},
			},
		}
		specSecuritySchemes := createSecuritySchemes(securitySchemes)
		Convey("When securitySchemeExists method is called with an existing security definition", func() {
			secDef := newAPIKeyHeaderSecurityDefinition("secDef1", "secDef1")
			exists := specSecuritySchemes.securitySchemeExists(secDef)
			Convey("Then the specSecuritySchemes should be true", func() {
				So(exists, ShouldBeTrue)
			})
		})
	})

	Convey("Given a list of specSecuritySchemes", t, func() {
		securitySchemes := []map[string][]string{
			{
				"secDef1": {},
			},
		}
		specSecuritySchemes := createSecuritySchemes(securitySchemes)
		Convey("When securitySchemeExists method is called with a NON existing security definition", func() {
			secDef := newAPIKeyHeaderSecurityDefinition("secDefNonExisting", "secDefNonExisting")
			exists := specSecuritySchemes.securitySchemeExists(secDef)
			Convey("Then the specSecuritySchemes should be false", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}

func TestCreateSecuritySchemes(t *testing.T) {
	Convey("Given a map of securitySchemes with multi auth AND support", t, func() {
		securitySchemes := []map[string][]string{
			{
				"secDef1": {},
				"secDef2": {},
			},
		}
		Convey("When createSecuritySchemes method is called with the securitySchemes", func() {
			specSecuritySchemes := createSecuritySchemes(securitySchemes)
			Convey("Then the specSecuritySchemes should not be empty", func() {
				So(specSecuritySchemes, ShouldNotBeEmpty)
			})
			Convey("Then the specSecuritySchemes name contain the expected items", func() {
				So(specSecuritySchemes, ShouldContain, SpecSecurityScheme{Name: "secDef1"})
				So(specSecuritySchemes, ShouldContain, SpecSecurityScheme{Name: "secDef2"})
			})
		})
	})

	Convey("Given a map of securitySchemes with multi auth OR support", t, func() {
		securitySchemes := []map[string][]string{
			{
				"secDef1": {},
				"secDef2": {},
			},
			{
				"secDef3": {},
			},
		}
		Convey("When createSecuritySchemes method is called with the securitySchemes", func() {
			specSecuritySchemes := createSecuritySchemes(securitySchemes)
			Convey("Then the specSecuritySchemes should not be empty", func() {
				So(specSecuritySchemes, ShouldNotBeEmpty)
			})
			Convey("Then the specSecuritySchemes name contain the expected items which are the first one in the array (by design these take preference)", func() {
				So(specSecuritySchemes, ShouldContain, SpecSecurityScheme{Name: "secDef1"})
				So(specSecuritySchemes, ShouldContain, SpecSecurityScheme{Name: "secDef2"})
			})
			Convey("Then the specSecuritySchemes should not contain anything else", func() {
				So(specSecuritySchemes, ShouldNotContain, SpecSecurityScheme{Name: "secDef3"})
			})
		})
	})
}
