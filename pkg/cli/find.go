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
	"errors"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
)

// FindDomain uses the user's input (name/UUID) to find a specific domain within the token scope.
func FindDomain(userInput string) (*Domain, error) {
	identityV3, _ := getServiceClients()
	d := new(Domain)

	// check if userInput is a UUID
	tmpD, err := domains.Get(identityV3, userInput).Extract()
	if err == nil {
		d.ID = tmpD.ID
		d.Name = tmpD.Name
	} else {
		// userInput appears to be a name so we do domain listing restricted to the name
		page, err := domains.List(identityV3, domains.ListOpts{Name: userInput}).AllPages()
		if err != nil {
			return nil, err
		}
		dList, err := domains.ExtractDomains(page)
		if err != nil {
			return nil, err
		}
		// no need to continue, if there are multiple domains in the list
		if len(dList) > 1 {
			return nil, errors.New("more than one domain exists with the name " + userInput)
		}

		for _, dInList := range dList {
			d.ID = dInList.ID
			d.Name = dInList.Name
		}
	}
	if d.ID == "" {
		return nil, errors.New("domain not found")
	}

	return d, nil
}

// FindProject uses the user's input (name/UUID) to find a specific project within the token scope.
func FindProject(userInputProject, userInputDomain string) (*Project, error) {
	identityV3, _ := getServiceClients()
	p := new(Project)

	// check if userInputProject is a UUID
	tmpP, err := projects.Get(identityV3, userInputProject).Extract()
	if err == nil {
		p.ID = tmpP.ID
		p.Name = tmpP.Name
		p.DomainID = tmpP.DomainID
	} else {
		// userInputProject appears to be a name so we do project listing
		// restricted to the name and domain ID (if given)
		var page pagination.Page
		if userInputDomain != "" {
			d, err := FindDomain(userInputDomain)
			if err != nil {
				return nil, err
			}
			p.DomainName = d.Name

			page, err = projects.List(identityV3, projects.ListOpts{
				Name:     userInputProject,
				DomainID: d.ID,
			}).AllPages()
			if err != nil {
				return nil, err
			}
		} else {
			page, err = projects.List(identityV3, projects.ListOpts{Name: userInputProject}).AllPages()
			if err != nil {
				return nil, err
			}
		}

		pList, err := projects.ExtractProjects(page)
		if err != nil {
			return nil, err
		}
		// no need to continue, if there are multiple projects in the list
		if len(pList) > 1 {
			return nil, errors.New("more than one project exists with the name " + userInputProject)
		}

		for _, pInList := range pList {
			p.ID = pInList.ID
			p.Name = pInList.Name
			p.DomainID = pInList.DomainID
		}
	}
	if p.ID == "" {
		return nil, errors.New("project not found")
	}

	// this is needed in case the user did not gave a domain ID at input
	// which means we still don't have the domain name
	if p.DomainName == "" {
		d, err := domains.Get(identityV3, p.DomainID).Extract()
		if err != nil {
			return nil, err
		}
		p.DomainName = d.Name
	}

	return p, nil
}
