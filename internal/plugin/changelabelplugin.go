/*
Copyright 2018, 2019 the Velero contributors.

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

package plugin

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// ChangeLabelPlugin is a restore item action plugin for Velero
type ChangeLabelPlugin struct {
	log             logrus.FieldLogger
	configMapClient corev1client.ConfigMapInterface
}

// NewChangeLabelPlugin instantiates a ChangeLabelPlugin.
func NewChangeLabelPlugin(log logrus.FieldLogger, configMapClient corev1client.ConfigMapInterface) *ChangeLabelPlugin {
	return &ChangeLabelPlugin{
		log:             log,
		configMapClient: configMapClient,
	}
}

// AppliesTo returns information about which resources this action should be invoked for.
// A RestoreItemAction's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.g
func (p *ChangeLabelPlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		//IncludedResources: []string{"persistentvolumeclaims", "persistentvolumes"},
	}, nil
}

// Execute allows the ChangeLabelPlugin to perform arbitrary logic with the item being restored,
// in this case, setting a custom annotation on the item being restored.
func (p *ChangeLabelPlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.log.Info("Executing ChangeLabelPlugin!")

	metadata, err := meta.Accessor(input.Item)
	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}
	labels := metadata.GetLabels()

	p.log.Debug("Getting plugin config")
	config, err := getPluginConfig(framework.PluginKindRestoreItemAction, "velero.io/change-label", p.configMapClient)
	if err != nil {
		return nil, err
	}

	p.log.Debug("configmap", config)

	if config.Data != nil {
		for label, value := range config.Data {
			p.log.Debug("configmap data - ", label, ": ", value)

			existingLabel, ok := labels[label]
			if ok {
				p.log.Debug("Existing value for label: ", existingLabel)
			}
			labels[label] = value
		}
		metadata.SetLabels(labels)
	} else {
		p.log.Info("No configurations found in configmap: ", config.Name)
	}
	return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
}

func getPluginConfig(kind framework.PluginKind, name string, client corev1client.ConfigMapInterface) (*corev1.ConfigMap, error) {
	opts := metav1.ListOptions{
		// velero.io/plugin-config: true
		// velero.io/restic: RestoreItemAction
		LabelSelector: fmt.Sprintf("velero.io/plugin-config,%s=%s", name, kind),
	}

	list, err := client.List(opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(list.Items) == 0 {
		return nil, nil
	}

	if len(list.Items) > 1 {
		var items []string
		for _, item := range list.Items {
			items = append(items, item.Name)
		}
		return nil, errors.Errorf("found more than one ConfigMap matching label selector %q: %v", opts.LabelSelector, items)
	}

	return &list.Items[0], nil
}
