# MemFS

MemFS that implements http.FileSystem and http.File interface for quick mock up a in-memory http server

```golang

fs, _ := memfs.New(map[string][]byte{
	"/root/home/foo/1.txt": []byte("@/root/home/foo/1.txt"),
	"/root/home/bar/1.txt": []byte("@/root/home/bar/1.txt"),
	"/root/home/1.txt":     []byte("@/root/home/1.txt"),
	"/root/home/xyz/":      nil,
}, nil)

srv := &http.Server{Addr: ":8080"}
http.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
srv.ListenAndServe();
```
