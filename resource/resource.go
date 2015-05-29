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

package resource

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/resource/internal"
)

const (
	ResourcesVersion = "1.8.6"
	vanillaURL       = "https://s3.amazonaws.com/Minecraft.Download/versions/%[1]s/%[1]s.jar"
)

var (
	lock           sync.RWMutex
	packs          []*pack
	errMissingFile = errors.New("file not found")
)

type pack struct {
	name string

	files map[string]opener
}

type opener func() (io.ReadCloser, error)

// Open searches through all open resource packs for the requested file.
// If a file exists but fails to open that error will be returned instead
// of the standard 'file not found' but only if another file couldn't be
// found to be used in its place.
func Open(plugin, name string) (io.ReadCloser, error) {
	lock.RLock()
	defer lock.RUnlock()
	var lastErr error
	for i := len(packs) - 1; i >= 0; i-- {
		pck := packs[i]
		if f, ok := pck.files[fmt.Sprintf("assets/%s/%s", plugin, name)]; ok {
			r, err := f()
			if err != nil {
				lastErr = err
				continue
			}
			return r, nil
		}
	}
	if lastErr == nil {
		return nil, errMissingFile
	}
	return nil, lastErr
}

// OpenAll searches through all open resource packs for the requested file.
// If a file exists but fails to open that error will be returned instead
// of the standard 'file not found' but only if another file couldn't be
// found to be used in its place. Unlike Open this will return all matching
// files
func OpenAll(plugin, name string) ([]io.ReadCloser, error) {
	lock.RLock()
	defer lock.RUnlock()
	var lastErr error
	var out []io.ReadCloser
	for i := len(packs) - 1; i >= 0; i-- {
		pck := packs[i]
		if f, ok := pck.files[fmt.Sprintf("assets/%s/%s", plugin, name)]; ok {
			r, err := f()
			if err != nil {
				lastErr = err
				continue
			}
			out = append(out, r)
		}
	}
	if lastErr == nil && len(out) == 0 {
		return nil, errMissingFile
	}
	return out, lastErr
}

// Search searches for files that belong to the passed plugin and exist
// the passed path with the passed extension. This searches all open packs.
func Search(plugin, path, ext string) []string {
	lock.RLock()
	defer lock.RUnlock()
	found := map[string]bool{}
	var lst []string
	search := fmt.Sprintf("assets/%s/%s", plugin, path)
	base := fmt.Sprintf("assets/%s/", plugin)
	for _, pck := range packs {
		for k := range pck.files {
			if !found[k] && strings.HasPrefix(k, search) && strings.HasSuffix(k, ext) {
				found[k] = true
				lst = append(lst, k[len(base):])
			}
		}
	}
	return lst
}

func IsActive(name string) bool {
	lock.RLock()
	defer lock.RUnlock()
	for _, pck := range packs {
		if pck.name == name {
			return true
		}
	}
	return false
}

func RemovePack(name string) {
	lock.RLock()
	defer lock.RUnlock()
	for i, pck := range packs {
		if pck.name == name {
			packs = append(packs[:i], packs[i+1:]...)
			return
		}
	}
}

type TickFunc func(progress float64, done bool)

// TODO(Think) Ideally this package has a way to start the download instead of
// being an init thing. Also should have a way to get progress information.

func Init(tick TickFunc, sync chan<- func()) {
	fromInternal()
	defLocation := fmt.Sprintf("./resources-%s", ResourcesVersion)
	_, err := os.Stat(defLocation)
	if os.IsNotExist(err) {
		go func() {
			sync <- func() { tick(0, false) }
			downloadDefault(tick, sync, defLocation)
			sync <- func() {
				if err := fromDir(defLocation); err != nil {
					panic(err)
				}
				tick(1, true)
			}
		}()
	} else {
		if err := fromDir(defLocation); err != nil {
			panic(err)
		}
	}
}

type dummyCloser struct {
	*bytes.Reader
}

func (dummyCloser) Close() error { return nil }

func fromInternal() {
	lock.Lock()
	defer lock.Unlock()
	p := &pack{
		name:  "$internal",
		files: map[string]opener{},
	}
	for _, name := range internal.AssetNames() {
		name := name
		p.files[name] = func() (io.ReadCloser, error) {
			data, err := internal.Asset(name)
			return dummyCloser{bytes.NewReader(data)}, err
		}
	}
	packs = append(packs, p)
}

func fromDir(d string) error {
	lock.Lock()
	defer lock.Unlock()
	p := &pack{
		name:  d,
		files: map[string]opener{},
	}

	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(d, path)
		if err != nil {
			return err
		}
		rel = strings.Replace(rel, string(filepath.Separator), "/", -1)
		p.files[rel] = func() (io.ReadCloser, error) {
			return os.Open(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	packs = append([]*pack{packs[0], p}, packs[1:]...)
	return nil
}

func LoadZip(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	return fromFile(f, name)
}

func fromFile(f *os.File, name string) error {
	lock.Lock()
	defer lock.Unlock()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	z, err := zip.NewReader(f, s.Size())
	if err != nil {
		return err
	}
	p := &pack{
		name:  name,
		files: map[string]opener{},
	}
	for _, f := range z.File {
		f := f
		p.files[f.Name] = f.Open
	}
	packs = append(packs, p)
	return nil
}

type progressRead struct {
	max  int64
	n    int64
	tick TickFunc
	sync chan<- func()

	r io.Reader
}

func (p *progressRead) Read(buf []byte) (n int, err error) {
	n, err = p.r.Read(buf)
	if n > 0 {
		p.n += int64(n)
		p.sync <- func() {
			p.tick(float64(p.n)/float64(p.max), false)
		}
	}
	return
}

func downloadDefault(tick TickFunc, sync chan<- func(), target string) {
	console.Text("Obtaining vanilla resources for %s, please wait...", ResourcesVersion)
	resp, err := http.Get(fmt.Sprintf(vanillaURL, ResourcesVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	os.MkdirAll("./", 0777)
	f, err := os.Create(target + ".tmp")
	if err != nil {
		panic(err)
	}
	defer os.Remove(target + ".tmp")
	defer f.Close()
	_, err = io.Copy(f, &progressRead{
		max:  resp.ContentLength,
		tick: tick,
		sync: sync,
		r:    resp.Body,
	})
	if err != nil {
		panic(err)
	}

	f.Seek(0, 0) // Go back to the start
	fr, err := zip.NewReader(f, resp.ContentLength)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(target, 0777)

	// Copy the assets (not the classes) in the new zip
	for _, f := range fr.File {
		if !strings.HasPrefix(f.Name, "assets/") {
			continue
		}
		func() {
			path := filepath.Join(target, f.Name)
			os.MkdirAll(filepath.Dir(path), 0777)
			w, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer w.Close()
			r, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer r.Close()
			_, err = io.Copy(w, r)
			if err != nil {
				panic(err)
			}
		}()
	}
}
