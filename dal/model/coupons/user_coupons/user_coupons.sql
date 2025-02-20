-- 用户优惠券关联表
CREATE TABLE `user_coupons`
(
    `id`        int unsigned NOT NULL AUTO_INCREMENT,
    `user_id`   int unsigned NOT NULL COMMENT '用户ID',
    `coupon_id` varchar(36)  NOT NULL COMMENT '优惠券ID',
    `status`    tinyint      NOT NULL DEFAULT 0 COMMENT '状态：0-未使用 1-已使用 2-已过期',
    `order_id`  varchar(36)           DEFAULT NULL COMMENT '使用的订单ID',

    `used_at`   TIMESTAMP             DEFAULT NULL COMMENT '使用时间',
    created_at  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_user_coupon` (`user_id`, `coupon_id`),
    UNIQUE KEY `uniq_order_id` (`order_id`),
    KEY `idx_user_status` (`user_id`, `status`),
    KEY `idx_order` (`order_id`)
) COMMENT ='用户持有优惠券';