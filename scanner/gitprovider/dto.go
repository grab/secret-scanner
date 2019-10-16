/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

// Repository is a universal struct for holding repo info fields
type Repository struct {
	Owner         string
	ID            string
	Name          string
	FullName      string
	CloneURL      string
	URL           string
	DefaultBranch string
	Description   string
	Homepage      string
}
