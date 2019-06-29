# Packer #

[https://github.com/hashicorp/packer](https://github.com/hashicorp/packer)
[官方文档](https://www.packer.io/docs/index.html)

Packer是HashiCorp公司（very cool）开源软件族中的一个，是一个从单一配置文件为多平台创建一致镜像的轻量级的开源工具。能够运行在主流的操作系统上，并行高效的创建多平台的虚拟机镜像，它为代码即基础结构提供了坚实的基础，通过packer较大的降低了了创建用户自定义镜像的难度，并且将创建镜像的过程从人工的随机过程变成可以配置管理代码，可测试的过程，从而减少了用户应用上云的障碍之一。持的平台包括Amazon EC2、DigitalOcean、OpenStack、VirtualBox和VMware等。

## Packer的组成 ##

[https://yq.aliyun.com/articles/72724?t=t1](https://yq.aliyun.com/articles/72724?t=t1)

Packer包含构建器(Builder),（派生器）Provisioner,(后处理器)Post-Processor三个组件，通过JSON格式的模板文件，可以灵活的组合这三种组件并行的、自动化的创建多平台一致的镜像文件。为单个平台生成镜像的单个任务称为构建，而单个构建的结果也称为工件(Artifact)，多个构建可以并行运行。

- **Builder**(构建器)  能够为单个平台创建镜像。构建器读取一些配置并使用它来运行和生成镜像。作为构建的一部分调用构建器以创建实际生成的镜像。构建器可以以插件的形式创建并添加到Packer中。
- **Provisioner**(派生器）  这一组件在Buider创建的运行的机器中安装和配置软件。他们执行使镜像包含有用软件的主要工作。常见的派生器包括shell脚本，Chef，Puppet等。
- **Post-Processors**(后处理器） 它使用构建器或另一个后处理器的结果来创建新工件的过程。例如压缩后处理器压缩工件，上传后处理器上传工件等。

## 示例 ##

```
{
  "builders": [
    {
      "type": "amazon-ebs",
      "access_key": "...",
      "secret_key": "...",
      "region": "us-east-1",
      "source_ami": "ami-fce3c696",
      "instance_type": "t2.micro",
      "ssh_username": "ubuntu",
      "ami_name": "packer {{timestamp}}"
    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "script": "setup_things.sh"
    }
  ]
}
```