CREATE TABLE checkouts
(
    pre_order_id    VARCHAR(36)  NOT NULL COMMENT '预订单ID',
    user_id         INT UNSIGNED NOT NULL COMMENT '用户ID',
    coupon_id       VARCHAR(36) COMMENT '使用优惠券ID列表（支持多券）',
    original_amount BIGINT       NOT NULL COMMENT '原始金额（单位：分）',
    final_amount    BIGINT       NOT NULL COMMENT '实付金额（单位：分）',
    status          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '状态：0-预占中 1-已确认 2-已取消 3-已过期',
    expire_time     BIGINT       NOT NULL COMMENT '过期时间戳（秒）',
    created_at      TIMESTAMP(3)          DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间（毫秒精度）',
    updated_at      TIMESTAMP(3)          DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (pre_order_id),
    INDEX idx_user_status (user_id, status),
    INDEX idx_expire (expire_time)
) COMMENT ='预订单主表';

CREATE TABLE checkout_items
(
    pre_order_id VARCHAR(64)  NOT NULL COMMENT '预订单ID',
    product_id   INT UNSIGNED NOT NULL COMMENT '商品ID',
    quantity     INT UNSIGNED NOT NULL COMMENT '数量',
    price        BIGINT       NOT NULL COMMENT '当时单价（分）',
    snapshot     JSON         NOT NULL COMMENT '商品快照（名称、规格等）',
    created_at   TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (pre_order_id),
    INDEX idx_product (product_id),
    INDEX idx_preorder (pre_order_id)
) COMMENT ='预订单商品明细';