# Ansible笔记 #

## 介绍 ##

ansible是基于paramiko开发的,并且基于模块化工作，本身没有批量部署的能力。真正具有批量部署的是ansible所运行的模块，ansible只是提供一种框架。

架构图：

![](img/ansible_struc.png)

- Ansible：Ansible核心程序。
- HostInventory：记录由Ansible管理的主机信息，包括端口、密码、ip等。
- Playbooks：“剧本”YAML格式文件，多个任务定义在一个文件中，定义主机需要调用哪些模块来完成的功能。
- CoreModules：核心模块，主要操作是通过调用核心模块来完成管理任务。
- CustomModules：自定义模块，完成核心模块无法完成的功能，支持多种语言。
- ConnectionPlugins：连接插件，Ansible和Host通信使用

## 剧本编写 ##

**tags**

ansible 2.5以后内置的tags有以下几个（当然可以为task自定义tag）：

- always: 指定这个tag后，task任务将永远被执行，而不用去考虑是否使用了--skip-tags标记
- tagged: 当 --tags 指定为它时，则只要有tags标记的task都将被执行,--skip-tags效果相反
- untagged: 当 --tags 指定为它时，则所有没有tag标记的task 将被执行,--skip-tags效果相反
- all: 这个标记无需指定，ansible-playbook默认执行的时候就是这个标记.所有task都被执行

**set_fact**

ansible有一个模块叫setup，用于获取远程主机的相关信息，并可以将这些信息作为变量在playbook里进行调用。而setup模块获取这些信息的方法就是依赖于fact。

- 手动设置fact

ansible除了能获取到预定义的fact的内容,还支持手动为某个主机定制fact。称之为本地fact。本地fact默认存放于被控端的/etc/ansible/facts.d目录下，如果文件为ini格式或者json格式，ansible会自动识别。

- 使用set_fact模块定义新的变量

set_fact模块可以自定义facts，这些自定义的facts可以通过template或者变量的方式在playbook中使用。

```
- name: set_fact example
  hosts: test
  tasks:
    - name: Calculate InnoDB buffer pool size
      set_fact: innodb_buffer_pool_size_mb="{{ ansible_memtotal_mb / 2 |int }}"
      
    - debug: var=innodb_buffer_pool_size_mb
```