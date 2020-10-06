package securitycenter

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
)

func ResourceArmSecurityCenterContactV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alert_notifications": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"alerts_to_admins": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func ResourceArmSecurityCenterContactUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Println("[DEBUG] Migrating ResourceType from v0 to v1 format")
	oldId := rawState["id"].(string)

	id, err := parse.SecurityCenterContactID(oldId)
	if err != nil {
		return rawState, err
	}

	rawState["name"] = id.Name

	return rawState, nil
}
