package main

import "os"

func remove(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}
