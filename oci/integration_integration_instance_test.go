// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v45/common"
	oci_integration "github.com/oracle/oci-go-sdk/v45/integration"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

type IdcsAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}

var (
	IntegrationInstanceRequiredOnlyResource = IntegrationInstanceResourceDependencies +
		generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Required, Create, integrationInstanceRepresentation)

	IntegrationInstanceResourceConfig = IntegrationInstanceResourceDependencies +
		generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Update, integrationInstanceRepresentation)

	integrationInstanceSingularDataSourceRepresentation = map[string]interface{}{
		"integration_instance_id": Representation{repType: Required, create: `${oci_integration_integration_instance.test_integration_instance.id}`},
	}

	integrationInstanceDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{repType: Required, create: `${var.compartment_id}`},
		"display_name":   Representation{repType: Optional, create: `displayName`, update: `displayName2`},
		"filter":         RepresentationGroup{Required, integrationInstanceDataSourceFilterRepresentation},
		"state":          Representation{repType: Optional, create: `ACTIVE`},
	}
	integrationInstanceDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_integration_integration_instance.test_integration_instance.id}`}},
	}

	integrationInstanceRepresentation = map[string]interface{}{
		"compartment_id":            Representation{repType: Required, create: `${var.compartment_id}`},
		"display_name":              Representation{repType: Required, create: `displayName`, update: `displayName2`},
		"integration_instance_type": Representation{repType: Required, create: `STANDARD`, update: `ENTERPRISE`},
		"is_byol":                   Representation{repType: Required, create: `false`, update: `true`},
		"message_packs":             Representation{repType: Required, create: `1`, update: `2`},
		// Not supported yet
		// "alternate_custom_endpoints": RepresentationGroup{Optional, integrationInstanceAlternateCustomEndpointsRepresentation},
		"consumption_model":         Representation{repType: Optional, create: `UCM`},
		"custom_endpoint":           RepresentationGroup{Optional, integrationInstanceCustomEndpointRepresentation},
		"defined_tags":              Representation{repType: Optional, create: `${map("${var.oci_identity_tag_namespace}.${var.oci_identity_tag}", "value")}`, update: `${map("${var.oci_identity_tag_namespace}.${var.oci_identity_tag}", "updatedValue")}`},
		"freeform_tags":             Representation{repType: Optional, create: map[string]string{"bar-key": "value"}, update: map[string]string{"Department": "Accounting"}},
		"idcs_at":                   Representation{repType: Required, create: `${var.idcs_access_token}`},
		"is_file_server_enabled":    Representation{repType: Optional, create: `false`, update: `true`},
		"is_visual_builder_enabled": Representation{repType: Optional, create: `false`, update: `true`},
		"network_endpoint_details":  RepresentationGroup{Optional, integrationInstanceNetworkEndpointDetailsRepresentation},
	}
	integrationInstanceAlternateCustomEndpointsRepresentation = map[string]interface{}{
		"hostname":              Representation{repType: Required, create: `althostname.com`, update: `althostname2.com`},
		"certificate_secret_id": Representation{repType: Optional, create: `${var.oci_vault_secret_id}`},
	}
	integrationInstanceCustomEndpointRepresentation = map[string]interface{}{
		"hostname":              Representation{repType: Required, create: `hostname2.com`, update: `hostname2-updated.com`},
		"certificate_secret_id": Representation{repType: Optional, create: `${var.oci_vault_secret_id}`},
	}
	integrationInstanceNetworkEndpointDetailsRepresentation = map[string]interface{}{
		"network_endpoint_type":          Representation{repType: Required, create: `PUBLIC`},
		"allowlisted_http_ips":           Representation{repType: Optional, create: []string{`172.16.0.239/32`}},
		"allowlisted_http_vcns":          RepresentationGroup{Optional, integrationInstanceNetworkEndpointDetailsAllowlistedHttpVcnsRepresentation},
		"is_integration_vcn_allowlisted": Representation{repType: Optional, create: `false`},
	}

	integrationInstanceVcnRepresentation = `resource "oci_core_vcn" "vcn" {
cidr_blocks    = ["10.0.0.0/16"]
dns_label      = "vcn"
compartment_id = var.compartment_id
display_name   = "vcn"
}`

	integrationInstanceNetworkEndpointDetailsAllowlistedHttpVcnsRepresentation = map[string]interface{}{
		"id":              Representation{repType: Required, create: `oci_core_vcn.vcn.id`},
		"allowlisted_ips": Representation{repType: Optional, create: []string{`172.16.0.239/32`}},
	}

	IntegrationInstanceResourceDependencies = DefinedTagsDependencies + KmsVaultIdVariableStr + integrationInstanceVcnRepresentation
)

// issue-routing-tag: integration/default
func TestIntegrationIntegrationInstanceResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestIntegrationIntegrationInstanceResource_basic")
	defer httpreplay.SaveScenario()

	if strings.Contains(getEnvSettingWithBlankDefault("suppressed_tests"), "TestIntegrationIntegrationInstanceResource_basic") {
		t.Skip("Skipping suppressed TestIntegrationIntegrationInstanceResource_basic")
	}

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	vaultSecretId := getEnvSettingWithBlankDefault("oci_vault_secret_id")
	vaultSecretIdStr := fmt.Sprintf("variable \"oci_vault_secret_id\" { default = \"%s\" }\n", vaultSecretId)

	resourceName := "oci_integration_integration_instance.test_integration_instance"
	datasourceName := "data.oci_integration_integration_instances.test_integration_instances"
	singularDatasourceName := "data.oci_integration_integration_instance.test_integration_instance"

	var resId, resId2 string
	// Save TF content to create resource with optional properties. This has to be exactly the same as the config part in the "create with optionals" step in the test.
	saveConfigContent(config+compartmentIdVariableStr+IntegrationInstanceResourceDependencies+
		generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Create, integrationInstanceRepresentation), "integration", "integrationInstance", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckIntegrationIntegrationInstanceDestroy,
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config + compartmentIdVariableStr + idcsAccessTokenVariableStr() + IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Required, Create, integrationInstanceRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttr(resourceName, "integration_instance_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "is_byol", "false"),
					resource.TestCheckResourceAttr(resourceName, "message_packs", "1"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + idcsAccessTokenVariableStr(),
			},
			// verify create with optionals
			{
				Config: config + compartmentIdVariableStr +
					tagVariablesStr() +
					idcsAccessTokenVariableStr() +
					vaultSecretIdStr +
					IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap(
						"oci_integration_integration_instance",
						"test_integration_instance",
						Optional,
						Create,
						integrationInstanceRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					// resource.TestCheckResourceAttr(resourceName, "alternate_custom_endpoints.#", "1"),
					// CheckResourceSetContainsElementWithProperties(resourceName, "alternate_custom_endpoints", map[string]string{
					// 	"hostname": "hostname",
					// },
					// 	[]string{
					// 		"certificate_secret_id",
					// 	}),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "consumption_model", "UCM"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "consumption_model", "UCM"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_endpoint.0.certificate_secret_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.0.hostname", "hostname2.com"),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "idcs_at"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_url"),
					resource.TestCheckResourceAttr(resourceName, "integration_instance_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "is_byol", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_file_server_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_visual_builder_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "message_packs", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.allowlisted_http_vcns.#", "1"),
					CheckResourceSetContainsElementWithProperties(resourceName, "network_endpoint_details.0.allowlisted_http_vcns", map[string]string{
						"id": getEnvSettingWithBlankDefault("allow_listed_http_vcn"),
					},
						[]string{}),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.is_integration_vcn_allowlisted", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.network_endpoint_type", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),

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
				Config: config + compartmentIdVariableStr +
					tagVariablesStr() +
					idcsAccessTokenVariableStr() +
					compartmentIdUVariableStr +
					vaultSecretIdStr +
					IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap(
						"oci_integration_integration_instance",
						"test_integration_instance",
						Optional,
						Create,
						representationCopyWithNewProperties(integrationInstanceRepresentation, map[string]interface{}{
							"compartment_id": Representation{repType: Required, create: `${var.compartment_id_for_update}`},
						})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					// resource.TestCheckResourceAttr(resourceName, "alternate_custom_endpoints.#", "1"),
					// CheckResourceSetContainsElementWithProperties(resourceName, "alternate_custom_endpoints", map[string]string{
					// 	"hostname": "hostname",
					// },
					// 	[]string{
					// 		"certificate_secret_id",
					// 	}),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttr(resourceName, "consumption_model", "UCM"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_endpoint.0.certificate_secret_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.0.hostname", "hostname2.com"),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "idcs_at"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_url"),
					resource.TestCheckResourceAttr(resourceName, "integration_instance_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "is_byol", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_file_server_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_visual_builder_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "message_packs", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.allowlisted_http_vcns.#", "1"),
					CheckResourceSetContainsElementWithProperties(resourceName, "network_endpoint_details.0.allowlisted_http_vcns", map[string]string{
						"id": getEnvSettingWithBlankDefault("allow_listed_http_vcn"),
					},
						[]string{}),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.is_integration_vcn_allowlisted", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.network_endpoint_type", "PUBLIC"),

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
				Config: config + compartmentIdVariableStr + tagVariablesStr() + idcsAccessTokenVariableStr() + vaultSecretIdStr + IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Update, integrationInstanceRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					// resource.TestCheckResourceAttr(resourceName, "alternate_custom_endpoints.#", "1"),
					// CheckResourceSetContainsElementWithProperties(resourceName, "alternate_custom_endpoints", map[string]string{
					// 	"hostname": "hostname2",
					// },
					// 	[]string{
					// 		"certificate_secret_id",
					// 	}),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "consumption_model", "UCM"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_endpoint.0.certificate_secret_id"),
					resource.TestCheckResourceAttr(resourceName, "custom_endpoint.0.hostname", "hostname2-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "idcs_at"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_url"),
					resource.TestCheckResourceAttr(resourceName, "integration_instance_type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(resourceName, "is_byol", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_file_server_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_visual_builder_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "message_packs", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.allowlisted_http_vcns.#", "1"),
					CheckResourceSetContainsElementWithProperties(resourceName, "network_endpoint_details.0.allowlisted_http_vcns", map[string]string{
						"id": getEnvSettingWithBlankDefault("allow_listed_http_vcn"),
					},
						[]string{}),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.is_integration_vcn_allowlisted", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_endpoint_details.0.network_endpoint_type", "PUBLIC"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),

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
					generateDataSourceFromRepresentationMap("oci_integration_integration_instances", "test_integration_instances", Optional, Update, integrationInstanceDataSourceRepresentation) +
					compartmentIdVariableStr + tagVariablesStr() + idcsAccessTokenVariableStr() + vaultSecretIdStr + IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Update, integrationInstanceRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),

					resource.TestCheckResourceAttr(datasourceName, "integration_instances.#", "1"),
					// resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.alternate_custom_endpoints.#", "1"),
					// CheckResourceSetContainsElementWithProperties(datasourceName, "integration_instances.0.alternate_custom_endpoints", map[string]string{
					// 	"hostname": "hostname2",
					// },
					// 	[]string{
					// 		"certificate_secret_id",
					// 		"certificate_secret_version",
					// 	}),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.consumption_model", "UCM"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.custom_endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.custom_endpoint.0.certificate_secret_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.custom_endpoint.0.certificate_secret_version"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.custom_endpoint.0.hostname", "hostname2-updated.com"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.instance_url"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.integration_instance_type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.is_byol", "true"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.is_file_server_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.is_visual_builder_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.message_packs", "2"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.network_endpoint_details.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.network_endpoint_details.0.allowlisted_http_vcns.#", "1"),
					CheckResourceSetContainsElementWithProperties(datasourceName, "integration_instances.0.network_endpoint_details.0.allowlisted_http_vcns", map[string]string{
						"id": getEnvSettingWithBlankDefault("allow_listed_http_vcn"),
					},
						[]string{}),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.network_endpoint_details.0.is_integration_vcn_allowlisted", "false"),
					resource.TestCheckResourceAttr(datasourceName, "integration_instances.0.network_endpoint_details.0.network_endpoint_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.time_created"),
					resource.TestCheckResourceAttrSet(datasourceName, "integration_instances.0.time_updated"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Update, integrationInstanceSingularDataSourceRepresentation) +
					compartmentIdVariableStr + tagVariablesStr() + idcsAccessTokenVariableStr() + vaultSecretIdStr + IntegrationInstanceResourceDependencies +
					generateResourceFromRepresentationMap("oci_integration_integration_instance", "test_integration_instance", Optional, Update, integrationInstanceRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					// resource.TestCheckResourceAttr(singularDatasourceName, "alternate_custom_endpoints.#", "1"),
					// CheckResourceSetContainsElementWithProperties(singularDatasourceName, "alternate_custom_endpoints", map[string]string{
					// 	"hostname": "hostname2",
					// },
					// 	[]string{
					// 		"certificate_secret_version",
					// 	}),
					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "consumption_model", "UCM"),
					resource.TestCheckResourceAttr(singularDatasourceName, "custom_endpoint.#", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "custom_endpoint.0.certificate_secret_version"),
					resource.TestCheckResourceAttr(singularDatasourceName, "custom_endpoint.0.hostname", "hostname2-updated.com"),
					resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "instance_url"),
					resource.TestCheckResourceAttr(singularDatasourceName, "integration_instance_type", "ENTERPRISE"),
					resource.TestCheckResourceAttr(singularDatasourceName, "is_byol", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "is_file_server_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "is_visual_builder_enabled", "true"),
					resource.TestCheckResourceAttr(singularDatasourceName, "message_packs", "2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "network_endpoint_details.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "network_endpoint_details.0.allowlisted_http_vcns.#", "1"),
					CheckResourceSetContainsElementWithProperties(singularDatasourceName, "network_endpoint_details.0.allowlisted_http_vcns", map[string]string{
						"id": getEnvSettingWithBlankDefault("allow_listed_http_vcn"),
					},
						[]string{}),
					resource.TestCheckResourceAttr(singularDatasourceName, "network_endpoint_details.0.allowlisted_http_vcns.0.id",
						getEnvSettingWithBlankDefault("allow_listed_http_vcn")),
					resource.TestCheckResourceAttr(singularDatasourceName, "network_endpoint_details.0.is_integration_vcn_allowlisted", "false"),
					resource.TestCheckResourceAttr(singularDatasourceName, "network_endpoint_details.0.network_endpoint_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_updated"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + tagVariablesStr() + idcsAccessTokenVariableStr() + vaultSecretIdStr + IntegrationInstanceResourceConfig,
			},
			// verify resource import
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"idcs_at",
				},
				ResourceName: resourceName,
			},
		},
	})
}

func testAccCheckIntegrationIntegrationInstanceDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).integrationInstanceClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_integration_integration_instance" {
			noResourceFound = false
			request := oci_integration.GetIntegrationInstanceRequest{}

			tmp := rs.Primary.ID
			request.IntegrationInstanceId = &tmp

			request.RequestMetadata.RetryPolicy = getRetryPolicy(true, "integration")

			response, err := client.GetIntegrationInstance(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_integration.IntegrationInstanceLifecycleStateDeleted): true,
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
	if !inSweeperExcludeList("IntegrationIntegrationInstance") {
		resource.AddTestSweepers("IntegrationIntegrationInstance", &resource.Sweeper{
			Name:         "IntegrationIntegrationInstance",
			Dependencies: DependencyGraph["integrationInstance"],
			F:            sweepIntegrationIntegrationInstanceResource,
		})
	}
}

func sweepIntegrationIntegrationInstanceResource(compartment string) error {
	integrationInstanceClient := GetTestClients(&schema.ResourceData{}).integrationInstanceClient()
	integrationInstanceIds, err := getIntegrationInstanceIds(compartment)
	if err != nil {
		return err
	}
	for _, integrationInstanceId := range integrationInstanceIds {
		if ok := SweeperDefaultResourceId[integrationInstanceId]; !ok {
			deleteIntegrationInstanceRequest := oci_integration.DeleteIntegrationInstanceRequest{}

			deleteIntegrationInstanceRequest.IntegrationInstanceId = &integrationInstanceId

			deleteIntegrationInstanceRequest.RequestMetadata.RetryPolicy = getRetryPolicy(true, "integration")
			_, error := integrationInstanceClient.DeleteIntegrationInstance(context.Background(), deleteIntegrationInstanceRequest)
			if error != nil {
				fmt.Printf("Error deleting IntegrationInstance %s %s, It is possible that the resource is already deleted. Please verify manually \n", integrationInstanceId, error)
				continue
			}
			waitTillCondition(testAccProvider, &integrationInstanceId, integrationInstanceSweepWaitCondition, time.Duration(3*time.Minute),
				integrationInstanceSweepResponseFetchOperation, "integration", true)
		}
	}
	return nil
}

func getIntegrationInstanceIds(compartment string) ([]string, error) {
	ids := getResourceIdsToSweep(compartment, "IntegrationInstanceId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	integrationInstanceClient := GetTestClients(&schema.ResourceData{}).integrationInstanceClient()

	listIntegrationInstancesRequest := oci_integration.ListIntegrationInstancesRequest{}
	listIntegrationInstancesRequest.CompartmentId = &compartmentId
	listIntegrationInstancesRequest.LifecycleState = oci_integration.ListIntegrationInstancesLifecycleStateActive
	listIntegrationInstancesResponse, err := integrationInstanceClient.ListIntegrationInstances(context.Background(), listIntegrationInstancesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting IntegrationInstance list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, integrationInstance := range listIntegrationInstancesResponse.Items {
		id := *integrationInstance.Id
		resourceIds = append(resourceIds, id)
		addResourceIdToSweeperResourceIdMap(compartmentId, "IntegrationInstanceId", id)
	}
	return resourceIds, nil
}

func integrationInstanceSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if integrationInstanceResponse, ok := response.Response.(oci_integration.GetIntegrationInstanceResponse); ok {
		return integrationInstanceResponse.LifecycleState != oci_integration.IntegrationInstanceLifecycleStateDeleted
	}
	return false
}

func integrationInstanceSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.integrationInstanceClient().GetIntegrationInstance(context.Background(), oci_integration.GetIntegrationInstanceRequest{
		IntegrationInstanceId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}

func idcsAccessTokenVariableStr() string {
	return fmt.Sprintf("variable \"idcs_access_token\" { default = \"%s\" }\n", readIdcsAccessToken())
}

func tagVariablesStr() string {
	vars := []string{
		fmt.Sprintf("variable \"oci_identity_tag_namespace\" { default = \"%s\" }\n", getEnvSettingWithBlankDefault("oci_identity_tag_namespace")),
		fmt.Sprintf("variable \"oci_identity_tag\" { default = \"%s\" }\n", getEnvSettingWithBlankDefault("oci_identity_tag")),
		fmt.Sprintf("variable \"allow_listed_http_vcn\" { default = \"%s\" }\n", getEnvSettingWithBlankDefault("allow_listed_http_vcn")),
	}

	return strings.Join(vars[:], "\n")
}

func readIdcsAccessToken() string {
	return getEnvSettingWithBlankDefault("idcs_access_token")
}
