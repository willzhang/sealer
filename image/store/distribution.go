// Copyright © 2021 Alibaba Group Holding Ltd.
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

package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alibaba/sealer/utils"

	"github.com/opencontainers/go-digest"
)

type DistributionMetadataItem struct {
	SourceRepository   string        `json:"source_repository"`
	LayerContentDigest digest.Digest `json:"layer_content_digest"`
}

// distributionMetadata is the data from {layerdb}/distribution_layer_digest
// which indicate that digest of compressedlayerStream in specific registry and repository
type DistributionMetadata []DistributionMetadataItem

func (ls LayerStorage) LoadDistributionMetadata(layerID LayerID) (map[string]digest.Digest, error) {
	var (
		layerDBPath = ls.LayerDBDir(layerID.ToDigest())
		metadatas   = DistributionMetadata{}
		res         = map[string]digest.Digest{}
	)
	distributionMetadataFile, err := os.Open(filepath.Join(layerDBPath, "distribution_layer_digest"))
	if err != nil {
		return res, nil
	}
	defer distributionMetadataFile.Close()

	err = json.NewDecoder(distributionMetadataFile).Decode(&metadatas)
	if err != nil {
		return res, err
	}

	for _, item := range metadatas {
		res[item.SourceRepository] = item.LayerContentDigest
	}

	return res, nil
}

func (ls LayerStorage) addDistributionMetadata(layerID LayerID, newMetadatas map[string]digest.Digest) error {
	// load from distribution_layer_digest
	metadataMap, err := ls.LoadDistributionMetadata(layerID)
	if err != nil {
		return err
	}
	// override metadata items, and add new metadata
	for key, value := range newMetadatas {
		metadataMap[key] = value
	}

	distributionMetadatas := DistributionMetadata{}
	for key, value := range metadataMap {
		distributionMetadatas = append(distributionMetadatas, DistributionMetadataItem{
			SourceRepository:   key,
			LayerContentDigest: value,
		})
	}

	distributionMetadatasJSON, err := json.Marshal(&distributionMetadatas)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(ls.LayerDBDir(layerID.ToDigest()), "distribution_layer_digest"), distributionMetadatasJSON)
}
