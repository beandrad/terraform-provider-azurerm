package tests

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
)

func TestAzureRMSecurityCenterContactMigrateState(t *testing.T) {
	inputAttributes := map[string]interface{}{
		"id": "/subscriptions/00000000-0000-0000-0000-000000000000/securityContact/default",
	}
	expectedName := "default"

	rawState, _ := securitycenter.ResourceArmSecurityCenterContactUpgradeV0ToV1(inputAttributes, nil)
	if rawState["name"].(string) != expectedName {
		t.Fatalf("ResourceType migration failed, expected %q, got: %q", expectedName, rawState["name"].(string))
	}
}
