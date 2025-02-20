

# 用户下单地址快照
CREATE TABLE order_addresses
(
    address_id       BIGINT UNSIGNED AUTO_INCREMENT,
    recipient_name   VARCHAR(100) NOT NULL COMMENT '收件人姓名',
    phone_number     VARCHAR(50)  DEFAULT NULL COMMENT '联系电话',
    province         VARCHAR(100) DEFAULT NULL COMMENT '州/省',
    city             VARCHAR(100) NOT NULL COMMENT '城市',
    detailed_address VARCHAR(255) NOT NULL COMMENT '详细地址',
    created_at       TIMESTAMP    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at       TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (address_id),
    INDEX idx_recipient_name (recipient_name)
)