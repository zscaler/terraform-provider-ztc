package ztc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-ztc/ztc/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/provisioning/provisioning_url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

var (
	sweeperLogger   hclog.Logger
	sweeperLogLevel hclog.Level
)

func init() {
	sweeperLogLevel = hclog.Warn
	if os.Getenv("TF_LOG") != "" {
		sweeperLogLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	sweeperLogger = hclog.New(&hclog.LoggerOptions{
		Level:      sweeperLogLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

func logSweptResource(kind, id, nameOrLabel string) {
	sweeperLogger.Warn(fmt.Sprintf("sweeper found dangling %q %q %q", kind, id, nameOrLabel))
}

type testClient struct {
	sdkV3Client *zscaler.Client
}

var (
	testResourcePrefix   = "tf-acc-test-"
	updateResourcePrefix = "tf-updated-"
)

func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("ZTW_VCR_TF_ACC") != "" {
		t.Skip("forced sweeper is live and will never be run within VCR")
		return
	}
	if os.Getenv("ZTC_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "ZTC_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	provider := ZTCProvider()
	c := terraform.NewResourceConfigRaw(nil)
	diag := provider.Configure(context.Background(), c)
	if diag.HasError() {
		t.Skipf("sweeper's provider configuration failed: %v", diag)
		return
	}

	sdkClient, err := sdkV3ClientForTest()
	if err != nil {
		t.Fatalf("Failed to get SDK client: %s", err)
	}

	testClient := &testClient{
		sdkV3Client: sdkClient,
	}

	// sweepTestSourceIPGroup(testClient)
	// sweepTestDestinationIPGroup(testClient)
	sweepTestNetworkServices(testClient)
	sweepTestNetworkServicesGroup(testClient)
	sweepTestIPPoolGroup(testClient)
	sweepTestTrafficForwardingRule(testClient)
	sweepTestZIAForwardingGateway(testClient)
	sweepTestDNSForwardingGateway(testClient)
	sweepTestLocationTemplate(testClient)
	sweepTestProvisioningURL(testClient)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			// Retrieve the client and handle the error
			sdkClient, err := sdkV3ClientForTest()
			if err != nil {
				return fmt.Errorf("failed to initialize SDK V3 client for sweeper: %w", err)
			}

			// Pass the client to the deleter function
			return del(&testClient{sdkV3Client: sdkClient})
		},
	})
}

/*
	func sweepTestSourceIPGroup(client *testClient) error {
		var errorList []error

		service := &zscaler.Service{
			Client: client.sdkV3Client,
		}

		ipSourceGroup, err := ipsourcegroups.GetAll(context.Background(), service)
		if err != nil {
			return err
		}
		// Logging the number of identified resources before the deletion loop
		sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipSourceGroup)))
		for _, b := range ipSourceGroup {
			// Check if the resource name has the required prefix before deleting it
			if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
				if _, err := ipsourcegroups.Delete(context.Background(), service, b.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.IPSourceGroup, fmt.Sprintf("%d", b.ID), b.Name)
			}
		}
		// Log errors encountered during the deletion process
		if len(errorList) > 0 {
			for _, err := range errorList {
				sweeperLogger.Error(err.Error())
			}
		}
		return condenseError(errorList)
	}

	func sweepTestDestinationIPGroup(client *testClient) error {
		var errorList []error

		service := &zscaler.Service{
			Client: client.sdkV3Client,
		}

		ipDestGroup, err := ipdestinationgroups.GetAll(context.Background(), service)
		if err != nil {
			return err
		}
		// Logging the number of identified resources before the deletion loop
		sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipDestGroup)))
		for _, b := range ipDestGroup {
			// Check if the resource name has the required prefix before deleting it
			if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
				if _, err := ipdestinationgroups.Delete(context.Background(), service, b.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.IPDestinationGroup, fmt.Sprintf("%d", b.ID), b.Name)
			}
		}
		// Log errors encountered during the deletion process
		if len(errorList) > 0 {
			for _, err := range errorList {
				sweeperLogger.Error(err.Error())
			}
		}
		return condenseError(errorList)
	}
*/
func sweepTestIPPoolGroup(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	ipDestGroup, err := ipgroups.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipDestGroup)))
	for _, b := range ipDestGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := ipgroups.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.IPPoolGroup, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestNetworkServices(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	services, err := networkservices.GetAllNetworkServices(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(services)))
	for _, b := range services {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservices.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.NetworkServices, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestNetworkServicesGroup(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	groups, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservicegroups.DeleteNetworkServiceGroups(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.NetworkServiceGroups, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestTrafficForwardingRule(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	rule, err := forwarding_rules.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := forwarding_rules.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.TrafficForwardingRule, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestZIAForwardingGateway(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	rule, err := zia_forwarding_gateway.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := zia_forwarding_gateway.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ZIAForwardingGateway, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestDNSForwardingGateway(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	rule, err := dns_forwarding_gateway.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dns_forwarding_gateway.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.DNSForwardingGateway, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestLocationTemplate(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	rule, err := locationtemplate.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := locationtemplate.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.LocationTemplate, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestProvisioningURL(client *testClient) error {
	var errorList []error

	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	rule, err := provisioning_url.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := provisioning_url.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ProvisioningURL, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}
