/*
Copyright 2020 The Kubernetes Authors.

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
	"os"

	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // for rest client metric registration
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	_ "github.com/tanjunchen/tanjunchen-scheduler/apis/config/scheme"
	"github.com/tanjunchen/tanjunchen-scheduler/pkg/dynamic"
	"github.com/tanjunchen/tanjunchen-scheduler/pkg/example"
	"github.com/tanjunchen/tanjunchen-scheduler/pkg/names"
)

func main() {
	// Register custom plugins to the scheduler framework.
	// Later they can consist of scheduler profile(s) and hence
	// used by various kinds of workloads.
	command := app.NewSchedulerCommand(
		app.WithPlugin(names.DynamicName, dynamic.NewDynamicPlugin),
		app.WithPlugin(names.ExampleName, example.NewExamplePlugin),
	)

	code := cli.Run(command)
	os.Exit(code)
}
