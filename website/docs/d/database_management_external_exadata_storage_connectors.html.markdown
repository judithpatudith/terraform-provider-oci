---
subcategory: "Database Management"
layout: "oci"
page_title: "Oracle Cloud Infrastructure: oci_database_management_external_exadata_storage_connectors"
sidebar_current: "docs-oci-datasource-database_management-external_exadata_storage_connectors"
description: |-
  Provides the list of External Exadata Storage Connectors in Oracle Cloud Infrastructure Database Management service
---

# Data Source: oci_database_management_external_exadata_storage_connectors
This data source provides the list of External Exadata Storage Connectors in Oracle Cloud Infrastructure Database Management service.

Lists the Exadata storage server connectors for the specified Exadata infrastructure.


## Example Usage

```hcl
data "oci_database_management_external_exadata_storage_connectors" "test_external_exadata_storage_connectors" {
	#Required
	compartment_id = var.compartment_id
	external_exadata_infrastructure_id = oci_database_management_external_exadata_infrastructure.test_external_exadata_infrastructure.id

	#Optional
	display_name = var.external_exadata_storage_connector_display_name
}
```

## Argument Reference

The following arguments are supported:

* `compartment_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment.
* `display_name` - (Optional) The optional single value query filter parameter on the entity display name.
* `external_exadata_infrastructure_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Exadata infrastructure.


## Attributes Reference

The following attributes are exported:

* `external_exadata_storage_connector_collection` - The list of external_exadata_storage_connector_collection.

### ExternalExadataStorageConnector Reference

The following attributes are exported:

* `additional_details` - The additional details of the resource defined in `{"key": "value"}` format. Example: `{"bar-key": "value"}` 
* `agent_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the agent for the Exadata storage server.
* `connection_uri` - The unique string of the connection. For example, "https://<storage-server-name>/MS/RESTService/".
* `display_name` - The name of the Exadata resource. English letters, numbers, "-", "_" and "." only.
* `exadata_infrastructure_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Exadata infrastructure.
* `id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Exadata resource.
* `internal_id` - The internal ID of the Exadata resource.
* `lifecycle_details` - The details of the lifecycle state of the Exadata resource.
* `resource_type` - The type of Exadata resource.
* `state` - The current lifecycle state of the database resource.
* `status` - The status of the Exadata resource.
* `storage_server_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Exadata storage server.
* `time_created` - The timestamp of the creation of the Exadata resource.
* `time_updated` - The timestamp of the last update of the Exadata resource.
* `version` - The version of the Exadata resource.

