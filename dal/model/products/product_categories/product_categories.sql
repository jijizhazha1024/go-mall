CREATE TABLE product_categories (
    id INT AUTO_INCREMENT COMMENT '自增主键',
    product_id INT COMMENT '商品id',
    category_id INT COMMENT '分类id',
    PRIMARY KEY (id),
    UNIQUE KEY uk_product_category (product_id, category_id) COMMENT '商品与分类的唯一约束'
);