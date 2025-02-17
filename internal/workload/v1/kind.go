// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package v1

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var ErrInvalidKind = errors.New("unrecognized workload kind in workload config")

// WorkloadKind indicates which of the supported workload kinds are being used.
type WorkloadKind int32

const (
	WorkloadKindUnknown WorkloadKind = iota
	WorkloadKindStandalone
	WorkloadKindCollection
	WorkloadKindComponent
)

func (wk WorkloadKind) String() string {
	kinds := map[WorkloadKind]string{
		WorkloadKindStandalone: "StandaloneWorkload",
		WorkloadKindCollection: "WorkloadCollection",
		WorkloadKindComponent:  "ComponentWorkload",
		WorkloadKindUnknown:    "Unknown Workload Type",
	}

	return kinds[wk]
}

func (wk *WorkloadKind) UnmarshalJSON(data []byte) error {
	kind, ok := workloadKindsMap()[string(data)]
	if !ok {
		return ErrInvalidKind
	}

	*wk = kind

	return nil
}

func (wk *WorkloadKind) UnmarshalYAML(node *yaml.Node) error {
	kind, ok := workloadKindsMap()[node.Value]
	if !ok {
		return ErrInvalidKind
	}

	*wk = kind

	return nil
}

func workloadKindsMap() map[string]WorkloadKind {
	return map[string]WorkloadKind{
		"StandaloneWorkload": WorkloadKindStandalone,
		"WorkloadCollection": WorkloadKindCollection,
		"ComponentWorkload":  WorkloadKindComponent,
	}
}

func decodeStandalone(dc *yaml.Decoder) (*StandaloneWorkload, error) {
	v := &StandaloneWorkload{}
	if err := dc.Decode(v); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return v, nil
}

func decodeCollection(dc *yaml.Decoder) (*WorkloadCollection, error) {
	v := &WorkloadCollection{}
	if err := dc.Decode(v); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return v, nil
}

func decodeComponent(dc *yaml.Decoder) (*ComponentWorkload, error) {
	v := &ComponentWorkload{}
	if err := dc.Decode(v); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return v, nil
}
