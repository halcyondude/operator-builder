// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CollectionSetNames(t *testing.T) {
	t.Parallel()

	sharedNameInput := WorkloadShared{
		Name: "shared-name",
		Kind: WorkloadKindCollection,
	}

	sharedNameExpected := WorkloadShared{
		Name:        "shared-name",
		PackageName: "sharedname",
		Kind:        WorkloadKindCollection,
	}

	for _, tt := range []struct {
		name     string
		input    *WorkloadCollection
		expected *WorkloadCollection
	}{
		{
			name: "workload collection missing root command",
			input: &WorkloadCollection{
				WorkloadShared: sharedNameInput,
			},
			expected: &WorkloadCollection{
				WorkloadShared: sharedNameExpected,
				Spec: WorkloadCollectionSpec{
					CompanionCliRootcmd: CliCommand{},
				},
			},
		},
		{
			name: "workload collection with root command missing description",
			input: &WorkloadCollection{
				WorkloadShared: sharedNameInput,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name: "hasrootcommand",
					},
				},
			},
			expected: &WorkloadCollection{
				WorkloadShared: sharedNameExpected,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage workloadcollectiontest collection and components",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
					CompanionCliSubcmd: CliCommand{
						Name:         "collection",
						Description:  "Manage workloadcollectiontest workload",
						VarName:      "Collection",
						FileName:     "collection",
						IsSubcommand: true,
					},
				},
			},
		},
		{
			name: "workload collection with root command",
			input: &WorkloadCollection{
				WorkloadShared: sharedNameInput,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:        "hasrootcommand",
						Description: "Manage workloadcollectiontest collection and components custom",
					},
				},
			},
			expected: &WorkloadCollection{
				WorkloadShared: sharedNameExpected,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage workloadcollectiontest collection and components custom",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
					CompanionCliSubcmd: CliCommand{
						Name:         "collection",
						Description:  "Manage workloadcollectiontest workload",
						VarName:      "Collection",
						FileName:     "collection",
						IsSubcommand: true,
					},
				},
			},
		},
		{
			name: "workload collection with full subcommand",
			input: &WorkloadCollection{
				WorkloadShared: sharedNameInput,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name: "hasrootcommand",
					},
					CompanionCliSubcmd: CliCommand{
						Name:        "collection",
						Description: "Manage workloadcollectiontest workload custom",
						VarName:     "Collection",
						FileName:    "collection",
					},
				},
			},
			expected: &WorkloadCollection{
				WorkloadShared: sharedNameExpected,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage workloadcollectiontest collection and components",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
					CompanionCliSubcmd: CliCommand{
						Name:         "collection",
						Description:  "Manage workloadcollectiontest workload custom",
						VarName:      "Collection",
						FileName:     "collection",
						IsSubcommand: true,
					},
				},
			},
		},
		{
			name: "workload collection with full subcommand missing description",
			input: &WorkloadCollection{
				WorkloadShared: sharedNameInput,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name: "hasrootcommand",
					},
					CompanionCliSubcmd: CliCommand{
						Name:     "collection",
						VarName:  "Collection",
						FileName: "collection",
					},
				},
			},
			expected: &WorkloadCollection{
				WorkloadShared: sharedNameExpected,
				Spec: WorkloadCollectionSpec{
					API: WorkloadAPISpec{
						Kind: "WorkloadCollectionTest",
					},
					CompanionCliRootcmd: CliCommand{
						Name:          "hasrootcommand",
						Description:   "Manage workloadcollectiontest collection and components",
						VarName:       "Hasrootcommand",
						FileName:      "hasrootcommand",
						IsRootcommand: true,
					},
					CompanionCliSubcmd: CliCommand{
						Name:         "collection",
						Description:  "Manage workloadcollectiontest workload",
						VarName:      "Collection",
						FileName:     "collection",
						IsSubcommand: true,
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
