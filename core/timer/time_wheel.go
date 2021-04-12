package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
	此时间轮请勿多服务复用,建议为单服务开启自己的时间轮操作事件
*/

type TimeWheel struct {
	element_cnt_per_wheel []uint32   //每个时间轮的槽(元素)数量。在 256+64+64+64+64 = 512 个槽中，表示的范围为 2^32
	right_shift_per_wheel []uint32   //当指针指向当前时间轮最后一位数，再走一位就需要向上进位。每个时间轮进位的时候，使用右移的方式，最快实现进位。这里是每个轮的进位二进制位数
	base_per_wheel        []uint32   //记录每个时间轮指针当前指向的位置
	mutex                 sync.Mutex //加锁
	rwmutex               sync.RWMutex
	newest                []uint32  //每个时间轮当前指针所指向的位置
	timewheels            [][]*Node //定义时间轮
	TimerMap              map[string]*Node
	bRun                  bool
	cRun                  int32
	Mutex                 sync.Mutex
	running               bool
	stop_signl            chan bool
}

func (this *TimeWheel) SetTimer(name string, inteval uint32, handler func(interface{}), args interface{}) {
	if inteval <= 0 {
		return
	}

	var bucket_no uint8 = 0
	var offset uint32 = inteval
	var left uint32 = inteval
	for offset >= this.element_cnt_per_wheel[bucket_no] { //偏移量大于当前时间轮容量，则需要向高位进位
		offset >>= this.right_shift_per_wheel[bucket_no] //计算高位的值。偏移量除以低位的进制。比如低位当前是256，则右移8个二进制位，就是除以256，得到的结果是高位的值。
		var tmp uint32 = 1
		if bucket_no == 0 {
			tmp = 0
		}
		left -= this.base_per_wheel[bucket_no] * (this.element_cnt_per_wheel[bucket_no] - this.newest[bucket_no] - tmp)
		bucket_no++
	}
	if offset < 1 {
		return
	}
	if inteval < this.base_per_wheel[bucket_no]*offset {
		return
	}
	left -= this.base_per_wheel[bucket_no] * (offset - 1)
	pos := (this.newest[bucket_no] + offset) % this.element_cnt_per_wheel[bucket_no] //通过类似hash的方式，找到在时间轮上的插入位置

	var node Node
	node.SetData(Timer{name, left, handler, args})

	this.rwmutex.RLock()
	this.TimerMap[name] = this.timewheels[bucket_no][pos].InsertHead(node) //插入定时器
	this.rwmutex.RUnlock()
}

func (this *TimeWheel) step() {
	//var dolist list.List
	{
		this.rwmutex.RLock()
		//遍历所有桶
		var bucket_no uint8 = 0
		for bucket_no = 0; bucket_no < wheel_cnt; bucket_no++ {
			this.newest[bucket_no] = (this.newest[bucket_no] + 1) % this.element_cnt_per_wheel[bucket_no] //当前指针递增1
			//fmt.Println(newest)
			var head *Node = this.timewheels[bucket_no][this.newest[bucket_no]] //返回当前指针指向的槽位置的表头
			var firstElement *Node = head.Next()
			for firstElement != nil { //链表不为空
				if value, ok := firstElement.Data().(Timer); ok { //如果element里面确实存储了Timer类型的数值，那么ok返回true，否则返回false。
					inteval := value.Inteval
					doSomething := value.DoSomething
					args := value.Args
					if nil != doSomething { //有遇到函数为nil的情况，所以这里判断下非nil
						if 0 == bucket_no || 0 == inteval {
							//dolist.PushBack(value) //执行自定义处理函数
							go doSomething(args)
						} else {
							SetTimer(value.Name, inteval, doSomething, args) //重新插入计时器
						}
					}
					Delete(firstElement) //删除定时器
				}
				firstElement = head.Next() //重新定位到链表第一个元素头
			}
			if 0 != this.newest[bucket_no] { //指针不是0，还未转回到原点，跳出。如果回到原点，则说明转完了一圈，需要向高位进位1，则继续循环入高位步进一步。
				break
			}
		}
		this.rwmutex.RUnlock()
	}
}

func (this *TimeWheel) run() {
	this.running = true
	var i int = 0
	for this.bRun {
		go this.step()
		i++
		//log.Printf("第%ds", i)
		//间隔时间inteval=1s
		time.Sleep(1 * time.Second)
	}
	this.running = false
	this.stop_signl <- true
}

func (this *TimeWheel) Start() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	atomic.AddInt32(&(this.cRun), 1)
	if !this.bRun {
		this.bRun = true
	} else {
		return
	}

	this.stop_signl = make(chan bool)
	go this.run()
}

func (this *TimeWheel) Stop() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	newRun := atomic.AddInt32(&(this.cRun), -1)
	if newRun <= 0 {
		this.bRun = false
		if this.running {
			<-this.stop_signl
		}
		this.Init()
	}
}

func (this *TimeWheel) Init() {
	this.running = false
	this.cRun = 0
	this.newest = []uint32{0, 0, 0, 0, 0}
	this.element_cnt_per_wheel = []uint32{256, 64, 64, 64, 64}                          //每个时间轮的槽(元素)数量。在 256+64+64+64+64 = 512 个槽中，表示的范围为 2^32
	this.right_shift_per_wheel = []uint32{8, 6, 6, 6, 6}                                //当指针指向当前时间轮最后一位数，再走一位就需要向上进位。每个时间轮进位的时候，使用右移的方式，最快实现进位。这里是每个轮的进位二进制位数
	this.base_per_wheel = []uint32{1, 256, 256 * 64, 256 * 64 * 64, 256 * 64 * 64 * 64} //记录每个时间轮指针当前指向的位置
	this.TimerMap = make(map[string]*Node)
	this.timewheels = make([][]*Node, wheel_cnt)

	var bucket_no uint8 = 0
	for bucket_no = 0; bucket_no < wheel_cnt; bucket_no++ {
		this.timewheels[bucket_no] = make([]*Node, 0)
		var i uint32 = 0
		for ; i < this.element_cnt_per_wheel[bucket_no]; i++ {
			this.timewheels[bucket_no] = append(this.timewheels[bucket_no], new(Node))
		}
	}
}

func NewTimeWheel() *TimeWheel {
	timeWheel := &TimeWheel{}
	timeWheel.Init()
	return timeWheel
}
