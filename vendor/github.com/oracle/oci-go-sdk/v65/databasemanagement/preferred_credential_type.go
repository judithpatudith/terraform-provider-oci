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

// PreferredCredentialTypeEnum Enum with underlying type: string
type PreferredCredentialTypeEnum string

// Set of constants representing the allowable values for PreferredCredentialTypeEnum
const (
	PreferredCredentialTypeBasic PreferredCredentialTypeEnum = "BASIC"
)

var mappingPreferredCredentialTypeEnum = map[string]PreferredCredentialTypeEnum{
	"BASIC": PreferredCredentialTypeBasic,
}

var mappingPreferredCredentialTypeEnumLowerCase = map[string]PreferredCredentialTypeEnum{
	"basic": PreferredCredentialTypeBasic,
}

// GetPreferredCredentialTypeEnumValues Enumerates the set of values for PreferredCredentialTypeEnum
func GetPreferredCredentialTypeEnumValues() []PreferredCredentialTypeEnum {
	values := make([]PreferredCredentialTypeEnum, 0)
	for _, v := range mappingPreferredCredentialTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetPreferredCredentialTypeEnumStringValues Enumerates the set of values in String for PreferredCredentialTypeEnum
func GetPreferredCredentialTypeEnumStringValues() []string {
	return []string{
		"BASIC",
	}
}

// GetMappingPreferredCredentialTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingPreferredCredentialTypeEnum(val string) (PreferredCredentialTypeEnum, bool) {
	enum, ok := mappingPreferredCredentialTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
