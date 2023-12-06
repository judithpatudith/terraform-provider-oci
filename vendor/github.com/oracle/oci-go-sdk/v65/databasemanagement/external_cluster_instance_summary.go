// Copyright (c) 2016, 2018, 2023, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Database Management API
//
// Use the Database Management API to perform tasks such as obtaining performance and resource usage metrics
// for a fleet of Managed Databases or a specific Managed Database, creating Managed Database Groups, and
// running a SQL job on a Managed Database or Managed Database Group.
//

package databasemanagement

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"strings"
)

// ExternalClusterInstanceSummary The summary of an external cluster instance.
type ExternalClusterInstanceSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external cluster instance.
	Id *string `mandatory:"true" json:"id"`

	// The user-friendly name for the cluster instance. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The name of the external cluster instance.
	ComponentName *string `mandatory:"true" json:"componentName"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external cluster that the cluster instance belongs to.
	ExternalClusterId *string `mandatory:"true" json:"externalClusterId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external DB system that the cluster instance is a part of.
	ExternalDbSystemId *string `mandatory:"true" json:"externalDbSystemId"`

	// The current lifecycle state of the external cluster instance.
	LifecycleState ExternalClusterInstanceLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external DB node.
	ExternalDbNodeId *string `mandatory:"false" json:"externalDbNodeId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external connector.
	ExternalConnectorId *string `mandatory:"false" json:"externalConnectorId"`

	// The name of the host on which the cluster instance is running.
	HostName *string `mandatory:"false" json:"hostName"`

	// The role of the cluster node.
	NodeRole ExternalClusterInstanceNodeRoleEnum `mandatory:"false" json:"nodeRole,omitempty"`

	// The Oracle base location of Cluster Ready Services (CRS).
	CrsBaseDirectory *string `mandatory:"false" json:"crsBaseDirectory"`

	// The Automatic Diagnostic Repository (ADR) home directory for the cluster instance.
	AdrHomeDirectory *string `mandatory:"false" json:"adrHomeDirectory"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The date and time the external cluster instance was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The date and time the external cluster instance was last updated.
	TimeUpdated *common.SDKTime `mandatory:"false" json:"timeUpdated"`
}

func (m ExternalClusterInstanceSummary) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ExternalClusterInstanceSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingExternalClusterInstanceLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetExternalClusterInstanceLifecycleStateEnumStringValues(), ",")))
	}

	if _, ok := GetMappingExternalClusterInstanceNodeRoleEnum(string(m.NodeRole)); !ok && m.NodeRole != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for NodeRole: %s. Supported values are: %s.", m.NodeRole, strings.Join(GetExternalClusterInstanceNodeRoleEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}
