# 项目名 (gitlab-version)

用于识别gitlab版本号

## 安装

您可以通过以下步骤安装该项目：

1. 打开终端并导航到您希望安装该项目的目录中。
2. 使用以下命令克隆该仓库：

```sh
git clone https://github.com/XTeam-Wing/gitlab-version.git
cd gitlab-version
go mod download
go run main.go -h
```
   
## 使用
    
```sh
  go run main.go -i https://target.com
  go run main.go -i target.txt
  go run main.go -i https://target1.com,https://target2.com
```
   ## 贡献
欢迎贡献代码！