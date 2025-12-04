#!/bin/bash

# Docker Hub 发布脚本
# 使用方式: ./publish.sh [tag]
# 示例: ./publish.sh v1.0.0

set -e

# 配置
DOCKER_USER="${DOCKER_USER:-your-username}"
IMAGE_NAME="monitor"
TAG="${1:-latest}"

FULL_IMAGE="${DOCKER_USER}/${IMAGE_NAME}"

echo "=========================================="
echo "  Docker Hub 发布脚本"
echo "=========================================="
echo "镜像: ${FULL_IMAGE}:${TAG}"
echo ""

# 检查是否登录
if ! docker info 2>/dev/null | grep -q "Username"; then
    echo "请先登录 Docker Hub..."
    docker login
fi

# 构建镜像
echo "正在构建镜像..."
docker build -t "${FULL_IMAGE}:${TAG}" .

# 如果不是 latest，同时打上 latest 标签
if [ "${TAG}" != "latest" ]; then
    docker tag "${FULL_IMAGE}:${TAG}" "${FULL_IMAGE}:latest"
fi

# 推送镜像
echo "正在推送镜像..."
docker push "${FULL_IMAGE}:${TAG}"

if [ "${TAG}" != "latest" ]; then
    docker push "${FULL_IMAGE}:latest"
fi

echo ""
echo "✓ 发布成功!"
echo "  ${FULL_IMAGE}:${TAG}"
[ "${TAG}" != "latest" ] && echo "  ${FULL_IMAGE}:latest"
