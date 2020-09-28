package securitycenter

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceArmJitNetworkAccessPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmJitNetworkAccessPoliciesCreateOrUpdate,
		Read:   resourceArmJitNetworkAccessPoliciesRead,
		Update: resourceArmJitNetworkAccessPoliciesCreateOrUpdate,
		Delete: resourceArmJitNetworkAccessPoliciesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// Timeouts: &schema.ResourceTimeout{
		// 	Create: schema.DefaultTimeout(60 * time.Minute),
		// 	Read:   schema.DefaultTimeout(5 * time.Minute),
		// 	Update: schema.DefaultTimeout(60 * time.Minute),
		// 	Delete: schema.DefaultTimeout(60 * time.Minute),
		// },

		Schema: map[string]*schema.Schema{
			"asc_location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_machines": virtualMachinesPolicySchema(),
		},
	}
}

func virtualMachinesPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"ports": portsPolicySchema(),
			},
		},
	}
}

func portsPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port": {
					Type:         schema.TypeInt,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.PortNumber,
				},
				"protocol": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Default:  string(security.TCP),
					ValidateFunc: validation.StringInSlice([]string{
						string(security.TCP),
						string(security.UDP),
					}, false),
				},
				"allowed_source_address_prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"max_request_access_duration": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validate.ISO8601Duration,
				},
			},
		},
	}
}

func buildJitNetworkAccessPortRules(portData map[string]interface{}) *[]security.JitNetworkAccessPortRule {
	ports := make([]security.JitNetworkAccessPortRule, 0)
	if v, ok := portData["ports"].([]interface{}); ok { //&& len(v.List()) > 0 {
		for _, portConfig := range v {
			data := portConfig.(map[string]interface{})
			number := int32(data["port"].(int))
			protocol := security.Protocol(data["protocol"].(string))
			allowed_source_address_prefix := data["allowed_source_address_prefix"].(string)
			max_request_access_duration := data["max_request_access_duration"].(string)
			port := security.JitNetworkAccessPortRule{
				Number:                     &number,
				Protocol:                   protocol,
				AllowedSourceAddressPrefix: &allowed_source_address_prefix,
				MaxRequestAccessDuration:   &max_request_access_duration,
			}
			ports = append(ports, port)
		}
	}
	return &ports
}

func resourceArmJitNetworkAccessPoliciesCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.JitNetworkAccessPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachinesConfig := d.Get("virtual_machines").([]interface{})
	virtualMachines := make([]security.JitNetworkAccessPolicyVirtualMachine, 0)

	// number := int32(22)
	// protocol := security.Protocol("TCP")
	// allowed_source_address_prefix := "*"
	// max_request_access_duration := "PT3H"
	// ports := []security.JitNetworkAccessPortRule{
	// 	{
	// 		Number:                     &number,
	// 		Protocol:                   protocol,
	// 		AllowedSourceAddressPrefix: &allowed_source_address_prefix,
	// 		MaxRequestAccessDuration:   &max_request_access_duration,
	// 	},
	// }

	for _, virtualMachineConfig := range virtualMachinesConfig {
		data := virtualMachineConfig.(map[string]interface{})
		id := data["id"].(string)
		virtualMachine := security.JitNetworkAccessPolicyVirtualMachine{
			ID:    &id,
			Ports: buildJitNetworkAccessPortRules(data), // &ports,
		}
		virtualMachines = append(virtualMachines, virtualMachine)
	}

	accessRequests := make([]security.JitNetworkAccessRequest, 0)
	// TODO: build access request
	jitNetworkAccessPolicyProperties := security.JitNetworkAccessPolicyProperties{
		VirtualMachines: &virtualMachines,
		Requests:        &accessRequests,
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceGroup := d.Get("resource_group_name").(string)
	kind := d.Get("kind").(string)
	location := "uksouth"
	name := "default"
	policyId := fmt.Sprintf(
		"/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Security/locations/%v/jitNetworkAccessPolicies/%v",
		subscriptionId, resourceGroup, location, name,
	)

	policyType := "Microsoft.Security/locations/jitNetworkAccessPolicies"

	jitNetworkAccessPolicy := security.JitNetworkAccessPolicy{
		ID:                               &policyId,
		Name:                             &name,
		Type:                             &policyType,
		Kind:                             &kind,
		Location:                         &location,
		JitNetworkAccessPolicyProperties: &jitNetworkAccessPolicyProperties,
	}

	ascLocation := d.Get("asc_location").(string)

	log.Printf("[DEBUG] Ports %v", *jitNetworkAccessPolicyProperties.VirtualMachines)

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, name, ascLocation, jitNetworkAccessPolicy)
	if err != nil {
		return fmt.Errorf("Error creating/updating Security Center Subscription pricing: %+v", err)
	}

	// resp, err := client.GetSubscriptionPricing(ctx, name)
	// if err != nil {
	// 	return fmt.Errorf("Error reading Security Center Subscription pricing: %+v", err)
	// }
	// if resp.ID == nil {
	// 	return fmt.Errorf("Security Center Subscription pricing ID is nil")
	// }

	d.SetId(*resp.ID)

	return resourceArmJitNetworkAccessPoliciesRead(d, meta)
}

func resourceArmJitNetworkAccessPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}

func resourceArmJitNetworkAccessPoliciesDelete(_ *schema.ResourceData, _ interface{}) error {
	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}
