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

// ExternalDbHomeSummary The summary of an external database home.
type ExternalDbHomeSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external DB home.
	Id *string `mandatory:"true" json:"id"`

	// The user-friendly name for the external DB home. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the external DB system that the DB home is a part of.
	ExternalDbSystemId *string `mandatory:"true" json:"externalDbSystemId"`

	// The current lifecycle state of the external DB home.
	LifecycleState ExternalDbHomeLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the external DB home was created.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The date and time the external DB home was last updated.
	TimeUpdated *common.SDKTime `mandatory:"true" json:"timeUpdated"`

	// The name of the external DB home.
	ComponentName *string `mandatory:"false" json:"componentName"`

	// The location of the DB home.
	HomeDirectory *string `mandatory:"false" json:"homeDirectory"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`
}

func (m ExternalDbHomeSummary) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ExternalDbHomeSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingExternalDbHomeLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetExternalDbHomeLifecycleStateEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}
