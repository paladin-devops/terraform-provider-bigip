package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_DNS_NAME = fmt.Sprintf("/%s/test-dns", TEST_PARTITION)

var TEST_DNS_RESOURCE = `
resource "bigip_dns" "test-dns" {
   description = "` + TEST_DNS_NAME + `"
   name_servers = ["1.1.1.1"]
   numberof_dots = 2
   search = ["f5.com"]
}

`

func TestBigipLtmdns_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//CheckDestroy: testCheckdnssDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
					resource.TestCheckResourceAttr("bigip_dns.test-dns", "description", TEST_DNS_NAME),
					resource.TestCheckResourceAttr("bigip_dns.test-dns", "numberof_dots", "2"),
					resource.TestCheckResourceAttr("bigip_dns.test-dns",
						fmt.Sprintf("name_servers.%d", schema.HashString("1.1.1.1")),
						"1.1.1.1"),
					resource.TestCheckResourceAttr("bigip_dns.test-dns",
						fmt.Sprintf("search.%d", schema.HashString("f5.com")),
						"f5.com"),
				),
			},
		},
	})
}

func TestBigipLtmdns_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers: testAccProviders,
		//	CheckDestroy: testCheckdnssDestroyed, ( No Delet API support)
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TEST_DNS_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckdnsExists(TEST_DNS_NAME, true),
				),
				ResourceName:      TEST_DNS_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

//var TEST_NODE_IN_POOL_RESOURCE = `
//resource "bigip_ltm_pool" "test-pool" {
//	name = "` + TEST_POOL_NAME + `"
//  	load_balancing_mode = "round-robin"
//  	nodes = ["${formatlist("%s:80", bigip_ltm_node.*.name)}"]
//  	allow_snat = false
//}
//`
//func TestBigipLtmNode_removeNode(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAcctPreCheck(t)
//		},
//		Providers: testAccProviders,
//		CheckDestroy: testCheckNodesDestroyed,
//		Steps: []resource.TestStep{
//			resource.TestStep{
//				Config: TEST_NODE_RESOURCE + TEST_NODE_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckNodeExists(TEST_NODE_NAME, true),
//					testCheckPoolExists(TEST_POOL_NAME, true),
//					testCheckPoolMember(TEST_POOL_NAME, TEST_NODE_NAME),
//				),
//			},
//			resource.TestStep{
//				Config: TEST_NODE_IN_POOL_RESOURCE,
//				Check: resource.ComposeTestCheckFunc(
//					testCheckNodeExists(fmt.Sprintf("%s:%s", TEST_NODE_NAME, "80"), false),
//					testCheckEmptyPool(TEST_POOL_NAME),
//				),
//			},
//		},
//	})
//}

func testCheckdnsExists(description string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		dns, err := client.DNSs()
		if err != nil {
			return err
		}
		if exists && dns == nil {
			return fmt.Errorf("dns ", description, " was not created.")

		}
		if !exists && dns != nil {
			return fmt.Errorf("dns ", description, " still exists.")

		}
		return nil
	}
}

func testCheckdnssDestroyed(s *terraform.State) error {
	/* client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_dns" {
			continue
		}

		description := rs.Primary.ID
		dns, err := client.dnss()
		if err != nil {
			return err
		}
		if dns != nil {
			return fmt.Errorf("dns ", description, " not destroyed.")

		}
	}*/
	return nil
}