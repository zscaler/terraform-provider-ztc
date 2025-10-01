package variable

// Firewall Filtering IP Destination Group resource/datasource
const (
	IPDSTGroupName         = "this is an acceptance test"
	IPDSTGroupDescription  = "this is an acceptance test"
	IPDSTGroupTypeDSTNFQDN = "DSTN_FQDN"
)

// Firewall Filtering IP Source Group resource/datasource
const (
	IPSRCGroupName        = "this is an acceptance test"
	IPSRCGroupDescription = "this is an acceptance test"
)

// Firewall network services groups resource/datasource
const (
	NetworkServicesGroupName        = "this is an acceptance test"
	NetworkServicesGroupDescription = "this is an acceptance test"
	// IPNetworkServices = "this is an acceptance test"
)

// Firewall network services resource/datasource
const (
	NetworkServicesName        = "this is an acceptance test"
	NetworkServicesDescription = "this is an acceptance test"
	NetworkServicesType        = "CUSTOM"
)

const (
	IPPoolGroupName        = "this is an acceptance test"
	IPPoolGroupDescription = "this is an acceptance test"
)

const (
	ForwardGWDescription   = "this is an acceptance test"
	ForwardGWFailClose     = true
	ForwardGWPrimaryType   = "AUTO"
	ForwardGWSecondaryType = "AUTO"
	ForwardGWType          = "ZIA"
)

const (
	ForwardControlRuleDescription = "this is an acceptance test"
	ForwardControlRuleType        = "EC_RDR"
	ForwardControlMethod          = "DIRECT"
	ForwardControlState           = "ENABLED"
)
