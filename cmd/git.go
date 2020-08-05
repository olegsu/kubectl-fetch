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
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olegsu/kubectl-fetch-yaml/pkg/downloader"
	"github.com/olegsu/kubectl-fetch-yaml/pkg/logger"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var gitCmdOptions struct {
	sshKeyPath string
	user       string
	repo       string
	branch     string
	revision   string
	path       string
	token      string
	verbose    bool
}

var gitCmd = &cobra.Command{
	Use:  "git",
	Long: "Download file or directory from git repository",
	Run: func(cmd *cobra.Command, args []string) {
		lgr := logger.New(logger.Options{
			Verbose: gitCmdOptions.verbose,
		})
		ctx, cancel := context.WithCancel(context.Background())
		startSignalHandler(lgr, cancel)
		var singer ssh.Signer
		var err error
		if gitCmdOptions.sshKeyPath != "" {
			singer, err = readPrivateKey(gitCmdOptions.sshKeyPath)
			dieOnError("Failed to read key file", err)
		}

		downloader := downloader.NewGitDownloader(downloader.GitOptions{
			Target:   os.Stdout,
			Singer:   singer,
			User:     gitCmdOptions.user,
			Branch:   gitCmdOptions.branch,
			Revision: gitCmdOptions.revision,
			Repo:     gitCmdOptions.repo,
			Path:     gitCmdOptions.path,
			Token:    gitCmdOptions.token,
			Logger:   lgr,
		})

		err = downloader.Download(ctx, os.Stdout)
		dieOnError("Failed to download file", err)
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	gitCmd.Flags().StringVar(&gitCmdOptions.sshKeyPath, "key-file", "", "Path to ssh key")
	gitCmd.Flags().StringVar(&gitCmdOptions.token, "token", "", "Authentication token")
	gitCmd.Flags().StringVar(&gitCmdOptions.user, "user", "git", "User to be used to authenticated")
	gitCmd.Flags().StringVar(&gitCmdOptions.branch, "branch", "master", "Branch to clone")
	gitCmd.Flags().StringVar(&gitCmdOptions.repo, "repo", "", "Repository to clone")
	gitCmd.Flags().StringVar(&gitCmdOptions.revision, "revision", "", "Revision to clone, default is HEAD of branch")
	gitCmd.Flags().StringVar(&gitCmdOptions.path, "path", "", "Path inside the repo")
	gitCmd.Flags().BoolVar(&gitCmdOptions.verbose, "verbose", false, "Print more logs")
}

func readPrivateKey(path string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key: %w", err)
	}
	return signer, nil
}
