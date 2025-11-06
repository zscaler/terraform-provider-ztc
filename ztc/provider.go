package ztc

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

func ZTCProvider() *schema.Provider {
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
			"ztc_activation_status":           resourceActivationStatus(),
			"ztc_location_template":           resourceLocationTemplate(),
			"ztc_provisioning_url":            resourceProvisioningURL(),
			"ztc_traffic_forwarding_rule":     resourceTrafficForwardingRule(),
			"ztc_traffic_forwarding_dns_rule": resourceTrafficForwardingDNSRule(),
			"ztc_traffic_forwarding_log_rule": resourceTrafficForwardingLogRule(),
			"ztc_forwarding_gateway":          resourceForwardingGateway(),
			"ztc_dns_forwarding_gateway":      resourceDNSForwardingGateway(),
			"ztc_ip_destination_groups":       resourceIPDestinationGroups(),
			"ztc_ip_source_groups":            resourceIPSourceGroups(),
			"ztc_ip_pool_groups":              resourceIPPoolGroups(),
			"ztc_network_services":            resourceNetworkServices(),
			"ztc_network_service_groups":      resourceNetworkServiceGroups(),
			"ztc_account_groups":              resourceAccountGroup(),
			"ztc_public_cloud_info":           resourcePublicCloudInfo(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ztc_activation_status":           dataSourceActivationStatus(),
			"ztc_location_template":           dataSourceLocationTemplate(),
			"ztc_provisioning_url":            dataSourceProvisioningURL(),
			"ztc_location_management":         dataSourceLocationManagement(),
			"ztc_edge_connector_group":        dataSourceEdgeConnectorGroup(),
			"ztc_traffic_forwarding_rule":     dataSourceTrafficForwardingRule(),
			"ztc_traffic_forwarding_dns_rule": dataSourceTrafficForwardingDNSRule(),
			"ztc_traffic_forwarding_log_rule": dataSourceTrafficForwardingLogRule(),
			"ztc_forwarding_gateway":          dataSourceForwardingGateway(),
			"ztc_dns_forwarding_gateway":      dataSourceDNSForwardingGateway(),
			"ztc_ip_destination_groups":       dataSourceIPDestinationGroups(),
			"ztc_ip_source_groups":            dataSourceIPSourceGroups(),
			"ztc_ip_pool_groups":              dataSourceIPPoolGroups(),
			"ztc_network_services":            dataSourceNetworkServices(),
			"ztc_network_service_groups":      dataSourceNetworkServiceGroups(),
			"ztc_account_groups":              dataSourceAccountGroup(),
			"ztc_public_cloud_info":           dataSourcePublicCloudInfo(),
			"ztc_supported_regions":           dataSourceSupportedRegions(),
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
