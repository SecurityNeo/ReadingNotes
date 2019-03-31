# Go语言基础 #


- 常量计数器iota
	iota是常量计数器,只能在常量的表达式中使用。iota在const关键字出现时将被重置为0(const内部的第一行之前)，const中每新增一行常量声明将使iota计数一次(即其值自动加1)。使用iota能简化定义，在定义枚举时很有用。
	
	例：
	1、
	```
	const (
	a = iota    //a=0
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