
CREATE TABLE order_items
(
    item_id      BIGINT UNSIGNED AUTO_INCREMENT ,
    order_id     VARCHAR(64)  NOT NULL COMMENT '关联订单号',

    -- 商品快照
    product_id   INT          NOT NULL COMMENT '商品ID',
    quantity     INT          NOT NULL COMMENT '购买数量',
    product_name VARCHAR(255) NOT NULL COMMENT '商品名称',
    product_desc TEXT COMMENT '规格描述',
    unit_price   BIGINT       NOT NULL COMMENT '单价(分)',

#     FOREIGN KEY (order_id) REFERENCES orders (order_id),
    PRIMARY KEY (item_id),
    INDEX idx_order_product (order_id, product_id)
) COMMENT ='订单商品快照';