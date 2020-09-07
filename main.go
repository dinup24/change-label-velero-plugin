/*
Copyright 2017, 2019 the Velero contributors.

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

package main

import (
	"fmt"
	"os"

	"github.com/dinup24/change-label-velero-plugin/internal/plugin"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/client"
	"github.com/vmware-tanzu/velero/pkg/plugin/framework"
)

func main() {
	config, err := client.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error reading config file: %v\n", err)
	}

	f := client.NewFactory("velero", config)

	framework.NewServer().
		RegisterRestoreItemAction("velero.io/change-label", newChangeLabelPlugin(f)).
		Serve()
}

func newChangeLabelPlugin(f client.Factory) framework.HandlerInitializer {
	return func(logger logrus.FieldLogger) (interface{}, error) {
		client, err := f.KubeClient()
		if err != nil {
			return nil, err
		}

		return plugin.NewChangeLabelPlugin(logger, client.CoreV1().ConfigMaps(f.Namespace())), nil
	}
}
