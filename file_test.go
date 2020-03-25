package memfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	fs, err := New(fsDef, nil)

	assert.NoError(t, err)
	assert.NotNil(t, fs)

	// regular file
	f, err := fs.Open("/root/home/1.txt")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	fi, err := f.Stat()
	assert.NoError(t, err)
	assert.False(t, fi.IsDir())
	assert.Equal(t, "1.txt", fi.Name())
	assert.EqualValues(t, 17, fi.Size())

	_, err = f.Readdir(10)
	assert.Error(t, err)

	buf := make([]byte, 6)
	n, err := f.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	assert.Equal(t, []byte("@/root"), buf)

	pos, err := f.Seek(6, 1)
	assert.NoError(t, err)
	assert.EqualValues(t, 12, pos)

	buf = make([]byte, 5)
	n, err = f.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("1.txt"), buf)

	pos, err = f.Seek(7, 0)
	assert.NoError(t, err)
	assert.EqualValues(t, 7, pos)

	buf = make([]byte, 4)
	n, err = f.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("home"), buf)

	pos, err = f.Seek(-3, 2)
	assert.NoError(t, err)
	assert.EqualValues(t, 14, pos)

	buf = make([]byte, 3)
	n, err = f.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("txt"), buf)

	// dir
	d, err := fs.Open("/root/home/")
	assert.NoError(t, err)
	assert.NotNil(t, d)

	fi, err = d.Stat()
	assert.NoError(t, err)
	assert.True(t, fi.IsDir())
	assert.Equal(t, "home", fi.Name())
	assert.EqualValues(t, 0, fi.Size())

	fis, err := d.Readdir(10)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(fis))

	buf = make([]byte, 6)
	_, err = d.Read(buf)
	assert.Error(t, err)

	_, err = d.Seek(6, 1)
	assert.Error(t, err)
}
