package google

import (
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func replaceVarsForTest(config *transport_tpg.Config, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider config
		if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}

// This function isn't a test of transport.go; instead, it is used as an alternative
// to ReplaceVars inside tests.
func replaceVarsForFrameworkTest(prov *frameworkProvider, rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")
	var project, region, zone string

	if strings.Contains(linkTmpl, "{{project}}") {
		project = rs.Primary.Attributes["project"]
	}

	if strings.Contains(linkTmpl, "{{region}}") {
		region = GetResourceNameFromSelfLink(rs.Primary.Attributes["region"])
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
		zone = GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])
	}

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}

		if v, ok := rs.Primary.Attributes[m]; ok {
			return v
		}

		// Attempt to draw values from the provider
		if f := reflect.Indirect(reflect.ValueOf(prov)).FieldByName(m); f.IsValid() {
			return f.String()
		}

		return ""
	}

	return re.ReplaceAllStringFunc(linkTmpl, replaceFunc), nil
}
