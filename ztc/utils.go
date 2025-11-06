package ztc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
)

func intPtr(n int) *int {
	return &n
}

func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func SetToStringList(d *schema.ResourceData, key string) []string {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []string{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []string{}
	}
	return SetToStringSlice(set)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

func getIntFromResourceData(d *schema.ResourceData, key string) (int, bool) {
	obj, isSet := d.GetOk(key)
	val, isInt := obj.(int)
	return val, isSet && isInt && val > 0
}

func SetToIntList(d *schema.ResourceData, key string) []int {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []int{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []int{}
	}

	intList := make([]int, set.Len())
	for i, v := range set.List() {
		intList[i] = v.(int)
	}
	return intList
}

var failFastErrorCodes = []string{
	"INVALID_INPUT_ARGUMENT",
	"TRIAL_EXPIRED",
	"EDIT_LOCK_NOT_AVAILABLE",
	"DUPLICATE_ITEM",
	// Add more codes here as needed
}

// failFastOnErrorCodes detects known fatal API error codes and returns the original error to fail immediately.
func failFastOnErrorCodes(err error) error {
	if err == nil {
		return nil
	}

	// Case 1: SDK's structured ErrorResponse (preferred path)
	var apiErr *errorx.ErrorResponse
	if errors.As(err, &apiErr) {
		code := extractErrorCodeFromBody(apiErr.Message)
		for _, c := range failFastErrorCodes {
			if code == c {
				log.Printf("[ERROR] Failing immediately due to API error code '%s': %s", c, apiErr.Message)
				return err
			}
		}
	}

	// Case 2: fallback for unstructured errors
	errMsg := err.Error()
	for _, code := range failFastErrorCodes {
		match := fmt.Sprintf(`"code":"%s"`, code)
		if strings.Contains(errMsg, match) {
			log.Printf("[WARN] Failing due to fallback match for code '%s': %s", code, errMsg)
			return err
		}
	}

	return nil
}

func extractErrorCodeFromBody(body string) string {
	type apiErrorBody struct {
		Code string `json:"code"`
	}
	var parsed apiErrorBody
	if err := json.Unmarshal([]byte(body), &parsed); err == nil {
		return parsed.Code
	}
	return ""
}

func processCountries(countries []string) []string {
	processedCountries := make([]string, len(countries))
	for i, country := range countries {
		if country != "ANY" && country != "NONE" && len(country) == 2 { // Assuming the 2 letter code is an ISO Alpha-2 Code
			processedCountries[i] = "COUNTRY_" + country
		} else {
			processedCountries[i] = country
		}
	}
	return processedCountries
}

func DetachRuleIDNameExtensions(ctx context.Context, client *Client, id int, resource string, getResources func(*forwarding_rules.ForwardingRules) []common.IDNameExtensions, setResources func(*forwarding_rules.ForwardingRules, []common.IDNameExtensions)) error {
	service := client.Service

	log.Printf("[INFO] Detaching filtering rule from %s: %d\n", resource, id)
	rules, err := forwarding_rules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[error] Error while getting filtering rule")
		return err
	}

	for _, rule := range rules {
		ids := []common.IDNameExtensions{}
		shouldUpdate := false
		for _, destGroup := range getResources(&rule) {
			if destGroup.ID != id {
				ids = append(ids, destGroup)
			} else {
				shouldUpdate = true
			}
		}
		if shouldUpdate {
			setResources(&rule, ids)
			time.Sleep(time.Second * 5)
			_, err = forwarding_rules.Get(ctx, service, rule.ID)
			if err == nil {
				_, err = forwarding_rules.Update(ctx, service, rule.ID, &rule)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
