// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package v1

import (
	"errors"
	"fmt"

	"sigs.k8s.io/kubebuilder/v3/pkg/model/resource"

	"github.com/vmware-tanzu-labs/operator-builder/internal/utils"
)

const (
	defaultCollectionSubcommandName         = "collection"
	defaultCollectionSubcommandDescription  = `Manage %s workload`
	defaultCollectionRootcommandDescription = `Manage %s collection and components`
)

var ErrMissingRequiredFields = errors.New("missing required fields")

// WorkloadCollectionSpec defines the attributes for a workload collection.
type WorkloadCollectionSpec struct {
	API                 WorkloadAPISpec      `json:"api" yaml:"api"`
	CompanionCliRootcmd CliCommand           `json:"companionCliRootcmd,omitempty" yaml:"companionCliRootcmd,omitempty" validate:"omitempty"`
	CompanionCliSubcmd  CliCommand           `json:"companionCliSubcmd,omitempty" yaml:"companionCliSubcmd,omitempty" validate:"omitempty"`
	ComponentFiles      []string             `json:"componentFiles" yaml:"componentFiles"`
	Components          []*ComponentWorkload `json:",omitempty" yaml:",omitempty" validate:"omitempty"`
	WorkloadSpec        `yaml:",inline"`
}

// WorkloadCollection defines a workload collection.
type WorkloadCollection struct {
	WorkloadShared `yaml:",inline"`
	Spec           WorkloadCollectionSpec `json:"spec" yaml:"spec" validate:"required"`
}

func NewWorkloadCollection(
	name string,
	spec WorkloadAPISpec,
	componentFiles []string,
) *WorkloadCollection {
	return &WorkloadCollection{
		WorkloadShared: WorkloadShared{
			Kind: WorkloadKindCollection,
			Name: name,
		},
		Spec: WorkloadCollectionSpec{
			API:            spec,
			ComponentFiles: componentFiles,
		},
	}
}

func (c *WorkloadCollection) Validate() error {
	missingFields := []string{}

	// required fields
	if c.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if c.Spec.API.Domain == "" {
		missingFields = append(missingFields, "spec.api.domain")
	}

	if c.Spec.API.Group == "" {
		missingFields = append(missingFields, "spec.api.group")
	}

	if c.Spec.API.Version == "" {
		missingFields = append(missingFields, "spec.api.version")
	}

	if c.Spec.API.Kind == "" {
		missingFields = append(missingFields, "spec.api.kind")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("%w: %s", ErrMissingRequiredFields, missingFields)
	}

	return nil
}

func (c *WorkloadCollection) GetWorkloadKind() WorkloadKind {
	return c.Kind
}

// methods that implement WorkloadInitializer.
func (c *WorkloadCollection) GetDomain() string {
	return c.Spec.API.Domain
}

func (c *WorkloadCollection) HasRootCmdName() bool {
	return c.Spec.CompanionCliRootcmd.hasName()
}

func (c *WorkloadCollection) HasRootCmdDescription() bool {
	return c.Spec.CompanionCliRootcmd.hasDescription()
}

func (c *WorkloadCollection) HasSubCmdName() bool {
	return c.Spec.CompanionCliSubcmd.hasName()
}

func (c *WorkloadCollection) HasSubCmdDescription() bool {
	return c.Spec.CompanionCliSubcmd.hasDescription()
}

// methods that implement WorkloadAPIBuilder.
func (c *WorkloadCollection) GetName() string {
	return c.Name
}

func (c *WorkloadCollection) GetPackageName() string {
	return c.PackageName
}

func (c *WorkloadCollection) GetAPIGroup() string {
	return c.Spec.API.Group
}

func (c *WorkloadCollection) GetAPIVersion() string {
	return c.Spec.API.Version
}

func (c *WorkloadCollection) GetAPIKind() string {
	return c.Spec.API.Kind
}

func (c *WorkloadCollection) IsClusterScoped() bool {
	return c.Spec.API.ClusterScoped
}

func (c *WorkloadCollection) IsStandalone() bool {
	return false
}

func (c *WorkloadCollection) IsComponent() bool {
	return false
}

func (c *WorkloadCollection) IsCollection() bool {
	return true
}

func (c *WorkloadCollection) SetResources(workloadPath string) error {
	err := c.Spec.processManifests(FieldMarkerType, CollectionMarkerType)
	if err != nil {
		return err
	}

	for _, cpt := range c.Spec.Components {
		for _, csr := range cpt.Spec.Resources {
			// add to spec fields if not present
			err := c.Spec.processMarkers(csr, CollectionMarkerType)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *WorkloadCollection) GetDependencies() []*ComponentWorkload {
	return []*ComponentWorkload{}
}

func (c *WorkloadCollection) SetComponents(components []*ComponentWorkload) error {
	c.Spec.Components = components

	return nil
}

func (c *WorkloadCollection) HasChildResources() bool {
	return len(c.Spec.Resources) > 0
}

func (c *WorkloadCollection) GetCollection() *WorkloadCollection {
	return c.Spec.Collection
}

func (c *WorkloadCollection) GetComponents() []*ComponentWorkload {
	return c.Spec.Components
}

func (c *WorkloadCollection) GetSourceFiles() *[]SourceFile {
	return c.Spec.SourceFiles
}

func (c *WorkloadCollection) GetFuncNames() (createFuncNames, initFuncNames []string) {
	return getFuncNames(*c.GetSourceFiles())
}

func (c *WorkloadCollection) GetAPISpecFields() *APIFields {
	return c.Spec.APISpecFields
}

func (c *WorkloadCollection) GetRBACRules() *[]RBACRule {
	var rules []RBACRule = *c.Spec.RBACRules

	return &rules
}

func (*WorkloadCollection) GetOwnershipRules() *[]OwnershipRule {
	return &[]OwnershipRule{}
}

func (c *WorkloadCollection) GetComponentResource(domain, repo string, clusterScoped bool) *resource.Resource {
	var namespaced bool
	if clusterScoped {
		namespaced = false
	} else {
		namespaced = true
	}

	api := resource.API{
		CRDVersion: "v1",
		Namespaced: namespaced,
	}

	return &resource.Resource{
		GVK: resource.GVK{
			Domain:  domain,
			Group:   c.Spec.API.Group,
			Version: c.Spec.API.Version,
			Kind:    c.Spec.API.Kind,
		},
		Plural: resource.RegularPlural(c.Spec.API.Kind),
		Path: fmt.Sprintf(
			"%s/apis/%s/%s",
			repo,
			c.Spec.API.Group,
			c.Spec.API.Version,
		),
		API:        &api,
		Controller: true,
	}
}

func (c *WorkloadCollection) SetNames() {
	c.PackageName = utils.ToPackageName(c.Name)

	// only set the names if we have specified the root command name else none
	// of the following values will matter as the code for the cli will not be
	// generated
	if c.HasRootCmdName() {
		// set the root command values
		c.Spec.CompanionCliRootcmd.setCommonValues(c, false)

		// set the subcommand values
		c.Spec.CompanionCliSubcmd.setCommonValues(c, true)
	}
}

func (c *WorkloadCollection) GetRootCommand() *CliCommand {
	return &c.Spec.CompanionCliRootcmd
}

func (c *WorkloadCollection) GetSubCommand() *CliCommand {
	return &c.Spec.CompanionCliSubcmd
}

func (c *WorkloadCollection) LoadManifests(workloadPath string) error {
	resources, err := expandResources(workloadPath, c.Spec.Resources)
	if err != nil {
		return err
	}

	c.Spec.Resources = resources
	for _, r := range c.Spec.Resources {
		if err := r.loadContent(c.IsCollection()); err != nil {
			return err
		}
	}

	return nil
}
