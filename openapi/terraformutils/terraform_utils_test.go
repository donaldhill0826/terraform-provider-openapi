package terraformutils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTerraformUtilsGetTerraformPluginsVendorDir(t *testing.T) {
	Convey("Given an TerraformUtils set up with a homeDir and darwin platform", t, func() {
		t := TerraformUtils{
			HomeDir:  "/Users/username",
			Platform: "darwin",
		}
		Convey("When GetTerraformPluginsVendorDir is called", func() {
			vendorDir, err := t.GetTerraformPluginsVendorDir()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And vendor dir should be the default one as no os was specified`", func() {
				So(vendorDir, ShouldEqual, "/Users/username/.terraform.d/plugins")
			})
		})
	})
	Convey("Given an TerraformUtils set up with a homeDir and linux platform", t, func() {
		t := TerraformUtils{
			HomeDir:  "/Users/username",
			Platform: "linux",
		}
		Convey("When GetTerraformPluginsVendorDir is called", func() {
			vendorDir, err := t.GetTerraformPluginsVendorDir()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And vendor dir should be the default one as no os was specified`", func() {
				So(vendorDir, ShouldEqual, "/Users/username/.terraform.d/plugins")
			})
		})
	})
	Convey("Given an TerraformUtils set up with a homeDir and windows platform", t, func() {
		homeDir := "C:\\Users\\username\\"
		t := TerraformUtils{
			Platform: "windows",
			HomeDir:  homeDir,
		}
		Convey("When GetTerraformPluginsVendorDir is called", func() {
			vendorDir, err := t.GetTerraformPluginsVendorDir()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And vendor dir should be the windows specific path`", func() {
				So(vendorDir, ShouldEqual, fmt.Sprintf("%s\\AppData\\terraform.d\\plugins", homeDir))
			})
		})
	})
	Convey("Given an TerraformUtils missing the homeDir configuration", t, func() {
		t := TerraformUtils{
			HomeDir:  "",
			Platform: "darwin",
		}
		Convey("When GetTerraformPluginsVendorDir is called", func() {
			_, err := t.GetTerraformPluginsVendorDir()
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And error message should match the expected one`", func() {
				So(err.Error(), ShouldEqual, "mandatory HomeDir value missing")
			})
		})
	})
	Convey("Given an TerraformUtils missing the platform configuration", t, func() {
		t := TerraformUtils{
			HomeDir:  "/Users/username",
			Platform: "",
		}
		Convey("When GetTerraformPluginsVendorDir is called", func() {
			_, err := t.GetTerraformPluginsVendorDir()
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And error message should match the expected one`", func() {
				So(err.Error(), ShouldEqual, "mandatory platform information is missing")
			})
		})
	})
}

func TestConvertToTerraformCompliantFieldName(t *testing.T) {
	Convey("Given a name that is terraform name compliant", t, func() {
		propertyName := "some_prop_name_that_is_terraform_field_name_compliant"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that is NOT terraform name compliant", t, func() {
		propertyName := "thisIsACamelCaseNameWhichIsNotTerraformNameCompliant"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, "this_is_a_camel_case_name_which_is_not_terraform_name_compliant")
			})
		})
	})

	Convey("Given a name that is terraform name compliant but with numbers with no _ between number and next word", t, func() {
		propertyName := "cdns_v1id"
		expected := "cdns_v1_id"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, expected)
			})
		})
	})

	Convey("Given a name that is terraform name compliant with one number in the middle", t, func() {
		propertyName := "cdns_v1_id"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that is terraform name compliant with multiple numbers", t, func() {
		propertyName := "cdns_v1_firewall_v2_id"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that is terraform name compliant with one number at the end", t, func() {
		propertyName := "cdns_v1"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that ends with double underscore ( __ )", t, func() {
		propertyName := "cdns__"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The double underscore is left intact", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that has v_1 on purpose", t, func() {
		propertyName := "cdns_v_1"
		expected := "cdns_v1"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			//TODO possible bug?
			Convey("The underscore between v and 1 is removed", func() {
				So(fieldName, ShouldEqual, expected)
			})
		})
	})

	Convey("Given a name that ends with __1", t, func() {
		propertyName := "cdns__1"
		expected := "cdns_1"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			//TODO possible bug?
			Convey("The double underscore is changed to single underscore, ", func() {
				So(fieldName, ShouldEqual, expected)
			})
		})
	})

	Convey("Given a name that ends with v1_", t, func() {
		propertyName := "cdns_v1_"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is left intact, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name that ends with _1", t, func() {
		propertyName := "cdns_1"
		expectedName := "cdns1"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			//TODO possible bug?
			Convey("The string returned has the underscore removed, ", func() {
				So(fieldName, ShouldEqual, expectedName)
			})
		})
	})

	Convey("Given a name with underscore at the end", t, func() {
		propertyName := "cdns_"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is left intact, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name with leading and trailing underscores", t, func() {
		propertyName := "_cdns_"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			Convey("The string returned is left intact, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a name 1", t, func() {
		propertyName := "1"
		Convey("ConvertToTerraformCompliantName panics", func() {
			//TODO possible bug?
			So(func() { ConvertToTerraformCompliantName(propertyName) }, ShouldPanic)
		})
	})

	Convey("Given a name with a number and an underscore at the end", t, func() {
		propertyName := "cdns_1_"
		expectedName := "cdns1_"
		Convey("When ConvertToTerraformCompliantName method is called", func() {
			fieldName := ConvertToTerraformCompliantName(propertyName)
			//TODO possible bug?
			Convey("The string returned has the first underscore remove and the trailing one left intact, ", func() {
				So(fieldName, ShouldEqual, expectedName)
			})
		})
	})
}

func TestCreateSchema(t *testing.T) {
	Convey("Given an environment variable, schemaType of type string, required property and an empty default value", t, func() {
		propertyName := "propertyName"
		envVariableValue := "someValue"
		defaultValue := ""
		os.Setenv(strings.ToUpper(propertyName), envVariableValue)
		schemaType := schema.TypeString
		required := true
		Convey("When createSchema method is called", func() {
			schema := createSchema(propertyName, schemaType, required, defaultValue)
			Convey("Then the schema returned should be of type string", func() {
				So(schema.Type, ShouldEqual, schemaType)
			})
			Convey("And the schema returned should be required", func() {
				So(schema.Required, ShouldEqual, required)
			})
			Convey("And the schema default function should return the value set for te environment variable", func() {
				value, err := schema.DefaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldEqual, envVariableValue)
			})
		})
		os.Unsetenv(strings.ToUpper(propertyName))
	})
	Convey("Given a schemaType of type bool, an optional property and an empty default value", t, func() {
		schemaType := schema.TypeBool
		required := false
		defaultValue := ""
		Convey("When createSchema method is called", func() {
			schema := createSchema("propertyName", schemaType, required, defaultValue)
			Convey("Then the schema returned should be of type bool", func() {
				So(schema.Type, ShouldEqual, schemaType)
			})
			Convey("And the schema returned should be optional", func() {
				So(schema.Optional, ShouldEqual, !required)
			})
			Convey("And the schema default function should return nil as there's no env variable matching the property name", func() {
				value, err := schema.DefaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldBeNil)
			})
		})
	})

	Convey("Given a schemaType of type string, required property and a NON empty default value", t, func() {
		propertyName := "propertyName"
		defaultValue := "defaultValue"
		schemaType := schema.TypeString
		required := true
		Convey("When createSchema method is called", func() {
			schema := createSchema(propertyName, schemaType, required, defaultValue)
			Convey("Then the schema returned should be of type string", func() {
				So(schema.Type, ShouldEqual, schemaType)
			})
			Convey("And the schema returned should be required", func() {
				So(schema.Required, ShouldEqual, required)
			})
			Convey("And the schema default function should return the value set for te environment variable", func() {
				value, err := schema.DefaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldEqual, defaultValue)
			})
		})
		os.Unsetenv(strings.ToUpper(propertyName))
	})
}

func TestCreateStringSchema(t *testing.T) {
	Convey("Given a required property of type string", t, func() {
		required := true
		Convey("When CreateStringSchemaProperty method is called", func() {
			s := CreateStringSchemaProperty("propertyName", required, "")
			Convey("Then the schema returned should be of type string", func() {
				So(s.Type, ShouldEqual, schema.TypeString)
			})
			Convey("And the schema returned should be required", func() {
				So(s.Required, ShouldEqual, required)
			})
			Convey("And the schema default function should return nil as there's no env variable matching the property name", func() {
				value, err := s.DefaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldBeNil)
			})
		})
	})
}

func TestEnvDefaultFunc(t *testing.T) {
	Convey("Given a property name that has an environment variable set up and nil default value", t, func() {
		propertyName := "propertyName"
		envVariableValue := "someValue"
		os.Setenv(strings.ToUpper(propertyName), envVariableValue)
		Convey("When envDefaultFunc method is called", func() {
			defaultFunc := envDefaultFunc(propertyName, nil)
			Convey("And the returned defaultFunc is invoked the value returned should be the value of the environment variable", func() {
				value, err := defaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldEqual, envVariableValue)
			})
		})
		os.Unsetenv(strings.ToUpper(propertyName))
	})
	Convey("Given a property name that DOES NOT have an environment variable set up and a default value is configured", t, func() {
		defaultValue := "someDefaultValue"
		Convey("When envDefaultFunc method is called", func() {
			defaultFunc := envDefaultFunc("propertyName", defaultValue)
			Convey("And the returned defaultFunc is invoked the value returned should be the defaultValue", func() {
				value, err := defaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldEqual, defaultValue)
			})
		})
	})
}
