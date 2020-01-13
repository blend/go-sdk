package main

import (
	"fmt"
	"time"
	selector "github.com/blend/go-sdk/selector"
	k8s "k8s.io/apimachinery/pkg/labels"
)

func benchSelector(sel string, labels []map[string]string, binder func(string, map[string]string) (bool, error)) (d time.Duration, err error) {
	start := time.Now()
	defer func() {
		d = time.Since(start)
	}()
	var result bool
	for i := 0; i < 1024; i++ {
		for _, labelSet := range labels {
			result, err = binder(sel, labelSet)
			if err != nil {
				return
			}
			if !result {
				err = fmt.Errorf("selector failed")
				return
			}
		}
	}
	return
}

func kubeRunner(sel string, labels map[string]string) (bool, error) {
	s, err := k8s.Parse(sel)
	if err != nil {
		return false, err
	}

	return s.Matches(k8s.Set(labels)), nil
}

func blendRunner(sel string, labels map[string]string) (bool, error) {
	s, err := selector.Parse(sel)
	if err != nil {
		return false, err
	}
	return s.Matches(labels), nil
}

func main() {
	sel := "foo==bar,foo!=baz,moo in (foo, bar, baz, buzz),!thing"
	labels := []map[string]string{
		{"foo": "bar", "thing1": "", "moo": "foo"},
		{"foo": "bar", "thing1": "", "moo": "bar"},
		{"foo": "bar", "thing1": "", "moo": "baz"},
		{"foo": "bar", "thing1": "", "moo": "buzz"},
	}

	fmt.Println("starting bench")

	fmt.Println("k8s starting")
	k8s, err := benchSelector(sel, labels, kubeRunner)
	if err != nil {
		fmt.Printf("k8s failed: %v\n", err)
	} else {
		fmt.Println("k8s complete")
		fmt.Printf("%v\n", k8s)
	}
	fmt.Println("blend starting")
	blend, err := benchSelector(sel, labels, blendRunner)
	if err != nil {
		fmt.Printf("blend failed: %v\n", err)
	} else {
		fmt.Println("blend complete")
		fmt.Printf("%v\n", blend)
	}
}
