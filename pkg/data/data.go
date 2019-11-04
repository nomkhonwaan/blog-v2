// Code generated by bindata. DO NOT EDIT.
// sources:
// data/facebook-open-graph-template.html
// data/graphql-playground.html

package data

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type gzipAsset struct {
	bytes []byte
	info  gzipFileInfoEx
}

type gzipFileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type gzipBindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi gzipBindataFileInfo) Name() string {
	return fi.name
}
func (fi gzipBindataFileInfo) Size() int64 {
	return fi.size
}
func (fi gzipBindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi gzipBindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi gzipBindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi gzipBindataFileInfo) IsDir() bool {
	return false
}
func (fi gzipBindataFileInfo) Sys() interface{} {
	return nil
}

var _gzipBindataDataFacebookopengraphtemplatehtml = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\x51\x4b\xf3\x30\x14\x86\xef\xfb\x2b\xf2\xe5\xfa\xcb\xa6\x77\x22" +
		"\xc9\x40\x36\x07\x82\xa0\xc8\x0a\x7a\x19\x9b\x63\x7b\x20\x4d\x62\x7a\xda\x51\x4a\xff\xbb\xb4\x1b\xb4\x1d\xce\xab" +
		"\x90\xf3\xbe\xcf\xc3\x81\x44\xfe\xdb\xbd\x6c\x0f\x1f\xaf\x8f\xac\xa0\xd2\x6e\x12\x39\x1c\xcc\x6a\x97\x2b\x0e\x8e" +
		"\x0f\x03\xd0\x66\x93\x30\x26\x4b\x20\xcd\xb2\x42\xc7\x0a\x48\xf1\xf4\xb0\x17\x77\x7c\x0a\x9c\x2e\x41\xf1\x06\xe1" +
		"\x18\x7c\x24\xce\x32\xef\x08\x1c\x29\x7e\x44\x43\x85\x32\xd0\x60\x06\x62\xbc\xfc\x67\xe8\x90\x50\x5b\x51\x65\xda" +
		"\x82\xba\x5d\xdd\xcc\x44\x05\x51\x10\xf0\x5d\x63\xa3\xf8\xbb\x48\x1f\xc4\xd6\x97\x41\x13\x7e\x5a\x98\x59\x11\x14" +
		"\x98\x1c\x66\x5c\x88\x3e\x40\xa4\x56\x71\x9f\xdf\xd7\xd1\xce\xca\x5d\xb7\x4a\xdf\x9e\xfb\xfe\x5a\x9b\xda\x00\xcb" +
		"\xfa\xa1\x0d\xf0\x47\x1f\xc9\x5e\x02\xc3\xe8\x3a\x61\xa0\xca\x22\x06\x42\xef\x96\xdc\x6e\x0a\xce\x74\xd7\xe1\x17" +
		"\x5b\xed\x41\x53\x1d\xc1\x3c\x95\x3a\x87\xbe\x4f\x18\xfb\xd5\x8b\x43\xbc\x34\x5e\x90\x67\x27\x38\x33\x5a\xe4\xb8" +
		"\xfb\x66\xda\x58\xae\x4f\x93\x44\xae\x4f\x2f\x2d\xd7\xe3\x4f\xf8\x09\x00\x00\xff\xff\x19\x26\xaf\xd3\x19\x02\x00" +
		"\x00")

func gzipBindataDataFacebookopengraphtemplatehtml() (*gzipAsset, error) {
	bytes := _gzipBindataDataFacebookopengraphtemplatehtml
	info := gzipBindataFileInfo{
		name:        "data/facebook-open-graph-template.html",
		size:        537,
		md5checksum: "",
		mode:        os.FileMode(420),
		modTime:     time.Unix(1572828204, 0),
	}

	a := &gzipAsset{bytes: bytes, info: info}

	return a, nil
}

var _gzipBindataDataGraphqlplaygroundhtml = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5c\x69\x6f\xdb\xb8\xd6\xfe\x9e\x5f\xc1\x57\x45\xd1\x14\x88\x68\x92" +
		"\xda\x3d\x76\xf0\x76\xc9\x4d\x8b\xc9\xdc\x66\x16\xb4\x77\xee\x37\x45\xa2\x6c\x26\x32\xe5\x4a\xf4\xd6\xc1\xfc\xf7" +
		"\x0b\x92\xb2\xbc\xc8\x8b\xe2\xc4\x99\x00\xad\x81\xc4\x32\x79\x78\x96\xe7\xe1\x91\xb9\x98\xea\xfc\xdf\xfb\x4f\xef" +
		"\xfe\xf8\xf3\xfa\x02\xf4\xc5\x20\x3d\x3f\x39\xe9\x54\xef\x34\x8c\xcf\x4f\x00\xe8\x0c\xa8\x08\x41\xd4\x0f\xf3\x82" +
		"\x8a\xee\x48\x24\xa6\x0f\x5a\x8b\x0a\x1e\x0e\x68\xd7\x18\x33\x3a\x19\x66\xb9\x30\x40\x94\x71\x41\xb9\xe8\x1a\xa3" +
		"\x82\xe6\x66\x11\x85\x69\x78\x93\xd2\x2e\xcf\xce\x00\xe3\x4c\xb0\x30\x55\x85\xb4\x8b\x21\x3a\x03\x03\xc6\xd9\x60" +
		"\x34\x58\x29\x0a\xa7\xb5\x22\x29\x15\xa6\xe6\x88\x19\xca\xb0\x60\x22\xa5\xe7\x97\x79\x38\xec\xff\x7a\x05\xae\xd3" +
		"\x70\xd6\xcb\xb3\x11\x8f\x3b\x2d\x5d\x23\x65\x52\xc6\xef\x40\x4e\xd3\xae\x51\x88\x59\x4a\x8b\x3e\xa5\xc2\x00\xfd" +
		"\x9c\x26\x5d\xa3\xd5\x8a\x62\x0e\x6f\x8b\x98\xa6\x6c\x9c\x43\x4e\x45\x8b\x0f\x07\xad\x9e\x54\xf8\x35\x35\x87\x95" +
		"\x42\x33\xa7\x61\x24\x5a\x37\x23\x96\xc6\xad\x42\x84\x82\x45\xad\xa8\x28\x5a\x8c\xc7\x74\x0a\xa3\xa2\x30\x4a\x28" +
		"\x96\xac\xf5\xb3\x5c\x44\x23\x01\x58\x94\xf1\x87\x19\x4c\xc2\xb1\x54\x02\x87\xbc\x37\xb7\x53\x44\x39\x1b\x0a\x50" +
		"\xe4\xd1\x03\xa3\xb8\x2d\x5a\x03\x16\xc7\x29\x9d\x84\x39\x85\xb7\x85\x71\xde\x69\x69\xe5\x92\xfb\x96\x26\xff\xa4" +
		"\x73\x93\xc5\xb3\xf3\x93\x8e\x82\x10\x88\xd9\x90\x76\x0d\x41\xa7\x42\x82\xa0\xa8\x90\x9d\x05\xfc\x75\x02\x00\x00" +
		"\x49\xc6\x85\x99\x84\x03\x96\xce\xda\xc0\xf8\x34\xa4\x1c\xfc\x1e\xf2\xc2\x38\x03\x45\xc8\x0b\xb3\xa0\x39\x4b\x7e" +
		"\x52\x92\xd9\x98\xe6\x49\x9a\x4d\xda\xa0\xcf\xe2\x98\x72\x59\xfa\xf7\xc9\x09\x00\xd2\x5a\xa9\x6d\x10\xe6\x3d\xc6" +
		"\xdb\x00\xe9\x26\x37\x61\x74\xa7\x83\x69\x83\x17\xd8\x23\xa1\x15\x56\xad\xe0\x22\xd2\x8f\xbc\x6c\x6e\x4e\xe8\xcd" +
		"\x1d\x13\x66\x28\x3b\x8e\x60\x19\x6f\x83\x15\x29\x04\x9d\x02\xd0\xb0\xa0\x66\x36\x12\x20\xc9\xf2\x49\x98\xc7\x85" +
		"\xb6\x75\xff\x36\xca\x8d\xff\x9f\xdb\xbc\xa3\xb3\x24\x0f\x07\xb4\x00\x1b\x1c\x4b\xf2\x6c\x50\x5e\x02\x90\x0d\xc3" +
		"\x88\x89\x59\x15\xe5\xc2\x6f\x91\x87\xbc\x48\xb2\x7c\xd0\x06\xea\x32\x0d\x05\xfd\xf3\x14\xa3\xe1\xf4\x75\x25\x3a" +
		"\x28\x9a\x88\xed\x15\xf9\x5b\xfd\x17\x59\xdd\x2b\xdc\xcc\x2b\xd4\xc0\x25\xb4\xc7\x1f\xb4\xe4\x8c\x46\xf3\x07\x8a" +
		"\x0f\x45\xb1\xd3\x52\x69\x2b\xd3\x78\x6b\xfe\xc2\x24\x8c\xe9\xa7\x91\xd8\x9e\x36\x73\x81\x86\x19\xb3\x5f\x7c\x5b" +
		"\xb2\xc8\x96\x8f\x43\xb0\xd9\x90\x61\xb3\x01\xc5\xe6\xf3\xe6\x78\x35\x53\x7e\x40\x78\x20\x84\x9b\x7b\xe3\x22\x2f" +
		"\x36\x63\xf9\x0f\xb9\xbc\x19\xbf\x67\x40\xec\xa6\xee\xf8\x03\xc3\x43\x30\xac\xf7\xc7\x70\x38\xa4\x61\xfe\x38\xc9" +
		"\xdd\xcc\xfb\xfd\xbe\x3f\xfb\xb4\xfe\x01\xdf\x23\xdf\x15\xd5\x9c\x6c\x33\x8a\x1b\x1c\x54\xd2\xdb\x7d\x5b\xaf\xde" +
		"\x5a\x55\x83\x67\xab\x2d\xbc\xdb\x16\xde\x6e\x0b\xef\xe8\x3c\xdf\x53\xd4\x75\xce\x19\xe7\x34\x7f\x9f\x87\x93\x0a" +
		"\x01\xf4\xb2\xf2\xa9\x10\x79\x76\x47\xcd\x38\x2c\xfa\x59\x92\x14\x54\xb4\x81\x87\x96\xfd\x77\x76\x0b\x63\x7b\x45" +
		"\x1a\xa3\xdd\xe2\x04\xa3\xed\x34\x3d\x73\x47\xeb\xc8\x66\x23\x71\x3f\x87\xdd\xfb\xb8\x80\x1d\xb2\x1d\xab\xa7\x36" +
		"\x0d\xfb\x1f\xbe\xdc\xde\x8d\xd7\xa6\x19\x55\x7f\x34\xb3\x9c\xe9\xe9\xfe\x70\x2a\xff\x74\xeb\x95\xde\xbc\x59\x64" +
		"\x9f\x86\x7d\xe9\xb9\x33\x39\xb7\xb7\xaa\x4f\x93\xe6\xf7\x08\x04\x89\x53\x80\x94\x71\x1a\xe6\xd5\xb4\x47\x96\xae" +
		"\xbd\xea\x33\xa7\x83\x34\x68\x70\x7b\xef\xde\x7f\xfa\x16\x7f\xbf\xe0\xda\x0f\x06\x77\xa3\x86\xb2\xe7\x0e\xde\x45" +
		"\x53\xf6\xfd\x82\xeb\x3e\x18\xdc\x8d\x1a\x34\xb8\xf4\x43\x38\xf8\xf5\x3b\x06\xd7\x5f\x45\xc6\xba\x3f\xb8\x1b\x35" +
		"\x68\x70\x6f\x66\xfd\xde\xe5\xe8\x7b\x05\x17\x43\xb4\x0a\x0d\xbe\x2f\xb8\x5b\x34\x68\x70\xd3\xf4\xcd\xcf\xd7\xdf" +
		"\x2f\xb6\xe4\x81\x1d\x77\x8b\x86\xb2\xe3\xf6\xd2\x8f\x97\xbf\xec\x05\xd7\xb5\x87\x53\x40\xfc\xdd\xf0\xae\x09\xed" +
		"\xd7\xf2\x4c\x20\x3e\xf2\x90\xe1\xae\x98\xfe\xf6\xf1\xe3\x5e\x88\x03\x07\x06\xbe\x83\x10\x72\x31\xb2\x1c\xec\xb8" +
		"\xc3\x29\xb0\x5d\xe8\x60\x84\x10\x22\xc4\xf7\x7d\xec\x5b\xee\x6e\x06\xee\xa7\xe3\xc1\x3e\x3c\x1b\xfe\x8e\x3a\x2a" +
		"\x89\xbe\xe4\x6f\x07\x37\x8d\xf8\xf3\xb0\x4b\x5c\xe2\x63\xcf\xf2\x65\x37\x07\xbe\x05\xed\x00\x21\x84\x5d\xcf\xb7" +
		"\x5d\x77\x4f\xfe\xdc\x47\xc3\x03\xed\x3f\x1b\xe6\x8e\x3a\xe4\xf9\xc2\x47\xc5\x7e\xe2\xd4\x5d\x09\x23\x0c\x03\x2f" +
		"\x08\x02\xc7\x23\x9e\x63\x05\x68\x4f\xa6\xed\x69\x74\x7f\x2b\xcf\x86\x91\xe3\x8e\x93\x92\xeb\xaf\xc9\x5e\x46\x2c" +
		"\x02\x91\x15\xf8\x04\xbb\xc4\xf6\x1c\xc7\x77\xef\x9d\x4a\xf7\xd1\xf0\x40\xfb\xcf\x84\xb8\x23\x8f\xc1\x68\xfc\xdb" +
		"\xbb\x3f\xfe\xdd\x8c\x39\xcb\x71\x08\x76\x65\x3f\xf7\x91\x73\xc0\xb7\xd8\x3d\x95\x3c\xdc\x8b\x67\x43\x21\x79\x60" +
		"\xee\xed\x1c\xea\xb1\x8b\xcb\xcf\x5f\xe6\x6b\x50\xeb\x2b\xf1\x4b\x4b\x4b\x61\x9e\x87\xb3\xc5\x3a\xd4\x86\x30\x96" +
		"\x97\xb4\x36\xef\x3e\x03\x04\xad\xb5\x57\x71\xb6\xd8\x1a\x40\x10\x37\x6c\x54\x03\xe0\xc9\x6d\xd7\xe2\x37\x99\xa0" +
		"\xb9\xbe\x8a\xb2\x11\x17\x6d\x80\xcf\xe6\xbb\x07\x0d\x84\xca\x1b\x61\x91\x45\xf1\xf4\xc9\xc8\x70\x0e\x01\xa4\xd6" +
		"\xe8\x30\x32\x1e\xd3\xf6\x71\xc8\xb8\x7d\xf3\xdf\xff\x0c\xae\x9f\x8c\x0c\x6f\x35\x34\xbb\x11\x20\xb5\x46\x87\x91" +
		"\xf1\x98\xb6\x8f\x43\x46\xff\x77\xfa\x26\x7f\xba\xcc\x08\x0e\xe9\x9d\xb5\x46\x87\x91\xf1\x98\xb6\x8f\x74\x9b\xfa" +
		"\xdc\xfb\x7a\x79\xf7\x44\x64\x60\x88\xef\x0f\xc8\x86\x46\x87\x90\xf1\xb8\xb6\x8f\x94\x19\x17\xff\xfa\xfa\x56\x3c" +
		"\x19\x19\xd6\x21\x80\xd4\x1a\x1d\x46\xc6\x63\xda\x3e\x0e\x19\xf1\xb7\x8b\x9f\xdf\xfd\xd2\x9c\x0c\xb4\x95\x8c\xe5" +
		"\xcd\xd4\x79\x64\x8c\xd7\x82\x73\x57\x5e\x5e\x0d\x91\xfa\x58\xb0\xd6\xa4\xc6\xc5\x53\x9b\x6e\x40\x05\xe3\x09\xe3" +
		"\x4c\xd0\x26\x8c\xac\xc9\x6a\x62\xde\xff\xc9\xaf\xef\xf1\xf5\xf1\x50\x5e\xd6\x06\x2a\xa4\x01\x38\xb5\x26\x07\xf2" +
		"\xf2\x78\xa6\x9f\x82\x97\xfe\xed\xf5\xc5\x9b\x5f\x9f\x8e\x18\x0f\xad\xbc\x9a\xa0\x53\x6b\x72\x20\x31\x8f\x67\xfa" +
		"\x09\x88\x79\x91\x66\x61\xcc\x78\xcf\x9c\xe4\xd2\xc9\xbc\x64\x68\x98\x15\x4c\xc7\x1c\xde\x14\x59\x3a\x12\x54\x6b" +
		"\x9f\xb0\x58\xf4\xdb\x00\x23\x34\x9e\xe8\x92\x3e\x65\xbd\xbe\xd0\x45\x7d\x5d\x14\xb3\x62\x98\x4a\x16\xe7\xfe\xdf" +
		"\x64\xd3\x2d\x35\x49\x4a\x6b\x55\x83\x42\x15\xd7\x1b\x2d\x84\x2b\x60\x52\xd6\x53\x81\x0e\x8a\x36\x88\x28\x17\x34" +
		"\x5f\x15\xb8\xc9\xa6\x5a\x68\xad\xba\xb4\xb1\xa9\xae\x91\xce\x61\x18\xdd\x6d\xae\xbd\x1d\x15\x82\x25\x33\xb3\x3c" +
		"\x9e\xb3\xc5\x6e\xbd\xfd\xee\x76\x4b\x70\x99\x31\xcb\x69\xa4\xc9\x89\xb2\x74\x34\xe0\x6b\xba\xb7\xd4\x6f\xaf\x2b" +
		"\xf7\xd9\xb2\xde\xfc\x27\x50\x25\xcb\x9e\x33\x5f\x23\x99\x93\xbc\x28\xd1\xc7\x46\xcc\x9b\x4c\x88\x6c\xd0\x06\xa4" +
		"\xda\x39\x5b\xcf\xed\x2d\xbf\x74\x6f\x7e\x34\x64\xaf\xb4\x76\x5f\xd0\xa9\x58\x3e\x20\x53\xb0\x6f\xb4\x0d\x2c\x32" +
		"\xf7\x4b\x15\x4e\xca\x38\x08\x2a\xbd\x93\xad\x36\x75\x82\x28\x4b\xb3\xbc\x0d\xf2\xde\x4d\x78\x4a\x1c\xe7\x0c\x2c" +
		"\xfe\x21\xe8\xbe\xfe\x67\x43\x8d\x2f\x93\x0f\x49\xb4\x1c\xec\x3c\x2e\x5b\xc7\xb5\x74\x24\xa1\x13\xb3\x31\x60\x71" +
		"\xd7\x58\xcb\x73\x7d\xbe\xab\x18\xf7\x40\x94\x86\x45\x21\xeb\x7b\x99\x01\xc6\x8c\x4e\xde\x66\xd3\xae\x81\x00\x02" +
		"\x98\xf8\xf2\xcf\x00\xd3\x41\xca\x8b\xf6\x34\x65\xfc\xae\x6b\xf4\x85\x18\xb6\x5b\xad\xc9\x64\x02\x27\x16\xcc\xf2" +
		"\x5e\x0b\x07\x41\xd0\x52\xb5\x4a\xeb\x8e\x73\x63\xe0\x2a\xeb\x65\x4b\x87\xc7\x00\xe8\xc4\x34\x29\xf4\xa5\x3e\xdd" +
		"\x45\xc3\xfc\x32\x0f\x63\x46\xb9\xd0\x7e\xaf\x14\x99\xd8\x00\x53\xdc\x35\x6c\xe8\xbb\x2f\x0d\x30\x25\x5d\x23\x70" +
		"\x21\xc1\x2f\x0d\x30\xc3\x5d\x03\xc9\x77\x59\x16\x40\xd7\x7d\x69\xcc\xf5\xca\x48\x45\x36\x04\xf2\x9f\xa9\xa8\xed" +
		"\x1a\x2f\x2e\x10\x42\x3e\x31\x74\x61\x49\x65\xd7\x80\xbe\x01\xf4\xaf\xab\x94\xb6\x73\x09\x64\x36\x6c\xa4\x68\xde" +
		"\x0c\xa3\x7a\xc3\x4e\x6b\x35\x8e\x32\xf8\xd6\x22\xfa\x4e\xaf\x92\x95\x29\xaa\x62\x9f\x4b\x1b\x3a\x23\xbb\x06\x26" +
		"\x1e\x0c\x5c\xa3\xcc\xc7\xc5\xe7\x59\xd7\xc0\x06\x48\x58\x9a\x76\x8d\x51\x9e\x9e\xbe\x58\x47\xed\xb5\x01\xf2\x69" +
		"\xd7\xb0\xa5\x5f\x52\x7d\x65\x6b\x18\x8a\xbe\xb2\xf5\x36\xcb\x63\x9a\xcf\x95\x54\x41\xc9\x8f\x66\x3e\x4a\x69\xd7" +
		"\xe0\x19\xff\x46\xf3\xcc\x00\x71\xd7\xf8\xc5\x86\x1e\x20\xd0\xb7\x23\x13\x43\xc7\x07\xc8\x24\xd0\x77\x01\x86\xc4" +
		"\xd7\x57\x04\xfa\xce\x18\x63\x17\x3a\x5e\x84\xe4\xd0\xc3\x53\x95\xaa\x8d\xaa\x54\x57\xfd\x52\x42\xd5\x23\x55\x64" +
		"\x62\x48\x5c\x7d\x45\xa0\x6f\x7d\x76\xa0\xeb\x45\x48\x5a\x71\x54\x95\x2a\x5d\xfc\xfb\x60\x43\xd7\xfb\xa6\xdc\x41" +
		"\x5a\x9b\x1f\x59\x10\xdb\x00\x01\x07\xba\xd2\x9e\xe3\xe8\x2b\x07\x7a\xe3\x52\x00\x01\x29\x62\x12\xe8\xd8\xaa\xce" +
		"\x2c\x05\x5c\x5f\xaa\xf3\x23\xd3\x82\xd8\x02\x48\x15\x2b\x29\xb3\x92\x92\xfe\xf8\xef\x4c\x2c\x15\xcb\x78\x1d\x07" +
		"\x20\xa0\xac\x7f\x93\xd8\x4a\x38\x57\xb1\x2d\xd3\x4b\xef\xb4\x1b\x60\xda\x35\x5c\x5b\x11\x26\x13\xab\x04\x3b\x49" +
		"\x12\x8d\xaa\x6b\x03\xcb\x8d\x4c\x1b\xda\x04\x20\xd3\x37\x2d\xe8\xf8\xa6\x6f\xfa\x85\xbe\x00\xea\x0f\xc8\x0f\x40" +
		"\x7e\xd0\x17\xb2\x4c\x76\xe3\x99\x24\x69\xd3\x4f\x8e\x4f\x31\x42\xc3\xe9\x19\x50\x6f\xaf\x7f\xda\xe9\xa7\xde\xae" +
		"\x56\x7e\xd6\xb6\x83\x95\xdb\xf5\x85\xf4\x7a\x18\x7e\x00\x91\x0d\x1c\x04\x1d\x12\x99\x04\x12\xd3\x92\x74\xc2\xc0" +
		"\xf4\xa1\x67\x01\x02\x03\xdb\xc4\x08\x06\x2e\xb0\x34\x8d\x04\xf8\xd0\x23\x26\x0c\x80\x2c\x76\x94\x04\x90\xc5\xb2" +
		"\x1d\x0c\x64\xad\x14\x0b\x6c\x55\xef\x4a\x75\x52\x88\x48\x7d\x2e\x0c\x94\x32\x4f\x09\x18\x55\xa2\xea\x71\xe5\xa3" +
		"\x60\x12\xa9\x2d\xe0\x0a\x93\x95\x6d\x56\x85\xc9\xfa\xee\x50\x1d\x11\x8c\x88\x8c\xc3\x83\x8e\x02\x44\x05\x66\x7a" +
		"\x10\x03\x07\x62\x47\xf9\x6f\x6b\x5c\x2c\xdd\xf1\x89\xe9\xc8\x1e\xea\x41\x4c\xcc\x25\xbc\xe6\x58\x02\x55\x21\xdb" +
		"\x6a\xc4\x34\x38\x96\x4e\x3c\x62\x49\xad\xae\x92\x01\x0b\xd4\x8e\x02\x8d\xda\x63\x5d\xee\xd5\xf5\x3d\xce\x8d\xbd" +
		"\x1c\x63\xa4\xba\xb9\x55\x76\x73\x57\xf6\x72\x88\x64\xaf\xb7\xa1\x2d\x43\x71\x3c\x5d\xa0\xcb\x8b\x45\xaf\x87\x88" +
		"\x44\x32\xe5\x24\x56\x8e\xa7\x3e\x9b\xba\xf8\x28\x11\xaa\x3d\x4b\x15\x61\x6d\x5f\xb0\x21\xf7\xc4\x81\x78\x41\x7d" +
		"\x2d\x17\xac\xc7\xcb\x05\xbb\xcc\x05\xe7\x78\xb9\xa0\x77\x02\x17\x78\xac\xec\xb3\x35\xbe\x41\x58\xbe\x8c\x78\x71" +
		"\x83\xa8\xf2\x81\x2c\x12\xc2\xa9\x12\x82\xcc\x13\x82\x54\x09\x41\x9a\x26\x84\x55\x25\x84\xfd\x24\x09\xa1\xb7\xd9" +
		"\x8c\x72\x4a\x3b\x8f\xb9\x9c\xe0\x96\xdf\xe2\x76\x55\x20\xbf\xa6\xa3\x70\xd8\x35\xd4\xe8\x68\xa5\xf8\x36\x63\xbc" +
		"\x2a\x57\x59\x63\xc9\x2f\x1b\xe2\x41\x27\xb5\xe4\xcd\x04\x07\xa6\x7c\x37\x71\xb0\xe7\x9b\x47\xed\x35\x1d\xc3\xa3" +
		"\xc0\x05\xb6\x3b\xb6\x7c\xd3\xf2\x77\xbb\xa0\x77\x58\x8e\xe3\x02\xb4\x1d\xe0\xdb\xd0\x49\xcd\x12\x15\xd0\x08\x15" +
		"\xbd\xcf\x70\x14\x9e\x6c\xe9\x12\x46\xd6\xdc\x27\xb3\xf4\x09\xec\x65\x4a\x2d\xb7\x1f\xc3\x27\x8b\x00\xdf\xfe\xac" +
		"\xc8\xda\x03\x8b\x5a\x64\x3e\x8a\x0b\x6a\xac\x24\x6f\x0e\x69\x09\x8a\xd9\x08\x14\xbd\xd2\x6a\xa8\x31\xea\x1f\x39" +
		"\x0b\x79\x2f\xa5\xe6\x5b\x35\xed\xdc\xeb\xa6\xb2\x8b\x80\x6f\xf7\x3d\xb4\xcd\xe9\x5d\xb6\xd5\x62\xe2\x9a\xe9\x2b" +
		"\x9a\xec\xc7\x47\xf5\x03\x07\x10\xf7\x4a\x9a\xf7\x0e\xb1\xad\x17\xcc\xd6\x8c\xff\x26\x87\xfd\x8d\xac\x07\x3e\xf0" +
		"\xbd\x2b\xd7\x02\xc4\x6d\x64\xbd\xd3\x52\xf3\x8f\x4e\xab\x18\xeb\x0b\x39\x6f\x2c\x5d\x91\x93\x64\xe3\xfc\x4a\x4f" +
		"\x20\xb5\x74\x31\x0c\x79\xc5\x90\x9a\x8e\x1a\x1b\x1f\x13\x22\xe5\xb4\xde\x98\x8d\xcf\x4f\xca\xb7\xc5\xac\x34\xcf" +
		"\x32\xa1\x1e\xb9\x31\x7f\xe0\xc6\xd2\xd1\xf9\xdb\x70\x1c\xea\xd2\x72\x22\x37\x61\x3c\xce\x26\x30\x8c\xe3\x8b\x31" +
		"\xe5\xe2\x8a\x15\x82\x72\x9a\x9f\xbe\x92\x73\xdb\x57\x67\x20\x19\x71\xb5\xb6\x01\x4e\xa9\xac\x7f\x0d\xfe\x3a\xa9" +
		"\x6e\xf1\x51\xc6\x0b\x01\xca\x49\xf0\x97\x72\xad\xab\x0b\xe2\x2c\x1a\x0d\x28\x17\xb0\x47\xc5\x45\x4a\xe5\xe5\xdb" +
		"\xd9\xc7\x58\xab\x5c\x9a\x2e\xbf\xaa\x8e\x76\x81\x35\x25\x50\xa1\x20\x7d\x91\x8e\x9d\xbe\x2a\xcf\xe2\xca\x06\xeb" +
		"\xd6\x65\xb0\xbb\x6c\xca\xfa\x65\x43\xf2\xf3\xba\xfa\xe5\x67\x34\x28\x1b\x73\xe1\x12\xfd\x05\xf8\x90\x71\x26\x4e" +
		"\xa5\x8e\xb3\xea\x08\x91\x7e\x51\x1e\x0f\x33\xc6\x45\x1b\xbc\x9a\x3f\xc4\xe4\xd5\x59\x25\xf1\xf7\xfc\xa8\xda\xeb" +
		"\x9f\x4e\x16\x8f\x2a\xe9\xb4\xca\x27\x94\xb4\xd4\x63\x6b\xfe\x17\x00\x00\xff\xff\x7d\xa0\x6e\xee\xce\x46\x00\x00" +
		"")

func gzipBindataDataGraphqlplaygroundhtml() (*gzipAsset, error) {
	bytes := _gzipBindataDataGraphqlplaygroundhtml
	info := gzipBindataFileInfo{
		name:        "data/graphql-playground.html",
		size:        18126,
		md5checksum: "",
		mode:        os.FileMode(420),
		modTime:     time.Unix(1567041319, 0),
	}

	a := &gzipAsset{bytes: bytes, info: info}

	return a, nil
}

// GzipAsset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func GzipAsset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _gzipbindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("GzipAsset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// MustGzipAsset is like GzipAsset but panics when GzipAsset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
func MustGzipAsset(name string) []byte {
	a, err := GzipAsset(name)
	if err != nil {
		panic("asset: GzipAsset(" + name + "): " + err.Error())
	}

	return a
}

// GzipAssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
func GzipAssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _gzipbindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("GzipAssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

// GzipAssetNames returns the names of the assets.
// nolint: deadcode
func GzipAssetNames() []string {
	names := make([]string, 0, len(_gzipbindata))
	for name := range _gzipbindata {
		names = append(names, name)
	}
	return names
}

//
// _gzipbindata is a table, holding each asset generator, mapped to its name.
//
var _gzipbindata = map[string]func() (*gzipAsset, error){
	"data/facebook-open-graph-template.html": gzipBindataDataFacebookopengraphtemplatehtml,
	"data/graphql-playground.html":           gzipBindataDataGraphqlplaygroundhtml,
}

// GzipAssetDir returns the file names below a certain
// directory embedded in the file by bindata.
// For example if you run bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then GzipAssetDir("data") would return []string{"foo.txt", "img"}
// GzipAssetDir("data/img") would return []string{"a.png", "b.png"}
// GzipAssetDir("foo.txt") and GzipAssetDir("notexist") would return an error
// GzipAssetDir("") will return []string{"data"}.
func GzipAssetDir(name string) ([]string, error) {
	node := _gzipbintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op:   "open",
					Path: name,
					Err:  os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type gzipBintree struct {
	Func     func() (*gzipAsset, error)
	Children map[string]*gzipBintree
}

var _gzipbintree = &gzipBintree{Func: nil, Children: map[string]*gzipBintree{
	"data": {Func: nil, Children: map[string]*gzipBintree{
		"facebook-open-graph-template.html": {Func: gzipBindataDataFacebookopengraphtemplatehtml, Children: map[string]*gzipBintree{}},
		"graphql-playground.html":           {Func: gzipBindataDataGraphqlplaygroundhtml, Children: map[string]*gzipBintree{}},
	}},
}}
