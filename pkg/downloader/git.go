package downloader

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/olegsu/kubectl-fetch-yaml/pkg/logger"
	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
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

const yamlSeperator = "---"

type (
	git struct {
		target io.Writer
		singer ssh.Signer
		user   string
		repo   string
		branch string
		path   string
		logger logger.Logger
		token  string
	}
)

func (g git) Download(context context.Context, dest io.Writer) error {
	cloneOptions := &gogit.CloneOptions{
		URL:           g.repo,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(g.branch),
		Auth:          g.buildAuthentication(),
	}
	g.logger.Debug("Cloning", "repo", g.repo, "branch", g.branch)
	m := memory.NewStorage()
	repo, err := gogit.CloneContext(context, m, nil, cloneOptions)
	if err != nil {
		return fmt.Errorf("Failed to clone repo: %w", err)
	}
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("Failed to get HEAD: %w", err)
	}
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return fmt.Errorf("Failed create commit iterator: %w", err)
	}
	commitIter.ForEach(func(c *object.Commit) error {
		if c.Hash != head.Hash() {
			return nil
		}
		g.logger.Debug("Iterating over commit tree", "committer", c.Committer.String(), "message", c.Message)

		filesIter, err := c.Files()
		if err != nil {
			return fmt.Errorf("Failed commit file iterator: %w", err)
		}
		return filesIter.ForEach(func(f *object.File) error {
			if !strings.HasPrefix(f.Name, g.path) {
				return nil
			}
			str, err := f.Contents()
			if err != nil {
				return err
			}

			fmt.Fprintf(g.target, "\n%s\n", str)
			if path.Ext(f.Name) == "yaml" || path.Ext(f.Name) == "yml" {
				fmt.Fprintf(g.target, "\n%s\n", yamlSeperator)
			}
			return nil
		})
	})

	return nil
}

func (g git) buildAuthentication() transport.AuthMethod {
	if g.singer != nil {
		g.logger.Debug("Authentication ssh key")
		auth := &gitssh.PublicKeys{
			User:   g.user,
			Signer: g.singer,
		}
		auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		return auth
	}
	if g.token != "" {
		g.logger.Debug("Authentication with user and token", "user", g.user, "token", g.token)
		auth := &http.BasicAuth{
			Username: g.user,
			Password: g.token,
		}
		return auth
	}

	return nil
}
