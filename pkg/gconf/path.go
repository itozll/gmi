package gconf

import "github.com/itozll/gmi/pkg/env"

const _pathRoot = "conf/"

func PathAll() []string {
	return []string{
		_pathRoot,
		_pathRoot + env.Name(),
	}
}

func PathRoot() string {
	return _pathRoot
}

func PathEnv() string {
	return _pathRoot + env.Name()
}
