package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/accessanalyzer"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

// AccessAnalyzer is limited to one per region, so run serially
// locally and in TeamCity.
func TestAccAWSAccessAnalyzer_serial(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"Analyzer": {
			"basic":             testAccAWSAccessAnalyzerAnalyzer_basic,
			"disappears":        testAccAWSAccessAnalyzerAnalyzer_disappears,
			"Tags":              testAccAWSAccessAnalyzerAnalyzer_Tags,
			"Type_Organization": testAccAWSAccessAnalyzerAnalyzer_Type_Organization,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccPreCheckAWSAccessAnalyzer(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AccessAnalyzerConn

	input := &accessanalyzer.ListAnalyzersInput{}

	_, err := conn.ListAnalyzers(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
