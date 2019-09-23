# Golang方法 #

[https://blog.csdn.net/xiangxianghehe/article/details/78790601](https://blog.csdn.net/xiangxianghehe/article/details/78790601)

一个方法就是一个包含了接受者的函数，接受者可以是命名类型或者结构体类型的一个值或者是一个指针,但是接收者不能是一个接口类型。方法的声明和普通函数的声明类似，只是在函数名称前面多了一个参数，这个参数把这个方法绑定到这个参数对应的类型上。方法常见的语法：

```golang
func (variable_name variable_data_type) function_name() [return_type]{
   /* 函数体*/
}
```

方法是特殊的函数，定义在某一特定的类型上，通过类型的实例来进行调用，这个实例被叫接收者(receiver)。
函数将变量作为参数：`Function1(recv)`
方法在变量上被调用：`recv.Method1()`
接收者必须有一个显式的名字，这个名字必须在方法中被使用。`receiver_type`叫做（接收者）基本类型，这个类型必须在和方法同样的包中被声明。

注意：Go语言不允许为简单的内置类型添加方法。例如这样做是不合法的：

```golang
func (a int) Add (b int){    //方法非法！不能是内置数据类型
  fmt.Println(a+b)
}
```
