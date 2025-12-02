package ztc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/ecgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/provisioning/provisioning_url"
)

func dataSourceProvisioningURL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProvisioningURLRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prov_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prov_url_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prov_url_data": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zs_cloud_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"config_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registration_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pac_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_template": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"editable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"last_mod_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"template": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_prefix": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"xff_forward_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"auth_required": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"caution_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"aup_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"aup_timeout_in_days": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"ofw_enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"ips_control": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"enforce_bandwidth_control": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"up_bandwidth": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"dn_bandwidth": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"idle_time_in_minutes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"surrogate_ip_enforced_for_known_browsers": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"surrogate_refresh_time_in_minutes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"surrogate_ip": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"cloud_provider_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"form_factor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_scale_details": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_scale": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"cell_edge_deploy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"release_channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_mod_uid": UIDNameSchema(),
			"last_mod_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceProvisioningURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *provisioning_url.ProvisioningURL
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for provisioning url id: %d\n", id)
		res, err := provisioning_url.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for provisioning url name: %s\n", name)
		res, err := provisioning_url.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("desc", resp.Desc)
		_ = d.Set("prov_url", resp.ProvUrl)
		_ = d.Set("prov_url_type", resp.ProvUrlType)
		_ = d.Set("status", resp.Status)
		_ = d.Set("last_mod_time", resp.LastModTime)

		if err := d.Set("prov_url_data", flattenProvURLData(&resp.ProvUrlData)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_mod_uid", flattenIDExtensionsList(resp.LastModUid)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.Errorf("couldn't find any provisioning url with name '%s'", name)
	}

	return nil
}

func flattenProvURLData(provUrlData *provisioning_url.ProvUrlData) []map[string]interface{} {
	if provUrlData == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"zs_cloud_domain":     provUrlData.ZsCloudDomain,
			"org_id":              provUrlData.OrgID,
			"config_server":       provUrlData.ConfigServer,
			"registration_server": provUrlData.RegistrationServer,
			"api_server":          provUrlData.ApiServer,
			"pac_server":          provUrlData.PacServer,
			"cloud_provider_type": provUrlData.CloudProviderType,
			"form_factor":         provUrlData.FormFactor,
			"cell_edge_deploy":    provUrlData.CellEdgeDeploy,
			"release_channel":     provUrlData.ReleaseChannel,
			"location_template":   flattenLocationTemplateFromProvURL(&provUrlData.LocationTemplate),
			"auto_scale_details":  flattenAutoScaleDetails(&provUrlData.AutoScaleDetails),
		},
	}
}

func flattenLocationTemplateFromProvURL(locTemplate *locationtemplate.LocationTemplate) []map[string]interface{} {
	if locTemplate == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":            locTemplate.ID,
			"name":          locTemplate.Name,
			"desc":          locTemplate.Description,
			"editable":      locTemplate.Editable,
			"last_mod_time": locTemplate.LastModTime,
			"template":      flattenTemplateDetailsFromProvURL(locTemplate.LocationTemplateDetails),
		},
	}
}

func flattenTemplateDetailsFromProvURL(template *locationtemplate.LocationTemplateDetails) []map[string]interface{} {
	if template == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"template_prefix":           template.TemplatePrefix,
			"xff_forward_enabled":       template.XFFForwardEnabled,
			"auth_required":             template.AuthRequired,
			"caution_enabled":           template.CautionEnabled,
			"aup_enabled":               template.AupEnabled,
			"aup_timeout_in_days":       template.AupTimeoutInDays,
			"ofw_enabled":               template.OFWEnabled,
			"ips_control":               template.IPSControl,
			"enforce_bandwidth_control": template.EnforceBandwidthControl,
			"up_bandwidth":              template.UpBandwidth,
			"dn_bandwidth":              template.DnBandwidth,
			"idle_time_in_minutes":      template.IdleTimeInMinutes,
			"surrogate_ip_enforced_for_known_browsers": template.SurrogateIPEnforcedForKnownBrowsers,
			"surrogate_refresh_time_in_minutes":        template.SurrogateRefreshTimeInMinutes,
			"surrogate_ip":                             template.SurrogateIP,
		},
	}
}

func flattenAutoScaleDetails(autoScaleDetails *provisioning_url.AutoScaleDetails) []map[string]interface{} {
	if autoScaleDetails == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"auto_scale": autoScaleDetails.AutoScale,
		},
	}
}

func flattenLocationTemplate(locTemplate *locationtemplate.LocationTemplate) []map[string]interface{} {
	if locTemplate == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":            locTemplate.ID,
			"name":          locTemplate.Name,
			"desc":          locTemplate.Description,
			"editable":      locTemplate.Editable,
			"last_mod_time": locTemplate.LastModTime,
			"template":      flattenTemplateDetails(locTemplate.LocationTemplateDetails),
			"last_mod_uid":  flattenCommonIDNameExternalID(locTemplate.LastModUid),
		},
	}
}

func flattenUIDNamesListIds(uidNames []common.IDNameExtensions) []map[string]interface{} {
	if len(uidNames) == 0 {
		return nil
	}
	result := []map[string]interface{}{}
	ids := []int{}
	for _, uidname := range uidNames {
		ids = append(ids, uidname.ID)
	}
	result = append(result, map[string]interface{}{
		"id": ids,
	})
	return result
}

func flattenECInstance(ecInstance *common.ECInstances) []map[string]interface{} {
	if ecInstance == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ec_instance_type": ecInstance.ECInstanceType,
			"out_gw_ip":        ecInstance.OutGwIp,
			"nat_ip":           ecInstance.NatIP,
			"dns_ip":           ecInstance.DNSIP,
			"service_nw":       flattenNW(ecInstance.ServiceNw),
			"virtual_nw":       flattenNW(ecInstance.VirtualNw),
		},
	}
}

func flattenECInstances(ecInstances []common.ECInstances) []map[string]interface{} {
	if ecInstances == nil {
		return nil
	}
	result := []map[string]interface{}{}
	for _, ecInstance := range ecInstances {
		result = append(result, flattenECInstance(&ecInstance)...)
	}
	return result
}

func flattenECVm(ecvm *common.ECVMs) []map[string]interface{} {
	if ecvm == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":                 ecvm.ID,
			"name":               ecvm.Name,
			"form_factor":        ecvm.FormFactor,
			"management_nw":      flattenNW(ecvm.ManagementNw),
			"ec_instances":       flattenECInstances(ecvm.ECInstances),
			"city_geo_id":        ecvm.CityGeoId,
			"nat_ip":             ecvm.NATIP,
			"zia_gateway":        ecvm.ZiaGateway,
			"zpa_broker":         ecvm.ZpaBroker,
			"build_version":      ecvm.BuildVersion,
			"last_upgrade_time":  ecvm.LastUpgradeTime,
			"upgrade_status":     ecvm.UpgradeStatus,
			"upgrade_start_time": ecvm.UpgradeStartTime,
			"upgrade_end_time":   ecvm.UpgradeEndTime,
		},
	}
}

func flattenECVms(ecvms []common.ECVMs) []map[string]interface{} {
	if len(ecvms) == 0 {
		return nil
	}
	result := []map[string]interface{}{}
	for _, ecvm := range ecvms {
		result = append(result, flattenECVm(&ecvm)...)
	}
	return result
}

func flattenBCGroup(bcGroup *ecgroup.EcGroup) []map[string]interface{} {
	if bcGroup == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":                      bcGroup.ID,
			"name":                    bcGroup.Name,
			"desc":                    bcGroup.Description,
			"deploy_type":             bcGroup.DeployType,
			"platform":                bcGroup.Platform,
			"aws_availability_zone":   bcGroup.AWSAvailabilityZone,
			"azure_availability_zone": bcGroup.AzureAvailabilityZone,
			"max_ec_count":            bcGroup.MaxEcCount,
			"tunnel_mode":             bcGroup.TunnelMode,
			"location":                flattenGeneralPurpose(bcGroup.Location),
			"prov_template":           flattenGeneralPurpose(bcGroup.ProvTemplate),
			"ec_vms":                  flattenECVms(bcGroup.ECVMs),
		},
	}
}
