CREATE TABLE inventory
(
    product_id INT, -- 与商品服务共享同一ID
    total      INT NOT NULL,
    sold       INT NOT NULL,
    primary key (product_id)
);