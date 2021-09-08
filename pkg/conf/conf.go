package conf

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var ErrNotFound = errors.New("not found file")

// ListOf list all file at the path
func ListOf(path string, opts ...FuncOption) ([]string, error) {
	opt := defaultOptions()
	for _, fn := range opts {
		fn(opt)
	}

	return List(path, opt)
}

// List list all file at the path
func List(path string, opt *Options) ([]string, error) {
	if len(path) == 0 {
		path = "./"
	} else if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return listPath(path, opt, 1)
}

// listPath list all file at the path
func listPath(path string, opt *Options, depth int) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		name := f.Name()
		fa := path + f.Name()

		if f.IsDir() {
			if opt.IncludeDir {
				appendFile(opt, name, fa, &res)
			}

			if opt.Recursion {
				if opt.Depth == 0 || opt.Depth > depth {
					r, err := listPath(fa+"/", opt, depth+1)
					if err != nil {
						return nil, err
					}

					if len(r) > 0 {
						res = append(res, r...)
					}
				}
			}
		} else if f.Mode().IsRegular() {
			appendFile(opt, name, fa, &res)
		}
	}

	return res, nil
}

func appendFile(opt *Options, name string, fa string, res *[]string) {
	if len(opt.ExcludeFileName) > 0 && opt.ExcludeFileName == name {
		return
	}

	if len(opt.FileName) > 0 {
		if opt.FileName == name {
			*res = append(*res, fa)
		}

		return
	}

	if len(opt.Prefix) != 0 && !strings.HasPrefix(name, opt.Prefix) {
		return
	}

	if len(opt.Suffix) != 0 && !strings.HasSuffix(name, opt.Suffix) {
		return
	}

	*res = append(*res, fa)
}

// loadPath load files, and decode to out
func loadPath(path string, out interface{}, loadFunc func(interface{}, string) error, opt *Options) error {
	l, err := List(path, opt)
	if err != nil {
		return err
	}

	for _, f := range l {
		if err = loadFunc(out, f); err != nil {
			return err
		}
	}

	return nil
}

func AutoLoad(out interface{}, file string, loadFunc func(interface{}, string) error) (f string, err error) {
	var (
		fi     os.FileInfo
		loaded bool
	)

	for _, p := range PathAll() {
		f = strings.TrimRight(p, "/") + "/" + file

	retry:
		fi, err = os.Stat(f)

		if err != nil {
			if os.IsNotExist(err) {
				err = nil
				if !strings.HasSuffix(f, yamlSuffix) {
					f += yamlSuffix
					goto retry
				}

				continue
			}

			return
		}

		if !fi.Mode().IsRegular() {
			continue
		}

		if err = loadFunc(out, f); err != nil {
			return
		}

		loaded = true
	}

	if loaded {
		return "", nil
	}

	return file, ErrNotFound
}
