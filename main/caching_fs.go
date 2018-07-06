package main

import (
	"os"
	"path/filepath"
	"sync"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type CachingFileSystem struct {
	fs boshsys.FileSystem

	readCache     map[string][]byte
	readCacheLock sync.RWMutex

	globCache     map[string][]string
	globCacheLock sync.RWMutex

	logTag string
	logger boshlog.Logger
}

var _ boshsys.FileSystem = &CachingFileSystem{}

func NewCachingFileSystem(fs boshsys.FileSystem, logger boshlog.Logger) *CachingFileSystem {
	return &CachingFileSystem{
		fs: fs,

		readCache: map[string][]byte{},
		globCache: map[string][]string{},

		logTag: "CachingFileSystem",
		logger: logger,
	}
}

func (f *CachingFileSystem) DropCache() {
	f.logger.Info(f.logTag, "Reloading data")

	f.readCacheLock.Lock()
	f.readCache = map[string][]byte{}
	f.readCacheLock.Unlock()

	f.globCacheLock.Lock()
	f.globCache = map[string][]string{}
	f.globCacheLock.Unlock()
}

func (f *CachingFileSystem) ChangeTempRoot(path string) error {
	return f.fs.ChangeTempRoot(path)
}

func (f *CachingFileSystem) HomeDir(username string) (path string, err error) {
	return f.fs.HomeDir(username)
}

func (f *CachingFileSystem) MkdirAll(path string, perm os.FileMode) (err error) {
	return f.fs.MkdirAll(path, perm)
}

func (f *CachingFileSystem) RemoveAll(fileOrDir string) (err error) {
	return f.fs.RemoveAll(fileOrDir)
}

func (f *CachingFileSystem) Chown(path, username string) (err error) {
	return f.fs.Chown(path, username)
}

func (f *CachingFileSystem) Chmod(path string, perm os.FileMode) (err error) {
	return f.fs.Chmod(path, perm)
}

func (f *CachingFileSystem) OpenFile(path string, flag int, perm os.FileMode) (boshsys.File, error) {
	return f.fs.OpenFile(path, flag, perm)
}

func (f *CachingFileSystem) WriteFileString(path, content string) (err error) {
	return f.fs.WriteFileString(path, content)
}

func (f *CachingFileSystem) WriteFile(path string, content []byte) (err error) {
	return f.fs.WriteFile(path, content)
}

func (f *CachingFileSystem) ExpandPath(path string) (string, error) {
	return f.fs.ExpandPath(path)
}

func (f *CachingFileSystem) ConvergeFileContents(path string, content []byte, opts ...boshsys.ConvergeFileContentsOpts) (written bool, err error) {
	return f.fs.ConvergeFileContents(path, content, opts...)
}

func (f *CachingFileSystem) ReadFileString(path string) (string, error) {
	bytes, err := f.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (f *CachingFileSystem) ReadFile(path string) ([]byte, error) {
	f.readCacheLock.Lock()
	defer f.readCacheLock.Unlock()

	if content, found := f.readCache[path]; found {
		f.logger.Debug(f.logTag, "hit: read[%s]", path)
		return content, nil
	} else {
		f.logger.Debug(f.logTag, "miss: read[%s]", path)
	}

	content, err := f.fs.ReadFile(path)
	if err == nil {
		f.readCache[path] = content
	}

	return content, err
}

func (f *CachingFileSystem) ReadFileWithOpts(path string, opts boshsys.ReadOpts) (content []byte, err error) {
	return f.fs.ReadFileWithOpts(path, opts)
}

func (f *CachingFileSystem) FileExists(path string) bool {
	return f.fs.FileExists(path)
}

func (f *CachingFileSystem) Rename(oldPath, newPath string) (err error) {
	return f.fs.Rename(oldPath, newPath)
}

func (f *CachingFileSystem) Symlink(oldPath, newPath string) (err error) {
	return f.fs.Symlink(oldPath, newPath)
}

func (f *CachingFileSystem) ReadAndFollowLink(symlinkPath string) (targetPath string, err error) {
	return f.fs.ReadAndFollowLink(symlinkPath)
}

func (f *CachingFileSystem) Readlink(symlinkPath string) (targetPath string, err error) {
	return f.fs.Readlink(symlinkPath)
}

func (f *CachingFileSystem) CopyFile(srcPath, dstPath string) (err error) {
	return f.fs.CopyFile(srcPath, dstPath)
}

func (f *CachingFileSystem) CopyDir(srcPath, dstPath string) error {
	return f.fs.CopyDir(srcPath, dstPath)
}

func (f *CachingFileSystem) TempFile(prefix string) (boshsys.File, error) {
	return f.fs.TempFile(prefix)
}

func (f *CachingFileSystem) TempDir(prefix string) (path string, err error) {
	return f.fs.TempDir(prefix)
}

func (f *CachingFileSystem) Lstat(path string) (os.FileInfo, error) {
	return f.fs.Lstat(path)
}

func (f *CachingFileSystem) Stat(path string) (os.FileInfo, error) {
	return f.fs.Stat(path)
}

func (f *CachingFileSystem) StatWithOpts(path string, opts boshsys.StatOpts) (os.FileInfo, error) {
	return f.fs.StatWithOpts(path, opts)
}

func (f *CachingFileSystem) RecursiveGlob(pattern string) (matches []string, err error) {
	return f.fs.RecursiveGlob(pattern)
}

func (f *CachingFileSystem) WriteFileQuietly(path string, content []byte) error {
	return f.fs.WriteFileQuietly(path, content)
}

func (f *CachingFileSystem) Glob(pattern string) ([]string, error) {
	f.globCacheLock.Lock()
	defer f.globCacheLock.Unlock()

	if matches, found := f.globCache[pattern]; found {
		f.logger.Debug(f.logTag, "hit: glob[%s]", pattern)
		return matches, nil
	} else {
		f.logger.Debug(f.logTag, "miss: glob[%s]", pattern)
	}

	matches, err := f.fs.Glob(pattern)
	if err == nil {
		f.globCache[pattern] = matches
	}

	return matches, err
}

func (f *CachingFileSystem) Walk(root string, walkFunc filepath.WalkFunc) error {
	return f.fs.Walk(root, walkFunc)
}
