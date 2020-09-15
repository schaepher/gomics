// Copyright (c) 2013-2018 Utkan Güngördü <utkan@freeconsole.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package archive

import (
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
)

type Dir struct {
	files []string
}

func NewDir(name string) (*Dir, error) {
	var err error

	dir := new(Dir)

	files, err := ioutil.ReadDir(name)
	if err != nil {
		return nil, err
	}

	reg, _ := regexp.Compile(`\d+`)
	sort.Slice(files, func(i, j int) bool {
		intI, errI := strconv.Atoi(reg.FindString(files[i].Name()))
		intJ, errJ := strconv.Atoi(reg.FindString(files[j].Name()))

		if errI == nil && errJ == nil {
			return intI < intJ
		}

		return files[i].Name() < files[j].Name()
	})
	
	for _, fileInfo := range files {
		dir.files = append(dir.files, name + string(os.PathSeparator) + fileInfo.Name())
	}

	return dir, nil
}

func (dir *Dir) checkbounds(i int) error {
	if i < 0 || i >= len(dir.files) {
		return ErrBounds
	}
	return nil
}

func (dir *Dir) Load(i int, autorotate bool) (*gdk.Pixbuf, error) {
	if err := dir.checkbounds(i); err != nil {
		return nil, err
	}

	f, err := os.Open(dir.files[i])
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return LoadPixbuf(f, autorotate)
}

func (dir *Dir) Name(i int) (string, error) {
	if err := dir.checkbounds(i); err != nil {
		return "", err
	}

	return dir.files[i], nil
}

func (dir *Dir) Len() int {
	return len(dir.files)
}

func (dir *Dir) Close() error {
	return nil
}
