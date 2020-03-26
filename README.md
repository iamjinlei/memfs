[![Build Status](https://travis-ci.com/iamjinlei/memfs.svg?branch=master)](https://travis-ci.com/iamjinlei/memfs)

# MemFS

MemFS that implements http.FileSystem and http.File interface for quick mock up an in-memory http server

```golang

fs, _ := memfs.New(map[string][]byte{
	"/root/home/foo/1.txt": []byte("@/root/home/foo/1.txt"),
	"/root/home/bar/1.txt": []byte("@/root/home/bar/1.txt"),
	"/root/home/1.txt":     []byte("@/root/home/1.txt"),
	"/root/home/xyz/":      nil,
}, nil)

srv := &http.Server{Addr: ":8080"}
http.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
srv.ListenAndServe()
```

Check the example folder for a working implementation

```bash
go run example/fileserver.go
```
