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

package core

import (
	"encoding/json"

	"github.com/sapcc/limesctl/internal/errors"
)

type jsonData []byte

// getJSON returns the result body of a get/list/update operation.
func (c *Cluster) renderJSON() jsonData {
	return parseToJSON(c.Result.Body)
}

// getJSON returns the result body of a get/list/update operation.
func (d *Domain) renderJSON() jsonData {
	return parseToJSON(d.Result.Body)
}

// getJSON returns the result body of a get/list/update operation.
func (p *Project) renderJSON() jsonData {
	return parseToJSON(p.Result.Body)
}

func parseToJSON(data interface{}) jsonData {
	b, err := json.Marshal(data)
	errors.Handle(err, "could not marshal JSON")
	return jsonData(b)
}
