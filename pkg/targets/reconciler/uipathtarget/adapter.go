/*
Copyright 2022 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package uipathtarget

import (
	corev1 "k8s.io/api/core/v1"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/apis"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	commonv1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/common/v1alpha1"
	"github.com/triggermesh/triggermesh/pkg/apis/targets/v1alpha1"
	common "github.com/triggermesh/triggermesh/pkg/reconciler"
	"github.com/triggermesh/triggermesh/pkg/reconciler/resource"
)

// adapterConfig contains properties used to configure the target's adapter.
// Public fields are automatically populated by envconfig.
type adapterConfig struct {
	// Configuration accessor for logging/metrics/tracing
	obsConfig source.ConfigAccessor
	// Container image
	Image string `default:"gcr.io/triggermesh/uipathtarget-adapter"`
}

// Verify that Reconciler implements common.AdapterBuilder.
var _ common.AdapterBuilder[*servingv1.Service] = (*Reconciler)(nil)

// BuildAdapter implements common.AdapterBuilder.
func (r *Reconciler) BuildAdapter(trg commonv1alpha1.Reconcilable, _ *apis.URL) (*servingv1.Service, error) {
	typedTrg := trg.(*v1alpha1.UiPathTarget)

	return common.NewAdapterKnService(trg, nil,
		resource.Image(r.adapterCfg.Image),
		resource.EnvVars(makeAppEnv(typedTrg)...),
		resource.EnvVars(r.adapterCfg.obsConfig.ToEnvVars()...),
	), nil
}

func makeAppEnv(o *v1alpha1.UiPathTarget) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "UIPATH_ROBOT_NAME",
			Value: o.Spec.RobotName,
		}, {
			Name: "UIPATH_USER_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: o.Spec.UserKey.SecretKeyRef,
			},
		}, {
			Name:  "UIPATH_PROCESS_NAME",
			Value: o.Spec.ProcessName,
		}, {
			Name:  "UIPATH_TENANT_NAME",
			Value: o.Spec.TenantName,
		}, {
			Name:  "UIPATH_ACCOUNT_LOGICAL_NAME",
			Value: o.Spec.AccountLogicalName,
		},
		{
			Name:  "UIPATH_CLIENT_ID",
			Value: o.Spec.ClientID,
		},
		{
			Name:  "UIPATH_ORGANIZATION_UNIT_ID",
			Value: o.Spec.OrganizationUnitID,
		},
	}
}
