// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"encoding/json"

	"github.com/goharbor/harbor/src/api/artifact"
	"github.com/goharbor/harbor/src/api/artifact/abstractor/resolver"
	"github.com/goharbor/harbor/src/api/scan"
	"github.com/goharbor/harbor/src/pkg/scan/report"
	v1 "github.com/goharbor/harbor/src/pkg/scan/rest/v1"
)

func boolValue(v *bool) bool {
	if v != nil {
		return *v
	}

	return false
}

func resolveVulnerabilitiesAddition(ctx context.Context, artifact *artifact.Artifact) (*resolver.Addition, error) {
	art := &v1.Artifact{
		NamespaceID: artifact.ProjectID,
		Repository:  artifact.RepositoryName,
		Digest:      artifact.Digest,
		MimeType:    artifact.ManifestMediaType,
	}

	reports, err := scan.DefaultController.GetReport(art, []string{v1.MimeTypeNativeReport})
	if err != nil {
		return nil, err
	}

	vulnerabilities := make(map[string]interface{})
	for _, rp := range reports {
		// Resolve scan report data only when it is ready
		if len(rp.Report) == 0 {
			continue
		}

		vrp, err := report.ResolveData(rp.MimeType, []byte(rp.Report))
		if err != nil {
			return nil, err
		}

		vulnerabilities[rp.MimeType] = vrp
	}

	content, _ := json.Marshal(vulnerabilities)

	return &resolver.Addition{
		Content:     content,
		ContentType: "application/json",
	}, nil
}
