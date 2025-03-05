CREATE TABLE `carts` (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键 自增',
                         `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         `user_id` int(11) DEFAULT NULL COMMENT '用户ID',
                         `product_id` int(11) DEFAULT NULL COMMENT '商品ID',
                         `quantity` int(11) DEFAULT NULL COMMENT '商品数量',
                         `checked` tinyint(1) DEFAULT NULL COMMENT '商品是否选中',
                         PRIMARY KEY (`id`),
                         KEY `idx_carts_user_id` (`user_id`),
                         KEY `idx_carts_product_id` (`product_id`)
)