// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v45/common"

	oci_kms "github.com/oracle/oci-go-sdk/v45/keymanagement"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	KeyRequiredOnlyResource = KeyResourceDependencies +
		generateResourceFromRepresentationMap("oci_kms_key", "test_key", Required, Create, keyRepresentation)

	KeyResourceConfig = KeyResourceDependencies +
		generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Update, keyRepresentation)

	keySingularDataSourceRepresentation = map[string]interface{}{
		"key_id":              Representation{repType: Required, create: `${oci_kms_key.test_key.id}`},
		"management_endpoint": Representation{repType: Required, create: `${data.oci_kms_vault.test_vault.management_endpoint}`},
	}

	keyDataSourceRepresentation = map[string]interface{}{
		"compartment_id":      Representation{repType: Required, create: `${var.compartment_id}`},
		"management_endpoint": Representation{repType: Required, create: `${data.oci_kms_vault.test_vault.management_endpoint}`},
		"protection_mode":     Representation{repType: Optional, create: `SOFTWARE`},
		"algorithm":           Representation{repType: Optional, create: `AES`},
		"length":              Representation{repType: Optional, create: `16`},
		"filter":              RepresentationGroup{Required, keyDataSourceFilterRepresentation}}
	keyDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{repType: Required, create: `id`},
		"values": Representation{repType: Required, create: []string{`${oci_kms_key.test_key.id}`}},
	}

	deletionTime = time.Now().UTC().AddDate(0, 0, 8).Truncate(time.Millisecond)

	keyRepresentation = map[string]interface{}{
		"compartment_id":      Representation{repType: Required, create: `${var.compartment_id}`},
		"display_name":        Representation{repType: Required, create: `Key C`, update: `displayName2`},
		"key_shape":           RepresentationGroup{Required, keyKeyShapeRepresentation},
		"management_endpoint": Representation{repType: Required, create: `${data.oci_kms_vault.test_vault.management_endpoint}`},
		"desired_state":       Representation{repType: Optional, create: `ENABLED`, update: `DISABLED`},
		"freeform_tags":       Representation{repType: Optional, create: map[string]string{"Department": "Finance"}, update: map[string]string{"Department": "Accounting"}},
		"protection_mode":     Representation{repType: Optional, create: `SOFTWARE`},
		"time_of_deletion":    Representation{repType: Required, create: deletionTime.Format(time.RFC3339Nano)},
	}
	keyKeyShapeRepresentation = map[string]interface{}{
		"algorithm": Representation{repType: Required, create: `AES`},
		"length":    Representation{repType: Required, create: `16`},
	}

	kmsVaultId                = getEnvSettingWithBlankDefault("kms_vault_ocid")
	KmsVaultIdVariableStr     = fmt.Sprintf("variable \"kms_vault_id\" { default = \"%s\" }\n", kmsVaultId)
	kmsKeyIdForCreate         = getEnvSettingWithBlankDefault("key_ocid_for_create")
	kmsKeyIdCreateVariableStr = fmt.Sprintf("variable \"kms_key_id_for_create\" { default = \"%s\" }\n", kmsKeyIdForCreate)

	kmsKeyIdForUpdate         = getEnvSettingWithBlankDefault("key_ocid_for_update")
	kmsKeyIdUpdateVariableStr = fmt.Sprintf("variable \"kms_key_id_for_update\" { default = \"%s\" }\n", kmsKeyIdForUpdate)

	kmsKeyCompartmentId            = getEnvSettingWithBlankDefault("compartment_ocid")
	kmsKeyCompartmentIdVariableStr = fmt.Sprintf("variable \"kms_key_compartment_id\" { default = \"%s\" }\n", kmsKeyCompartmentId)

	// Should deprecate use of tenancy level resources
	KeyResourceDependencies = KmsVaultIdVariableStr + `
	data "oci_kms_vault" "test_vault" {
		#Required
		vault_id = "${var.kms_vault_id}"
	}
	`
	KeyResourceDependencyConfig = KeyResourceDependencies + `
	data "oci_kms_keys" "test_keys_dependency" {
		#Required
		compartment_id = "${var.tenancy_ocid}"
		management_endpoint = "${data.oci_kms_vault.test_vault.management_endpoint}"
		algorithm = "AES"

		filter {
    		name = "state"
    		values = ["ENABLED", "UPDATING"]
        }
	}
	data "oci_kms_keys" "test_keys_dependency_RSA" {
		#Required
		compartment_id = "${var.tenancy_ocid}"
		management_endpoint = "${data.oci_kms_vault.test_vault.management_endpoint}"
		algorithm = "RSA"

		filter {
    		name = "state"
    		values = ["ENABLED", "UPDATING"]
        }
	}
	`

	KeyResourceDependencyConfig2 = KeyResourceDependencies + kmsKeyCompartmentIdVariableStr + `
	data "oci_kms_keys" "test_keys_dependency" {
		#Required
		compartment_id = "${var.kms_key_compartment_id}"
		management_endpoint = "${data.oci_kms_vault.test_vault.management_endpoint}"
		algorithm = "AES"

		filter {
    		name = "state"
    		values = ["ENABLED", "UPDATING"]
        }
	}
	data "oci_kms_keys" "test_keys_dependency_RSA" {
		#Required
		compartment_id = "${var.kms_key_compartment_id}"
		management_endpoint = "${data.oci_kms_vault.test_vault.management_endpoint}"
		algorithm = "RSA"

		filter {
    		name = "state"
    		values = ["ENABLED", "UPDATING"]
        }
	}
	`
)

// issue-routing-tag: kms/default
func TestKmsKeyResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestKmsKeyResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_kms_key.test_key"
	datasourceName := "data.oci_kms_keys.test_keys"
	singularDatasourceName := "data.oci_kms_key.test_key"

	var resId, resId2 string
	// Save TF content to create resource with optional properties. This has to be exactly the same as the config part in the "create with optionals" step in the test.
	saveConfigContent(config+compartmentIdVariableStr+KeyResourceDependencies+
		generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Create, keyRepresentation), "keymanagement", "key", t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKMSKeyDestroy,
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create
			{
				Config: config + compartmentIdVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Required, Create, keyRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Key C"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.algorithm", "AES"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.length", "16"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next create
			{
				Config: config + compartmentIdVariableStr + KeyResourceDependencies,
			},
			// verify create with optionals
			{
				Config: config + compartmentIdVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Create, keyRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "current_key_version"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Key C"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.algorithm", "AES"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.length", "16"),
					resource.TestCheckResourceAttrSet(resourceName, "management_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vault_id"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// verify update to the compartment (the compartment will be switched back in the next step)
			{
				Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Create,
						representationCopyWithNewProperties(keyRepresentation, map[string]interface{}{
							"compartment_id": Representation{repType: Required, create: `${var.compartment_id_for_update}`},
						})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttrSet(resourceName, "current_key_version"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Key C"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.algorithm", "AES"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.length", "16"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vault_id"),

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
				Config: config + compartmentIdVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Update, keyRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "current_key_version"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.algorithm", "AES"),
					resource.TestCheckResourceAttr(resourceName, "key_shape.0.length", "16"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),
					resource.TestCheckResourceAttrSet(resourceName, "vault_id"),

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
					generateDataSourceFromRepresentationMap("oci_kms_keys", "test_keys", Optional, Update, keyDataSourceRepresentation) +
					compartmentIdVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Update, keyRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(datasourceName, "algorithm", "AES"),
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "management_endpoint"),
					resource.TestCheckResourceAttr(datasourceName, "protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttr(datasourceName, "length", "16"),

					resource.TestCheckResourceAttr(datasourceName, "keys.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "keys.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "keys.0.display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "keys.0.freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "keys.0.protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.time_created"),
					resource.TestCheckResourceAttrSet(datasourceName, "keys.0.vault_id"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					generateDataSourceFromRepresentationMap("oci_kms_key", "test_key", Required, Create, keySingularDataSourceRepresentation) +
					compartmentIdVariableStr + KeyResourceConfig + DefinedTagsDependencies,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "key_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "current_key_version"),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "is_primary"),
					resource.TestCheckResourceAttr(singularDatasourceName, "key_shape.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "key_shape.0.algorithm", "AES"),
					resource.TestCheckResourceAttr(singularDatasourceName, "key_shape.0.length", "16"),
					resource.TestCheckResourceAttr(singularDatasourceName, "protection_mode", "SOFTWARE"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "vault_id"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + KeyResourceConfig + DefinedTagsDependencies,
			},
			// revert the updates
			{
				Config: config + compartmentIdVariableStr + KeyResourceDependencies + DefinedTagsDependencies +
					generateResourceFromRepresentationMap("oci_kms_key", "test_key", Optional, Create, keyRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Key C"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify resource import
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: keyImportId,
				ImportStateVerifyIgnore: []string{
					"desired_state",
					"time_of_deletion",
					"replica_details",
				},
				ResourceName: resourceName,
			},
		},
	})
}

func testAccCheckKMSKeyDestroy(s *terraform.State) error {
	noResourceFound := true
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_kms_key" {
			client, err := testAccProvider.Meta().(*OracleClients).KmsManagementClient(rs.Primary.Attributes["management_endpoint"])
			if err != nil {
				return err
			}

			noResourceFound = false
			request := oci_kms.GetKeyRequest{}

			tmp := rs.Primary.ID
			request.KeyId = &tmp

			request.RequestMetadata.RetryPolicy = getRetryPolicy(true, "kms")

			response, err := client.GetKey(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_kms.KeyLifecycleStateSchedulingDeletion): true,
					string(oci_kms.KeyLifecycleStatePendingDeletion):    true,
				}
				if _, ok := deletedLifecycleStates[string(response.LifecycleState)]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.LifecycleState)
				}

				if !response.TimeOfDeletion.Equal(deletionTime) && !httpreplay.ModeRecordReplay() {
					return fmt.Errorf("resource time_of_deletion: %s is not set to %s", response.TimeOfDeletion.Format(time.RFC3339Nano), deletionTime.Format(time.RFC3339Nano))
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

func keyImportId(state *terraform.State) (string, error) {
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "oci_kms_key" {
			return fmt.Sprintf("managementEndpoint/%s/keys/%s", rs.Primary.Attributes["management_endpoint"], rs.Primary.ID), nil
		}
	}

	return "", fmt.Errorf("unable to create import id as no resource of type oci_kms_key in state")
}
