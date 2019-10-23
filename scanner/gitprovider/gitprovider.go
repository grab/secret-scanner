/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

// GitProvider defines interface for interacting with remote Git services
type GitProvider interface {
	Initialize(baseURL, token string, additionalParams map[string]string) error
	GetAdditionalParam(key string) string
	ValidateAdditionalParams(additionalParams map[string]string) bool
	GetRepository(opt map[string]string) (*Repository, error)
	Name() string
}
