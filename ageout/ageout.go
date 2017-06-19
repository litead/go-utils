package ageout

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Folder struct {
	Path         string
	MaxAge       time.Duration
	At           time.Duration
	Exts         []string
	TopLevelOnly bool
}

func (f *Folder) ageOutFile(now time.Time) {
	first := true

	fx := func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() {
			if f.TopLevelOnly && !first {
				return filepath.SkipDir
			}
			first = false
			return nil
		}
		if now.Sub(info.ModTime()) < f.MaxAge {
			return nil
		}
		remove := len(f.Exts) == 0
		for _, ext := range f.Exts {
			if strings.HasSuffix(path, ext) {
				remove = true
				break
			}
		}
		if remove {
			os.Remove(path)
		}
		return nil
	}

	filepath.Walk(f.Path, fx)
}

type folderSlice []*Folder

func (f folderSlice) Len() int {
	return len(f)
}

func (f folderSlice) Less(i, j int) bool {
	return f[i].At < f[j].At
}

func (f folderSlice) Swap(i, j int) {
	t := f[i]
	f[i] = f[j]
	f[j] = t
}

func CleanFolder(path string, age time.Duration, topLevelOnly bool) {
	f := Folder{Path: path, MaxAge: age, TopLevelOnly: topLevelOnly}
	f.ageOutFile(time.Now())
}

var folders folderSlice
var lock sync.Mutex

func AddFolder(folder *Folder) {
	lock.Lock()
	defer lock.Unlock()
	folders = append(folders, folder)
	sort.Sort(folders)
}

func doAgeOut() {
	var fs folderSlice
	var day time.Time

	for {
		lock.Lock()
		if len(fs) != len(folders) {
			if len(fs) == 0 {
				y, m, d := time.Now().Date()
				day = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			}
			fs = make(folderSlice, len(folders))
			copy(fs, folders)
		}
		lock.Unlock()

		if len(fs) == 0 {
			time.Sleep(time.Minute)
			continue
		}

		for _, f := range fs {
			now := time.Now()
			duration := now.Sub(day)
			if duration < f.At {
				time.Sleep(f.At - duration)
			}
			f.ageOutFile(now)
		}

		day = day.Add(time.Hour * 24)
	}
}

func Start() {
	go doAgeOut()
}
