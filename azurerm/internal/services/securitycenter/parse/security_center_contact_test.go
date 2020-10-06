package parse

import (
	"testing"
)

func TestSecurityCenterContactID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SecurityCenterContactId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No SecurityContact Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No SecurityContact Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/securityContact/",
			Error: true,
		},
		{
			Name:  "Security Center Subscription Pricing ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/securityContact/default",
			Expect: &SecurityCenterContactId{
				Name: "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SecurityCenterContactID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
