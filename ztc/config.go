package ztc

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-ztc/ztc/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
)

type (
	// Config contains our provider schema values and Zscaler clients.
	Config struct {
		clientID           string
		clientSecret       string
		vanityDomain       string
		cloud              string
		privateKey         string
		httpProxy          string
		retryCount         int
		parallelism        int
		backoff            bool
		minWait            int
		maxWait            int
		logLevel           int
		requestTimeout     int
		useLegacyClient    bool
		zscalerSDKClientV3 *zscaler.Client
		logger             hclog.Logger
		TerraformVersion   string // New field for Terraform version
		ProviderVersion    string // New field for Provider version

		// Options for Legacy V2 SDK
		Username   string
		Password   string
		APIKey     string
		ZTCBaseURL string
		UserAgent  string
	}
)

type Client struct {
	Service *zscaler.Service
}

func NewConfig(d *schema.ResourceData) *Config {
	// defaults
	config := Config{
		backoff:        true,
		minWait:        30,
		maxWait:        300,
		retryCount:     30,
		parallelism:    1,
		logLevel:       int(hclog.Error),
		requestTimeout: 0,
	}
	logLevel := hclog.Level(config.logLevel)
	if os.Getenv("TF_LOG") != "" {
		logLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	config.logger = hclog.New(&hclog.LoggerOptions{
		Level:      logLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})

	if val, ok := d.GetOk("use_legacy_client"); ok {
		config.useLegacyClient = val.(bool)
	} else if os.Getenv("ZSCALER_USE_LEGACY_CLIENT") != "" {
		config.useLegacyClient = strings.ToLower(os.Getenv("ZSCALER_USE_LEGACY_CLIENT")) == "true"
	}

	if val, ok := d.GetOk("client_id"); ok {
		config.clientID = val.(string)
	}
	if config.clientID == "" && os.Getenv("ZSCALER_CLIENT_ID") != "" {
		config.clientID = os.Getenv("ZSCALER_CLIENT_ID")
	}

	if val, ok := d.GetOk("client_secret"); ok {
		config.clientSecret = val.(string)
	}
	if config.clientSecret == "" && os.Getenv("ZSCALER_CLIENT_SECRET") != "" {
		config.clientSecret = os.Getenv("ZSCALER_CLIENT_SECRET")
	}

	if val, ok := d.GetOk("private_key"); ok {
		config.privateKey = val.(string)
	}
	if config.privateKey == "" && os.Getenv("ZSCALER_PRIVATE_KEY") != "" {
		config.privateKey = os.Getenv("ZSCALER_PRIVATE_KEY")
	}

	if val, ok := d.GetOk("vanity_domain"); ok {
		config.vanityDomain = val.(string)
	}
	if config.vanityDomain == "" && os.Getenv("ZSCALER_VANITY_DOMAIN") != "" {
		config.vanityDomain = os.Getenv("ZSCALER_VANITY_DOMAIN")
	}

	if val, ok := d.GetOk("zscaler_cloud"); ok {
		config.cloud = val.(string)
	}
	if config.cloud == "" && os.Getenv("ZSCALER_CLOUD") != "" {
		config.cloud = os.Getenv("ZSCALER_CLOUD")
	}

	if val, ok := d.GetOk("username"); ok {
		config.Username = val.(string)
	}
	if config.Username == "" {
		config.Username = os.Getenv("ZTC_USERNAME")
	}

	if val, ok := d.GetOk("password"); ok {
		config.Password = val.(string)
	}
	if config.Password == "" {
		config.Password = os.Getenv("ZTC_PASSWORD")
	}

	if val, ok := d.GetOk("api_key"); ok {
		config.APIKey = val.(string)
	}
	if config.APIKey == "" {
		config.APIKey = os.Getenv("ZTC_API_KEY")
	}

	if val, ok := d.GetOk("ZTC_cloud"); ok {
		config.ZTCBaseURL = val.(string)
	}
	if config.ZTCBaseURL == "" {
		config.ZTCBaseURL = os.Getenv("ZTC_CLOUD")
	}

	if val, ok := d.GetOk("max_retries"); ok {
		config.retryCount = val.(int)
	}

	if val, ok := d.GetOk("parallelism"); ok {
		config.parallelism = val.(int)
	}

	if val, ok := d.GetOk("backoff"); ok {
		config.backoff = val.(bool)
	}

	if val, ok := d.GetOk("min_wait_seconds"); ok {
		config.minWait = val.(int)
	}

	if val, ok := d.GetOk("max_wait_seconds"); ok {
		config.maxWait = val.(int)
	}

	if val, ok := d.GetOk("log_level"); ok {
		config.logLevel = val.(int)
	}

	if val, ok := d.GetOk("request_timeout"); ok {
		config.requestTimeout = val.(int)
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.httpProxy = httpProxy
	}
	if config.httpProxy == "" && os.Getenv("ZSCALER_HTTP_PROXY") != "" {
		config.httpProxy = os.Getenv("ZSCALER_HTTP_PROXY")
	}

	return &config
}

// loadClients initializes SDK clients based on configuration
func (c *Config) loadClients() diag.Diagnostics {
	if c.useLegacyClient {
		log.Println("[INFO] Initializing ZTC V2 (Legacy) client")
		v2Client, err := zscalerSDKV2Client(c)
		if err != nil {
			return diag.Errorf("failed to initialize SDK V2 client: %v", err)
		}
		c.zscalerSDKClientV3 = v2Client.Client
		return nil
	}

	log.Println("[INFO] Initializing ZTC V3 client")
	v3Client, err := zscalerSDKV3Client(c)
	if err != nil {
		return diag.Errorf("failed to initialize SDK V3 client: %v", err)
	}
	c.zscalerSDKClientV3 = v3Client

	return nil
}

// SelectClient returns the appropriate client based on authentication type or other factors.
func (c *Config) SelectClient() (*zscaler.Client, *ztw.Client, error) {
	if !c.useLegacyClient && c.zscalerSDKClientV3 != nil {
		return c.zscalerSDKClientV3, nil, nil
	}

	return nil, nil, fmt.Errorf("no valid client configuration provided")
}

// generateUserAgent constructs the user agent string with all required details
func generateUserAgent(terraformVersion string) string {
	// Fetch the provider version dynamically from common.Version()
	providerVersion := common.Version()

	return fmt.Sprintf("(%s %s) Terraform/%s Provider/%s",
		runtime.GOOS,
		runtime.GOARCH,
		terraformVersion,
		providerVersion,
	)
}

func zscalerSDKV2Client(c *Config) (*zscaler.Service, error) {
	customUserAgent := generateUserAgent(c.TerraformVersion)

	// Start with base configuration setters
	setters := []ztw.ConfigSetter{
		ztw.WithCache(false),
		ztw.WithHttpClientPtr(http.DefaultClient),
		ztw.WithRateLimitMaxRetries(int32(c.retryCount)),
		ztw.WithRequestTimeout(time.Duration(c.requestTimeout) * time.Second),
		ztw.WithUserAgentExtra(customUserAgent), // Set the custom user agent
	}

	// Apply credentials and mandatory parameters
	setters = append(
		setters,
		ztw.WithZtwUsername(c.Username),
		ztw.WithZtwPassword(c.Password),
		ztw.WithZtwAPIKey(c.APIKey),
		ztw.WithZtwCloud(c.ZTCBaseURL),
	)

	// Configure HTTP proxy if provided
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}
		setters = append(setters, ztw.WithProxyHost(_url.Hostname()))

		// Default to port 80 if not provided
		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		// Parse the port as a 32-bit integer
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}

		// Optionally, you can also check the port range if needed
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		// Safe cast to int32
		port32 := int32(port64)
		setters = append(setters, ztw.WithProxyPort(port32))
	}

	// Initialize ZTC configuration
	ztwCfg, err := ztw.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZTC configuration: %v", err)
	}
	ztwCfg.UserAgent = customUserAgent
	// Initialize ZTC client
	wrappedV2Client, err := zscaler.NewLegacyZtwClient(ztwCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZTC client: %v", err)
	}

	log.Println("[INFO] Successfully initialized ZTC V2 client")
	return wrappedV2Client, nil
}

func zscalerSDKV3Client(c *Config) (*zscaler.Client, error) {
	customUserAgent := generateUserAgent(c.TerraformVersion)

	// Start with base configuration setters
	setters := []zscaler.ConfigSetter{
		zscaler.WithCache(false),
		zscaler.WithHttpClientPtr(http.DefaultClient),
		zscaler.WithRateLimitMaxRetries(int32(c.retryCount)),
		zscaler.WithRequestTimeout(time.Duration(c.requestTimeout) * time.Second),
		zscaler.WithUserAgentExtra(customUserAgent),
	}

	// Configure HTTP proxy if provided
	if c.httpProxy != "" {
		_url, err := url.Parse(c.httpProxy)
		if err != nil {
			return nil, err
		}
		setters = append(setters, zscaler.WithProxyHost(_url.Hostname()))

		sPort := _url.Port()
		if sPort == "" {
			sPort = "80"
		}
		// Parse the port as a 32-bit integer
		port64, err := strconv.ParseInt(sPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}

		// Optionally, you can also check the port range if needed
		if port64 < 1 || port64 > 65535 {
			return nil, fmt.Errorf("invalid port number: must be between 1 and 65535, got: %d", port64)
		}
		// Safe cast to int32
		port32 := int32(port64)
		setters = append(setters, zscaler.WithProxyPort(port32))
	}

	// Main switch for OAuth2 authentication
	switch {
	case c.clientID != "" && c.clientSecret != "" && c.vanityDomain != "":
		setters = append(setters,
			zscaler.WithClientID(c.clientID),
			zscaler.WithClientSecret(c.clientSecret),
			zscaler.WithVanityDomain(c.vanityDomain),
		)

		if c.cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(c.cloud))
		}

	case c.clientID != "" && c.privateKey != "" && c.vanityDomain != "":
		setters = append(setters,
			zscaler.WithClientID(c.clientID),
			zscaler.WithPrivateKey(c.privateKey),
			zscaler.WithVanityDomain(c.vanityDomain),
		)

		if c.cloud != "" {
			setters = append(setters, zscaler.WithZscalerCloud(c.cloud))
		}

	default:
		return nil, fmt.Errorf("invalid authentication configuration: missing required parameters")
	}

	// Create configuration for OAuth2 authentication
	config, err := zscaler.NewConfiguration(setters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK V3 configuration: %v", err)
	}

	config.UserAgent = customUserAgent

	// Initialize the client with the configuration
	v3Client, err := zscaler.NewOneAPIClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zscaler API client: %v", err)
	}

	return v3Client.Client, nil
}

func (c *Config) Client() (*Client, error) {
	// Legacy client logic
	if c.useLegacyClient {
		wrappedV2Client, err := zscalerSDKV2Client(c)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize legacy v2 client: %w", err)
		}
		return &Client{
			Service: zscaler.NewService(wrappedV2Client.Client, nil),
		}, nil
	}

	// Fallback to V3 client logic
	v3Client, err := zscalerSDKV3Client(c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize v3 client: %w", err)
	}
	return &Client{
		Service: zscaler.NewService(v3Client, nil),
	}, nil
}
