// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package packageclient

import (
	esVersion "github.com/vmware-tanzu/tanzu-framework/apis/externalsecrets/v1beta1"

	"github.com/vmware-tanzu/tanzu-framework/packageclients/pkg/packagedatamodel"
)

func (p *pkgClient) ListExternalSecrets(o *packagedatamodel.ExternalSecretOptions) (*esVersion.ExternalSecretList, error) {
	// list := &esVersion.ExternalSecretList{}

	list, err := p.kappClient.ListExternalSecrets(o.Namespace)

	return list, err
}
