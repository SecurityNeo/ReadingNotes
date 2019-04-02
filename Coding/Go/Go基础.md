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
	