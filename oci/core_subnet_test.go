// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v45/common"
	oci_core "github.com/oracle/oci-go-sdk/v45/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	SubnetRequiredOnlyResource = SubnetResourceDependencies +
		generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation)

	SubnetResourceConfig = SubnetResourceDependencies +
		generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Update, subnetRepresentation)

	subnetSingularDataSourceRepresentation = map[string]interface{}{
		"subnet_id": Representation{repType: Required, create: `${oci_core_subnet.test_subnet.id}`},
	}

	subnetDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"display_name":   Representation{repType: Optional, create: `MySubnet`, update: `displayName2`},
		"state":          Representation{repType: Optional, create: `AVAILABLE`},
		"vcn_id":         Representation{repType: Optional, create: `${oci_core_vcn.test_vcn.id}`},
		"filter":         RepresentationGroup{Required, subnetDataSourceFilterRepresentation}}
	subnetDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_core_subnet.test_subnet.id}`}},
	}

	subnetRepresentation = map[string]interface{}{
		"cidr_block":                 Representation{repType: Required, create: `10.0.0.0/24`, update: "10.0.0.0/16"},
		"compartment_id":             Representation{repType: Required, create: `${var.compartment_id}`},
		"vcn_id":                     Representation{repType: Required, create: `${oci_core_vcn.test_vcn.id}`},
		"availability_domain":        Representation{repType: Optional, create: `${lower("${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}")}`},
		"defined_tags":               Representation{repType: Optional, create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"dhcp_options_id":            Representation{repType: Optional, create: `${oci_core_vcn.test_vcn.default_dhcp_options_id}`, update: `${oci_core_dhcp_options.test_dhcp_options.id}`},
		"display_name":               Representation{repType: Optional, create: `MySubnet`, update: `displayName2`},
		"dns_label":                  Representation{repType: Optional, create: `dnslabel`},
		"freeform_tags":              Representation{repType: Optional, create: map[string]string{"Department": "Finance"}, update: map[string]string{"Department": "Accounting"}},
		"prohibit_public_ip_on_vnic": Representation{repType: Optional, create: `false`},
		"prohibit_internet_ingress":  Representation{repType: Optional, create: `false`},
		"route_table_id":             Representation{repType: Optional, create: `${oci_core_vcn.test_vcn.default_route_table_id}`, update: `${oci_core_route_table.test_route_table.id}`},
		"security_list_ids":          Representation{repType: Optional, create: []string{`${oci_core_vcn.test_vcn.default_security_list_id}`}, update: []string{`${oci_core_security_list.test_security_list.id}`}},
	}

	SubnetResourceDependencies = generateResourceFromRepresentationMap("oci_core_dhcp_options", "test_dhcp_options", Required, Create, dhcpOptionsRepresentation) +
		generateResourceFromRepresentationMap("oci_core_internet_gateway", "test_internet_gateway", Required, Create, internetGatewayRepresentation) +
		generateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Create, routeTableRepresentation) +
		generateResourceFromRepresentationMap("oci_core_security_list", "test_security_list", Required, Create, securityListRepresentation) +
		generateDataSourceFromRepresentationMap("oci_core_services", "test_services", Required, Create, serviceDataSourceRepresentation) +
		generateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Create, representationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
			"dns_label":      Representation{repType: Required, create: `dnslabel`},
			"is_ipv6enabled": Representation{repType: Optional, create: `true`},
		})) +
		AvailabilityDomainConfig +
		DefinedTagsDependencies
	AnotherSecurityListRequiredOnlyResource = generateResourceFromRepresentationMap("oci_core_security_list", "test_security_list", Required, Create, securityListRepresentation)
	SubnetRequiredOnlyResourceDependencies  = AvailabilityDomainConfig + VcnResourceConfig
)

// issue-routing-tag: core/virtualNetwork
func TestCoreSubnetResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreSubnetResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_core_subnet.test_subnet"
	datasourceName := "data.oci_core_subnets.test_subnets"
	singularDatasourceName := "data.oci_core_subnet.test_subnet"

	// Get subnet CIDR block based on its VCN CIDR Block
	// For example: VCN CIDR Block: 2607:9b80:9a0f:0100::/56, Subnet CIDR Block: 2607:9b80:9a0f:0100::/64
	subnetCidrBlock := `${substr(oci_core_vcn.test_vcn.ipv6cidr_blocks[0], 0, length(oci_core_vcn.test_vcn.ipv6cidr_blocks[0]) - 2)}${64}`
	var resId, resId2 string
	// Save TF content to create resource with optional properties. This has to be exactly the same as the config part in the "create with optionals" step in the test.
	saveConfigContent(config+compartmentIdVariableStr+SubnetResourceDependencies+
		generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Create, subnetRepresentation), "core", "subnet", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckCoreSubnetDestroy,
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config + compartmentIdVariableStr + SubnetResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + SubnetResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + compartmentIdVariableStr + SubnetResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Create, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{
						"ipv6cidr_block": Representation{repType: Optional, create: subnetCidrBlock},
					})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "availability_domain"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "MySubnet"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6cidr_block"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_public_ip_on_vnic", "false"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_internet_ingress", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_mac"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "true")); isEnableExportCompartment {
							if errExport := testExportCompartmentWithResourceName(&resId, &compartmentId, resourceName); errExport != nil {
								return errExport
							}
						}
						return err
					},
				),
			},

			// verify update to the compartment (the compartment will be switched back in the next step)
			{
				Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + SubnetResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Create,
						representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{
							"compartment_id": Representation{repType: Required, create: `${var.compartment_id_for_update}`},
						})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "availability_domain"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "MySubnet"),
					resource.TestCheckResourceAttr(resourceName, "dns_label", "dnslabel"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6cidr_block"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_public_ip_on_vnic", "false"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_internet_ingress", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_mac"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("resource recreated when it was supposed to be updated")
						}
						return err
					},
				),
			},

			// verify updates to updatable parameters
			{
				Config: config + compartmentIdVariableStr + SubnetResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Update, representationCopyWithNewProperties(subnetRepresentation, map[string]interface{}{
						"ipv6cidr_block": Representation{repType: Optional, update: subnetCidrBlock},
					})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "availability_domain"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "dhcp_options_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "dns_label", "dnslabel"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6cidr_block"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_public_ip_on_vnic", "false"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_internet_ingress", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "route_table_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_router_mac"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_core_subnets", "test_subnets", Optional, Update, subnetDataSourceRepresentation) +
					compartmentIdVariableStr + SubnetResourceDependencies +
					generateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Optional, Update, subnetRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),
					resource.TestCheckResourceAttrSet(datasourceName, "vcn_id"),

					resource.TestCheckResourceAttr(datasourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.availability_domain"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.defined_tags.%", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.dhcp_options_id"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.dns_label", "dnslabel"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.ipv6cidr_block"),
					resource.TestCheckResourceAttr(datasourceName, "subnets.0.prohibit_public_ip_on_vnic", "false"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_internet_ingress", "false"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.route_table_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.subnet_domain_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.time_created"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.vcn_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.virtual_router_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "subnets.0.virtual_router_mac"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetSingularDataSourceRepresentation) +
					compartmentIdVariableStr + SubnetResourceConfig,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "subnet_id"),

					resource.TestCheckResourceAttrSet(singularDatasourceName, "availability_domain"),
					resource.TestCheckResourceAttr(singularDatasourceName, "cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "ipv6cidr_block"),
					resource.TestCheckResourceAttr(singularDatasourceName, "prohibit_public_ip_on_vnic", "false"),
					resource.TestCheckResourceAttr(resourceName, "prohibit_internet_ingress", "false"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "subnet_domain_name"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "virtual_router_ip"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "virtual_router_mac"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + SubnetResourceConfig,
			},
			// verify resource import
			{
				Config:                  config,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ResourceName:            resourceName,
			},
		},
	})
}

func testAccCheckCoreSubnetDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).virtualNetworkClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_subnet" {
			noResourceFound = false
			request := oci_core.GetSubnetRequest{}

			tmp := rs.Primary.ID
			request.SubnetId = &tmp

			request.RequestMetadata.RetryPolicy = getRetryPolicy(true, "core")

			response, err := client.GetSubnet(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.SubnetLifecycleStateTerminated): true,
				}
				if _, ok := deletedLifecycleStates[string(response.LifecycleState)]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.LifecycleState)
				}
				//resource lifecycle state is in expected deleted lifecycle states. continue with next one.
				continue
			}

			//Verify that exception is for '404 not found'.
			if failure, isServiceError := common.IsServiceError(err); !isServiceError || failure.GetHTTPStatusCode() != 404 {
				return err
			}
		}
	}
	if noResourceFound {
		return fmt.Errorf("at least one resource was expected from the state file, but could not be found")
	}

	return nil
}

func init() {
	if DependencyGraph == nil {
		initDependencyGraph()
	}
	if !inSweeperExcludeList("CoreSubnet") {
		resource.AddTestSweepers("CoreSubnet", &resource.Sweeper{
			Name:         "CoreSubnet",
			Dependencies: DependencyGraph["subnet"],
			F:            sweepCoreSubnetResource,
		})
	}
}

func sweepCoreSubnetResource(compartment string) error {
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()
	subnetIds, err := getSubnetIds(compartment)
	if err != nil {
		return err
	}
	for _, subnetId := range subnetIds {
		if ok := SweeperDefaultResourceId[subnetId]; !ok {
			deleteSubnetRequest := oci_core.DeleteSubnetRequest{}

			deleteSubnetRequest.SubnetId = &subnetId

			deleteSubnetRequest.RequestMetadata.RetryPolicy = getRetryPolicy(true, "core")
			_, error := virtualNetworkClient.DeleteSubnet(context.Background(), deleteSubnetRequest)
			if error != nil {
				fmt.Printf("Error deleting Subnet %s %s, It is possible that the resource is already deleted. Please verify manually \n", subnetId, error)
				continue
			}
			waitTillCondition(testAccProvider, &subnetId, subnetSweepWaitCondition, time.Duration(3*time.Minute),
				subnetSweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getSubnetIds(compartment string) ([]string, error) {
	ids := getResourceIdsToSweep(compartment, "SubnetId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()

	listSubnetsRequest := oci_core.ListSubnetsRequest{}
	listSubnetsRequest.CompartmentId = &compartmentId
	listSubnetsRequest.LifecycleState = oci_core.SubnetLifecycleStateAvailable
	listSubnetsResponse, err := virtualNetworkClient.ListSubnets(context.Background(), listSubnetsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting Subnet list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, subnet := range listSubnetsResponse.Items {
		id := *subnet.Id
		resourceIds = append(resourceIds, id)
		addResourceIdToSweeperResourceIdMap(compartmentId, "SubnetId", id)
	}
	return resourceIds, nil
}

func subnetSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if subnetResponse, ok := response.Response.(oci_core.GetSubnetResponse); ok {
		return subnetResponse.LifecycleState != oci_core.SubnetLifecycleStateTerminated
	}
	return false
}

func subnetSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.virtualNetworkClient().GetSubnet(context.Background(), oci_core.GetSubnetRequest{
		SubnetId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
