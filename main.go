package main

import (
	"fmt"
	"net"

	"k8s.io/kubernetes/pkg/controller/nodeipam/ipam/cidrset"
)

var cidrs = []*net.IPNet{}

func main() {
	_, nodeCidr, err := net.ParseCIDR("10.44.0.0/16")
	if err != nil {
		panic(err)
	}
	fmt.Println("Node CIDR:", nodeCidr)

	// 4096 possible networks in a /16 when slicing into /28's
	cidrSet, err := cidrset.NewCIDRSet(nodeCidr, 28)
	if err != nil {
		panic(err)
	}

	// occupy an example cidr
	{
		_, cidr, err := net.ParseCIDR("10.44.0.16/28")
		if err != nil {
			panic(err)
		}

		if err := cidrSet.Occupy(cidr); err != nil {
			panic(err)
		}

		cidrs = append(cidrs, cidr)
		fmt.Println("Example CIDR marked already Occupied:", cidr)
	}

	// allocate a cidr
	{
		cidr, err := cidrSet.AllocateNext()
		if err != nil {
			panic(err)
		}

		cidrs = append(cidrs, cidr)
		fmt.Println("CIDR Allocated:", cidr)
	}

	// allocate another cidr
	{
		cidr, err := cidrSet.AllocateNext()
		if err != nil {
			panic(err)
		}

		cidrs = append(cidrs, cidr)
		fmt.Println("CIDR Allocated:", cidr)
	}

	// release all cidrs
	for _, cidr := range cidrs {
		if err := cidrSet.Release(cidr); err != nil {
			panic(err)
		}
		fmt.Println("Released CIDR:", cidr)

		cidrs = []*net.IPNet{}
	}

	// allocate another cidr
	{
		cidr, err := cidrSet.AllocateNext()
		if err != nil {
			panic(err)
		}

		fmt.Println("CIDR Allocated:", cidr)
	}
}
