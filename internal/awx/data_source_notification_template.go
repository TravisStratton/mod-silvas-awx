package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/josh-silvas/terraform-provider-awx/tools/goawx"
	"github.com/josh-silvas/terraform-provider-awx/tools/utils"
)

func dataSourceNotificationTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationTemplatesRead,
		Description: "Data source for AWX Notification Template",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the Notification Template",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Notification Template",
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceNotificationTemplatesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	notificationTemplates, _, err := client.NotificationTemplatesService.List(params)
	if err != nil {
		return utils.DiagFetch(diagInventoryTitle, params, err)
	}
	if len(notificationTemplates) > 1 {
		return utils.Diagf(
			"Get: find more than one Element",
			"The Query Returns more than one Group, %d",
			len(notificationTemplates),
		)
	}
	if len(notificationTemplates) == 0 {
		return utils.Diagf(
			"Get: Notification Template does not exist",
			"The Query Returns no Notification Template matching filter %v",
			params,
		)
	}

	notificationTemplate := notificationTemplates[0]
	_ = setNotificationTemplateResourceData(d, notificationTemplate)
	return diags
}
