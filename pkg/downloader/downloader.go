package downloader

import (
	"context"
	"io"

	"github.com/olegsu/kubectl-fetch/pkg/logger"
	"golang.org/x/crypto/ssh"
)

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

type (
	// Downloader downloads files
	Downloader interface {
		Download(context.Context, io.Writer) error
	}

	GitOptions struct {
		Singer ssh.Signer
		User   string
		Repo   string
		Branch string
		Path   string
		Logger logger.Logger
		Target io.Writer
		Token  string
	}
)

// NewGitDownloader download files from git repos
func NewGitDownloader(opt GitOptions) Downloader {
	return &git{
		singer: opt.Singer,
		user:   opt.User,
		branch: opt.Branch,
		path:   opt.Path,
		repo:   opt.Repo,
		logger: opt.Logger,
		target: opt.Target,
		token:  opt.Token,
	}
}
