CREATE TABLE payments
(
    payment_id      VARCHAR(36)  NOT NULL COMMENT '支付单ID（UUID）',
    pre_order_id    VARCHAR(64)  NOT NULL COMMENT '预订单ID（外键）',
    order_id        VARCHAR(36)  DEFAULT NULL COMMENT '关联订单ID（支付成功后更新）',

    original_amount BIGINT       NOT NULL COMMENT '订单原价（单位：分）',
    paid_amount     BIGINT       DEFAULT NULL COMMENT '实付金额（分）',

    -- 支付信息
    payment_method  VARCHAR(20)  NOT NULL COMMENT '支付渠道（wx_pay/alipay）',
    transaction_id  VARCHAR(64)  DEFAULT NULL COMMENT '支付平台交易号',
    pay_url         Text         NOT NULL COMMENT '支付跳转链接',
    expire_time     BIGINT       NOT NULL COMMENT '支付链接过期时间戳（秒）',

    status          TINYINT      NOT NULL COMMENT '支付状态（0-未定义 1-待支付 2-已支付...）',
    idempotency_key VARCHAR(64)  NOT NULL COMMENT '幂等键（唯一索引）',

    created_at      TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间（毫秒精度）',
    updated_at      TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    paid_at         BIGINT       DEFAULT NULL COMMENT '支付成功时间戳（秒）',

    PRIMARY KEY (payment_id),
    UNIQUE KEY uk_idempotency (idempotency_key),
    INDEX idx_pre_order (pre_order_id),
    INDEX idx_order (order_id),
    INDEX idx_status_method (status, payment_method),
    INDEX idx_create_time (created_at)
)