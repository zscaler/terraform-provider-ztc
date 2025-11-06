---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZTC Config Activation"
subcategory: "Activation"
---

# ZIA Activator Configuration

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/activation"
)

func getEnvVarOrFail(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	log.Fatalf("[ERROR] Couldn't find environment variable %s\n", k)
	return ""
}

func main() {
	log.Printf("[INFO] Initializing ZTW activation client")

	useLegacy := strings.ToLower(os.Getenv("ZSCALER_USE_LEGACY_CLIENT")) == "true"

	var (
		service *zscaler.Service
		err     error
	)

	if useLegacy {
		log.Printf("[INFO] Using Legacy Client mode")

		username := getEnvVarOrFail("ZTW_USERNAME")
		password := getEnvVarOrFail("ZTW_PASSWORD")
		apiKey := getEnvVarOrFail("ZTW_API_KEY")
		cloud := getEnvVarOrFail("ZTW_CLOUD")

		ztwCfg, err := ztw.NewConfiguration(
			ztw.WithZtwUsername(username),
			ztw.WithZtwPassword(password),
			ztw.WithZtwAPIKey(apiKey),
			ztw.WithZtwCloud(cloud),
			ztw.WithUserAgentExtra(fmt.Sprintf("(%s %s) cli/ztwActivator", runtime.GOOS, runtime.GOARCH)),
		)
		if err != nil {
			log.Fatalf("Error creating ZTC configuration: %v", err)
		}

		service, err = zscaler.NewLegacyZtwClient(ztwCfg)
		if err != nil {
			log.Fatalf("Error creating ZTC legacy client: %v", err)
		}
	} else {
		log.Printf("[INFO] Using OneAPI Client mode")

		clientID := getEnvVarOrFail("ZSCALER_CLIENT_ID")
		clientSecret := getEnvVarOrFail("ZSCALER_CLIENT_SECRET")
		vanityDomain := getEnvVarOrFail("ZSCALER_VANITY_DOMAIN")
		cloud := getEnvVarOrFail("ZSCALER_CLOUD")

		cfg, err := zscaler.NewConfiguration(
			zscaler.WithClientID(clientID),
			zscaler.WithClientSecret(clientSecret),
			zscaler.WithVanityDomain(vanityDomain),
			zscaler.WithZscalerCloud(cloud),
			zscaler.WithUserAgentExtra(fmt.Sprintf("(%s %s) cli/ztcActivator", runtime.GOOS, runtime.GOARCH)),
		)
		if err != nil {
			log.Fatalf("[ERROR] Failed to build OneAPI configuration: %v", err)
		}

		service, err = zscaler.NewOneAPIClient(cfg)
		if err != nil {
			log.Fatalf("[ERROR] Failed to initialize OneAPI client: %v", err)
		}
	}

	ctx := context.Background()

	resp, err := activation.UpdateActivationStatus(ctx, service, activation.ECAdminActivation{
		OrgEditStatus:         "org_edit_status",
		OrgLastActivateStatus: "org_last_activate_status",
		// AdminStatusMap:        "admin_status_map",
		AdminActivateStatus: "admin_activate_status",
	})
	if err != nil {
		log.Fatalf("[ERROR] Activation Failed: %v", err)
	}

	log.Printf("[INFO] Activation succeeded: %#v\n", resp)

	// Perform logout if using Legacy Client
	if useLegacy && service.LegacyClient != nil && service.LegacyClient.ZtwClient != nil {
		log.Printf("[INFO] Destroying session...\n")
		if err := service.LegacyClient.ZtwClient.Logout(ctx); err != nil {
			log.Printf("[WARN] Logout failed: %v\n", err)
		}
	}
}
```