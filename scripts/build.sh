#!/bin/bash

# 检查当前目录是否为 go-mall
CURRENT_DIR=$(basename "$PWD")
if [ "$CURRENT_DIR" != "go-mall" ]; then
  echo "Error: Please run this script from the go-mall directory."
  exit 1
fi

# 设置基础路径
BASE_DIR="/app"
APIS_DIR="$BASE_DIR/apis"
SERVICES_DIR="$BASE_DIR/services"
ETC_DIR="$BASE_DIR/etc"

# 创建必要的目录
echo "building..."
mkdir -p "$APIS_DIR" "$SERVICES_DIR" "$ETC_DIR"

# 构建APIs
build_module() {
  local module_dir=$1
  local target_dir=$2
  pushd "$module_dir" || { echo "Failed to enter directory $module_dir"; exit 1; }
  for i in *; do
    echo "building $i"
    capName=$(echo "${i^}")
    cd "$i" || { echo "Failed to enter directory $i"; exit 1; }
    go build -ldflags="-s -w" -o "$target_dir/${capName}Api" || { echo "Failed to build $i"; exit 1; }
    cp ./etc/* "$ETC_DIR/" 2>/dev/null ||  { echo "Failed to copy etc files"; exit 1;}
    cd ..
  done
  popd || { echo "Failed to leave directory $module_dir"; exit 1; }
}

echo "building apis"
build_module "apis" "$APIS_DIR"
echo "build apis done"

echo "building services"
build_module "services" "$SERVICES_DIR"
echo "build services done"

echo "build done"