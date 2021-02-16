/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBigipLtmNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipLtmNodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the node.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User defined description of the node.",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address of the node of the node.",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node connection limit.",
			},
			"dynamic_ratio": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dynamic ratio number for the node.",
			},
			"monitor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the health monitors the system currently uses to monitor this node.",
			},
			"rate_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Node rate limit.",
			},
			"ratio": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Node ratio weight.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the node.",
			},
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The FQDN node's address family.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The fully qualified domain name of the node.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The amount of time before sending the next DNS query.",
						},
						"downinterval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The number of attempts to resolve a domain name.",
						},
						"autopopulate": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies if the node should scale to the IP address set returned by DNS.",
						},
					},
				},
			},
		},
	}
}
func dataSourceBigipLtmNodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	d.SetId("")
	name := fmt.Sprintf("/%s/%s", d.Get("partition").(string), d.Get("name").(string))
	log.Println("[DEBUG] Reading Node : " + name)
	node, err := client.GetNode(name)
	if err != nil {
		return fmt.Errorf("Error retrieving node %s: %v", name, err)
	}
	if node == nil {
		log.Printf("[DEBUG] Node %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", node.Name)
	_ = d.Set("partition", node.Partition)
	_ = d.Set("address", node.Address)
	_ = d.Set("connection_limit", node.ConnectionLimit)
	_ = d.Set("dynamic_ratio", node.DynamicRatio)
	_ = d.Set("monitor", node.Monitor)
	_ = d.Set("rate_limit", node.RateLimit)
	_ = d.Set("ratio", node.Ratio)
	_ = d.Set("state", node.State)
	_ = d.Set("fqdn.0.interval", node.FQDN.Interval)
	_ = d.Set("fqdn.0.downinterval", node.FQDN.DownInterval)
	_ = d.Set("fqdn.0.autopopulate", node.FQDN.AutoPopulate)
	_ = d.Set("fqdn.0.address_family", node.FQDN.AddressFamily)
	_ = d.Set("fqdn.0.name", node.FQDN.Name)
	d.SetId(node.Name)

	return nil

}
