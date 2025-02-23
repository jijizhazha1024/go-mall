CREATE TABLE `inventory_lock` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `order_id` VARCHAR(64) NOT NULL COMMENT '唯一订单ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  
    primary key (`id`),
    UNIQUE KEY `uniq_order_user` (`order_id`, `user_id`)
 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='幂等锁表';