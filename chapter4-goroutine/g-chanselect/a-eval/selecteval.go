package main

import "fmt"

var intChan1 chan int
var intChan2 chan int
var channels = []chan int{intChan1, intChan2}

var numbers = []int{1, 2, 3, 4, 5}


//分支选择规则
//(1)在开始执行select语句的时候，所有跟在case关键字右边的发送语句或接收语句中的通道表达式和元素表达式都会先求值(求值的顺序是从左到右、自上而下的)，无论它们所在的case是否有可能被选择都会是这样
//(2)通道intChan1和intChan2都未被初始化，向他们发送元素值的操作会被永久阻塞，所以select被执行时最终选择了default case,因为其他两个case走不通
//(3)在执行select语句的时候，运行时系统会自上而下地判断每个case中的发送或接收操作是否可以立即进行。这里的"立即进行"，指的是当前goroutine不会因此操作而被阻塞。
//   这个判断还需要依据通道的具体特性(缓冲或非缓冲)以及那一时刻的具体情况来进行。只要发现有个case上的判断是肯定的，该case就会被选中
//(4)当有一个case被选中时，运行时系统就会执行该case及其包含的语句，而其他case会被忽略。如果同时有多个case满足条件，那么运行时系统会通过一个伪随机的算法选中一个case。
func main() {
	select {
	case getChan(0) <- getNumber(0): //把切片numbers中的第0个元素发送到切片channels中的第0个元素中
		fmt.Println("The 1th case is selected.")
	case getChan(1) <- getNumber(1):  //把切片numbers中的第1个元素发送到切片channels中的第1个元素中
		fmt.Println("The 2nd case is selected.")
	default:
		fmt.Println("Default!")
	}
}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]  ", i)
	return channels[i]
}
