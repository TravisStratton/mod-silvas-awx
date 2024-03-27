package awx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
)

func resourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource `awx_notification_template` manages notification templates within an AWX organization.",
		CreateContext: resourceNotificationTemplateCreate,
		ReadContext:   resourceNotificationTemplateRead,
		UpdateContext: resourceNotificationTemplateUpdate,
		DeleteContext: resourceNotificationTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the notification template.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The organization ID to associate with the notification template.",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of notification template.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the notification template.",
			},
			"notification_configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The configuration of the notification template.",
			},
		},
	}
}

func resourceNotificationTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService

	notificationConfigurationStr := d.Get("notification_configuration").(string)
	notificationConfigurationMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(notificationConfigurationStr), &notificationConfigurationMap)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("error while unmarshal notification_configuration: %s", err.Error()),
		})
		return diags
	}

	result, err := awxService.Create(map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": notificationConfigurationMap,
	}, map[string]string{})
	if err != nil {
		log.Printf("Fail to Create notification_template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("NotificationTemplate failed to create %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
	id, diags := convertStateIDToNummeric("Update NotificationTemplate", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("notification_template", id, err)
	}

	notificationConfigurationStr := d.Get("notification_configuration").(string)
	notificationConfigurationMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(notificationConfigurationStr), &notificationConfigurationMap)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create NotificationTemplate",
			Detail:   fmt.Sprintf("error while unmarshal notification_configuration: %s", err.Error()),
		})
		return diags
	}

	_, err = awxService.Update(id, map[string]interface{}{
		"name":                       d.Get("name").(string),
		"description":                d.Get("description").(string),
		"organization":               d.Get("organization_id").(string),
		"notification_type":          d.Get("notification_type").(string),
		"notification_configuration": notificationConfigurationMap,
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update NotificationTemplate",
			Detail:   fmt.Sprintf("notification_template with name %s failed to update %s", d.Get("name").(string), err.Error()),
		})
		return diags
	}

	return resourceNotificationTemplateRead(ctx, d, m)
}

func resourceNotificationTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
	id, diags := convertStateIDToNummeric("Read notification_template", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("notification_template", id, err)

	}
	d = setNotificationTemplateResourceData(d, res)
	return nil
}

func resourceNotificationTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.NotificationTemplatesService
	id, diags := convertStateIDToNummeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.Delete(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setNotificationTemplateResourceData(d *schema.ResourceData, r *awx.NotificationTemplate) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	d.Set("notification_type", r.NotificationType)
	d.Set("notification_configuration", r.NotificationConfiguration)
	d.SetId(strconv.Itoa(r.ID))
	return d
}