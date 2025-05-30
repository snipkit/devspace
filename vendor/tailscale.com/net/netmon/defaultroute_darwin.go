// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

//go:build darwin || ios

package netmon

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"tailscale.com/syncs"
)

const (
	KHULNASOFT_ADMIN_HOST = "admin.khulnasoft.com"
)

var (
	lastKnownDefaultRouteIfName syncs.AtomicValue[string]
)

// UpdateLastKnownDefaultRouteInterface is called by ipn-go-bridge in the iOS app when
// our NWPathMonitor instance detects a network path transition.
func UpdateLastKnownDefaultRouteInterface(ifName string) {
	if ifName == "" {
		return
	}
	if old := lastKnownDefaultRouteIfName.Swap(ifName); old != ifName {
		log.Printf("defaultroute_darwin: update from Swift, ifName = %s (was %s)", ifName, old)
	}
}

// resolveHostname resolves net.IP of given hostname
func resolveHostname(hostname string) (net.IP, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}

	// Prefer IPv4 addresses
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4, nil
		}
	}

	return nil, fmt.Errorf("no IPv4 address found for %s", hostname)
}

// interfaceTo returns string name of network interface used to reach target IP
func interfaceTo(ip net.IP) (string, error) {
	cmd := exec.Command("route", "get", ip.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "interface:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("interface to %s not found", ip)
}

// getInterfaceByRoute gets default interface by checking route to specific host
func getInterfaceByRoute(host string) (*net.Interface, error) {
	ip, err := resolveHostname(host)
	if err != nil {
		return nil, err
	}

	ifaceName, err := interfaceTo(ip)
	if err != nil {
		return nil, err
	}

	intf, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}

	return intf, nil
}

func defaultRoute() (d DefaultRouteDetails, err error) {
	// We cannot rely on the delegated interface data on darwin. The NetworkExtension framework
	// seems to set the delegate interface only once, upon the *creation* of the VPN tunnel.
	// If a network transition (e.g. from Wi-Fi to Cellular) happens while the tunnel is
	// connected, it will be ignored and we will still try to set Wi-Fi as the default route
	// because the delegated interface is not updated by the NetworkExtension framework.
	//
	// We work around this on the Swift side with a NWPathMonitor instance that observes
	// the interface name of the first currently satisfied network path. Our Swift code will
	// call into `UpdateLastKnownDefaultRouteInterface`, so we can rely on that when it is set.
	//
	// If for any reason the Swift machinery didn't work and we don't get any updates, we will
	// fallback to the BSD logic.

	// Start by getting all available interfaces.
	interfaces, err := netInterfaces()
	if err != nil {
		log.Printf("defaultroute_darwin: could not get interfaces: %v", err)
		return d, ErrNoGatewayIndexFound
	}

	getInterfaceByName := func(name string) *Interface {
		for _, ifc := range interfaces {
			if ifc.Name != name {
				continue
			}

			if !ifc.IsUp() {
				log.Printf("defaultroute_darwin: %s is down", name)
				return nil
			}

			addrs, _ := ifc.Addrs()
			if len(addrs) == 0 {
				log.Printf("defaultroute_darwin: %s has no addresses", name)
				return nil
			}
			return &ifc
		}
		return nil
	}

	// Try to get interface by running route against Khulnasoft admin
	iface, err := getInterfaceByRoute(KHULNASOFT_ADMIN_HOST)
	if err == nil {
		d.InterfaceIndex = iface.Index
		d.InterfaceName = iface.Name
		return d, nil
	}

	// If that fails (like in air-gapped environments) fallback to default TS logic

	// Did Swift set lastKnownDefaultRouteInterface? If so, we should use it and don't bother
	// with anything else. However, for sanity, do check whether Swift gave us with an interface
	// that exists, is up, and has an address.
	if swiftIfName := lastKnownDefaultRouteIfName.Load(); swiftIfName != "" {
		ifc := getInterfaceByName(swiftIfName)
		if ifc != nil {
			d.InterfaceName = ifc.Name
			d.InterfaceIndex = ifc.Index
			return d, nil
		}
	}

	// Fallback to the BSD logic
	idx, err := DefaultRouteInterfaceIndex()
	if err != nil {
		return d, err
	}
	iface, err = net.InterfaceByIndex(idx)
	if err != nil {
		return d, err
	}
	d.InterfaceName = iface.Name
	d.InterfaceIndex = idx
	return d, nil
}
