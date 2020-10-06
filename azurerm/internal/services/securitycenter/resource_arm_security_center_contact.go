package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/azuresdkhacks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// default name to keep resource backwards compatible
const securityCenterContactName = "default1"

func resourceArmSecurityCenterContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterContactCreateUpdate,
		Read:   resourceArmSecurityCenterContactRead,
		Update: resourceArmSecurityCenterContactCreateUpdate,
		Delete: resourceArmSecurityCenterContactDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceArmSecurityCenterContactV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceArmSecurityCenterContactUpgradeV0ToV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      securityCenterContactName,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
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

func resourceArmSecurityCenterContactCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Checking for presence of existing Security Center Contact: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_contact", *existing.ID)
		}
	}

	contact := security.Contact{
		ContactProperties: &security.ContactProperties{
			Email: utils.String(d.Get("email").(string)),
			Phone: utils.String(d.Get("phone").(string)),
		},
	}

	if alertNotifications := d.Get("alert_notifications").(bool); alertNotifications {
		contact.AlertNotifications = security.On
	} else {
		contact.AlertNotifications = security.Off
	}

	if alertNotifications := d.Get("alerts_to_admins").(bool); alertNotifications {
		contact.AlertsToAdmins = security.AlertsToAdminsOn
	} else {
		contact.AlertsToAdmins = security.AlertsToAdminsOff
	}

	if d.IsNewResource() {
		if _, err := azuresdkhacks.CreateSecurityCenterContact(client, ctx, name, contact); err != nil {
			return fmt.Errorf("Creating Security Center Contact: %+v", err)
		}
	} else if _, err := client.Update(ctx, name, contact); err != nil {
		return fmt.Errorf("Updating Security Center Contact: %+v", err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Reading Security Center Contact: %+v", err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Security Center Contact ID is nil")
	}
	d.SetId(*resp.ID)

	return resourceArmSecurityCenterContactRead(d, meta)
}

func resourceArmSecurityCenterContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SecurityCenterContactID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Reading Security Center Contact: %+v", err)
	}

	if properties := resp.ContactProperties; properties != nil {
		d.Set("email", properties.Email)
		d.Set("phone", properties.Phone)
		d.Set("alert_notifications", properties.AlertNotifications == security.On)
		d.Set("alerts_to_admins", properties.AlertsToAdmins == security.AlertsToAdminsOn)
	}
	d.Set("name", resp.Name)

	return nil
}

func resourceArmSecurityCenterContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.ContactsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SecurityCenterContactID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Contact was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Deleting Security Center Contact: %+v", err)
	}

	return nil
}
