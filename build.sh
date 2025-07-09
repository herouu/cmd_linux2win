#!/bin/bash

# 启用错误捕获
set -e
# 检测是否为 Cygwin 环境
if [ -n "$(uname -a | grep -i cygwin)" ]; then
    echo "当前处于 Cygwin 环境，继续执行构建流程。"
else
    echo "当前不是 Cygwin 环境, 请使用 Cygwin 环境执行脚本。"
    exit 1
fi
# 定义coreutils目录路径
coreutilsDir="src/coreutils"

# 检查coreutils目录是否存在
if [ ! -d "$coreutilsDir" ]; then
    echo "coreutils目录不存在，请检查路径。" >&2
    exit 1
fi

# 创建输出目录
outputDir="bin"
if [ ! -d "$outputDir" ]; then
    mkdir -p "$outputDir"
    echo "创建输出目录: $outputDir"
fi

# 初始化错误标志
buildErrors=0

# 遍历coreutils目录及其子目录下的所有.go文件
find "$coreutilsDir" -name "*.go" | while read -r goFile; do
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

# 根据构建结果返回适当的退出码
if [ $buildErrors -eq 1 ]; then
    echo "构建过程中存在错误，请检查日志。" >&2
    exit 1
else
    echo "所有Go文件构建成功!"
    exit 0
fi