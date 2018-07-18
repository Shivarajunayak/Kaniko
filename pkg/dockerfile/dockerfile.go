/*
Copyright 2018 Google LLC

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

package dockerfile

import (
	"bytes"
	"github.com/GoogleContainerTools/kaniko/pkg/constants"
	"github.com/GoogleContainerTools/kaniko/pkg/util"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// Parse parses the contents of a Dockerfile and returns a list of commands
func Parse(b []byte) ([]instructions.Stage, error) {
	p, err := parser.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	stages, _, err := instructions.Parse(p.AST)
	if err != nil {
		return nil, err
	}
	return stages, err
}

// ResolveStages resolves any calls to previous stages with names to indices
// Ex. --from=second_stage should be --from=1 for easier processing later on
func ResolveStages(stages []instructions.Stage) {
	nameToIndex := make(map[string]string)
	for i, stage := range stages {
		index := strconv.Itoa(i)
		if stage.Name != index {
			nameToIndex[stage.Name] = index
		}
		for _, cmd := range stage.Commands {
			switch c := cmd.(type) {
			case *instructions.CopyCommand:
				if c.From != "" {
					if val, ok := nameToIndex[c.From]; ok {
						c.From = val
					}
				}
			}
		}
	}
}

// ParseCommands parses an array of commands into an array of instructions.Command; used for onbuild
func ParseCommands(cmdArray []string) ([]instructions.Command, error) {
	var cmds []instructions.Command
	cmdString := strings.Join(cmdArray, "\n")
	ast, err := parser.Parse(strings.NewReader(cmdString))
	if err != nil {
		return nil, err
	}
	for _, child := range ast.AST.Children {
		cmd, err := instructions.ParseCommand(child)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

// Dependencies returns a list of files in this stage that will be needed in later stages
func Dependencies(index int, stages []instructions.Stage, image v1.Image, buildArgs *BuildArgs) ([]string, error) {
	var dependencies []string
	for stageIndex, stage := range stages {
		if stageIndex <= index {
			continue
		}
		var sourceImage v1.Image
		logrus.Infof("Comparing stage basename %s with stage name %s", stage.BaseName, stages[index].Name)
		if stage.BaseName == constants.NoBaseImage {
			sourceImage = empty.Image
		} else if stage.BaseName == stages[index].Name {
			sourceImage = image
		} else if util.FilepathExists(filepath.Join(constants.KanikoIntermediateStagesDir, stage.BaseName)) {
			var err error
			sourceImage, err = tarball.ImageFromPath(filepath.Join(filepath.Join(constants.KanikoIntermediateStagesDir, stage.BaseName), constants.StageTar), nil)
			if err != nil {
				return nil, err
			}
		} else {
			// Initialize source image
			logrus.Infof("trying to intiialize %s", stage.BaseName)
			ref, err := name.ParseReference(stage.BaseName, name.WeakValidation)
			if err != nil {
				return nil, err

			}
			auth, err := authn.DefaultKeychain.Resolve(ref.Context().Registry)
			if err != nil {
				return nil, err
			}
			sourceImage, err = remote.Image(ref, remote.WithAuth(auth), remote.WithTransport(http.DefaultTransport))
			if err != nil {
				return nil, err
			}
		}
		imageConfig, err := sourceImage.ConfigFile()
		if err != nil {
			return nil, err
		}
		for _, cmd := range stage.Commands {
			switch c := cmd.(type) {
			case *instructions.EnvCommand:
				replacementEnvs := buildArgs.ReplacementEnvs(imageConfig.Config.Env)
				if err := util.UpdateConfigEnv(c.Env, &imageConfig.Config, replacementEnvs); err != nil {
					return nil, err
				}
			case *instructions.ArgCommand:
				buildArgs.AddArg(c.Key, c.Value)
			case *instructions.CopyCommand:
				if c.From != strconv.Itoa(index) {
					continue
				}
				// First, resolve any environment replacement
				replacementEnvs := buildArgs.ReplacementEnvs(imageConfig.Config.Env)
				resolvedEnvs, err := util.ResolveEnvironmentReplacementList(c.SourcesAndDest, replacementEnvs, true)
				if err != nil {
					return nil, err
				}
				// Resolve wildcards and get a list of resolved sources
				srcs, err := util.ResolveSources(resolvedEnvs, constants.RootDir)
				if err != nil {
					return nil, err
				}
				for index, src := range srcs {
					if !filepath.IsAbs(src) {
						srcs[index] = filepath.Join(constants.RootDir, src)
					}
				}
				dependencies = append(dependencies, srcs...)
			}
		}
	}
	return dependencies, nil
}
