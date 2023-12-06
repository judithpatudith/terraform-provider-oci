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
	"strings"
)

// RoleEnum Enum with underlying type: string
type RoleEnum string

// Set of constants representing the allowable values for RoleEnum
const (
	RoleNormal RoleEnum = "NORMAL"
	RoleSysdba RoleEnum = "SYSDBA"
)

var mappingRoleEnum = map[string]RoleEnum{
	"NORMAL": RoleNormal,
	"SYSDBA": RoleSysdba,
}

var mappingRoleEnumLowerCase = map[string]RoleEnum{
	"normal": RoleNormal,
	"sysdba": RoleSysdba,
}

// GetRoleEnumValues Enumerates the set of values for RoleEnum
func GetRoleEnumValues() []RoleEnum {
	values := make([]RoleEnum, 0)
	for _, v := range mappingRoleEnum {
		values = append(values, v)
	}
	return values
}

// GetRoleEnumStringValues Enumerates the set of values in String for RoleEnum
func GetRoleEnumStringValues() []string {
	return []string{
		"NORMAL",
		"SYSDBA",
	}
}

// GetMappingRoleEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingRoleEnum(val string) (RoleEnum, bool) {
	enum, ok := mappingRoleEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
