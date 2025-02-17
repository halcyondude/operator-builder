// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StandaloneSetNames(t *testing.T) {
	t.Parallel()

	sharedNameInput := WorkloadShared{
		Name: "shared-name",
		Kind: WorkloadKindStandalone,
	}

	sharedNameExpected := WorkloadShared{
		Name:        "shared-name",
		PackageName: "sharedname",
		Kind:        WorkloadKindStandalone,
	}

	for _, tt := range []struct {
		name     string
		input    *StandaloneWorkload
		expected *StandaloneWorkload
	}{
		{
			name: "standalone workload missing root command",
			input: &StandaloneWorkload{
				WorkloadShared: sharedNameInput,
			},
			expected: &StandaloneWorkload{
				WorkloadShared: sharedNameExpected,
				Spec: StandaloneWorkloadSpec{
					CompanionCliRootcmd: CliCommand{},
				},
			},
		},
		{
			name: "standalone workload with root command missing description",
			input: &StandaloneWorkload{
				WorkloadShared: sharedNameInput,
				Spec: StandaloneWorkloadSpec{
					API: WorkloadAPISpec{
						Kind: "StandaloneWorkloadTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name: "hasrootcommand",
					},
				},
			},
			expected: &StandaloneWorkload{
				WorkloadShared: sharedNameExpected,
				Spec: StandaloneWorkloadSpec{
					API: WorkloadAPISpec{
						Kind: "StandaloneWorkloadTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage standaloneworkloadtest workload",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
				},
			},
		},
		{
			name: "standalone workload with root command",
			input: &StandaloneWorkload{
				WorkloadShared: sharedNameInput,
				Spec: StandaloneWorkloadSpec{
					API: WorkloadAPISpec{
						Kind: "StandaloneWorkloadTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:        "hasrootcommand",
						Description: "Manage standaloneworkloadtest workload custom",
					},
				},
			},
			expected: &StandaloneWorkload{
				WorkloadShared: sharedNameExpected,
				Spec: StandaloneWorkloadSpec{
					API: WorkloadAPISpec{
						Kind: "StandaloneWorkloadTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage standaloneworkloadtest workload custom",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
				},
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.input.SetNames()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}
