

-- 优惠券使用记录表
CREATE TABLE `coupon_usage`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT,
    `order_id`        varchar(32)              NOT NULL COMMENT '关联的订单ID',
    `coupon_id`       varchar(32)              NOT NULL COMMENT '优惠券ID',
    `user_id`         int unsigned             NOT NULL COMMENT '用户ID',
    `coupon_type`     tinyint                  NOT NULL COMMENT '当时优惠券类型：1-满减 2-折扣 3-立减',
    `origin_value`    decimal(10, 2)                      NOT NULL COMMENT '当时优惠值（根据类型：分/百分比）',
    `discount_amount` decimal(10, 2)                      NOT NULL COMMENT '抵扣金额（分）',
    `applied_at`      TIMESTAMP                NOT NULL COMMENT '应用时间',
    PRIMARY KEY (`id`),
    KEY `idx_order` (`order_id`),
    KEY `idx_user` (`user_id`)
) COMMENT ='优惠券使用明细';
