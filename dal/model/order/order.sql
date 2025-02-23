CREATE TABLE orders
(
    order_id        VARCHAR(36)  NOT NULL COMMENT '订单ID（业务主键）',
    pre_order_id    VARCHAR(36)  NOT NULL COMMENT '预订单ID（关联结算服务）',
    user_id         INT UNSIGNED NOT NULL COMMENT '用户ID',

    -- 支付信息
    payment_method  TINYINT COMMENT '支付方式（1-微信 2-支付宝）',
    transaction_id  VARCHAR(64) COMMENT '支付平台流水号',
    paid_at         BIGINT COMMENT '支付成功时间戳（秒）',

    -- 金额信息（与结算服务保持一致）
    original_amount BIGINT       NOT NULL COMMENT '订单原始金额（分）',
    discount_amount BIGINT       NOT NULL DEFAULT 0 COMMENT '优惠总金额（分）',
    payable_amount  BIGINT       NOT NULL COMMENT '应付金额（分）',
    paid_amount     BIGINT                DEFAULT NULL COMMENT '实收金额（分）',


    -- 状态管理
    order_status    TINYINT      NOT NULL COMMENT '订单状态：1-待支付 2-已支付 3-已发货 4-已完成 5-已取消...',
    payment_status  TINYINT      NOT NULL COMMENT '支付状态：0-未支付 1-支付中 2-已支付 3-已退款...',

    reason          VARCHAR(255) COMMENT '取消原因',
    expire_time     BIGINT       NOT NULL COMMENT '过期时间戳',
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (order_id),
    UNIQUE KEY idx_pre_order (pre_order_id),
    KEY idx_user_status (user_id, order_status),
    KEY idx_payment_time (paid_at)
) COMMENT ='主订单表';

CREATE TABLE order_items
(
    order_id     VARCHAR(36)  NOT NULL COMMENT '预订单ID',
    product_id   INT UNSIGNED NOT NULL COMMENT '商品ID',
    quantity     INT UNSIGNED NOT NULL COMMENT '数量',
    price        BIGINT       NOT NULL COMMENT '当时单价（分）',
    product_name VARCHAR(255) NOT NULL COMMENT '商品名称',
    product_desc VARCHAR(255) NOT NULL COMMENT '商品描述',
    created_at   TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (order_id),
    INDEX idx_product (product_id)
) COMMENT ='订单商品快照';
# 用户下单地址快照
CREATE TABLE order_addresses
(
    order_id         VARCHAR(36)  NOT NULL COMMENT '订单ID',
    address_id       BIGINT UNSIGNED AUTO_INCREMENT,
    recipient_name   VARCHAR(100) NOT NULL COMMENT '收件人姓名',
    phone_number     VARCHAR(50)  DEFAULT NULL COMMENT '联系电话',
    province         VARCHAR(100) DEFAULT NULL COMMENT '州/省',
    city             VARCHAR(100) NOT NULL COMMENT '城市',
    detailed_address VARCHAR(255) NOT NULL COMMENT '详细地址',
    created_at       TIMESTAMP    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at       TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (address_id),
    UNIQUE KEY idx_order_address (order_id),
    INDEX idx_recipient_name (recipient_name)
)