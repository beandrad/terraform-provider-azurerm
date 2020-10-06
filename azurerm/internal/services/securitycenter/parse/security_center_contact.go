package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterContactId struct {
	Name string
}

func SecurityCenterContactID(input string) (*SecurityCenterContactId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Security Center Contact ID %q: %+v", input, err)
	}

	contact := SecurityCenterContactId{}

	if contact.Name, err = id.PopSegment("securityContact"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &contact, nil
}
