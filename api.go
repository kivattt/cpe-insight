package main

import (
	"fmt"
	"strings"
)

const CPE_INSIGHT_API_VERSION = "42.4242"

// TODO: Sort the output
func PrintEndPoints() {
	longestKeyLength := 0
	longestValueLength := 0
	for k, v := range apiEndPoints {
		longestKeyLength = max(longestKeyLength, len(k))
		longestValueLength = max(longestValueLength, len(v))
	}

	fmt.Print("\x1b[1;4m") // Bold
	fmt.Print("NAME" + strings.Repeat(" ", longestKeyLength-len("NAME")))
	fmt.Println(" ENDPOINT" + strings.Repeat(" ", longestValueLength-len("ENDPOINT")))
	fmt.Print("\x1b[0m") // Reset

	// TODO: Use $PAGER with a fallback when not on Windows
	for k, v := range apiEndPoints {
		//fmt.Println("\x1b[0;34m" + k + "\x1b[0m " + strings.Repeat(" ", longestKeyLength - len(k)) + v)
		fmt.Println(k + " \x1b[0;34m" + strings.Repeat(" ", longestKeyLength-len(k)) + v + "\x1b[0m")
	}
}

// "${t}" is replaced with the username
var apiEndPoints = map[string]string{
	"alertGet":               "/alert",
	"customerConnectMeta":    "/connect",
	"referenceGet":           "/reference",
	"adminUrlGet":            "/${t}/admin-url",
	"bandwidthGet":           "/${t}/bandwidth",
	"bbtGet":                 "/${t}/bbt",
	"capabilitiesGet":        "/${t}/capabilities",
	"customerGet":            "/${t}/customer",
	"customerLocationMapGet": "/${t}/customer/location-map",
	"ddnsGet":                "/${t}/ddns",
	"ddnsStatus":             "/${t}/ddns/status",
	"defaultConfigGet":       "/${t}/default-config",
	"deviceFamily":           "/${t}/device-family",
	"dhcpGet":                "/${t}/dhcp",
	"dmzGet":                 "/${t}/dmz",
	"dnsGet":                 "/${t}/dns",
	"ethernetPortsGet":       "/${t}/ethernet-ports",
	"findReference":          "/${t}/find",
	"firewallRulesGet":       "/${t}/firewall-rules",
	"firewallStatusGet":      "/${t}/firewall-status",
	"firmwareVersionGet":     "/${t}/firmware-version",
	"firmwareVersionsGet":    "/${t}/firmware-versions",
	//"forgotPasswordUrlCheck": "/${t}/forgot-password/change/${e}/${r}", // TODO: What is ${e} and ${r} ?
	"lanGet":                   "/${t}/lan",
	"lanHosts":                 "/${t}/lan-hosts",
	"lanHostsAliasGet":         "/${t}/lan-hosts/alias",
	"lanIpv6Get":               "/${t}/lan-ipv6",
	"ledGet":                   "/${t}/led",
	"logGet":                   "/${t}/log",
	"mobileGet":                "/${t}/mobile",
	"checkBngL2Tunnel":         "/${t}/msan/check-bng-l2tunnel",
	"checkBngMcvlan":           "/${t}/msan/check-bng-mcvlan",
	"checkMsanL2Tunnel":        "/${t}/msan/check-msan-l2tunnel",
	"checkMsanMcvlan":          "/${t}/msan/check-msan-mcvlan",
	"ponInfo":                  "/${t}/msan/pon-info",
	"msanStatus":               "/${t}/msan/status",
	"verify":                   "/${t}/msan/verify", // Slightly inexplicit key name?
	"natRulesGet":              "/${t}/nat-rules",
	"orderSystemParametersGet": "/${t}/order-system/parameters",
	"orderSystemServicesGet":   "/${t}/order-system/services",
	"pairsGet":                 "/${t}/pair",
	"ping":                     "/${t}/ping",
	"policiesGet":              "/${t}/policies",
	"publicLanGet":             "/${t}/public-lan",
	"schedulesGet":             "/${t}/schedules",
	"services":                 "/${t}/services",
	"sessionsGet":              "/${t}/sessions",
	"sharedSecretAuth":         "/${t}/shared-secret-auth",
	"sikStatusGet":             "/${t}/sik-status",
	"skynetGet":                "/${t}/skynet",
	"smtpOutboundGet":          "/${t}/smtp-outbound",
	"staticDhcpLeasesGet":      "/${t}/static-dhcp-leases",
	"staticRoutesGet":          "/${t}/static-routes",
	"status":                   "/${t}/status",
	"statusFromSessions":       "/${t}/status-from-sessions",
	//"telemetryHistoricGet": "/${t}/telemetry-historic/${e}", // TODO: What is ${e} ?
	"upnpGet":        "/${t}/upnp",
	"upnpMappings":   "/${t}/upnp/mappings",
	"uptime":         "/${t}/uptime",
	"wanProtocolGet": "/${t}/wan-protocol",
	"wlanGet":        "/${t}/wlan",
	"wlanRadioGet":   "/${t}/wlan/radio",
	"wlanScan":       "/${t}/wlan/scan",
	"wlanWpsGet":     "/${t}/wlan/wps",
	"taskHistory":    "/task/history",
}
