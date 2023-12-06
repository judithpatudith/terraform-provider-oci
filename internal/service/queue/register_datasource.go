// Copyright (c) 2017, 2023, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package queue

import "github.com/oracle/terraform-provider-oci/internal/tfresource"

func RegisterDatasource() {
	tfresource.RegisterDatasource("oci_queue_queue", QueueQueueDataSource())
	tfresource.RegisterDatasource("oci_queue_queues", QueueQueuesDataSource())
}
