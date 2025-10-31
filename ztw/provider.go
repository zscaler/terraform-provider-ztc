package ztw

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiSemaphore chan struct{}

func ZTWProvider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "zpa client id",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa client secret",
				ConflictsWith: []string{"private_key"},
			},
			"private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Description:   "zpa private key",
				ConflictsWith: []string{"client_secret"},
			},
			"vanity_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Vanity Domain",
			},
			"zscaler_cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Zscaler Cloud Name",
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ztw_cloud": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"zscaler",
					"zscalerone",
					"zscalertwo",
					"zscalerthree",
					"zscloud",
					"zscalerbeta",
					"zscalergov",
					"zscalerten",
					"zspreview",
				}, false),
				Optional: true,
			},
			"use_legacy_client": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable ZTW API via Legacy Mode",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://help.zscaler.com/oneapi/understanding-rate-limiting.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ztw_activation_status":           resourceActivationStatus(),
			"ztw_location_template":           resourceLocationTemplate(),
			"ztw_provisioning_url":            resourceProvisioningURL(),
			"ztw_traffic_forwarding_rule":     resourceTrafficForwardingRule(),
			"ztw_traffic_forwarding_dns_rule": resourceTrafficForwardingDNSRule(),
			"ztw_traffic_forwarding_log_rule": resourceTrafficForwardingLogRule(),
			"ztw_forwarding_gateway":          resourceForwardingGateway(),
			"ztw_dns_forwarding_gateway":      resourceDNSForwardingGateway(),
			"ztw_ip_destination_groups":       resourceIPDestinationGroups(),
			"ztw_ip_source_groups":            resourceIPSourceGroups(),
			"ztw_ip_pool_groups":              resourceIPPoolGroups(),
			"ztw_network_services":            resourceNetworkServices(),
			"ztw_network_service_groups":      resourceNetworkServiceGroups(),
			"ztw_account_groups":              resourceAccountGroup(),
			"ztw_public_cloud_info":           resourcePublicCloudInfo(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ztw_activation_status":           dataSourceActivationStatus(),
			"ztw_location_template":           dataSourceLocationTemplate(),
			"ztw_provisioning_url":            dataSourceProvisioningURL(),
			"ztw_location_management":         dataSourceLocationManagement(),
			"ztw_edge_connector_group":        dataSourceEdgeConnectorGroup(),
			"ztw_traffic_forwarding_rule":     dataSourceTrafficForwardingRule(),
			"ztw_traffic_forwarding_dns_rule": dataSourceTrafficForwardingDNSRule(),
			"ztw_traffic_forwarding_log_rule": dataSourceTrafficForwardingLogRule(),
			"ztw_forwarding_gateway":          dataSourceForwardingGateway(),
			"ztw_dns_forwarding_gateway":      dataSourceDNSForwardingGateway(),
			"ztw_ip_destination_groups":       dataSourceIPDestinationGroups(),
			"ztw_ip_source_groups":            dataSourceIPSourceGroups(),
			"ztw_ip_pool_groups":              dataSourceIPPoolGroups(),
			"ztw_network_services":            dataSourceNetworkServices(),
			"ztw_network_service_groups":      dataSourceNetworkServiceGroups(),
			"ztw_account_groups":              dataSourceAccountGroup(),
			"ztw_public_cloud_info":           dataSourcePublicCloudInfo(),
			"ztw_supported_regions":           dataSourceSupportedRegions(),
		},
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, err := providerConfigure(d, terraformVersion)
		if err != nil {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", err),
					AttributePath: cty.Path{},
				},
			}
		}
		return r, nil
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create configuration from schema
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Load the correct SDK client (prioritizing V3)
	if diags := config.loadClients(); diags.HasError() {
		return nil, diags
	}

	// Return the configured client
	client, err := config.Client()
	if err != nil {
		return nil, diag.Errorf("failed to configure Zscaler client: %v", err)
	}

	// Initialize the global semaphore based on the configured parallelism
	if config.parallelism > 0 {
		apiSemaphore = make(chan struct{}, config.parallelism)
	} else {
		apiSemaphore = make(chan struct{}, 1)
	}

	return client, nil
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
