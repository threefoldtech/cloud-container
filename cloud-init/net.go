package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/vishvananda/netlink"
)

const (
	NetworkFile = "network-config"
)

type Route struct {
	To     string `yaml:"to"`
	Via    net.IP `yaml:"via"`
	Metric int    `yaml:"metric,omitempty"`
}

type Nameservers struct {
	Search    []string `yaml:"search,omitempty"`
	Addresses []string `yaml:"addresses,omitempty"`
}

type Ethernet struct {
	Match struct {
		Mac string `yaml:"macaddress"`
	} `yaml:"match"`
	DHCP4       bool        `yaml:"dhcp4"`
	Addresses   []string    `yaml:"addresses"`
	Gateway4    net.IP      `yaml:"gateway4,omitempty"`
	Gateway6    net.IP      `yaml:"gateway6,omitempty"`
	Routes      []Route     `yaml:"routes,omitempty"`
	Nameservers Nameservers `yaml:"nameservers"`
}

func ApplyNetwork(seed, root string) error {
	var network struct {
		Ethernets map[string]Ethernet `yaml:"ethernets"`
	}

	if err := load(filepath.Join(seed, NetworkFile), &network); err != nil {
		return fmt.Errorf("failed to load network config file: %w", err)
	}

	ns := make(map[string]struct{})

	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("failed to list available nics: %s", err)
	}
	nics := make(map[string]netlink.Link)
	for _, link := range links {
		log("found device with mac: %s", link.Attrs().HardwareAddr.String())
		nics[link.Attrs().HardwareAddr.String()] = link

		if link.Attrs().Name == "lo" {
			if err := netlink.LinkSetUp(link); err != nil {
				log("failed to set lo device up: %s", err)
			}
		}
	}

	for _, eth := range network.Ethernets {
		mac := eth.Match.Mac
		link, ok := nics[mac]
		if !ok {
			log("no nic found with mac: %s", mac)
			continue
		}
		log("setting up (%s)", mac)
		if err := netlink.LinkSetUp(link); err != nil {
			log("failed to set device '%s' up: %s", mac, err.Error())
			continue
		}

		for _, address := range eth.Addresses {
			ip, err := netlink.ParseAddr(address)
			if err != nil {
				log("failed to parse address '%s': %v", address, err)
				continue
			}

			if err := netlink.AddrAdd(link, ip); err != nil {
				log("failed to assign ip address '%s' to intreface '%s': %s", address, mac, err.Error())
			}
		}

		for _, route := range eth.Routes {
			_, to, err := net.ParseCIDR(route.To)
			if err != nil {
				log("failed to parse cidr %s: %s", route.To, err.Error())
				continue
			}

			if err := netlink.RouteAdd(&netlink.Route{
				Dst:       to,
				Gw:        route.Via,
				LinkIndex: link.Attrs().Index,
			}); err != nil {
				log("failed to set route %s via %s on interface %s: %s", route.To, route.Via.String(), mac, err.Error())
			}
		}

		if len(eth.Gateway4) != 0 {
			if err := netlink.RouteAdd(&netlink.Route{
				Dst: &net.IPNet{
					IP:   net.ParseIP("0.0.0.0"),
					Mask: net.CIDRMask(0, 8*net.IPv4len),
				},
				Gw:        eth.Gateway4,
				LinkIndex: link.Attrs().Index,
			}); err != nil {
				log("failed to set route default(4) via %s on interface %s: %s", eth.Gateway4.String(), mac, err.Error())
			}
		}

		if len(eth.Gateway6) != 0 {
			if err := netlink.RouteAdd(&netlink.Route{
				Dst: &net.IPNet{
					IP:   net.ParseIP("::"),
					Mask: net.CIDRMask(0, 8*net.IPv6len),
				},
				Gw:        eth.Gateway6,
				LinkIndex: link.Attrs().Index,
			}); err != nil {
				log("failed to set route default(6) via %s on interface %s: %s", eth.Gateway4.String(), mac, err.Error())
			}
		}

		for _, nameserver := range eth.Nameservers.Addresses {
			ns[nameserver] = struct{}{}
		}
	}

	// write down resolv.conf file
	path := filepath.Join(root, "etc")
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create /etc: %w", err)
	}

	res, err := os.Create(filepath.Join(path, "resolv.conf"))
	if err != nil {
		return fmt.Errorf("failed to create /etc/resolv.conf: %w", err)
	}

	defer res.Close()
	for server := range ns {
		if _, err := res.WriteString(fmt.Sprintf("nameserver %s\n", server)); err != nil {
			return fmt.Errorf("failed to write resolv.conf file: %w", err)
		}
	}

	return nil
}
