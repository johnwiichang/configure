```
   __________  _   ____________________  ______  ______
  / ____/ __ \/ | / / ____/  _/ ____/ / / / __ \/ ____/
 / /   / / / /  |/ / /_   / // / __/ / / / /_/ / __/   
/ /___/ /_/ / /|  / __/ _/ // /_/ / /_/ / _, _/ /___   
\____/\____/_/ |_/_/   /___/\____/\____/_/ |_/_____/                                     
```

Configure 是一个附带检查的统一配置提供器。

Configure is a unified configuration provider with checks.

## **功能** What you can get from this lib
- **全新结构配置** Brand-New Structured Configuration
    > 传统基于结构的配置将会提供确切类型数据，初始化时需要执行全般初始化，对于可空、可默认以及配置错误的场景可能出现意外的配置问题。现在，可以在初始化载入配置时自动检查。
    >
    > You might assigned accurate type to configuration type, and will initial them all before program running. If your configuration file is broken or just wrong, you might need lots of check functions to ensure the value is correct. You can finish the value validation via this lib nowadays.
    > 
    > 例如：针对 `HostAddress`，如果配置文件错误，那么默认将会给出空字符串。空字符串将无法在执行载入时报错，只能频繁检查。另，更新配置时，容易出现意外的值错误。
    > 
    > E.g.: There is an item named `HostAddress`, if your file goes wrong, the empty string will fill into the item. Empty string will not cause panic when it comes to you, so we must invoke a large number of check functions to void panic. However, the unexpected value might make you mad during updating the configuration.

- **多嵌套的配置** Muti-Layer Configuration Support
    > 非结构式配置通常使用全局映射（map）来解决。这类映射虽然看起来靠谱，但是在并发读写出现的时候容易诱发 panic。同时，扁平化的配置管理容易造成尴尬的局面——配置项超多造成的调用混乱甚至毫无层次感而言。
    > 
    > Non-Structured configuration will use map to solve the problem. This method seems robust, but if you do a write action during fetch the value, you might meet panic instead. And, the flat configuration management will make you embarrassed just because lots of items will let you code without structure and lead you into mess.
    > 
    > 例如，邮件配置中，发件人与 SMTP 服务器地址其实与其他参数无关。使用全局映射的话，可能构成 `infrastructure.confMap["MailConf.Sender"]` 或者将键常量化，出现 `infrastructure.confMap[const.Descriptor]` 一类的访问。通过本库，可以实现：`infrastructure.ConfigurationProvider.MailService.Sender` 的访问形式，不管使用 Visual Studio Code 还是其他工具，配合代码提示可以降低因为变量名接近带来的风险。
    > 
    > E.g.: We all know that mail service need SMTP server address and sender, but both two arguments have no relations to others. Under the circumstances that we fetch value from global map data structure, the access code might like `infrastructure.confMap["MailConf.Sender"]`, and we can make it shorter like `infrastructure.confMap[const.Descriptor]`. But you can do it like `infrastructure.ConfigurationProvider.MailService.Sender` by using this lib. And if you use tools like Visual Studio Code, you can get better code experience and make less bugs.

- **多文件的配置** Muti-File Configuration Support
    > 项目中若使用配置文件分离的办法解决配置过长或意义不明的场景时，使用基于结构的配置容易造成装配文件的多次检查以及手动引用。对于嵌套的类型无法做到部分+全局的一次性兼容，需要二选其一，没有灰度行为的可能。
    > 
    > If you use multiple files to solve the problem which you have lots of items, you need lots of check functions to ensure your project will not be terminated by unexpected value and you should manage reference manually. And structured configuration design will meet some problems when partial items and global items loaded in the flash. You need to choose one of them, and there is no possibility of grayscale behaviour.
    > 
    > 例如：某些配置不需要独立配置文件（例如数据库连接字符串），但是为了定义符合设计，需要嵌入到二级配置中。
    > 
    > E.g.: Some configuration items designed to be a sub-item, but it just a simple string without any dependencies just like your database connection string ought to put into database section and it’s a basic item you might never change it in independent file.

## **开始** QuickStart
一切的一切需要从引入 configure 开始：

If you wanna use configure, please execute:

```
go get -u github.com/johnwiichang/configure
```

Demo:
```
package main

import (
	"fmt"

	"github.com/johnwiichang/configure"
)

type server struct {
	Address    configure.Field
	DataCentre configure.Field `nilable:"yes"`
	Position   configure.Field `key:"Location"`
}

//Config コンフィギュレーションの基本構造。
var Config = struct {
	Name   configure.Field
	Server server `external:"yes"`
}{}

func main() {
	map1 := map[string]interface{}{
		"Name": "チョウソウイ",
	}
    // Let's load the first configuration entity.
	err := configure.Load(map1, &Config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Config.Server)
	map2 := map[string]interface{}{
		"Address":  "0.0.0.0",
		"Location": "CKG",
	}
    // The additional configurations.
	err = configure.Load(map2, &Config.Server)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Config.Server)
	map2["Location"] = "TYO"
    // Reload the latest configuration without dirty actions when error occurred.
	err = configure.Load(map2, &Config.Server)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Config.Server)
	delete(map2, "Location")
    // This action will not be applied to provider
	err = configure.Load(map2, &Config.Server)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Config.Server)
}
```

Output:
```
{<nil> <nil> <nil>}
{0.0.0.0 <nil> CKG}
{0.0.0.0 <nil> TYO}
config: the key Location is required but get nil
{0.0.0.0 <nil> TYO}
```