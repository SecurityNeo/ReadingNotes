# Go语言基础 #


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
	