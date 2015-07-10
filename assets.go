// Copyright 2015 Matthew Collins
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

package steven

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thinkofdeath/steven/console"
)

const (
	assetsVersion    = "1.8"
	assetIndexURL    = "https://s3.amazonaws.com/Minecraft.Download/indexes/%s.json"
	assetResourceURL = "http://resources.download.minecraft.net/%s"
)

var (
	assets assetIndex
)

type assetIndex struct {
	Objects map[string]struct {
		Hash string `json:"hash"`
		Size int    `json:"size"`
	} `json:"objects"`
}

func init() {
	defLocation := "./resources"
	loc := fmt.Sprintf("%s/%s.index", defLocation, assetsVersion)
	_, err := os.Stat(loc)
	if os.IsNotExist(err) {
		getAssetIndex()
	} else {
		f, err := os.Open(fmt.Sprintf("./resources/%s.index", assetsVersion))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&assets)
		if err != nil {
			panic(err)
		}
	}
	syncChan <- downloadAssets
}

func getAssetIndex() {
	console.Text("Getting asset index")
	defLocation := "./resources"
	resp, err := http.Get(fmt.Sprintf(assetIndexURL, assetsVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		panic(err)
	}
	os.MkdirAll("./resources", 0777)
	f, err := os.Create(fmt.Sprintf("%s/%s.index", defLocation, assetsVersion))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(assets)
	console.Text("Got asset index for %s", assetsVersion)
}

func downloadAssets() {
	defLocation := "./resources"
	total := newProgressBar()
	var totalSize, totalCount int64
	for _, v := range assets.Objects {
		totalSize += int64(v.Size)
	}
	limiter := make(chan struct{}, 4)
	for i := 0; i < 4; i++ {
		limiter <- struct{}{}
	}
	go func() {
		for file, v := range assets.Objects {
			v := v
			file := file
			path := hashPath(v.Hash)
			loc := fmt.Sprintf("%s/%s", defLocation, path)
			_, err := os.Stat(loc)
			if !os.IsNotExist(err) {
				continue
			}
			<-limiter
			go func() {
				var prog *progressBar
				wait := make(chan struct{})
				syncChan <- func() { prog = newProgressBar(); wait <- struct{}{} }
				<-wait
				func() {
					resp, err := http.Get(fmt.Sprintf(assetResourceURL, path))
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					os.MkdirAll(filepath.Dir(loc), 0777)
					f, err := os.Create(loc + ".tmp")
					if err != nil {
						panic(err)
					}
					defer f.Close()
					n, err := prog.watchCopy(f, resp.Body, resp.ContentLength, fmt.Sprintf("Downloading %s: %%v/100", file))
					if err != nil {
						panic(err)
					}
					if n != int64(v.Size) {
						panic(fmt.Sprintf("Got: %d, Wanted: %d for %s", n, v.Size, fmt.Sprintf(assetResourceURL, path)))
					}
					console.Text("Downloaded: %s", loc)
				}()
				os.Rename(loc+".tmp", loc)
				syncChan <- func() {
					prog.remove()
					totalCount += int64(v.Size)
					progress := float64(totalCount) / float64(totalSize)
					total.update(progress, fmt.Sprintf("Downloading assets: %v/100", int(100*progress)))
				}
				limiter <- struct{}{}
			}()
		}
		syncChan <- func() {
			total.remove()
		}
	}()
}

func hashPath(hash string) string {
	return hash[:2] + "/" + hash
}
