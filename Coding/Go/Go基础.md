# Go语言基础 #

[https://github.com/Unknwon/the-way-to-go_ZH_CN/blob/master/eBook/directory.md](https://github.com/Unknwon/the-way-to-go_ZH_CN/blob/master/eBook/directory.md)

- 常量计数器iota

	iota是常量计数器,只能在常量的表达式中使用。iota在const关键字出现时将被重置为0(const内部的第一行之前)，const中每新增一行常量声明将使iota计数一次(即其值自动加1)。使用iota能简化定义，在定义枚举时很有用。
	
	例：
	1、
	```
	const (
		a = iota	//a=0
		b = iota	//b=1
	)
	```
	可以简写为：
	```
	const (
		a = iota
		b
	)
	```
	2、位掩码表达式
	```
	type Allergen int

	const ( 
	    IgEggs Allergen = 1 << iota         // 1 << 0 which is 00000001 
	    IgChocolate                         // 1 << 1 which is 00000010 
	    IgNuts                              // 1 << 2 which is 00000100 
	    IgStrawberries                      // 1 << 3 which is 00001000 
	    IgShellfish                         // 1 << 4 which is 00010000 
	)
	```
	3、定义数量级
	```
	type ByteSize float64

	const (
	    _           = iota                   // ignore first value by assigning to blank identifier
	    KB ByteSize = 1 << (10 * iota) 		 // 1 << (10*1)
	    MB                                   // 1 << (10*2)
	    GB                                   // 1 << (10*3)
	    TB                                   // 1 << (10*4)
	    PB                                   // 1 << (10*5)
	    EB                                   // 1 << (10*6)
	    ZB                                   // 1 << (10*7)
	    YB                                   // 1 << (10*8)
	)
	```
	
- JSON Tag

	在定义struct的时候，可以在字段后面添加tag，来控制encode/decode：是否要decode/encode某个字段，JSON中的字段名称是什么。可以选择的控制字段有三种：
	`-`：不要解析这个字段
	`omitempty`：当字段为空（默认值）时，不要解析这个字段。比如false;0;nil;长度为0的array、map、slice、string
	`FieldName`：当解析json的时候，使用这个名字

	例1、
	```
	type Product struct {
	Name      string  `json:"name"`
	ProductID int64   `json:"-"` 					
	Number    int     `json:"number,omitempty"`     
	Price     float64 `json:"price"`
	IsOnSale  bool    `json:"is_on_sale,string"`
	}
	
	func main() {
		p := &Product{}
		p.Name = "Apple"
		p.IsOnSale = true
		p.Number = 0
		p.Price = 8999.00
		p.ProductID = 1
		data, _ := json.Marshal(p)
		fmt.Println(string(data))
	}
	```
	序列化之后的结果为：
	`{"name":"Apple","price":8999,"is_on_sale":"true"}`

	2、在某些特殊场景下，结构体中包括一个布尔类型，需要实现当有传递值时进行序列化，没有传递时不进行序列化。
	
	例：需要实现当布尔类型值“is_boot”有赋值时序列化，没有赋值时忽略。
	```
	type JsonType struct {
		UUID	string 	`json:"uuid"`
		IsBoot	*bool	`json:"is_boot,omitempty"`
	}
	
	func main() {
		jsonData := &JsonType{}
		jsonData.UUID = "12345qwer"
		data, _ := json.Marshal(jsonData)
		fmt.Println(string(data))
		Boot := false
		jsonData.IsBoot = &Boot
		data2, _ := json.Marshal(jsonData)
		fmt.Println(string(data2))
	}
	```
	序列化之后的结果为：
	```
	{"uuid":"12345qwer"}
	{"uuid":"12345qwer","is_boot":false}
	```

- 值类型和引用类型

	**值类型**

	int、float、bool和string这些类型都属于值类型，使用这些类型的变量直接指向存在内存中的值，值类型的变量的值存储在栈中。当使用等号将一个变量的值赋给另一个变量时，如`j = i`,实际上是在内存中将i的值进行了拷贝。可以通过`&i`获取变量i`的内存地址.
	
	**引用类型**

	特指指针、slice、map、channel等预定义类型。一个引用类型的变量r1存储的是r1的值所在的内存地址（数字），或内存地址中第一个字所在的位置，这个内存地址被称之为指针，这个指针实际上也被存在另外的某一个字中。被引用的变量会存储在堆中，以便进行垃圾回收，且比栈拥有更大的内存空间。
	
	例：定义了一个数组a（值类型），复制给b，当b发生变化后a并不会发生任何变化
	```
	func main() {
		a :=[5]int{1,2,3,4,5}
		b := a
		b[2] = 8
		fmt.Println(a, b)
	}
	```
	打印结果为：`[1 2 3 4 5] [1 2 8 4 5]`

	例：定义一个切片a（引用类型），复制给b，当b发生变化后a也会跟着变化
	```
		func main() {
		a :=[]int{1,2,3,4,5}
		b := a
		b[2] = 8
		fmt.Println(a, b)
	}
	```
	打印结果为：`[1 2 8 4 5] [1 2 8 4 5]`

- 字符串

	Go中的字符串根据需要占用1至4个字节。
	字符串拼接：可以使用`+`来拼接字符串，`strings.Join()`函数拼接字符串更加高效，强烈推荐使用字节缓冲`bytes.Buffer`拼接。
	
	**strings和strconv包**
	
	- `HasPrefix`
		
		判断字符串s是否以prefix开头：
		`strings.HasPrefix(s, prefix string) bool`

	- `HasSuffix`
		
		判断字符串s是否以suffix结尾：
		`strings.HasSuffix(s, suffix string) bool`

	- `Contains`
	
		判断字符串s是否包含substr：
		`strings.Contains(s, substr string) bool`

	- `Index`

		返回字符串str在字符串s中的索引（str的第一个字符的索引），-1表示字符串s不包含字符串str:
		`strings.Index(s, str string) int`

	- `Replace`
		
		用于将字符串str中的前n个字符串old替换为字符串new，并返回一个新的字符串，如果n = -1则替换所有字符串old为字符串new：
		`strings.Replace(str, old, new, n) string`

	- `Count`

		用于计算字符串str在字符串s中出现的非重叠次数：
		`strings.Count(s, str string) int`

	- `Repeat`

		用于重复count次字符串s并返回一个新的字符串：
		`strings.Repeat(s, count int) string`

	- `ToLower`
	
		将字符串中的Unicode字符全部转换为相应的小写字符：
		`strings.ToLower(s) string`

	- `TrimSpace`

		剔除字符串开头和结尾的空白符号；如果想要剔除指定字符，则可以使用`strings.Trim(s, "cut")`来将开头和结尾的cut去除掉。该函数的第二个参数可以包含任何字符，如果只想剔除开头或者结尾的字符串，则可以使用`TrimLeft`或者`TrimRight`来实现。

	- `Fields`

		利用1个或多个空白符号来作为动态长度的分隔符将字符串分割成若干小块，并返回一个slice，如果字符串只包含空白符号，则返回一个长度为0的slice。`trings.Split(s, sep)`用于自定义分割符号来对指定字符串进行分割，同样返回slice。

	- `strconvb包`
		
		与字符串相关的类型转换都是通过strconv包实现。
		
		**数字类型转换到字符串**：
		
			`strconv.Itoa(i int) string`返回数字i所表示的字符串类型的十进制数。
			`strconv.FormatFloat(f float64, fmt byte, prec int, bitSize int) string`将64位浮点型的数字转换为字符串，其中fmt表示格式（其值可以是 'b'、'e'、'f' 或 'g'），prec表示精度，bitSize则使用32表示float32，用64表示float64。

		**从字符串类型转换为数字类型**:

			`strconv.Atoi(s string) (i int, err error)`将字符串转换为in 型。
			`strconv.ParseFloat(s string, bitSize int) (f float64, err error)`将字符串转换为float64型。

- 指针

	一个指针变量可以指向任何一个值的内存地址，它指向那个值的内存地址，在32位机器上占用4个字节，在64位机器上占用8个字节，并且与它所指向的值的大小无关。当一个指针被定义后没有分配到任何变量时，它的值为nil。不能获取一个文字或常量的地址，这是非法的。

	指针的一个高级应用是你可以传递一个变量的引用（如函数的参数），这样不会传递变量的拷贝。指针传递是很廉价的，只占用 4 个或 8 个字节。当程序在工作中需要占用大量的内存，或很多变量，或者两者都有，使用指针会减少内存占用和提高效率。被指向的变量也保存在内存中，直到没有任何指针指向它们，所以从它们被创建开始就具有相互独立的生命周期。

	```
	package main
	import "fmt"
	func main() {
		var i1 = 5
		fmt.Printf("An integer: %d, its location in memory: %p\n", i1, &i1)
		var intP *int
		intP = &i1
		fmt.Printf("The value at memory location %p is %d\n", intP, *intP)
	}
	```

- 标签与goto

	for、switch或select语句都可以配合标签（label）形式的标识符使用，即某一行第一个以冒号（:）结尾的单词（gofmt会将后续代码自动移至下一行）。标签的名称是大小写敏感的，一般建议使用全部大写字母。
	注意：不建议使用标签和goto语句。
	
	例1：

	```
	package main

	import "fmt"
	
	func main() {
	
	LABEL1:
		for i := 0; i <= 5; i++ {
			for j := 0; j <= 5; j++ {
				if j == 4 {
					continue LABEL1
				}
				fmt.Printf("i is: %d, and j is: %d\n", i, j)
			}
		}
	
	}
	```
	
	例2：

	```
	package main
	
	func main() {
		i:=0
		HERE:
			print(i)
			i++
			if i==5 {
				return
			}
			goto HERE
	}
	```

- 函数
	
	- 按值传递和按引用传递
	
		**按值传递**
	
			传递参数的副本，也是GO默认的传递方式。函数接收参数副本之后，在使用变量的过程中可能对副本的值进行更改，但不会影响到原来的变量。

		**按引用传递**

			将参数的地址（变量名前面添加&符号，比如&variable）传递给函数，此时传递给函数的是一个指针，函数可以直接修改参数的值，而不是对参数的副本进行操作。
			在函数调用时，像切片（slice）、字典（map）、接口（interface）、通道（channel）这样的引用类型都是默认使用引用传递（即使没有显式的指出指针）。

	- 传递变长参数

		如果函数的最后一个参数是采用`...type`的形式，那么这个函数就可以处理一个变长的参数，这个长度可以为0，这样的函数称为变参函数。
		```
		func Greeting(prefix string, who ...string)
		Greeting("hello:", "Joe", "Anna", "Eileen")
		```
		如果变长参数的类型并不是都相同,可以使用结构或空接口`interface{}`

	- defer关键字
		
		关键字`defer`允许程序推迟到函数返回之前（或任意位置执行return语句之后）一刻才执行某个语句或函数，类似finally。当有多个 defer行为被注册时，它们会以逆序执行（类似栈，即后进先出）。
		```
		func f() {
			for i := 0; i < 5; i++ {
				defer fmt.Printf("%d ", i)
			}
		}
		```
		输出结果为`4 3 2 1 0`
		defer关键字一般有如下几种使用场景：
		
		- 关闭文件流
		- 解锁一个加锁的资源
		- 打印最终报告
		- 关闭数据库链接
		- 使用defer语句实现代码追踪
		- 使用defer语句来记录函数的参数与返回值

	- 内置函数

		- close
		
			用于管道通信

		- len、cap
		
			len用于返回某个类型的长度或数量（字符串、数组、切片、map 和管道）；cap是容量的意思，用于返回某个类型的最大容量（只能用于切片和map）

		- new、make
		
			new和make均是用于分配内存。new用于值类型和用户定义的类型，如自定义结构，make用于内置引用类型（切片、map和管道）。它们的用法就像是函数，但是将类型作为参数：new(type)、make(type)。new(T) 分配类型T的零值并返回其地址，也就是指向类型 T 的指针。它也可以被用于基本类型：v := new(int)。make(T) 返回类型T的初始化之后的值，因此它比new进行更多的工作，new()是一个函数，不要忘记它的括号

		- copy、append
		
			用于复制和连接切片

		- panic、recover
		
			两者均用于错误处理机制

		- print、println
		
			底层打印函数，在部署环境中建议使用fmt包

		- complex、real imag
		
			用于创建和操作复数
	
	- 将函数作为参数

		函数可以作为其它函数的参数进行传递，然后在其它函数内调用执行，一般称之为回调。

		```
		package main

		import (
			"fmt"
		)
		
		func main() {
			callback(1, Add)
		}
		
		func Add(a, b int) {
			fmt.Printf("The sum of %d and %d is: %d\n", a, b, a+b)
		}
		
		func callback(y int, f func(int, int)) {
			f(y, 2) // this becomes Add(1, 2)
		}
		```

	- 闭包
	
		当不希望给函数起名字的时候，可以使用匿名函数，这样的函数不能够独立存在（编译器会返回错误：`non-declaration statement outside function body`），但可以被赋值于某个变量，即保存函数的地址到变量中：`fplus := func(x, y int) int { return x + y }`，然后通过变量名对函数进行调用：`fplus(3,4)`。也可以直接对匿名函数进行调用：`func(x, y int) int { return x + y } (3, 4)`

		```
		func() {
			sum := 0
			for i := 1; i <= 1e6; i++ {
				sum += i
			}
		}()
		```
			
		表示参数列表的第一对括号必须紧挨着关键字func。花括号{}涵盖着函数体，最后的一对括号表示对该匿名函数的调用。

		**将函数作为返回值**

		一个返回值为另一个函数的函数可以被称之为工厂函数。闭包函数保存并积累其中的变量的值，不管外部函数退出与否，它都能够继续操作外部函数中的局部变量。
	
		例1、：
		```
		package main

		import "fmt"
		
		func main() {
			var f = Adder()
			fmt.Print(f(1), " - ")
			fmt.Print(f(20), " - ")
			fmt.Print(f(300))
		}
		
		func Adder() func(int) int {
			var x int
			return func(delta int) int {
				x += delta
				return x
			}
		}
		```
		
		输出结果为：`1 - 21 - 321`

		例2：
		```
		func MakeAddSuffix(suffix string) func(string) string {
			return func(name string) string {
				if !strings.HasSuffix(name, suffix) {
					return name + suffix
				}
				return name
			}
		}

		addBmp := MakeAddSuffix(".bmp")
		addJpeg := MakeAddSuffix(".jpeg")

		addBmp("file") // returns: file.bmp
		addJpeg("file") // returns: file.jpeg
		```

		例3：使用闭包实现斐波拉切数列
		```
		package main

		import "fmt"
		
		// fibonacci is a function that returns
		// a function that returns an int.
		func fibonacci() func() int {
		    back1, back2:= 0, 1
		
		    return func() int {
		        
		        temp := back1
		        back1,back2 = back2,(back1 + back2)
		        return temp
		    }    
		}
		
		func main() {
		    f := fibonacci()
		    for i := 0; i < 10; i++ {
		        fmt.Println(f())
		    }
		}
		```

	- 内存缓存
	
		通过在内存中缓存和重复利用相同计算的结果，称之为内存缓存。
		斐波拉切数列的例子，将第n个数的值存在数组中索引为n的位置，然后在数组中查找是否已经计算过，如果没有找到，则再进行计算。

		```
		package main

		import (
			"fmt"
			"time"
		)
		
		const LIM = 41
		
		var fibs [LIM]uint64
		
		func main() {
			var result uint64 = 0
			start := time.Now()
			for i := 0; i < LIM; i++ {
				result = fibonacci(i)
				fmt.Printf("fibonacci(%d) is: %d\n", i, result)
			}
			end := time.Now()
			delta := end.Sub(start)
			fmt.Printf("longCalculation took this amount of time: %s\n", delta)
		}
		func fibonacci(n int) (res uint64) {
			// memoization: check if fibonacci(n) is already known in array:
			if fibs[n] != 0 {
				res = fibs[n]
				return
			}
			if n <= 1 {
				res = 1
			} else {
				res = fibonacci(n-1) + fibonacci(n-2)
			}
			fibs[n] = res
			return
		}
		```

- map

	map是引用类型：内存用make方法来分配，不要使用new，永远用make来构造map。map可以根据新增的key-value对动态的伸缩，因此它不存在固定长度或者最大限制。当map增长到容量上限的时候，如果再增加新的key-value对，map的大小会自动加1。所以出于性能的考虑，对于大的map或者会快速扩张的map，即使只是大概知道容量，也最好先标明。

	```
	package main
	import "fmt"
	
	func main() {
		mf := map[int]func() int{
			1: func() int { return 10 },
			2: func() int { return 20 },
			5: func() int { return 50 },
		}
		fmt.Println(mf)
	}
	```
	输出结果：`map[1:0x10903be0 5:0x10903ba0 2:0x10903bc0]`

- 标准库

	[gowalker](https://gowalker.org/search?q=gorepos)

	- 锁和sync包
	
		在Go语言中通过sync包中Mutex来实现锁的机制。`sync.Mutex`是一个互斥锁，它的作用是守护在临界区入口来确保同一时间只能有一个线程进入临界区。

		例1：
		```
		import  "sync"

		type Info struct {
			mu sync.Mutex
			// ... other fields, e.g.: Str string
		}
		```
		当有变量需要更新Info时，可以采用如下写法：
		```
		func Update(info *Info) {
			info.mu.Lock()
		    // critical section:
		    info.Str = // new value
		    // end critical section
		    info.mu.Unlock()
		}
		```
	
		例2
		通过Mutex来实现一个可以上锁的共享缓冲器:
		```
		type SyncedBuffer struct {
			lock 	sync.Mutex
			buffer  bytes.Buffer
		}
		```
	
		在sync包中还有一个RWMutex锁：他能通过RLock()来允许同一时间多个线程对变量进行读操作，但是只能一个线程进行写操作。如果使用 Lock()将和普通的Mutex作用相同。包中还有一个方便的Once类型变量的方法`once.Do(call)`，这个方法确保被调用函数只能被调用一次。
		
- Go的一些外部库
	
	- MySQL(GoMySQL), PostgreSQL(go-pgsql), MongoDB (mgo, gomongo), CouchDB (couch-go), ODBC (godbcl), Redis (redis.go) and SQLite3 (gosqlite) database drivers
	- SDL bindings
	- Google's Protocal Buffers(goprotobuf)
	- XML-RPC(go-xmlrpc)
	- Twitter(twitterstream)
	- OAuth libraries(GoAuth)


- 结构体

	- 结构体工厂

		Go语言不支持面向对象编程语言中那样的构造子方法，但是可以很容易的在Go中实现 “构造子工厂”方法。按惯例，工厂的名字以new或New开头。假设定义了如下的File结构体类型：

		```
		type File struct {
		    fd      int     // 文件描述符
		    name    string  // 文件名
		}
		```

		下面是这个结构体类型对应的工厂方法，它返回一个指向结构体实例的指针：

		```
		func NewFile(fd int, name string) *File {
		    if fd < 0 {
		        return nil
		    }
		
		    return &File{fd, name}
		}	
		```

		调用：

		```
		f := NewFile(10, "./test.txt")
		```

	- 匿名字段和内嵌结构体

		结构体可以包含一个或多个匿名（或内嵌）字段，即这些字段没有显式的名字，只有字段的类型是必须的，此时类型就是字段的名字，在一个结构体中对于每一种数据类型只能有一个匿名字段。匿名字段本身可以是一个结构体类型，即结构体可以包含内嵌结构体。

		```
		type innerS struct {
			in1 int
			in2 int
		}
		
		type outerS struct {
			b    int
			c    float32
			int  // anonymous field
			innerS //anonymous field
		}
		```

		**命名冲突**

		当两个字段拥有相同的名字（可能是继承来的名字）时该怎么办呢？

			- 外层名字会覆盖内层名字（但是两者的内存空间都保留），这提供了一种重载字段或方法的方式；
			- 如果相同的名字在同一级别出现了两次，如果这个名字被程序使用了，将会引发一个错误（不使用没关系）。没有办法来解决这种问题引起的二义性，必须由程序员自己修正。


- 方法

	Go语言中方法是作用在接收者（receiver）上的一个函数，接收者是某种类型的变量。因此方法是一种特殊类型的函数。任何类型都可以有方法，甚至可以是函数类型，可以是 int、bool、string 或数组的别名类型。但是接收者不能是一个接口类型，因为接口是一个抽象定义，但是方法却是具体实现，接收者也不能是一个指针类型，但是它可以是任何其他允许类型的指针。类型T（或 *T）上的所有方法的集合叫做类型T（或 *T）的方法集。因为方法是函数，所以同样的，不允许方法重载，即对于一个类型只能有一个给定名称的方法。
	
	定义方法的一般格式：
	`func (recv receiver_type) methodName(parameter_list) (return_value_list) { ... }`

	如果recv是receiver的实例，Method1是它的方法名，那么方法调用遵循传统的object.name选择器符号`recv.Method1()`。

	示例：
	
	```
	package main

	import "fmt"
	
	type TwoInts struct {
		a int
		b int
	}
	
	func main() {
		two1 := new(TwoInts)
		two1.a = 12
		two1.b = 10
	
		fmt.Printf("The sum is: %d\n", two1.AddThem())
		fmt.Printf("Add them to the param: %d\n", two1.AddToParam(20))
	
		two2 := TwoInts{3, 4}
		fmt.Printf("The sum is: %d\n", two2.AddThem())
	}
	
	func (tn *TwoInts) AddThem() int {
		return tn.a + tn.b
	}
	
	func (tn *TwoInts) AddToParam(param int) int {
		return tn.a + tn.b + param
	}
	```
	
	输出结果：
	```
	The sum is: 22
	Add them to the param: 42
	The sum is: 7
	```
	
	**函数和方法的区别**
	
	函数将变量作为参数：`Function1(recv)`，方法在变量上被调用：`recv.Method1()`；
	在接收者是指针时，方法可以改变接收者的值（或状态），这点函数也可以做到（当参数作为指针传递，即通过引用调用时，函数也可以改变参数的状态）；
	方法没有和数据定义（结构体）混在一起：它们是正交的类型；表示（数据）和行为（方法）是独立的。

	**指针或值作为接收者**

	鉴于性能的原因，recv最常见的是一个指向`receiver_type`的指针，跟函数中的引用传递类似。如果想要方法改变接收者的数据，就在接收者的指针类型上定义该方法。否则，就在普通的值类型上定义方法。
	对于类型T，如果在`*T`上存在方法`Meth()`，并且t是这个类型的变量，那么`t.Meth()`会被自动转换为`(&t).Meth()`。
	指针方法和值方法都可以在指针或非指针上被调用。
	示例：
	```
	package main

	import (
		"fmt"
	)
	
	type List []int
	
	func (l List) Len() int        { return len(l) }
	func (l *List) Append(val int) { *l = append(*l, val) }
	
	func main() {
		// 值
		var lst List
		lst.Append(1)
		fmt.Printf("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)
	
		// 指针
		plst := new(List)
		plst.Append(2)
		fmt.Printf("%v (len: %d)", plst, plst.Len()) // &[2] (len: 1)
	}
	```

	**方法和未导出字段**

	示例：
	```
	ackage person
	
	type Person struct {
		firstName string
		lastName  string
	}
	
	func (p *Person) FirstName() string {
		return p.firstName
	}
	
	func (p *Person) SetFirstName(newName string) {
		p.firstName = newName
	}
	```
	```
	package main

	import (
		"./person"
		"fmt"
	)
	
	func main() {
		p := new(person.Person)
		// p.firstName undefined
		// (cannot refer to unexported field or method firstName)
		// p.firstName = "Eric"
		p.SetFirstName("Eric")
		fmt.Println(p.FirstName()) // Output: Eric
	}
	```

	**内嵌类型的方法和继承**

	当一个匿名类型被内嵌在结构体中时，匿名类型的可见方法也同样被内嵌，这在效果上等同于外层类型继承了这些方法：将父类型放在子类型中来实现亚型。

	示例(内嵌结构体上的方法可以直接在外层类型的实例上调用)：
	```
	package main

	import (
		"fmt"
		"math"
	)
	
	type Point struct {
		x, y float64
	}
	
	func (p *Point) Abs() float64 {
		return math.Sqrt(p.x*p.x + p.y*p.y)
	}
	
	type NamedPoint struct {
		Point
		name string
	}
	
	func main() {
		n := &NamedPoint{Point{3, 4}, "Pythagoras"}
		fmt.Println(n.Abs()) // 打印5
	}
	```

	和内嵌类型方法具有同样名字的外层类型的方法会覆写内嵌类型对应的方法。结构体内嵌和自己在同一个包中的结构体时，可以彼此访问对方所有的字段和方法。

	**在类型中嵌入功能**

	实现在类型中嵌入功能：

		- 聚合（或组合）：包含一个所需功能类型的具名字段。
		
		- 内嵌：内嵌（匿名地）所需功能类型。

	示例1（聚合）：

	```
	package main

	import (
		"fmt"
	)
	
	type Log struct {
		msg string
	}
	
	type Customer struct {
		Name string
		log  *Log
	}
	
	func main() {
		c := new(Customer)
		c.Name = "Barak Obama"
		c.log = new(Log)
		c.log.msg = "1 - Yes we can!"
		// shorter
		c = &Customer{"Barak Obama", &Log{"1 - Yes we can!"}}
		// fmt.Println(c) &{Barak Obama 1 - Yes we can!}
		c.Log().Add("2 - After me the world will be a better place!")
		//fmt.Println(c.log)
		fmt.Println(c.Log())
	
	}
	
	func (l *Log) Add(s string) {
		l.msg += "\n" + s
	}
	
	func (l *Log) String() string {
		return l.msg
	}
	
	func (c *Customer) Log() *Log {
		return c.log
	}
	```
	
	输出结果：
	```
	1 - Yes we can!
	2 - After me the world will be a better place!
	```

	示例2（内嵌）：

	```
	package main

	import (
		"fmt"
	)
	
	type Log struct {
		msg string
	}
	
	type Customer struct {
		Name string
		Log
	}
	
	func main() {
		c := &Customer{"Barak Obama", Log{"1 - Yes we can!"}}
		c.Add("2 - After me the world will be a better place!")
		fmt.Println(c)
	
	}
	
	func (l *Log) Add(s string) {
		l.msg += "\n" + s
	}
	
	func (l *Log) String() string {
		return l.msg
	}
	
	func (c *Customer) String() string {
		return c.Name + "\nLog:" + fmt.Sprintln(c.Log)
	}
	```

	输出结果：
	```
	Barak Obama
	Log:{1 - Yes we can!
	2 - After me the world will be a better place!}
	```

	内嵌的类型不需要指针，Customer也不需要Add方法，它使用Log的Add方法，Customer有自己的String方法，并且在它里面调用了Log的String方法。

	**多重继承**

	多重继承指的是类型获得多个父类型行为的能力，它在传统的面向对象语言中通常是不被实现的（C++和Python例外）。因为在类继承层次中，多重继承会给编译器引入额外的复杂度。但是在Go语言中，通过在类型中嵌入所有必要的父类型，可以很简单的实现多重继承。

	示例：
	```
	package main

	import (
		"fmt"
	)
	
	type Camera struct{}
	
	func (c *Camera) TakeAPicture() string {
		return "Click"
	}
	
	type Phone struct{}
	
	func (p *Phone) Call() string {
		return "Ring Ring"
	}
	
	type CameraPhone struct {
		Camera
		Phone
	}
	
	func main() {
		cp := new(CameraPhone)
		fmt.Println("Our new CameraPhone exhibits multiple behaviors...")
		fmt.Println("It exhibits behavior of a Camera: ", cp.TakeAPicture())
		fmt.Println("It works like a Phone too: ", cp.Call())
	}
	```

	在Go中，类型就是类（数据和关联的方法）。Go不知道类似面向对象语言的类继承的概念。继承有两个好处：代码复用和多态。在Go中，代码复用通过组合和委托实现，多态通过接口的使用来实现：有时这也叫组件编程（Component Programming）。


	- 类型的`String()`方法和格式化描述符
	
	```
	package main

	import (
		"fmt"
		"strconv"
	)
	
	type TwoInts struct {
		a int
		b int
	}
	
	func main() {
		two1 := new(TwoInts)
		two1.a = 12
		two1.b = 10
		fmt.Printf("two1 is: %v\n", two1)
		fmt.Println("two1 is:", two1)
		fmt.Printf("two1 is: %T\n", two1)
		fmt.Printf("two1 is: %#v\n", two1)
	}
	
	func (tn *TwoInts) String() string {
		return "(" + strconv.Itoa(tn.a) + "/" + strconv.Itoa(tn.b) + ")"
	}
	```
	不要在`String()`方法里面调用涉及`String()`方法的方法，它会导致意料之外的错误。

	- 垃圾回收和SetFinalizer

	通过调用`runtime.GC()`函数可以显式的触发GC，但这只在某些罕见的场景下才有用，比如当内存资源不足时调用`runtime.GC()`，它会在此函数执行的点上立即释放一大片内存，此时程序可能会有短时的性能下降（因为GC进程在执行）。
	如果需要在一个对象 obj 被从内存移除前执行一些特殊操作，比如写到日志文件中，可以通过如下方式调用函数来实现：
	```runtime.SetFinalizer(obj, func(obj *typeObj))```
	`func(obj *typeObj)`需要一个`typeObj`类型的指针参数`obj，特殊操作会在它上面执行。func也可以是一个匿名函数。在对象被GC进程选中并从内存中移除以前，`SetFinalizer`都不会执行，即使程序正常结束或者发生错误。

	
- 接口（Interfaces）与反射（reflection）

	- **接口**
	
		类型（比如结构体）实现接口方法集中的方法，每一个方法的实现说明了此方法是如何作用于该类型的：即实现接口，同时方法集也构成了该类型的接口。类型不需要显式声明它实现了某个接口：接口被隐式地实现。多个类型可以实现同一个接口。一个类型可以实现多个接口。

		示例：
		```
		package main

		import "fmt"
		
		type Shaper interface {
			Area() float32
		}
		
		type Square struct {
			side float32
		}
		
		func (sq *Square) Area() float32 {
			return sq.side * sq.side
		}
		
		func main() {
			sq1 := new(Square)
			sq1.side = 5
		
			var areaIntf Shaper
			areaIntf = sq1
			// shorter,without separate declaration:
			// areaIntf := Shaper(sq1)
			// or even:
			// areaIntf := sq1
			fmt.Printf("The square has area: %f\n", areaIntf.Area())
		}
		```
		
		示例2：
		```
		package main

		import "fmt"
		
		type Shaper interface {
			Area() float32
		}
		
		type Square struct {
			side float32
		}
		
		func (sq *Square) Area() float32 {
			return sq.side * sq.side
		}
		
		type Rectangle struct {
			length, width float32
		}
		
		func (r Rectangle) Area() float32 {
			return r.length * r.width
		}
		
		func main() {
		
			r := Rectangle{5, 3} // Area() of Rectangle needs a value
			q := &Square{5}      // Area() of Square needs a pointer
			// shapes := []Shaper{Shaper(r), Shaper(q)}
			// or shorter
			shapes := []Shaper{r, q}
			fmt.Println("Looping through shapes for area ...")
			for n, _ := range shapes {
				fmt.Println("Shape details: ", shapes[n])
				fmt.Println("Area of this shape is: ", shapes[n].Area())
			}
		}
		```
		
		输出：
		```
		Looping through shapes for area ...
		Shape details:  {5 3}
		Area of this shape is:  15
		Shape details:  &{5}
		Area of this shape is:  25
		```

		示例3：
		
		```
		package main

		import "fmt"
		
		type stockPosition struct {
			ticker     string
			sharePrice float32
			count      float32
		}
		
		/* method to determine the value of a stock position */
		func (s stockPosition) getValue() float32 {
			return s.sharePrice * s.count
		}
		
		type car struct {
			make  string
			model string
			price float32
		}
		
		/* method to determine the value of a car */
		func (c car) getValue() float32 {
			return c.price
		}
		
		/* contract that defines different things that have value */
		type valuable interface {
			getValue() float32
		}
		
		func showValue(asset valuable) {
			fmt.Printf("Value of the asset is %f\n", asset.getValue())
		}
		
		func main() {
			var o valuable = stockPosition{"GOOG", 577.20, 4}
			showValue(o)
			o = car{"BMW", "M3", 66500}
			showValue(o)
		}
		```
		
		输出：
		```
		Value of the asset is 2308.800049
		Value of the asset is 66500.000000
		```

		**接口嵌套接口**

		一个接口可以包含一个或多个其他的接口，这相当于直接将这些内嵌接口的方法列举在外层接口中一样。
		示例：
		```
		type ReadWrite interface {
		    Read(b Buffer) bool
		    Write(b Buffer) bool
		}
		
		type Lock interface {
		    Lock()
		    Unlock()
		}
		
		type File interface {
		    ReadWrite
		    Lock
		    Close()
		}
		```
		
	- 类型断言

	一个接口类型的变量varI中可以包含任何类型的值，必须有一种方式来检测它的动态类型，即运行时在变量中存储的值的实际类型。在执行过程中动态类型可能会有所不同，但是它总是可以分配给接口变量本身的类型。通常我们可以使用 类型断言 来测试在某个时刻varI是否包含类型T的值：```v := varI.(T)       // unchecked type assertion```
	类型断言可能是无效的，虽然编译器会尽力检查转换是否有效，但是它不可能预见所有的可能性。如果转换在程序运行时失败会导致错误发生。更安全的方式是使用以下形式来进行类型断言：
	```
	if v, ok := varI.(T); ok {  // checked type assertion
	    Process(v)
	    return
	}
	// varI is not of type T
	```
	示例：
	```
	package main

	import (
		"fmt"
		"math"
	)
	
	type Square struct {
		side float32
	}
	
	type Circle struct {
		radius float32
	}
	
	type Shaper interface {
		Area() float32
	}
	
	func main() {
		var areaIntf Shaper
		sq1 := new(Square)
		sq1.side = 5
	
		areaIntf = sq1
		// Is Square the type of areaIntf?
		if t, ok := areaIntf.(*Square); ok {
			fmt.Printf("The type of areaIntf is: %T\n", t)
		}
		if u, ok := areaIntf.(*Circle); ok {
			fmt.Printf("The type of areaIntf is: %T\n", u)
		} else {
			fmt.Println("areaIntf does not contain a variable of type Circle")
		}
	}
	
	func (sq *Square) Area() float32 {
		return sq.side * sq.side
	}
	
	func (ci *Circle) Area() float32 {
		return ci.radius * ci.radius * math.Pi
	}
	```

	输出：
	```
	The type of areaIntf is: *main.Square
	areaIntf does not contain a variable of type Circle
	```

	- 类型判断：type-switch

	所有case语句中列举的类型（nil 除外）都必须实现对应的接口，如果被检测类型没有在case语句列举的类型中，就会执行default语句。
	
	示例：
	```
	func classifier(items ...interface{}) {
		for i, x := range items {
			switch x.(type) {
			case bool:
				fmt.Printf("Param #%d is a bool\n", i)
			case float64:
				fmt.Printf("Param #%d is a float64\n", i)
			case int, int64:
				fmt.Printf("Param #%d is a int\n", i)
			case nil:
				fmt.Printf("Param #%d is a nil\n", i)
			case string:
				fmt.Printf("Param #%d is a string\n", i)
			default:
				fmt.Printf("Param #%d is unknown\n", i)
			}
		}
	}
	```
	测试一个值是否实现了某个接口:
	```
	type Stringer interface {
	    String() string
	}
	
	if sv, ok := v.(Stringer); ok {
	    fmt.Printf("v implements String(): %s\n", sv.String()) // note: sv, not v
	}
	```

	- 使用方法集与接口

	在接口上调用方法时，必须有和方法定义时相同的接收者类型或者是可以从具体类型P直接可以辨识的：
	
		- 指针方法可以通过指针调用
		- 值方法可以通过值调用
		- 接收者是值的方法可以通过指针调用，因为指针会首先被解引用
		- 接收者是指针的方法不可以通过值调用，因为存储在接口中的值没有地址

	Go语言规范定义了接口方法集的调用规则：

		- 类型 *T 的可调用方法集包含接受者为 *T 或 T 的所有方法集
		- 类型 T 的可调用方法集包含接受者为 T 的所有方法
		- 类型 T 的可调用方法集不包含接受者为 *T 的方法

	示例1（使用Sorter接口排序）
	```
	//冒泡算法
	func Sort(data Sorter) {
	    for pass := 1; pass < data.Len(); pass++ {
	        for i := 0;i < data.Len() - pass; i++ {
	            if data.Less(i+1, i) {
	                data.Swap(i, i + 1)
	            }
	        }
	    }
	}
	```
	Sort 函数接收一个接口类型的参数：Sorter ，它声明了这些方法：
	```
	type Sorter interface {
	    Len() int
	    Less(i, j int) bool
	    Swap(i, j int)
	}
	```
	为数组定一个类型并在它上面实现 Sorter 接口的方法：
	```
	type IntArray []int
	func (p IntArray) Len() int           { return len(p) }
	func (p IntArray) Less(i, j int) bool { return p[i] < p[j] }
	func (p IntArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
	```

	示例2（io.Reader和io.Writer）
	```
	type Reader interface {
	    Read(p []byte) (n int, err error)
	}
	
	type Writer interface {
	    Write(p []byte) (n int, err error)
	}
	```
	只要类型实现了读写接口，提供Read()和Write方法，就可以从它读取数据，或向它写入数据。一个对象要是可读的，它必须实现io.Reader接口，这个接口只有一个签名是`Read(p []byte) (n int, err error)`的方法，它从调用它的对象上读取数据，并把读到的数据放入参数中的字节切片中，然后返回读取的字节数和一个 error 对象，如果没有错误发生返回 nil，如果已经到达输入的尾端，会返回`io.EOF("EOF")`，如果读取的过程中发生了错误，就会返回具体的错误信息。类似地，一个对象要是可写的，它必须实现io.Writer接口，这个接口也只有一个签名是`Write(p []byte) (n int, err error)`的方法，它将指定字节切片中的数据写入调用它的对象里，然后返回实际写入的字节数和一个 error 对象（如果没有错误发生就是 nil）。io 包里的Readers和Writers都是不带缓冲的，bufio包里提供了对应的带缓冲的操作，在读写UTF-8编码的文本文件时它们尤其有用。

	- 空接口
	
	空接口或者最小接口 不包含任何方法，它对实现不做任何要求：`type Any interface {}`。
	示例：
	```
	package main

	import "fmt"
	
	type specialString string
	
	var whatIsThis specialString = "hello"
	
	func TypeSwitch() {
		testFunc := func(any interface{}) {
			switch v := any.(type) {
			case bool:
				fmt.Printf("any %v is a bool type", v)
			case int:
				fmt.Printf("any %v is an int type", v)
			case float32:
				fmt.Printf("any %v is a float32 type", v)
			case string:
				fmt.Printf("any %v is a string type", v)
			case specialString:
				fmt.Printf("any %v is a special String!", v)
			default:
				fmt.Println("unknown type!")
			}
		}
		testFunc(whatIsThis)
	}
	
	func main() {
		TypeSwitch()
	}
	```
	输出：`any hello is a special String!`。

	- 使用空接口构建通用类型或包含不同类型变量的数组

	```
	type Element interface{}
	type Vector struct {
		a []Element
	}
	```

	Vector里能放任何类型的变量，因为任何类型都实现了空接口，实际上Vector里放的每个元素可以是不同类型的变量。

	- 复制数据切片至空接口切片

	```
	var dataSlice []myType = FuncReturnSlice()
	var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
	for i, d := range dataSlice {
	    interfaceSlice[i] = d
	}
	```
	注意：不能直接赋值，必须一个一个显式地复制。

	- 通用类型的节点数据结构

	示例（实现二叉树的部分代码）：
	```
	package main

	import "fmt"
	
	type Node struct {
		le   *Node
		data interface{}
		ri   *Node
	}
	
	func NewNode(left, right *Node) *Node {
		return &Node{left, nil, right}
	}
	
	func (n *Node) SetData(data interface{}) {
		n.data = data
	}
	
	func main() {
		root := NewNode(nil, nil)
		root.SetData("root node")
		// make child (leaf) nodes:
		a := NewNode(nil, nil)
		a.SetData("left node")
		b := NewNode(nil, nil)
		b.SetData("right node")
		root.le = a
		root.ri = b
		fmt.Printf("%v\n", root) // Output: &{0x125275f0 root node 0x125275e0}
	}
	```