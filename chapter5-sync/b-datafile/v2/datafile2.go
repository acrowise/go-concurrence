package v2

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Data 代表数据的类型。
type Data []byte

// DataFile 代表数据文件的接口类型。
type DataFile interface {

	Read() (rsn int64, d Data, err error) // Read 会读取一个数据块。

	Write(d Data) (wsn int64, err error)  // Write 会写入一个数据块。

	RSN() int64 // RSN 会获取最后读取的数据块的序列号。

	WSN() int64  // WSN 会获取最后写入的数据块的序列号。

	DataLen() uint32 // DataLen 会获取数据块的长度。

	Close() error // Close 会关闭数据文件。
}

// myDataFile 代表数据文件的实现类型。
type myDataFile struct {
	f       *os.File     // 文件。
	fmutex  sync.RWMutex // 被用于文件的读写锁。
	rcond   *sync.Cond   //读操作需要用到的条件变量
	woffset int64        // 写操作需要用到的偏移量。
	roffset int64        // 读操作需要用到的偏移量。
	wmutex  sync.Mutex   // 写操作需要用到的互斥锁。
	rmutex  sync.Mutex   // 读操作需要用到的互斥锁。
	dataLen uint32       // 数据块长度。
}

// NewDataFile 会新建一个数据文件的实例。
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: dataLen}
	//条件变量的创建
	df.rcond = sync.NewCond(df.fmutex.RLocker()) //让读写锁中的读锁与条件变量绑定
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量。
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块。
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	df.fmutex.RLock()//锁定读锁
	defer df.fmutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				//df.rcond.Wait()语句，添加这条语句的意义在于：当出现EOF错误时，让当前goroutine暂时"放弃"fmutex的读锁并等待通知的到来。
				//"放弃"fmutex的读锁，也就意味着Write方法中的数据块写操作不会受它的阻碍了。
				//一旦有新的写操作完成，应该及时向条件变量rcond发送通知，以唤醒为此而等待的goroutine。
				//请注意，在某个goroutine被换醒之后，应该再次检查需要满足的条件。
				//在这里，这个需要满足的条件是：在进行文件内容读取时不会出现EOF错误。
				//如果该条件满足，就可以进行后续的操作了；否则，再次"放弃"读锁并等待通知。这也是我依然保留for循环的原因。
				//todo 注意此处有两个小细节额，当进入等待的时候，会自动解开读锁的，并进入阻塞等待
				//todo  当收到通知，会自动锁定读锁，并往下执行到continue
				df.rcond.Wait()
				continue
			}
			return
		}
		d = bytes
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量。
	var offset int64
	df.wmutex.Lock()
	offset = df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()

	//写入一个数据块。
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	//todo
	df.rcond.Signal()
	return
}

func (df *myDataFile) RSN() int64 {
	df.rmutex.Lock()
	defer df.rmutex.Unlock()
	return df.roffset / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()
	return df.woffset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}
