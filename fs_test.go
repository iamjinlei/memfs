package memfs

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fsDef = map[string][]byte{
	"/root/home/foo/1.txt": []byte("@/root/home/foo/1.txt"),
	"/root/home/bar/1.txt": []byte("@/root/home/bar/1.txt"),
	"/root/home/1.txt":     []byte("@/root/home/1.txt"),
	"/root/home/xyz/":      nil,
	"/etc/1.txt":           []byte("@/etc/1.txt"),
	"/1.txt":               []byte("@/1.txt"),
}

func names(files []*File) []string {
	names := []string{}
	for _, f := range files {
		names = append(names, f.name)
	}
	sort.Strings(names)
	return names
}

func find(files []*File, name string) *File {
	for _, f := range files {
		if f.name == name {
			return f
		}
	}
	return nil
}

func TestFs(t *testing.T) {
	fs, err := New(fsDef, nil)

	assert.NoError(t, err)
	assert.NotNil(t, fs)

	n := fs.root.Load().(*File)
	assert.NotNil(t, n)
	assert.Equal(t, "/", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"1.txt", "etc", "root"}, names(n.children))

	n = find(fs.root.Load().(*File).children, "1.txt")
	assert.Equal(t, "1.txt", n.name)
	assert.Equal(t, []byte("@/1.txt"), n.bytes)
	assert.Nil(t, n.children)

	n = find(fs.root.Load().(*File).children, "etc")
	assert.Equal(t, "etc", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"1.txt"}, names(n.children))

	n = find(n.children, "1.txt")
	assert.Equal(t, "1.txt", n.name)
	assert.Equal(t, []byte("@/etc/1.txt"), n.bytes)
	assert.Nil(t, n.children)

	n = find(fs.root.Load().(*File).children, "root")
	assert.Equal(t, "root", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"home"}, names(n.children))

	n = find(n.children, "home")
	assert.Equal(t, "home", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"1.txt", "bar", "foo", "xyz"}, names(n.children))

	home := n
	n = find(n.children, "1.txt")
	assert.Equal(t, "1.txt", n.name)
	assert.Equal(t, []byte("@/root/home/1.txt"), n.bytes)
	assert.Nil(t, n.children)

	n = find(home.children, "bar")
	assert.Equal(t, "bar", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"1.txt"}, names(n.children))

	n = find(n.children, "1.txt")
	assert.Equal(t, "1.txt", n.name)
	assert.Equal(t, []byte("@/root/home/bar/1.txt"), n.bytes)
	assert.Nil(t, n.children)

	n = find(home.children, "foo")
	assert.Equal(t, "foo", n.name)
	assert.Nil(t, n.bytes)
	assert.Equal(t, []string{"1.txt"}, names(n.children))

	n = find(n.children, "1.txt")
	assert.Equal(t, "1.txt", n.name)
	assert.Equal(t, []byte("@/root/home/foo/1.txt"), n.bytes)
	assert.Nil(t, n.children)

	n = find(home.children, "xyz")
	assert.Equal(t, "xyz", n.name)
	assert.Nil(t, n.bytes)
	assert.Nil(t, n.children)

	// Open
	of, err := fs.Open("/1.txt")
	assert.NoError(t, err)
	f := of.(*File)
	assert.Equal(t, "1.txt", f.name)
	assert.Equal(t, []byte("@/1.txt"), f.bytes)
	assert.Nil(t, f.children)

	of, err = fs.Open("root/home/1.txt")
	assert.NoError(t, err)
	f = of.(*File)
	assert.Equal(t, "1.txt", f.name)
	assert.Equal(t, []byte("@/root/home/1.txt"), f.bytes)
	assert.Nil(t, f.children)

	of, err = fs.Open("/root/home/xyz")
	assert.NoError(t, err)
	f = of.(*File)
	assert.Equal(t, "xyz", f.name)
	assert.Nil(t, f.bytes)
	assert.Nil(t, f.children)

	of, err = fs.Open("/root/home/2.txt")
	assert.Error(t, err)
	assert.Nil(t, of)

	of, err = fs.Open("/")
	assert.NoError(t, err)
	d := of.(*File)
	assert.Equal(t, "/", d.name)
	assert.Nil(t, d.bytes)
	assert.Equal(t, []string{"1.txt", "etc", "root"}, names(d.children))

	of, err = fs.Open("root/home")
	assert.NoError(t, err)
	d = of.(*File)
	assert.Equal(t, "home", d.name)
	assert.Nil(t, d.bytes, nil)
	assert.Equal(t, []string{"1.txt", "bar", "foo", "xyz"}, names(d.children))

	fi, err := d.Stat()
	assert.NoError(t, err)
	assert.True(t, fi.IsDir())
	assert.EqualValues(t, 0, fi.Size())

	of, err = fs.Open("/root/home/")
	assert.NoError(t, err)
	d = of.(*File)
	assert.Equal(t, "home", d.name)
	assert.Nil(t, d.bytes, nil)
	assert.Equal(t, []string{"1.txt", "bar", "foo", "xyz"}, names(d.children))

	fi, err = d.Stat()
	assert.NoError(t, err)
	assert.True(t, fi.IsDir())
	assert.EqualValues(t, 0, fi.Size())

	fs.Update(map[string][]byte{
		"/0.txt": []byte("@/0.txt"),
	})

	of, err = fs.Open("/0.txt")
	assert.NoError(t, err)
	f = of.(*File)
	assert.Equal(t, "0.txt", f.name)
	assert.Equal(t, []byte("@/0.txt"), f.bytes)
	assert.Nil(t, f.children)

	_, err = fs.Open("/1.txt")
	assert.Error(t, err)
}
