package v1

import (
	"errors"
	"io"
	"os"
	"sync"
)

//-------------------------------------------------Data 代表数据的类型。 []byte的别名类型
type Data []byte

//----------------------------------------------DataFile 代表数据文件的接口类型。
type DataFile interface {
	Read() (rsn int64, d Data, err error) // Read 会读取一个数据块。
	Write(d Data) (wsn int64, err error)  // Write 会写入一个数据块。
	RSN() int64 // RSN 会获取最后读取的数据块的序列号。  Reading Serial Number  当前已读取的数据块的数量
	WSN() int64  // WSN 会获取最后写入的数据块的序列号。 Writing Serial Number  当前已写入的数据块的数量
	DataLen() uint32  // DataLen 会获取数据块的长度。
	Close() error  // Close 会关闭数据文件。
}

//-------------------------------------------myDataFile 代表数据文件的实现类型。
//todo  3把锁额
type myDataFile struct {
	f       *os.File     // 文件。
	fmutex  sync.RWMutex // 被用于文件的读写锁。
	woffset int64        // 写操作需要用到的偏移量。 用来记录写操作的进度的
	roffset int64        // 读操作需要用到的偏移量。 用来记录读操作的进度的
	wmutex  sync.Mutex   // 多个程序同时修改写操作进度的时候需要用到该互斥锁额
	rmutex  sync.Mutex   // 多个程序同时修改读操作进度的时候需要用到该互斥锁额
	dataLen uint32       // 数据块长度。每次读取的数据块长度
}




//在读之前已经读了80000个字节了,每次读取200个字节，那么此次的rsn为  80000/200=400

//todo 注意漏掉某个序列号的情形,比如我想读取第1000序列号的数据，
//todo 但是由于此时写的比较慢，压根没有这段数据，则下次再读的时候只会到了1001序列号的数据了，
//todo 当1000序列号写入了数据也没戏了,怎么办，不停的循环直到读到它呗
//读取指定序列号的数据块并返回(序列号由内部计算获得)
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// ①读取并更新读偏移量。
	var offset int64//用来保存读之前已经读至的偏移量
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//②根据读偏移量从文件中读取一块数据
	rsn = offset / int64(df.dataLen) // 已经读过的字节数/每次读取的字节数
	bytes := make([]byte, df.dataLen) //要读取的字节放置位置

	//todo 下面最大的弊端就是需要不停的读锁定读解锁无法通过使用 defer  df.fmutex.RUnlock()来集中解决
	//todo  一旦for语句发生恐慌，无法就无法执行下面的 df.fmutex.RUnlock() 就完蛋了
	for { //这是防止读取不到该序列号的数据块的
		df.fmutex.RLock()
		//ReadAt从指定的位置（相对于文件开始位置）读取len(bytes)字节数据并写入bytes。
		//它返回读取的字节数和可能遇到的任何错误。当返回字节数n<len(bytes)时，本方法总是会返回错误；如果是因为到达文件结尾，返回值err会是io.EOF。
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.fmutex.RUnlock()
				continue
			}
			df.fmutex.RUnlock()
			return
		}
		d = bytes
		df.fmutex.RUnlock()
		return
	}
}


//向文件中写入数据库 d
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
	//Write向文件中写入len(bytes)字节数据。它返回写入的字节数和可能遇到的任何错误。如果返回值n!=len(bytes)，本方法会返回一个非nil的错误。
	_, err = df.f.Write(bytes)
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

//--------------------------------------------NewDataFile 会新建一个数据文件的实例。
//todo 注意这里返回约束是实现了DataFile接口的类型
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	//Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）。
	//如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式。如果出错，错误底层类型是*PathError。
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	//只有*myDataFile类型才实现了接口DataFile，返回myDataFile类型没有实现接口DataFile，所以这里返回 &myDataFile{f: f, dataLen: dataLen}
	df := &myDataFile{f: f, dataLen: dataLen}
	return df, nil
}
