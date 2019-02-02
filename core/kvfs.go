package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type KVForFile struct {
	path     string
	readonly bool
	running  bool
}

func __ensurePath(p string) error {
	if fi, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(p, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
		}
	} else if !fi.IsDir() {
		return errors.New(fi.Name() + " exists is not directory")
	}
	return nil
}

func (f *KVForFile) Set(k string, v interface{}) error {
	val := f.Get(k)
	if val != nil && f.readonly {
		return errors.New("readonly")
	}
	if err := __ensurePath(f.path); err != nil {
		return err
	}
	js, err := json.Marshal([2]interface{}{k, v})
	if err != nil {
		return err
	}
	ioutil.WriteFile(f.path+"/"+f.__hash(k), js, os.ModePerm)
	return nil
}

func (f *KVForFile) __hash(k string) string {
	h := md5.New()
	h.Write([]byte(k))
	return hex.EncodeToString(h.Sum(nil)[4:12])
}

func (f *KVForFile) Get(k string) (ret interface{}) {
	if data, err := ioutil.ReadFile(f.path + "/" + f.__hash(k)); err == nil {
		var v [2]interface{}
		if json.Unmarshal(data, &v) == nil {
			return v[1]
		}
	}
	return
}

func (f *KVForFile) Increment(callback func(k string, v interface{})) {
	if files, err := ioutil.ReadDir(f.path); err == nil {
		for _, fi := range files {
			if fi.IsDir() && fi.Size() > 0 {
				continue
			}
			if data, err := ioutil.ReadFile(f.path + "/" + fi.Name()); err == nil {
				var v [2]interface{}
				if json.Unmarshal(data, &v) == nil {
					if k, ok := v[0].(string); ok {
						callback(k, v[1])
					}
				}
			}
		}
	}
}

func (f *KVForFile) Del(k string) {
	if !f.readonly {
		os.Remove(f.path + "/" + f.__hash(k))
	}
}

func (f *KVForFile) Cls() error {
	if f.readonly {
		return errors.New("readonly")
	}
	return os.RemoveAll(f.path)
}

func (f *KVForFile) IsRunning() bool {
	return f.running
}

func (f *KVForFile) Close() {
	f.running = false
}

func Open(path string, readonly bool) (*KVForFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	p := &KVForFile{}
	if !readonly {
		if err := __ensurePath(absPath); err != nil {
			return nil, err
		}
	} else {
		if fi, err := os.Stat(absPath); err != nil {
			return nil, err
		} else {
			if !fi.IsDir() {
				return nil, errors.New("directory exists but is not an empty directory")
			}
		}
	}
	p.readonly = readonly
	p.path = path
	p.running = true
	return p, nil
}
