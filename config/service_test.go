package config_test

import (
	"testing"

	"github.com/rebel-l/branma_be/config"
)

func TestService_GetPort(t *testing.T) {
	var service *config.Service
	if service.GetPort() != config.DefaultPort {
		t.Errorf(
			"failed to retrieve default value %d from nil struct, but got %d",
			config.DefaultPort,
			service.GetPort(),
		)
	}
}

type tcServiceMerge struct {
	name      string
	actual    *config.Service
	mergeWith *config.Service
	expected  *config.Service
}

func getTestCasesServiceMerge(t *testing.T) []tcServiceMerge { // nolint:funlen
	t.Helper()

	var testCases []tcServiceMerge

	port := 5000

	newPort := 6000

	// 1.
	tc := tcServiceMerge{
		name:      "config nil",
		mergeWith: &config.Service{Port: &port},
	}

	testCases = append(testCases, tc)

	// 2.
	tc = tcServiceMerge{
		name:     "parameter nil",
		actual:   &config.Service{Port: &port},
		expected: &config.Service{Port: &port},
	}

	testCases = append(testCases, tc)

	// 3.
	tc = tcServiceMerge{
		name:      "config has default values, parameter has values",
		actual:    &config.Service{},
		mergeWith: &config.Service{Port: &port},
		expected:  &config.Service{Port: &port},
	}

	testCases = append(testCases, tc)

	// 4.
	tc = tcServiceMerge{
		name:      "config has values, parameter has values",
		actual:    &config.Service{Port: &port},
		mergeWith: &config.Service{Port: &newPort},
		expected:  &config.Service{Port: &newPort},
	}

	testCases = append(testCases, tc)

	// 5.
	tc = tcServiceMerge{
		name:      "config has values, parameter has default values",
		actual:    &config.Service{Port: &port},
		mergeWith: &config.Service{},
		expected:  &config.Service{Port: &port},
	}

	testCases = append(testCases, tc)

	// 6.
	tc = tcServiceMerge{
		name:      "config has default values, parameter has default values",
		actual:    &config.Service{},
		mergeWith: &config.Service{},
		expected:  &config.Service{},
	}

	testCases = append(testCases, tc)

	return testCases
}

func TestService_Merge(t *testing.T) {
	for _, testCase := range getTestCasesServiceMerge(t) {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.actual.Merge(testCase.mergeWith)
			testService(t, testCase.expected, testCase.actual)
		})
	}
}

func testService(t *testing.T, expected, got *config.Service) {
	t.Helper()

	if expected.GetPort() != got.GetPort() {
		t.Errorf("failed to set service port: expected '%d' but got '%d'",
			expected.GetPort(), got.GetPort())
	}
}
