---
subcategory: "Database Management"
layout: "oci"
page_title: "Oracle Cloud Infrastructure: oci_database_management_external_asm"
sidebar_current: "docs-oci-resource-database_management-external_asm"
description: |-
  Provides the External Asm resource in Oracle Cloud Infrastructure Database Management service
---

# oci_database_management_external_asm
This resource provides the External Asm resource in Oracle Cloud Infrastructure Database Management service.

Updates the external ASM specified by `externalAsmId`.


## Example Usage

```hcl
resource "oci_database_management_external_asm" "test_external_asm" {
	#Required
	external_asm_id = oci_database_management_external_asm.test_external_asm.id

	#Optional
	external_connector_id = oci_database_management_external_connector.test_external_connector.id
}
```

## Argument Reference

The following arguments are supported:

* `external_asm_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external ASM.
* `external_connector_id` - (Optional) (Updatable) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external connector.


** IMPORTANT **
Any change to a property that does not support update will force the destruction and recreation of the resource with the new property values

## Attributes Reference

The following attributes are exported:

* `additional_details` - The additional details of the external ASM defined in `{"key": "value"}` format. Example: `{"bar-key": "value"}` 
* `compartment_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment.
* `component_name` - The name of the external ASM.
* `display_name` - The user-friendly name for the external ASM. The name does not have to be unique.
* `external_connector_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external connector.
* `external_db_system_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external DB system that the ASM is a part of.
* `grid_home` - The directory in which ASM is installed. This is the same directory in which Oracle Grid Infrastructure is installed.
* `id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external ASM.
* `is_cluster` - Indicates whether the ASM is a cluster ASM or not.
* `is_flex_enabled` - Indicates whether Oracle Flex ASM is enabled or not.
* `lifecycle_details` - Additional information about the current lifecycle state.
* `serviced_databases` - The list of databases that are serviced by the ASM.
	* `compartment_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment in which the external database resides.
	* `database_sub_type` - The subtype of Oracle Database. Indicates whether the database is a Container Database, Pluggable Database, Non-container Database, Autonomous Database, or Autonomous Container Database. 
	* `database_type` - The type of Oracle Database installation.
	* `db_unique_name` - The unique name of the external database.
	* `disk_groups` - The list of ASM disk groups used by the database.
	* `display_name` - The user-friendly name for the database. The name does not have to be unique.
	* `id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the external database.
	* `is_managed` - Indicates whether the database is a Managed Database or not.
* `state` - The current lifecycle state of the external ASM.
* `time_created` - The date and time the external ASM was created.
* `time_updated` - The date and time the external ASM was last updated.
* `version` - The ASM version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://registry.terraform.io/providers/oracle/oci/latest/docs/guides/changing_timeouts) for certain operations:
	* `create` - (Defaults to 20 minutes), when creating the External Asm
	* `update` - (Defaults to 20 minutes), when updating the External Asm
	* `delete` - (Defaults to 20 minutes), when destroying the External Asm


## Import

ExternalAsms can be imported using the `id`, e.g.

```
$ terraform import oci_database_management_external_asm.test_external_asm "id"
```

