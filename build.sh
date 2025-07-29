#!/bin/bash
# 启用错误捕获
set -e
# 判断是否是 Windows 下 shell 可运行环境
system_info=$(uname -a | tr '[:upper:]' '[:lower:]')
if [[ $system_info == *"cygwin"* || $system_info == *"mingw"* || $system_info == *"msys"* ]]; then
    echo "当前处于 Windows 下 Cygwin/MinGW/MSYS 环境，继续执行脚本。"
elif [[ $system_info == *"microsoft"* ]]; then
    echo "当前处于 Windows 下 WSL 环境，继续执行脚本。"
else
    echo "当前不是 Windows 下 shell 可运行环境，请使用 Cygwin、MinGW、MSYS 或 WSL 环境执行脚本。"
    exit 1
fi
# 定义coreutils目录路径
scriptDir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
directories=("$scriptDir/src/coreutils" "$scriptDir/src/tls-socks" "$scriptDir/src/procps" "$scriptDir/src/other" "$scriptDir/src/net-tools")

# 检查coreutils目录是否存在
for dir in "${directories[@]}"; do
if [ ! -d "$dir" ]; then
    echo "$dir目录不存在，请检查路径。" >&2
    exit 1
fi
done

# 创建输出目录
outputDir="bin"
if [ ! -d "$outputDir" ]; then
    mkdir -p "$outputDir"
    echo "创建输出目录: $outputDir"
fi

# 初始化错误标志
buildErrors=0

# 定义需要交叉编译的平台
platforms=("linux/amd64" "windows/amd64")


# 遍历所有指定目录及其子目录下的所有.go文件
for dir in "${directories[@]}"; do
    echo "正在处理目录: $dir"

   # 特殊处理 tls-socks 目录
  if [[ "$dir" == *"tls-socks"* ]]; then
        projectName="tls-socks"
        echo "特殊处理 $projectName 项目，进行交叉编译..."
        # 为每个平台编译
        for platform in "${platforms[@]}"; do
            OS=$(echo $platform | cut -d'/' -f1)
            ARCH=$(echo $platform | cut -d'/' -f2)

            # 设置输出文件名
            if [ "$OS" == "windows" ]; then
                outputFile="${projectName}.exe"
            else
                outputFile="${projectName}"
            fi

            echo "正在为 $OS/$ARCH 平台构建 $projectName..."

            # 执行交叉编译
            if ! GOOS=$OS GOARCH=$ARCH go build -o "$outputDir/$outputFile" "$dir/tls-socks.go"; then
                echo "为 $OS/$ARCH 平台构建 $projectName 时出错" >&2
                buildErrors=1
            else
                echo "成功构建 $outputDir/$outputFile"
            fi
        done
  else
      /usr/bin/find "$dir" -type f -name "*.go" | while read -r goFile; do
        # 获取文件名（不包含扩展名）
        fileName=$(basename "$goFile" .go)
        # 定义输出的二进制文件路径
        outputPath="$outputDir/${fileName}.exe"

        echo "正在构建 $goFile 到 $outputPath..."
        # 使用go build命令构建二进制文件
        if ! go build -o "$outputPath" "$goFile"; then
            echo "构建 $goFile 时出错: go build 返回非零退出码。" >&2
            buildErrors=1
        else
            echo "成功构建 $outputPath"
        fi
      done
  fi
done

# 根据构建结果返回适当的退出码
if [ $buildErrors -eq 1 ]; then
    echo "构建过程中存在错误，请检查日志。" >&2
    exit 1
else
    echo "所有Go文件构建成功!"
    exit 0
fi