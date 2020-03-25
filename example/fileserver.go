package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iamjinlei/memfs"
)

func main() {
	stopCh := make(chan bool)

	fs, err := memfs.New(
		map[string][]byte{
			"/root/home/foo/1.txt": []byte("@/root/home/foo/1.txt"),
			"/root/home/bar/1.txt": []byte("@/root/home/bar/1.txt"),
			"/root/home/1.txt":     []byte("@/root/home/1.txt"),
			"/root/home/xyz/":      nil,
			"/etc/1.txt":           []byte("@/etc/1.txt"),
			"/shutdown":            []byte("http server has shut down"),
		},
		map[string]func(path string){
			"Readdir": func(path string) {
				fmt.Printf("opening dir %v\n", path)
			},
			"Read": func(path string) {
				fmt.Printf("opening file %v\n", path)
				if path == "/shutdown" {
					select {
					case stopCh <- true:
					default:
					}
				}
			},
			"Close": func(path string) {
				fmt.Printf("closing %v\n", path)
			},
		})

	if err != nil {
		fmt.Printf("error building memfs %v\n", err)
		return
	}

	fmt.Printf("serving http://localhost:8080\n")

	srv := &http.Server{Addr: ":8080"}
	http.Handle("/", http.StripPrefix("/", http.FileServer(fs)))

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			fmt.Printf("ListenAndServe() err: %s\n", err)
		}
	}()

	<-stopCh

	if err := srv.Shutdown(context.TODO()); err != nil {
		fmt.Printf("error shutting down http server %v\n", err)
	}
	fmt.Printf("http server closed\n")
}
