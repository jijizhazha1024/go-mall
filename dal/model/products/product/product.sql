CREATE TABLE products (
    id INT AUTO_INCREMENT  COMMENT '主键，自增,商品id',
    name VARCHAR(255) NOT NULL  COMMENT '商品名称',
    description TEXT COMMENT '商品描述',
    picture VARCHAR(255)  COMMENT '商品图片信息',
    price DECIMAL(10, 2) NOT NULL COMMENT '商品价格',
    stock INT NOT NULL COMMENT '库存数量',
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
);
