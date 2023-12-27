/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package expvar

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
)

// Vars is a collection of expvars.
type Vars struct {
	vars		sync.Map	// map[string]Var
	varKeysMu	sync.RWMutex
	varKeys		[]string	// sorted
}

// Get retrieves a named exported variable. It returns nil if the name has
// not been registered.
func (v *Vars) Get(name string) Var {
	i, _ := v.vars.Load(name)
	val, _ := i.(Var)
	return val
}

// Publish declares a named exported variable. This should be called from a
// package's init function when it creates its Vars. If the name is already
// registered then this will log.Panic.
func (v *Vars) Publish(name string, val Var) error {
	if _, dup := v.vars.LoadOrStore(name, val); dup {
		return fmt.Errorf("reuse of exported var name: %s", name)
	}
	v.varKeysMu.Lock()
	defer v.varKeysMu.Unlock()
	v.varKeys = append(v.varKeys, name)
	sort.Strings(v.varKeys)
	return nil
}

// Forward forwards the vars contained in this set to another set with a given key prefix.
func (v *Vars) Forward(dst *Vars, keyPrefix string) error {
	var err error
	return v.Do(func(kv KeyValue) error {
		if err = dst.Publish(keyPrefix+kv.Key, kv.Value); err != nil {
			return err
		}
		return nil
	})
}

// Do calls f for each exported variable.
// The global variable map is locked during the iteration,
// but existing entries may be concurrently updated.
func (v *Vars) Do(f func(KeyValue) error) error {
	v.varKeysMu.RLock()
	defer v.varKeysMu.RUnlock()
	var err error
	for _, k := range v.varKeys {
		val, _ := v.vars.Load(k)
		err = f(KeyValue{k, val.(Var)})
		if err != nil {
			return err
		}
	}
	return nil
}

// Handler returns an http.HandlerFunc that renders the vars as json.
func (v *Vars) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = v.WriteTo(w)
}

// WriteTo writes the vars to a given writer as json.
//
// This is called by the Handler function to return the output.
func (v *Vars) WriteTo(wr io.Writer) (size int64, err error) {
	var n int
	n, err = fmt.Fprint(wr, "{")
	size += int64(n)
	if err != nil {
		return
	}
	first := true
	err = v.Do(func(kv KeyValue) error {
		if !first {
			n, err = fmt.Fprint(wr, ",")
			size += int64(n)
			if err != nil {
				return err
			}
		}
		first = false
		n, err = fmt.Fprintf(wr, "%q:%s", kv.Key, kv.Value)
		size += int64(n)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	n, err = fmt.Fprintln(wr, "}")
	size += int64(n)
	return
}
