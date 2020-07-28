package cmd

// Copyright © 2020 oleg2807@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"github.com/spf13/cobra"
)

const description = `
Download Kubernetes manifests from remote destinations to be applied via kubecl.`
const example = `git: kubectl fetch git --repo https://github.com/kubernetes/examples --path guestbook/all-in-one | kubectl apply -f -`

var version string

var rootCmdOptions struct {
	Verbose bool
}

var rootCmd = &cobra.Command{
	Use:     "kubectl-fetch",
	Long:    description,
	Example: example,
	Version: version,
}

// Execute - execute the root command
func Execute() {
	err := rootCmd.Execute()
	dieOnError("", err)
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.Verbose, "verbose", false, "Set to get more detailed output")
}
