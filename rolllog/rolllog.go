package rolllog

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type stFileConfig struct {
	maxSize  uint32
	maxCnt   uint32
	path     string
	baseName string
}

type RollWriter struct {
	sync.RWMutex
	fileConf stFileConfig
	fp       *os.File
	currSize uint32
}

func NewRollWriter(path string, baseName string) (*RollWriter, error) {
	w := &RollWriter{}
	err := w.Init(path, baseName)

	return w, err
}

func (w *RollWriter) Write(buf []byte) (n int, err error) {
	if w.currSize >= w.fileConf.maxSize {
		w.Lock()
		//fmt.Printf("curr_size:%d, maxSize:%d\n", w.currSize, w.fileConf.maxSize)
		if w.currSize >= w.fileConf.maxSize {
			w.roll()
			w.reOpen()
		}
		w.Unlock()
	}

	w.RLock()
	n, err = w.fp.Write(buf)
	atomic.AddUint32(&w.currSize, uint32(n))
	w.RUnlock()

	return n, err
}

func (w *RollWriter) Close(p []byte) error {
	return nil
}

func (w *RollWriter) Init(path string, baseName string) error {
	w.Lock()
	defer w.Unlock()
	w.fileConf.path = path
	w.fileConf.baseName = baseName
	w.fileConf.maxSize = 50 * 1024 * 1024
	w.fileConf.maxCnt = 10

	return w.reOpen()
}

func (w *RollWriter) SetMaxFileSize(size uint32) {
	w.fileConf.maxSize = size
}

func (w *RollWriter) SetMaxFileCnt(count uint32) {
	w.fileConf.maxCnt = count
}

func (w *RollWriter) reOpen() error {
	os.MkdirAll(w.fileConf.path, os.ModePerm)

	PathName := w.fileConf.path + "/" + w.fileConf.baseName + ".log"
	fmt.Printf("open file:%s\n", PathName)
	fp, err := os.OpenFile(PathName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to open file!err:%+v,file:%s\n", err, PathName)
		return err
	}

	fileInfo, err := w.fp.Stat()
	//w.Lock()
	if w.fp != nil {
		w.fp.Close()
	}
	w.fp = fp
	w.currSize = 0
	if err == nil {
		fmt.Printf("szie:%d\n", uint32(fileInfo.Size()))
		w.currSize = uint32(fileInfo.Size())
	}
	//w.Unlock()

	return nil
}

func (w *RollWriter) roll() {
	if w.fp != nil {
		w.fp.Close()
		w.fp = nil
	}

	basePathName := w.fileConf.path + w.fileConf.baseName

	for i := w.fileConf.maxCnt - 1; i > 0; i-- {
		var srcFile, dstFile string
		dstFile = basePathName + "_" + strconv.Itoa(int(i)) + ".log"
		if i != 1 {
			srcFile = basePathName + "_" + strconv.Itoa(int(i)-1) + ".log"
		} else {
			srcFile = basePathName + ".log"
		}
		fmt.Printf("mv %s %s\n", srcFile, dstFile)
		err := os.Rename(srcFile, dstFile)
		if err != nil {
			fmt.Printf("failed to mv file!srcFile:%s,dstFile:%s,err:%+v\n", srcFile, dstFile, err)
		}
	}
}
