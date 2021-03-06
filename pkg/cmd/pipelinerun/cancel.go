// Copyright © 2019 The Tekton Authors.
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

package pipelinerun

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/formatted"
	"github.com/tektoncd/cli/pkg/pipelinerun"
	"github.com/tektoncd/cli/pkg/validate"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	succeeded   = formatted.ColorStatus("Succeeded")
	failed      = formatted.ColorStatus("Failed")
	prCancelled = formatted.ColorStatus("Cancelled") + "(PipelineRunCancelled)"
)

func cancelCommand(p cli.Params) *cobra.Command {
	eg := `Cancel the PipelineRun named 'foo' from namespace 'bar':

    tkn pipelinerun cancel foo -n bar
`

	c := &cobra.Command{
		Use:          "cancel",
		Short:        "Cancel a PipelineRun in a namespace",
		Example:      eg,
		SilenceUsage: true,
		Annotations: map[string]string{
			"commandType": "main",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pr := args[0]

			s := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}

			if err := validate.NamespaceExists(p); err != nil {
				return err
			}

			return cancelPipelineRun(p, s, pr)
		},
	}

	_ = c.MarkZshCompPositionalArgumentCustom(1, "__tkn_get_pipelinerun")
	return c
}

func cancelPipelineRun(p cli.Params, s *cli.Stream, prName string) error {
	cs, err := p.Clients()
	if err != nil {
		return fmt.Errorf("failed to create tekton client")
	}

	pr, err := pipelinerun.GetV1beta1(cs, prName, metav1.GetOptions{}, p.Namespace())
	if err != nil {
		return fmt.Errorf("failed to find PipelineRun: %s", prName)
	}

	prCond := formatted.Condition(pr.Status.Conditions)
	if prCond == succeeded || prCond == failed || prCond == prCancelled {
		return fmt.Errorf("failed to cancel PipelineRun %s: PipelineRun has already finished execution", prName)
	}

	if _, err = pipelinerun.Patch(cs, prName, metav1.PatchOptions{}, p.Namespace()); err != nil {
		return fmt.Errorf("failed to cancel PipelineRun: %s: %v", prName, err)

	}

	fmt.Fprintf(s.Out, "PipelineRun cancelled: %s\n", pr.Name)
	return nil
}
