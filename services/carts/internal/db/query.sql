-- name: GetCartByUser :many
SELECT id, user_id, product_id, quantity, checked
FROM carts
WHERE user_id = $1;

-- name: UpdateCartItem :exec
UPDATE carts
SET checked = $1, quantity = $2
WHERE product_id = $3 AND user_id = $4;

-- name: DeleteCartItem :exec
DELETE FROM carts
WHERE product_id = $1 AND user_id = $2;

-- name: CreateOrUpdateCartItem :exec
-- 如果购物车中已存在该商品，则更新数量；否则插入新记录
INSERT INTO carts (user_id, product_id, quantity, checked)
VALUES ($1, $2, $3, false)
    ON DUPLICATE KEY UPDATE quantity = quantity + $3;