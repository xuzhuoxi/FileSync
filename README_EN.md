# FileSync  
File synchronization tool.  
Can be used for fixed or regular file copy, delete, move functions, and one-way or two-way synchronization of folders.  

[中文](README.md) | English  

## <span id="a1">1. Compatibility</span>  
go1.16.15  

## <span id="a2">2. Start</span>  

### <span id="a2.1">2.1. Download</span>  
- Download the release version [here](https://github.com/xuzhuoxi/FileSync/releases).  
- Download repository:  
```sh
  go get -u github.com/xuzhuoxi/FileSync
```

### <span id="a2.2">2.2. Build<span>  
- If you have downloaded the entire warehouse and related dependent warehouses, you can execute ([goxc_build.sh](/goxc_build/build.sh)) to build to get the executable program.  
- If necessary, you can modify ([goxc_build.sh](/goxc_build/build.sh)) for custom builds, build tools are described in [here](https://github.com/laher/goxc ).  

### <span id="a2.3">2.3. Run<span>  
- only supports command line operation  
- [Example](/demo)  

#### <span id="a2.3.1">2.3.1 Command Line Description<span>  
Two types of command-line actions are supported  

- run with the parameters in the loaded configuration file:  
  Command format: `execution file path -file=configuration file path -main=configuration task name/configuration task group name`  
  For configuration file format description, please refer to: [Configuration file description](#a3.2)  
  - Specify specific tasks to execute, such as:  
```sh
    FileSync -file=demo.yaml -main=copy
```
  - Do not specify a specific execution task, execute the default task in the configuration file, such as:  
```sh
    FileSync -file=demo.yaml
```

- Set parameters directly on the command line to run  
  Command format: `tool path -mode=execution mode -src=source information -tar=directory information -include=select processing settings -exclude=exclude processing settings -args=execution parameters`  
  E.g:  
```sh
  FileSync -mode=copy -src=/data/src1/*;/data/src2 -tar=/data/tar include=*.jpg exclude=*.txt args=/Lf/Lp/r/s
```

## <span id="a3">3. User Manual<span>  

### <span id="a3.1">3.1. Support mode<span>  

#### <span id="clear">3.1.1. clear<span>  

- Function Description  
  clear empty directories  

- required parameter  
  - mode  
    mode=clear  
  - src  
    only for directories  

- optional parameter  
  - include  
    Only the dir part is supported, the file part is ignored  
  - exclude  
    Only the dir part is supported, the file part is ignored  
  - args  
    - Support execution parameters: /Lf /Lp /r  
    - /r  
      - **Enable**: When **does not contain** files in **directories and subdirectories**, they will be cleared  
      - **OFF**: When the directory is **empty**, it will be cleared  
    - Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

- ignore arguments: tar  

#### <span id="copy">3.1.2. copy<span>  

- Function Description  
  Copy a file or directory to a specified directory  

- required parameter  
  - mode  
    mode=copy  
  - src  
    Copy the data source, please see [here](#src) for details   
  - tar  
    Target path, **only supports** directory path, **does not support** multiple paths, please see [here](#tar) for details  

- optional parameter  
  - include  
    See [here] for details (#include)  
    Special requirements: Ignore when /file in -args is enabled  
  - exclude  
    See [here] for details (#exclude)  
    Special requirements: Ignore when /file in -args is enabled     
  - args  
    - Support execution parameters: /i /Lf /Lp /r /s /size /time  
    - /i  
      - **Enable**: Hit directories **not added** to the processing list, resulting in **empty** directories **not** copied  
      - **OFF**: The hit directory **will be added** to the processing list, and the result will be **empty** directory **will be** copied  
    - /r  
      - **Enable**: The file specified by src and all files in **directory and subdirectory** will be scanned  
      - **OFF**: Only scan directories and files specified by src  
    - /s  
      - **Enable**: When copying files, they will be copied according to the **original directory structure**  
      - **Close**: All files or empty directories **tile copy** to the tar directory, **overwrite** if the names are the same  
    - /file  
      - **Enable**: Single file processing mode, src and tar must be file paths. The include and exclude parameters are ignored.   
      - **OFF**: Default  
    - /size  
      - **Enable**: Copy and overwrite only when the target file **does not exist** or the size of the source file** is larger than** the target file  
    - /time  
      - **Enable**: Copy and overwrite only if the target file **does not exist** or the **modified time of the source file is greater than** the target file  
    - /md5  
      - **Enable**: Copy and overwrite only when the target file **does not exist** or the **md5 value of the source file is not equa**l to the target file  
    - Notice  
      - Copy directly when /size, /time and /md5 are **closed**  
      - When /size, /time and /md5 are **enabled**, they will only be copied if they are satisfied at the same time  
    - Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

#### <span id="delete">3.1.3. delete<span>  

- Function Description  
  delete file or directory  

- required parameter  
  - mode  
    mode=delete  
  - src  
    Directory and file listing  

- optional parameter  
  - include  
    The dir part is only used for selection and will not be added to the processing list  
    See [here](#include) for more details   
  - exclude  
    No special requirements, please see [here](#exclude) for details   
  - args  
    - Support execution parameters: /Lf /Lp /r  
    - /r  
      - **Enable**: Search for files in subdirectories, and add them to the processing list if they are found  
      - **OFF**: Ignore subdirectories  
    - Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

- ignore arguments: tar  

#### <span id="move">3.1.4. move<span>  

- Function Description  
  Move a file or directory to a specified directory  

- required parameter  
  - mode  
    mode=move  
  - src  
    Mobile data sources, see [here](#src) for details  
  - tar  
    Target path, **only supports** directory path, **does not support** multiple paths, please see [here](#tar) for details  

- optional parameter  
  - include  
    See [here] for details (#include)  
    Special requirements: Ignore when /file in -args is enabled  
  - exclude  
    See [here] for details (#exclude)  
    Special requirements: Ignore when /file in -args is enabled   
  - args  
    - Support execution parameters: /i /Lf /Lp /r /s /size /time   
      -/i   
      - **Enable**: The hit directory **does not add** to the processing list, the result is **empty** directory **will not** be moved  
      - **OFF**: The hit directory **will be added** to the processing list, and the result will be **empty** directory **will be**moved  
    - /r  
      - **Enable**: The file specified by src and all files in **directory and subdirectory** will be scanned  
      - **OFF**: Only scan directories and files specified by src  
    - /s  
      - **Enable**: When copying files, they will be moved according to the **original directory structure**  
      - **Close**: All files or empty directories **tile copy** to the tar directory, **overwrite** if the names are the same  
    - /file  
      - **Enable**: Single file processing mode, src and tar must be file paths. The include and exclude parameters are ignored.   
      - **OFF**: Default  
    - /size  
      - **Enable**: Move and overwrite only when the target file **does not exist** or the size of the source file** is larger than the target file  
    - /time  
      - **Enable**: Move and overwrite only if the target file **does not exist** or the **modified time of the source file is greater than** the target file  
    - /md5  
      - **Enable**: Move and overwrite only when the target file **does not exist** or the **md5 value of the source file is not equal to the target file  
    - Notice  
      - When /size, /time and /md5 are **closed**, move directly  
      - When /size, /time and /md5 are **enabled**, they will only move when **simultaneous** are satisfied  
      - **Directory movement will only be performed if the source directory is empty after the file movement is complete**  
    - Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

#### <span id="sync">3.1.5. sync<span>

- Function Description  
  Two-way synchronization of two directories or one-way synchronization  

- required parameters  
  - mode  
    mode=sync  
  - src  
    Source directory, only supports directory path, does not support multiple, does not support wildcards  
  - tar  
    Target directory, only supports directory path, does not support multiple, does not support wildcards  

- optional parameter  
  - include  
    No special requirements, please see [here](#include) for details  
  - exclude  
    No special requirements, please see [here](#exclude) for details   
  - args  
    - Support execution parameters: /d /i /Lf /Lp /r /size /time  
    - /d  
      - **Enable**: 2-way sync  
      - **close**: singleton sync, src => tar  
    - /i  
      - **Enable**: The hit directory **is not added** to the processing list, and the result is **empty** directory **will not** be synchronized  
      - **OFF**: The hit directory **will be added** to the processing list, and the result will be **empty** directory **will be** synchronized  
    - /r  
      - **Enable**: The file specified by src and all files in **directory and subdirectory** will be scanned  
      - **OFF**: Only scan directories and files specified by src  
    - /size  
      - **Enable**: Synchronize and overwrite only when the target file **does not exist** or the size of the source file** is larger than** the target file  
    - /time  
      - **Enable**: Synchronize and overwrite only when the target file **does not exist** or the **modified time of the source file is greater than** the target file  
    - /md5  
      - **Enable**: Synchronize and overwrite only when the target file **does not exist** or the **md5 value of the source file is not equal to** the target file  
    - Notice  
      - When /d is enabled, **only one and only one of** /size and /time is enabled, and the /md5 parameter is not supported  
      - When /d is disabled, **at least one of** /size, /time and /md5 is enabled, and when both are enabled, it is **and**relation  
      - The sync function necessarily **maintains the directory structure**  
    - Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

### <span id="a3.2">3.2. Configuration file description<span>  
Using the configuration file in yaml format, the structure is as follows:  
```yaml
main: string //default, you can fill in task name or task group name
groups: //Array of task groups
  - {
    name: string //Task group name, used to distinguish each task or task group
    tasks: string //List of target tasks, each task is separated by a comma ","
    }
tasks: //task array
  - {
    name: string //Task name: used to distinguish each task or task group
    mode: string //task mode: used to distinguish the real execution behavior of the task
    src: string //Task source: The file or directory source of the task, wildcards are supported
    tar: string //task target: where the task's file or directory goes
    include: string //Include configuration: Hit include settings for task sources, support wildcards
    exclude: string //Sort configuration: Hit exclusion settings for task sources, wildcards are supported
    args: string //Execution parameters: Behavior parameter settings during execution
    }
```

### <span id="a3.3">3.3. Configuration file parameter description<span>  

#### <span id="name">3.3.1. name (task ID)<span>  
- Used to distinguish between different tasks and task groups  

#### <span id="mode">3.3.2. mode (execution mode)<span>  
- Now supported modes are: [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)  

#### <span id="src">3.3.3. src (source path)<span>  
- Source path, multiple paths are supported, separated by ";"  
- Supports wildcard paths, such as "\data\\*.png", "\*" matches [0,n) characters  
- **Note**: "\data" is the same as "\data\\", referring to the data directory, "\data\\*" refers to all files in the data directory  

#### <span id="tar">3.3.4.tar (target path)<span>  
- Target path, **only supports** directory paths, **does not support** multiple paths  
- **Note**: The current parameter is ignored in [clear](#clear), [delete](#delete) mode  

#### <span id="include">3.3.5.include (include configuration)<span>  
- Format: "file:\*.jpg,123.png;dir:folder1,fi\*er2,fi\*", the "file" part and the "dir" part are separated by ";"  
- **Note**: When "-include" or "-include=empty" is not configured, **match all src files that meet the requirements**  
- file part  
  Support specific file names, such as "123.png", etc.  
  Wildcards are supported, such as "\*.jpg", "a\*b.jpg", etc., where "\*" matches [0,n) characters  
  Use "," to separate multiple  
- dir part  
  Support specific file names, such as "folder1"  
  Wildcards are supported, such as "fi\*er2", "fi\*", etc., where "\*" matches [0,n) characters  
  Use "," to separate multiple  

#### <span id="exclude">3.3.6. exclude (exclude configuration)<span>  
- Format: "file:\*.jpg,123.png;dir:folder1,fi\*er2,fi\*", the "file" part and the "dir" part are separated by ";"  
- **Note**: When "-exclude" or "-exclude=empty" is not configured, **files are not excluded**  
- file part  
  Support specific directory names, such as "123.png", etc.  
  Wildcards are supported, such as "\*.jpg", "a\*b.jpg", etc., where "\*" matches [0,n) characters  
  Use "," to separate multiple  
- dir part  
  Support specific file names, such as "folder1"  
  Wildcards are supported, such as "fi\*er2", "fi\*", etc., where "\*" matches [0,n) characters  
  Use "," to separate multiple  

#### <span id="args">3.3.7. args (execution parameters)<span>  
- Execution parameters, supported as follows: **[/d](#/d)**, **[/i](#/i)**, **[/Lf](#/Lf)**, **[/Lp](#/Lp)**, **[/r](#/r)**, **[/s](#/s)**,  **[/file](#/file)**,  **[/size](#/size)**,  **[/time](#/time)**, **[/md5](#/md5)**   
- Multiple parameters can be directly spliced, such as "/d/i/Lf"  
- Please refer to [execution parameter description](#a3.4) for specific execution parameter description  

### <span id="a3.4">3.4. Execution parameter description<span>  

#### <span id="d">3.4.1. /d (two-way sync)<span>  
- Description: Enable two-way synchronization, the default is one-way  
- Scope: [sync](#FileSync)  

#### <span id="i">3.4.2. /i (ignore empty directories)<span>  
- Description: Ignore empty directories, the default is not to ignore  
- Scope: [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)  

#### <span id="Lf">3.4.3. /Lf (open file log)<span>  
- Description: Enable logging, no logging by default  
- Scope: [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)  

#### <span id="Lp">3.4.4. /Lp (enable printing)<span>  
- Description: the console prints information, the default does not print information  
- Scope: [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)  

#### <span id="r">3.4.5. /r (turn on recursion)<span>  
- Description: recursive, default not recursive  
- Scope: [clear](#clear), [copy](#copy), [delete](#delete), [move](#move), [sync](#sync)  

#### <span id="s">3.4.6. /s (keep directory structure)<span>  
- Description: keep the file directory structure, the default is not maintained  
- Scope: [copy](#copy), [move](#move)  

#### <span id="file">3.4.7. /file (single file processing mode)<span>
- Description: For single-file processing, src and tar must be file paths, and wildcards and multi-paths are not supported. The include and exclude parameters are ignored.  
- Scope: [copy](#copy), [move](#move)  

#### <span id="size">3.4.8. /size (processing according to file size difference)<span>  
- Description: Process by file size  
- Scope: [copy](#copy), [move](#move), [sync](#sync)  

#### <span id="time">3.4.9. /time (processing according to file modification time difference)<span>  
- Description: Process by time  
- Scope: [copy](#copy), [move](#move), [sync](#sync)  

#### <span id="time">3.4.10. /md5 (processing according to file md5 value difference)<span>  
- Description: Process according to md5 value  
- Scope: [copy](#copy), [move](#move), [sync](#sync)  

### <span id="a3.5">3.5. Hit (filter) description<span>  
The following is the **conventional logic** of the hit judgment of the file  
1. Wildcards in [src](#src)  
  - Table of contents  
  - document  
  - wildcard  
2. Wildcards in [exclude](#exclude) (if present)  
  - Complies with: Excluded  
  - Non-conformance: proceed to the next judgment  
3. Wildcards in [include](#include) (if present)  
  - Match: add to hit list  
  - Non-Compliant: Excluded  

**Notice**:  

- The dir parameter in [exclude](#exclude) **exists** and **matches**, the current directory and subdirectories will be excluded  
- The dir parameter in [include](#include) **exists** and **does not match**, the current directory and subdirectories will be excluded  
- The above is general logic, please refer to the relevant mode description for specific behavior: [Support mode](#a3.1)  

### <span id="a3.6">3.6. Command Line Description<span>  
- Supports **specified configuration operation** and **direct parameter operation** two types of command line functions  
- **NOTE PRIORITIES**: Specify Configuration Run > Direct Argument Run  

1. Specify the configuration to run:  
Command format: `tool path -file=configuration file path -main=configuration task name/configuration task group name`  
For the description of the configuration file format, please refer to: [Configuration file format](#a3.2)  
  - Form 1: Specify specific tasks to be executed  
    example:  
```sh
  FileSync -file=demo.yaml -main=copy
```
  - Form 2: Do not specify a specific execution task, execute the default task in the configuration file  
    example:  
```sh
  FileSync -file=demo.yaml
```
2. Direct parameter operation  
Command format: `tool path -mode=execution mode -src=source information -tar=directory information -include=select processing settings -exclude=exclude processing settings -args=execution parameters`  
example:  
```sh
  FileSync -mode=copy -src=/data/src1/*;/data/src2 -tar=/data/tar include=*.jpg exclude=*.txt args=/Lf/Lp/r/s
```

## <span id="a4">4. Core dependencies<span>  
- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)  
- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)  

## <span id="a5">5. Contact<span>  
xuzhuoxi  
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com>  

## <span id="a6">6. Open Source License<span>  
FileSync source code is available under the MIT [License](/LICENSE).  