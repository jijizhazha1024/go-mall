CREATE TABLE `carts` (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键 自增',
                         `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
                         `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
                         `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
                         `user_id` int(11) NOT NULL COMMENT '用户ID',
                         `product_id` int(11) NOT NULL COMMENT '商品ID',
                         `quantity` int(11) DEFAULT '0' COMMENT '商品数量',
                         `checked` tinyint(1) DEFAULT '0' COMMENT '商品是否选中',
                         PRIMARY KEY (`id`),
                         KEY `idx_carts_deleted_at` (`deleted_at`),
                         KEY `idx_carts_user_id` (`user_id`),
                         KEY `idx_carts_product_id` (`product_id`),
                         CONSTRAINT `fk_carts_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
                         CONSTRAINT `fk_carts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;