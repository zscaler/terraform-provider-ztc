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
			"ztc_cloud": {
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
			"skip_credentials_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Skip credentials validation and API client initialization entirely. " +
					"Intended for configurations where every ztc_* resource and data source is conditionally disabled " +
					"(e.g., count = 0) and no API call will ever be made, such as multi-environment deployments where " +
					"Zscaler is not present in every environment. Any resource or data source that does attempt an API " +
					"call will fail with an explanatory error. Can also be sourced from the " +
					"ZSCALER_SKIP_CREDENTIALS_VALIDATION environment variable.",
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
				Description:      "Maximum number of times a rate-limited or transiently failing request is retried before the operation fails. Default: 30. Maximum: 100. Does not affect request throughput; lowering it only makes runs fail sooner under sustained rate limiting.",
			},
			"parallelism": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "This attribute no longer has any effect and will be removed in a future major release. Remove it from the provider block. API rate limits are handled automatically: the provider honours the Retry-After header returned on a 429 response and retries transparently.",
				Description: "Deprecated and ignored. Previously limited the number of concurrent API requests. " +
					"Rate limiting is now handled automatically and this attribute has no effect.",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout in seconds for a single HTTP request to the Zscaler API. Default: 0, which selects the SDK's built-in 60-second timeout. Accepted values: 0-300. Does not affect request throughput.",
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
			"ztc_dns_gateway":                 resourceDNSGateway(),
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
			"ztc_workload_groups":             dataSourceWorkloadGroup(),
			"ztc_dns_gateway":                 dataSourceDNSGateway(),
		},
	}

	// Guard every resource and data source against the inert client returned
	// when skip_credentials_validation is enabled, so an accidental API call
	// yields a descriptive error instead of a nil-pointer panic.
	for _, r := range p.ResourcesMap {
		guardResourceAgainstInertClient(r)
	}
	for _, ds := range p.DataSourcesMap {
		guardResourceAgainstInertClient(ds)
	}

	p.ConfigureContextFunc = func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		r, diags := providerConfigure(d, terraformVersion)
		if diags.HasError() {
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "failed configuring the provider",
					Detail:        fmt.Sprintf("error:%v", diags),
					AttributePath: cty.Path{},
				},
			}
		}
		// Pass through non-error diagnostics (e.g. the
		// skip_credentials_validation warning).
		return r, diags
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Zscaler client")

	// Create configuration from schema
	config := NewConfig(d)
	config.TerraformVersion = terraformVersion

	// Skip mode: never construct the SDK client. The OneAPI SDK authenticates
	// inside its constructor, so building a client without valid credentials
	// would fail at configure time even when no resource will ever call the
	// API. Return an inert client instead; the CRUD guards installed in
	// ZTCProvider turn any actual API use into a descriptive error.
	if config.skipCredentialsValidation {
		log.Printf("[WARN] skip_credentials_validation is enabled; ZTC API client was not initialized")
		return &Client{skipCredentialsValidation: true}, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "ZTC credentials were not validated",
				Detail: "skip_credentials_validation is enabled, so the ZTC API client was not initialized. " +
					"Any ztc_* resource or data source that attempts an API call will fail. " +
					"This mode is intended for configurations where all ZTC resources are conditionally disabled (e.g., count = 0).",
			},
		}
	}

	// Load the correct SDK client (prioritizing V3)
	if diags := config.loadClients(); diags.HasError() {
		return nil, diags
	}

	// Return the configured client
	client, err := config.Client()
	if err != nil {
		return nil, diag.Errorf("failed to configure Zscaler client: %v", err)
	}

	return client, nil
}

// inertClientDiag is the error returned when a resource or data source is
// evaluated while the provider is in skip_credentials_validation mode.
func inertClientDiag() diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ZTC provider was configured with skip_credentials_validation",
			Detail: "This resource or data source attempted a ZTC API call, but the provider was configured with " +
				"skip_credentials_validation = true (or ZSCALER_SKIP_CREDENTIALS_VALIDATION=true), so no API client exists. " +
				"Either provide valid credentials and remove skip_credentials_validation, or ensure every ztc_* resource " +
				"and data source is disabled (e.g., count = 0) in this configuration.",
		},
	}
}

// isInertClient reports whether the provider meta is the inert client created
// in skip_credentials_validation mode.
func isInertClient(meta interface{}) bool {
	client, ok := meta.(*Client)
	return ok && client.skipCredentialsValidation
}

// guardResourceAgainstInertClient wraps a resource's CRUD and import
// functions so that, in skip_credentials_validation mode, they return a
// descriptive error instead of dereferencing the nil SDK service.
func guardResourceAgainstInertClient(r *schema.Resource) {
	wrap := func(f func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
		if f == nil {
			return nil
		}
		return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			if isInertClient(meta) {
				return inertClientDiag()
			}
			return f(ctx, d, meta)
		}
	}

	r.CreateContext = wrap(r.CreateContext)
	r.ReadContext = wrap(r.ReadContext)
	r.UpdateContext = wrap(r.UpdateContext)
	r.DeleteContext = wrap(r.DeleteContext)

	if r.Importer != nil && r.Importer.StateContext != nil {
		importer := r.Importer.StateContext
		r.Importer.StateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			if isInertClient(meta) {
				return nil, fmt.Errorf("cannot import: the ZTC provider was configured with skip_credentials_validation, so no API client exists")
			}
			return importer(ctx, d, meta)
		}
	}
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
