// Copyright 2016 The Linux Foundation
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

package schema_test

import (
	"strings"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/schema"
	"github.com/opencontainers/image-spec/specs-go/v1"
)

func TestManifest(t *testing.T) {
	for i, tt := range []struct {
		manifest string
		strict   bool
		fail     bool
	}{
		// expected failure: mediaType does not match pattern
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "invalid",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 148,
      "digest": "sha256:c57089565e894899735d458f0fd4bb17a0f1e0df8d72da392b85c9b35ee777cd"
    }
  ]
}
`,
			strict: true,
			fail:   true,
		},

		// expected failure: config.size is a string, expected integer
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": "1470",
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 148,
      "digest": "sha256:c57089565e894899735d458f0fd4bb17a0f1e0df8d72da392b85c9b35ee777cd"
    }
  ]
}
`,
			strict: true,
			fail:   true,
		},

		// expected failure: layers.size is string, expected integer
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": "675598",
      "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
    }
  ]
}
`,
			strict: true,
			fail:   true,
		},

		// expected failure: unrecognized layer media type and strict is true
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.other.layer",
      "size": "675598",
      "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
    }
  ]
}
`,
			strict: true,
			fail:   true,
		},

		// expected success: unrecognized layer media type, but strict is false
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.other.layer",
      "size": "675598",
      "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
    }
  ]
}
`,
			strict: false,
			fail:   true,
		},

		// valid manifest with optional fields
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 675598,
      "digest": "sha256:9d3dd9504c685a304985025df4ed0283e47ac9ffa9bd0326fddf4d59513f0827"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 156,
      "digest": "sha256:2b689805fbd00b2db1df73fae47562faac1a626d5f61744bfe29946ecff5d73d"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 148,
      "digest": "sha256:c57089565e894899735d458f0fd4bb17a0f1e0df8d72da392b85c9b35ee777cd"
    }
  ],
  "annotations": {
    "key1": "value1",
    "key2": "value2"
  }
}
`,
			strict: true,
			fail:   false,
		},

		// valid manifest with only required fields
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 675598,
      "digest": "sha256:9d3dd9504c685a304985025df4ed0283e47ac9ffa9bd0326fddf4d59513f0827"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 156,
      "digest": "sha256:2b689805fbd00b2db1df73fae47562faac1a626d5f61744bfe29946ecff5d73d"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 148,
      "digest": "sha256:c57089565e894899735d458f0fd4bb17a0f1e0df8d72da392b85c9b35ee777cd"
    }
  ]
}
`,
			fail: false,
		},

		// expected failure: empty layer, expected at least one
		{
			manifest: `
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 1470,
    "digest": "sha256:c86f7763873b6c0aae22d963bab59b4f5debbed6685761b5951584f6efb0633b"
  },
  "layers": []
}
`,
			fail: true,
		},
	} {
		manifestBytes := []byte(tt.manifest)
		reader := strings.NewReader(tt.manifest)
		descriptor := v1.Descriptor{
			MediaType: v1.MediaTypeImageManifest,
			Digest:    digest.FromBytes(manifestBytes).String(),
			Size:      int64(len(manifestBytes)),
		}
		err := schema.Validate(reader, &descriptor, tt.strict)
		if got := err != nil; tt.fail != got {
			t.Errorf("test %d: expected validation failure %t but got %t, err %v", i, tt.fail, got, err)
		}
	}
}
