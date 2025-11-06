package ztc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/partner_integrations/account_groups"
)

func resourceAccountGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountGroupCreate,
		ReadContext:   resourceAccountGroupRead,
		UpdateContext: resourceAccountGroupUpdate,
		DeleteContext: resourceAccountGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("group_id", idInt)
				} else {
					resp, err := account_groups.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("group_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AWS",
			},
			"public_cloud_accounts":  UIDNameSchema(),
			"cloud_connector_groups": UIDNameSchema(),
		},
	}
}

func resourceAccountGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandAccountGroup(d)
	log.Printf("[INFO] Creating zia ip groups\n%+v\n", req)

	resp, err := account_groups.CreateAccountGroups(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia ip groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("group_id", resp.ID)

	return resourceAccountGroupRead(ctx, d, meta)
}

func resourceAccountGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no ip groups id is set"))
	}
	resp, err := account_groups.GetAccountGroup(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if len(resp) == 0 {
		log.Printf("[WARN] No account group found with ID %d", id)
		d.SetId("")
		return nil
	}

	accountGroup := resp[0]
	log.Printf("[INFO] Getting zia ip source groups:\n%+v\n", accountGroup)

	d.SetId(fmt.Sprintf("%d", accountGroup.ID))
	_ = d.Set("group_id", accountGroup.ID)
	_ = d.Set("name", accountGroup.Name)
	_ = d.Set("description", accountGroup.Description)
	_ = d.Set("cloud_type", accountGroup.CloudType)

	if err := d.Set("public_cloud_accounts", flattenIDExtensionsListIDs(accountGroup.PublicCloudAccounts)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_connector_groups", flattenIDExtensionsListIDs(accountGroup.CloudConnectorGroups)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAccountGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia ip groups ID: %v\n", id)
	req := expandAccountGroup(d)
	if _, err := account_groups.GetAccountGroup(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := account_groups.UpdateAccountGroups(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountGroupRead(ctx, d, meta)
}

func resourceAccountGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "group_id")
	if !ok {
		log.Printf("[ERROR] ip groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip groups ID: %v\n", (d.Id()))

	if err := account_groups.DeleteAccountGroups(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia ip groups deleted")

	return nil
}

func expandAccountGroup(d *schema.ResourceData) account_groups.AccountGroups {
	id, _ := getIntFromResourceData(d, "group_id")
	return account_groups.AccountGroups{
		ID:                   id,
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		CloudType:            d.Get("cloud_type").(string),
		PublicCloudAccounts:  expandIDNameExtensionsSet(d, "public_cloud_accounts"),
		CloudConnectorGroups: expandIDNameExtensionsSet(d, "cloud_connector_groups"),
	}
}
