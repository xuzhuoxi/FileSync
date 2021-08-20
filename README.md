# FileSync
文件同步工具


## Compatibility(兼容性)
go 1.16.4

## Related Library(依赖)

- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)

- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc) 

## Getting Started(开始)

### Download Release(下载)

- 下载运行版本 [here](https://github.com/xuzhuoxi/ImageSplitter/releases).

- 下载仓库:

	```sh
	go get -u github.com/xuzhuoxi/FileSync
	```

### Build(构建)

Execution the construction file([build.sh](/build/build.sh)) to get the releases if you have already downloaded the repository.

You can modify the construction file([build.sh](/build/build.sh)) to achieve what you want if necessary. The command line description is [here](https://github.com/laher/goxc).

### Run(运行)

#### Demo(例子)

[Here](/demo/mac) is a running demo for MacOS platform.

The running command is consistent of all platforms.

Goto <a href="#command-line">Command Line Description</a>.

#### Command Line(命令行说明)

Supportted command line parameters as follow:

| -       | -            | -                                                            |
| :------ | :----------- | ------------------------------------------------------------ |
| -mode   | optional | The mode of the divisions.  1：小图使用固定尺寸；	2：小图使用平均尺寸|
| -order  | optional | The order of the divisions. 1：左上角为起始点；	2：左下角为起始点|
| -size   | **required**     | The size info of divisions. 格式：mxn。当mode为1时，m、n代表小图尺寸；当mode为2时，m、n代表分割数量|
| -in     | **required**     | Custom source file. |
| -out    | **required**     | Custom output files. 支持通配符（{n0},{N0},{n1},{N1},{x0},{X0},{x1},{X1},{y0},{Y0},{y1},{Y1}）|
| -format | optional     | The format of the generated image. Supported as follows: png, jpg, jpeg, jps |
| -ratio  | optional     | The quality of the generated image. Supported for jpg,jpeg,jps. |

E.g.:

-mode=1

-mode=2

-order=1

-order=2

-size=256x256

-in=./source/image.png

-in=/Users/aaa/image.jpg

-out=/Users/aaa/image_{n0}_{y1}_{x1}.png

-out=./out/image_{n0}_{y1}_{x1}.png

-format=jpeg

-format=jpg

-format=png

-ratio=85

## Manual(用户手册)

### 支持功能

#### clear

#### copy

#### delete

#### move

#### sync

### 参数说明

- /d

	- 说明：双向同步
	
	- 适用范围： [sync](#FileSync)


- /f

- /i

- /L

- /l

- /r

- /s

- /u

### 配置文件说明

### 命令行说明

## Contact(联系方式)

xuzhuoxi 

<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>

## License(开源许可证)
ImageSplitter source code is available under the MIT [License](/LICENSE).


