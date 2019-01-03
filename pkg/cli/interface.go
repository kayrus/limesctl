/*******************************************************************************
*
* Copyright 2018 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package cli

import (
	"os"

	"github.com/sapcc/gophercloud-limes/resources/v1/clusters"
	"github.com/sapcc/gophercloud-limes/resources/v1/domains"
	"github.com/sapcc/gophercloud-limes/resources/v1/projects"
	"github.com/sapcc/limes"
	"github.com/sapcc/limesctl/pkg/errors"
)

// Cluster contains information regarding a cluster(s).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Cluster.
type Cluster struct {
	ID     string
	Result clusters.CommonResult
	IsList bool
	Filter Filter
	Output Output
}

// Domain contains information regarding a domain(s).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Domain.
type Domain struct {
	ID     string
	Name   string
	Result domains.CommonResult
	IsList bool
	Filter Filter
	Output Output
}

// Project contains information regarding a project(s).
// As different methods are called on it, the fields within the structure are updated accordingly.
// Call its appropriate method to get/list/update a Project.
type Project struct {
	ID         string
	Name       string
	DomainID   string
	DomainName string
	Result     projects.CommonResult
	IsList     bool
	Filter     Filter
	Output     Output
}

// Filter contains different parameters for filtering a get/list/update operation.
type Filter struct {
	Cluster  string
	Area     string
	Service  string
	Resource string
}

// Output contains different options that affect the output of a get/list operation.
type Output struct {
	Names         bool
	Long          bool
	HumanReadable bool
}

// Renderer interface type contains different methods for rendering data in
// different formats.
type Renderer interface {
	renderJSON() *jsonData
	renderCSV() *csvData
}

// GetTask is the interface type that abstracts a get operation.
type GetTask interface {
	get()
	Renderer
}

// RunGetTask is the function that operates on a GetTask and shows the output in the respective
// format that is specified at the command line.
func RunGetTask(t GetTask, outputFmt string) {
	t.get()
	switch outputFmt {
	case "json":
		t.renderJSON().write(os.Stdout)
	case "csv":
		t.renderCSV().write(os.Stdout)
	default:
		t.renderCSV().writeTable(os.Stdout)
	}
}

// ListTask is the interface type that abstracts a list operation.
type ListTask interface {
	list()
	Renderer
}

// RunListTask is the function that operates on a ListTask and shows the output in the respective
// format that is specified at the command line.
func RunListTask(t ListTask, outputFmt string) {
	t.list()
	switch outputFmt {
	case "json":
		t.renderJSON().write(os.Stdout)
	case "csv":
		t.renderCSV().write(os.Stdout)
	default:
		t.renderCSV().writeTable(os.Stdout)
	}
}

// Resource contains quota information about a single resource.
type Resource struct {
	Name    string
	Value   int64
	Unit    limes.Unit
	Comment string
}

// Quotas is a map of service name to a list of resources. It contains the aggregate
// quota values used by the set methods to update a single cluster/domain/project.
type Quotas map[string][]Resource

// SetTask is the interface type that abstracts a put operation.
type SetTask interface {
	set(*Quotas)
}

// RunSetTask is the function that operates on a SetTask and shows the output in the respective
// format that is specified at the command line.
func RunSetTask(t SetTask, q *Quotas) {
	t.set(q)
}

// RunSyncTask schedules a sync job that pulls quota and usage data for a project from
// the backing services into Limes' local database.
func RunSyncTask(p *Project) {
	_, limesV1 := getServiceClients()

	err := projects.Sync(limesV1, p.DomainID, p.ID, projects.SyncOpts{
		Cluster: p.Filter.Cluster,
	})
	errors.Handle(err, "could not sync project")
}
