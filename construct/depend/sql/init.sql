/*
 Navicat Premium Dump SQL

 Source Server         : 124.71.72.124
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : 124.71.72.124:3306
 Source Schema         : mall

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 24/01/2025 14:50:23
*/


-- 创建一个单独的用户
CREATE USER IF NOT EXISTS 'jjzzchtt'@'%' IDENTIFIED BY 'jjzzchtt';
GRANT ALL PRIVILEGES ON *.* TO 'jjzzchtt'@'%';

-- 创建数据库
create database if not exists mall character set utf8mb4;


SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`
(
    `id`          int                                                           NOT NULL AUTO_INCREMENT COMMENT '主键，自增',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
    `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci         NULL COMMENT '分类描述',
    `created_at`  timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `idx_category_name` (`name` ASC) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------

-- ----------------------------
-- Table structure for product_categories
-- ----------------------------
DROP TABLE IF EXISTS `product_categories`;
CREATE TABLE `product_categories`
(
    `product_id`  int NOT NULL COMMENT '商品id',
    `category_id` int NOT NULL COMMENT '分类id',
    PRIMARY KEY (`product_id`, `category_id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_categories
-- ----------------------------

-- ----------------------------
-- Table structure for products
-- ----------------------------
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products`
(
    `id`          int                                                           NOT NULL AUTO_INCREMENT COMMENT '主键，自增',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '商品名称',
    `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci         NULL COMMENT '商品描述',
    `picture`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '商品图片信息',
    `price`       decimal(10, 2)                                                NOT NULL COMMENT '商品价格',
    `stock`       int                                                           NULL DEFAULT 0 COMMENT '库存数量',
    `created_at`  timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of products
-- ----------------------------

-- ----------------------------
-- Table structure for shopping_cart
-- ----------------------------
DROP TABLE IF EXISTS `shopping_cart`;
CREATE TABLE `shopping_cart`
(
    `id`         int        NOT NULL AUTO_INCREMENT COMMENT '主键 自增',
    `created_at` timestamp  NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp  NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp  NULL DEFAULT NULL COMMENT '删除时间',
    `is_deleted` tinyint(1) NULL DEFAULT NULL COMMENT '记录是否删除',
    `user_id`    int        NULL DEFAULT NULL COMMENT '用户ID',
    `goods_id`   int        NULL DEFAULT NULL COMMENT '商品ID',
    `nums`       int        NULL DEFAULT NULL COMMENT '商品数量',
    `checked`    tinyint(1) NULL DEFAULT NULL COMMENT '商品是否选中',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `idx_shopping_cart_goods` (`goods_id` ASC) USING BTREE,
    INDEX `idx_shopping_cart_user` (`user_id` ASC) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shopping_cart
-- ----------------------------

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `user_id`       int                                                           NOT NULL AUTO_INCREMENT COMMENT '主键，自增，用户 ID',
    `username`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名，可空',
    `email`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱，唯一',
    `password_hash` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '密码哈希值',
    `avatar_url`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像图片 URL',
    `created_at`    timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `user_deleted`  tinyint(1)                                                    NULL DEFAULT 0 COMMENT '用户是否已删除',
    `updated_at`    timestamp                                                     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`) USING BTREE,
    UNIQUE INDEX `email` (`email` ASC) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
