//
// Copyright 2021 The Sigstore Authors.
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

package mirroring

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/spf13/viper"
)

// TODO: Tests should not make network calls, mock out connections
func TestFetchLeavesByRange(t *testing.T) {
	viper.Set("rekorServerURL", "https://rekor.sigstore.dev")
	viper.Set("tree_file_dir", ".tree")
	viper.Set("metadata_file_dir", ".metadata")
	err := FetchLeavesByRange(0, 10)
	if err != nil {
		t.Errorf("%s\n", err)
	}
}

func TestComputeRoot(t *testing.T) {
	viper.Set("rekorServerURL", "https://rekor.sigstore.dev")
	viper.Set("tree_file_dir", ".tree")
	viper.Set("metadata_file_dir", ".metadata")

	// the .tree file is not an json array instead it have one json per line
	f, err := os.Open(".tree")
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	defer f.Close()

	var leaves []Artifact
	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		var leaf Artifact
		err = json.Unmarshal([]byte(line), &leaf)
		if err != nil {
			t.Errorf("%s\n", err)
			return
		}

		leaves = append(leaves, leaf)
	}

	STH, err := ComputeRootFromMemory(leaves)
	if err != nil {
		t.Errorf("%s\n", err)
		t.Log(STH)
		return
	}
}
