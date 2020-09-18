package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSecurityCenterSubscriptionPricing() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterSubscriptionPricingUpdate,
		Read:   resourceArmSecurityCenterSubscriptionPricingRead,
		Update: resourceArmSecurityCenterSubscriptionPricingUpdate,
		Delete: resourceArmSecurityCenterSubscriptionPricingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tier": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.Free),
					string(security.Standard),
				}, false),
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AppServices",
					"ContainerRegistry",
					"KeyVaults",
					"KubernetesService",
					"SqlServers",
					"SqlServerVirtualMachines",
					"StorageAccounts",
					"VirtualMachines",
				}, false),
			},
		},
	}
}

func resourceArmSecurityCenterSubscriptionPricingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// not doing import check as afaik it always exists (cannot be deleted)
	// all this resource does is flip a boolean

	pricing := security.Pricing{
		PricingProperties: &security.PricingProperties{
			PricingTier: security.PricingTier(d.Get("tier").(string)),
		},
	}

	resource_type := d.Get("resource_type").(string)

	if _, err := client.Update(ctx, resource_type, pricing); err != nil {
		return fmt.Errorf("Error creating/updating Security Center Subscription pricing: %+v", err)
	}

	resp, err := client.Get(ctx, resource_type)
	if err != nil {
		return fmt.Errorf("Error reading Security Center Subscription pricing: %+v", err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Security Center Subscription pricing ID is nil")
	}

	d.SetId(*resp.ID)

	return resourceArmSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceArmSecurityCenterSubscriptionPricingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resource_type := d.Get("resource_type").(string)

	resp, err := client.Get(ctx, resource_type)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %v Security Center Subscription was not found: %v", resource_type, err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading %v Security Center Subscription pricing: %+v", resource_type, err)
	}

	if properties := resp.PricingProperties; properties != nil {
		d.Set("tier", properties.PricingTier)
	}

	return nil
}

func resourceArmSecurityCenterSubscriptionPricingDelete(_ *schema.ResourceData, _ interface{}) error {
	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil // cannot be deleted.
}
