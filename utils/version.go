package utils

import "sync"

type version struct {
	version  string
	compTime string
}

var singleton *version
var once sync.Once

// Version gets the version object
func Version() *version {
	once.Do(func() {
		singleton = &version{}
	})
	return singleton
}

func (v *version) Get() string {
	return v.version
}

func (v *version) Set(s string) {
	v.version = s
}

func (v *version) GetCompTime() string {
	return v.compTime
}

func (v *version) SetCompTime(s string) {
	v.compTime = s
}
