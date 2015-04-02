package resource

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	resourcesVersion = "1.8.3"
	vanillaURL       = "https://s3.amazonaws.com/Minecraft.Download/versions/%[1]s/%[1]s.jar"
)

var (
	packs          []*pack
	errMissingFile = errors.New("file not found")
)

type pack struct {
	zip   *zip.Reader
	files map[string]*zip.File
}

// Open searches through all open resource packs for the requested file.
// If a file exists but fails to open that error will be returned instead
// of the standard 'file not found' but only if another file couldn't be
// found to be used in its place.
func Open(plugin, name string) (io.ReadCloser, error) {
	var lastErr error
	for i := len(packs) - 1; i >= 0; i-- {
		pck := packs[i]
		if f, ok := pck.files[fmt.Sprintf("assets/%s/%s", plugin, name)]; ok {
			r, err := f.Open()
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

// Search searches for files that belong to the passed plugin and exist
// the passed path with the passed extension. This searches all open packs.
func Search(plugin, path, ext string) []string {
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

// TODO(Think) Ideally this package has a way to start the download instead of
// being an init thing. Also should have a way to get process information.

func init() {
	defLocation := fmt.Sprintf(".steven/vanilla-%s.res", resourcesVersion)
	f, err := os.Open(defLocation)
	if os.IsNotExist(err) {
		f = downloadDefault(defLocation)
	}
	fromFile(f)

	if err := loadZip(".steven/pack.zip"); err != nil {
		fmt.Printf("Couldn't load pack.zip: %s\n", err)
	}
}

func loadZip(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	return fromFile(f)
}

func fromFile(f *os.File) error {
	s, err := f.Stat()
	if err != nil {
		return err
	}
	z, err := zip.NewReader(f, s.Size())
	if err != nil {
		return err
	}
	p := &pack{
		zip:   z,
		files: map[string]*zip.File{},
	}
	for _, f := range z.File {
		p.files[f.Name] = f
	}
	packs = append(packs, p)
	return nil
}

func downloadDefault(target string) *os.File {
	fmt.Printf("Obtaining vanilla resources for %s, please wait...\n", resourcesVersion)
	resp, err := http.Get(fmt.Sprintf(vanillaURL, resourcesVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	os.MkdirAll(".steven/", 0777)
	f, err := os.Create(target + ".tmp")
	if err != nil {
		panic(err)
	}
	defer os.Remove(target + ".tmp")
	defer f.Close()
	size, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	f.Seek(0, 0) // Go back to the start
	fr, err := zip.NewReader(f, size)

	t, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer t.Seek(0, 0) // Rollback to the start after writing the zip
	zt := zip.NewWriter(t)
	defer zt.Close()

	// Copy the assets (not the classes) in the new zip
	for _, f := range fr.File {
		if !strings.HasPrefix(f.Name, "assets/") {
			continue
		}
		w, err := zt.CreateHeader(&f.FileHeader)
		if err != nil {
			panic(err)
		}
		r, err := f.Open()
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(w, r)
		if err != nil {
			panic(err)
		}
		r.Close()
	}

	return t
}
