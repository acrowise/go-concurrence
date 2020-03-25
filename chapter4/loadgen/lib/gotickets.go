package lib

import (
	"errors"
	"fmt"
)




//==========================================我们定义这个goroutine票池的接口类型==========================
// GoTickets 表示Goroutine票池的接口。
// goroutine票池只负责增减票的数量，并以此真实地体现出正在运行的专用goroutine数量。
type GoTickets interface {
	Take() // 拿走一张票。
	Return() // 归还一张票。
	Active() bool // 票池是否已被激活。
	Total() uint32 // 票的总数。
	Remainder() uint32  // 剩余的票数。
}



//===============================让*myGoTickets成为实现GoTickets接口的类型=====================================================================
// myGoTickets 表示Goroutine票池的实现。
type myGoTickets struct {
	total    uint32        // 票的总数。
	ticketCh chan struct{} // 票的容器。
	active   bool          // 票池是否已被激活。即是否已正确初始化
}

// NewGoTickets 会新建一个Goroutine票池。
//返回一个接口类型
func NewGoTickets(total uint32) (GoTickets, error) {
	gt := myGoTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

//初始化工作交由包级私有的指针方法init来进行
func (gt *myGoTickets) init(total uint32) bool {
	if gt.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

func (gt *myGoTickets) Take() {
	<-gt.ticketCh
}

func (gt *myGoTickets) Return() {
	gt.ticketCh <- struct{}{}
}

func (gt *myGoTickets) Active() bool {
	return gt.active
}

func (gt *myGoTickets) Total() uint32 {
	return gt.total
}

func (gt *myGoTickets) Remainder() uint32 {
	return uint32(len(gt.ticketCh))
}
