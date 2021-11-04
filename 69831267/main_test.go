package main

import (
	"io"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func captureStdout(f func(w io.Writer)) string {
	originalStdout := os.Stdout
	defer func(){
		os.Stdout= originalStdout
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	f(w)
	if err := w.Close(); err != nil {
		panic(err)
	} else if out, err := ioutil.ReadAll(r); err != nil {
		panic(err)
	} else if err := r.Close(); err != nil {
		panic(err)
	} else {
		return string(out)
	}
}

func TestPrint(t *testing.T){
	exp := "Hello, world"
	for i := 0; i < 100000; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T){
			t.Parallel()
			if res := captureStdout(Print); res != exp  {
				t.Errorf("Expected %s, got %s", exp, res)
			}
		})
	}
}
