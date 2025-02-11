CREATE TABLE user_addresses (
    address_id INT AUTO_INCREMENT COMMENT '主键，自增，地址ID',
    user_id INT NOT NULL COMMENT '外键，关联到users表的user_id',
    detailed_address VARCHAR(255) NOT NULL COMMENT '详细地址',
    city VARCHAR(100) NOT NULL COMMENT '城市',
    province VARCHAR(100) DEFAULT NULL COMMENT '州/省',
    is_default BOOLEAN DEFAULT false COMMENT '是否默认地址',
    recipient_name VARCHAR(100) NOT NULL COMMENT '收件人姓名',
    phone_number VARCHAR(50) DEFAULT NULL COMMENT '联系电话',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (address_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);