# FileSync
文件同步工具

## <span id="a1">Compatibility(兼容性)</span>
go 1.16.4

## Getting Started(开始)

### 2.1. Download(下载)
- 下载运行版本 [这里](https://github.com/xuzhuoxi/FileSync/releases).
- 下载仓库:
	```sh
	go get -u github.com/xuzhuoxi/FileSync
	```

### 2.2. Build(构建)
Execution the construction file([goxc_build.sh](/goxc_build/build.sh)) to get the releases if you have already downloaded the repository.
You can modify the construction file([goxc_build.sh](/goxc_build/build.sh)) to achieve what you want if necessary. The command line description is [here](https://github.com/laher/goxc).

### 2.3. Run(运行)
仅支持命令行运行

#### 2.3.1. Command Line(命令行说明)
支持两类命令行行为
1. 使用加载配置文件中的参数运行:
	命令格式：`工具路径 -file=配置文件路径 -main=配置任务名/配置任务组名`
	配置文件格式说明请参考：[配置文件格式](#配置文件说明)
	- 形式一：指定具体执行任务执行
		例子：
		```sh
		FileSync -file=demo.yaml -main=copy
		```
	- 形式二：不指定具体执行任务，执行配置文件中默认任务 
		例子：
		```sh
		FileSync -file=demo.yaml
		```
2. 直接使用命令行设置参数运行
	命令格式：`工具路径 -mode=执行模式 -src=来源信息 -tar=目录信息 -include=选择处理设置 -exclude=排除处理设置 -args=执行参数`
	例子：
	```sh
	FileSync -mode=copy -src=/data/src1/*;/data/src2 -tar=/data/tar include=*.jpg exclude=*.txt args=/Lf/Lp/r/s
	```

## 3. Manual(用户手册)

### 3.1. 支持模式

#### 3.1.1. clear
- 功能说明：清除空目录
- 必要参数：
	- mode 	
		mode=clear	
	- src
		只针对目录处理
- 可选参数：
	- include
		只支持dir部分，忽略file部分
	- exclude
		只支持dir部分，忽略file部分
	- args
		- 支持执行参数： /Lf /Lp /r
		- /r
			- **启用**：**目录及子目录**中**不包含**文件时，会被清除	
			- **关闭**：目录为**空**时，会被清除
		- 具体执行参数说明请看[执行参数说明](#执行参数说明)
- 忽略参数：tar

#### 3.1.2. copy
- 功能说明：复制文件或目录到指定目录
- 必要参数：
	- mode 	
		mode=copy
	- src
		复制数据来源，详情请看[这里](#src)
	- tar
		目标路径，**只支持**目录路径，**不支持**多个路径，详情请看[这里](#tar)
- 可选参数：
	- include
		无特殊要求，详情请看[这里](#include)
	- exclude
		无特殊要求，详情请看[这里](#exclude)
	- args
		- 支持执行参数：/i /Lf /Lp /r /s /size /time
		- /i
			- **启用**：命中到的目录**不加入**到处理列表，结果为**空**目录**不会**被复制
			- **关闭**：命中到的目录**会加入**到处理列表，结果为**空**目录**会**被复制
		- /r
			- **启用**：会扫描src指定的文件，以及**目录及子目录**的全部文件
			- **关闭**：只扫描src指定的目录及文件
		- /s
			- **启用**：复制文件时会按照**原来的目录结构**进行复制
			- **关闭**：全部文件或空目录**平铺复制**到tar目录下，名字相同则**覆盖**
		- /size
			- **启用**：只有目标文件**不存在**或源文件的size**大于**目标文件时才进行复制并覆盖
		- /time
			- **启用**：只有目标文件**不存在**或源文件的**修改时间大于**目标文件时才进行复制并覆盖
		- 注意
			- /size与/time都**关闭**时，直接复制
			- /size与/time都**启用**时，**同时**满足才会复制
		- 具体执行参数说明请看[执行参数说明](#执行参数说明)

#### 3.1.3. delete
- 功能说明：删除文件或目录
- 必要参数：
	- mode 	
		mode=delete	
	- src
		目录及文件列表
- 可选参数：
	- include
		dir部分只用于选择，不会加入处理列表
		其它详情请看[这里](#include)
	- exclude
		无特殊要求，详情请看[这里](#exclude)
	- args
		- 支持执行参数： /Lf /Lp /r
		- /r
			- **启用**：对子目录中的文件进行查找，命中则加入到处理列表
			- **关闭**：忽略子目录
		- 具体执行参数说明请看[执行参数说明](#执行参数说明)
- 忽略参数：tar

#### 3.1.4. move
- 功能说明：移动文件或目录到指定目录
- 必要参数：
	- mode 	
		mode=move
	- src
		移动数据来源，详情请看[这里](#src)
	- tar
		目标路径，**只支持**目录路径，**不支持**多个路径，详情请看[这里](#tar)
- 可选参数：
	- include
		无特殊要求，详情请看[这里](#include)
	- exclude
		无特殊要求，详情请看[这里](#exclude)
	- args
		- 支持执行参数：/i /Lf /Lp /r /s /size /time
		- /i
			- **启用**：命中到的目录**不加入**到处理列表，结果为**空**目录**不会**被移动
			- **关闭**：命中到的目录**会加入**到处理列表，结果为**空**目录**会**被移动
		- /r
			- **启用**：会扫描src指定的文件，以及**目录及子目录**的全部文件
			- **关闭**：只扫描src指定的目录及文件
		- /s
			- **启用**：复制文件时会按照**原来的目录结构**进行移动
			- **关闭**：全部文件或空目录**平铺复制**到tar目录下，名字相同则**覆盖**
		- /size
			- **启用**：只有目标文件**不存在**或源文件的size**大于**目标文件时才进行移动并覆盖
		- /time
			- **启用**：只有目标文件**不存在**或源文件的**修改时间大于**目标文件时才进行移动并覆盖
		- 注意
			- /size与/time都**关闭**时，直接移动
			- /size与/time都**启用**时，**同时**满足才会移动
			- **只有文件移动完成后，源目录为空才会执行目录移动**
		- 具体执行参数说明请看[执行参数说明](#执行参数说明)

#### 3.1.5. sync
- 功能说明：双向同步两个目录 或 单向同步
- 必要参数：
	- mode 	
		mode=sync
	- src
		源目录，只支持目录路径，不支持多个，不支持通配符
	- tar
		目标目录，只支持目录路径，不支持多个，不支持通配符
- 可选参数：
	- include
		无特殊要求，详情请看[这里](#include)
	- exclude
		无特殊要求，详情请看[这里](#exclude)
	- args
		- 支持执行参数：/d /i /Lf /Lp /r /size /time
		- /d
			- **启用**：双向同步
			- **关闭**：单身同步，src => tar
		- /i
			- **启用**：命中到的目录**不加入**到处理列表，结果为**空**目录**不会**被同步
			- **关闭**：命中到的目录**会加入**到处理列表，结果为**空**目录**会**被同步
		- /r
			- **启用**：会扫描src指定的文件，以及**目录及子目录**的全部文件
			- **关闭**：只扫描src指定的目录及文件
		- /size
			- **启用**：只有目标文件**不存在**或源文件的size**大于**目标文件时才进行移动并覆盖
		- /time
			- **启用**：只有目标文件**不存在**或源文件的**修改时间大于**目标文件时才进行移动并覆盖
		- 注意
			- /d启用时，/size和/time**有且只有**一个启用
			- /d关闭时，/size和/time**至少**有一个启用，两个都启用时为**且**关系
			- 同步功能必然会**保持目录结构**
		- 具体执行参数说明请看[执行参数说明](#执行参数说明)

### 3.2. 配置文件说明
使用yaml格式的配置文件,结构如下：
```
main: string                    //默认，可填入 任务名 或 任务组名称
groups:                         //任务组数组
	- {
		name:       string      //任务组名称，用于区分每个任务或任务组
		targets:    string      //目标任务列表，各个任务间使用英文逗号“,”分隔
	  }
targets:                        //任务数组
	- {
		name:       string      //任务名称：	用于区分每个任务或任务组
		mode:       string      //任务模式：	用于区分任务真实的执行行为
		src:        string      //任务来源：	任务的文件或目录来源，支持通配符
		tar:        string      //任务目标：	任务的文件或目录的去处
		include:    string      //包含配置：	任务来源的命中包含设置，支持通配符
		exclude:    string      //排序配置：	任务来源的命中排除设置，支持通配符
		args:       string      //执行参数：	执行时的行为参数设置
	  }
```

### 3.3. 配置文件参数说明

#### 3.3.1. name 任务标识
用于区分不同任务和任务组

#### 3.3.2. mode 执行模式
现支持模式有：[clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)

#### 3.3.3. src
来源路径，支持多个路径，使用";"分隔
支持通配符路径，如”\data\\*.png“, "\*"匹配[0,n)个字符
**注意**："\data"和"\data\"相同，指的是data目录，"\data\\*"指的是data目录下的全部文件

#### 3.3.4. tar
目标路径，**只支持**目录路径，**不支持**多个路径
**注意**：[clear](#clear)、[delete](#delete)模式下忽略当前参数

#### 3.3.5. include
格式：“file:\*.jpg,123.png;dir:folder1,fi\*er2,fi\*”，"file"部分与"dir"部分使用";"分隔
**注意**：不配置"-include"或者"-include=空"时，**匹配全部src符合要求的文件**		
- file部分
	支持具体文件名，如 "123.png"等
	支持通配符，如 "\*.jpg"、"a\*b.jpg"等，其中"\*"匹配[0,n)个字符
	多个使用","分隔
- dir部分
	支持具体文件名，如 "folder1"
	支持通配符，如 "fi\*er2"、"fi\*"等，其中"\*"匹配[0,n)个字符
	多个使用","分隔

#### 3.3.6. exclude
格式：“file:\*.jpg,123.png;dir:folder1,fi\*er2,fi\*”，"file"部分与"dir"部分使用";"分隔
**注意**：不配置"-exclude"或者"-exclude=空"时，**不排除文件**		
- file部分
	支持具体目录名，如 "123.png"等
	支持通配符，如 "\*.jpg"、"a\*b.jpg"等，其中"\*"匹配[0,n)个字符
	多个使用","分隔
- dir部分
	支持具体文件名，如 "folder1"
	支持通配符，如 "fi\*er2"、"fi\*"等，其中"\*"匹配[0,n)个字符
	多个使用","分隔

#### 3.3.7. args 
- 执行参数，支持如下：**[/d](#/d)**,	**[/i](#/i)**,	**[/Lf](#/Lf)**,	**[/Lp](#/Lp)**,	**[/r](#/r)**,	**[/s](#/s)**,	**[/size](#/size)**,	**[/time](#/time)**
- 多个参数可直接拼接，如"/d/i/Lf"
- 具体执行参数说明请看[执行参数说明](#执行参数说明)

### 3.4. 执行参数说明

#### 3.4.1. /d (double)
- 说明：开启双向同步，默认为单向
- 适用范围： [sync](#FileSync)

#### 3.4.2. /i (ignore empty)
- 说明：忽略空目录，默认为不忽略
- 适用范围： [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)

#### 3.4.3. /Lf (log file)
- 说明：开启记录日志，默认不记录日志
- 适用范围： [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)

#### 3.4.4. /Lp (print)
- 说明：控制台打印信息，默认不打印信息
- 适用范围： [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)

#### 3.4.5. /r (recurse)
- 说明：递归，默认不递归
- 适用范围： [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)

#### 3.4.6. /s (stable)
- 说明：保持文件目录结构，默认不保持
- 适用范围： [copy](#copy), [move](#move)

#### 3.4.7. /size (update by size)
- 说明：按文件大小更新
- 适用范围： [copy](#copy), [move](#move), [sync](#sync)

#### 3.4.8. /time (update by time)
- 说明：按时间更新
- 适用范围： [copy](#copy), [move](#move), [sync](#sync)

### 3.5. 命中(过滤)说明

### 3.6. 命令行说明

## 4. Related Library(依赖)
- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)
- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc) 

## 5. Contact(联系方式)
xuzhuoxi 
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

[锚点测试](\a1)

## 6. License(开源许可证)
~~FileSync source code is available under the MIT [License](/LICENSE).~~
