-- 优惠券主表
CREATE TABLE `coupons`
(
    `id`              varchar(32)  NOT NULL COMMENT '优惠券ID',
    `name`            varchar(100) NOT NULL COMMENT '券名称',
    `type`            tinyint      NOT NULL COMMENT '类型：1-满减 2-折扣 3-立减',
    `value`           BIGINT       NOT NULL COMMENT '优惠值（根据类型：分/百分比）', -- 满减立减通用,折扣通用（例如：满200减50存50，8折存80）
    `min_amount`      BIGINT                DEFAULT 0 COMMENT '最低消费金额（分）', -- 满减立减通用,折扣通用
    `start_time`      timestamp    NOT NULL COMMENT '有效期开始',
    `end_time`        timestamp    NOT NULL COMMENT '有效期结束',
    `status`          tinyint      NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
    `total_count`     int          NOT NULL COMMENT '发行总量',
    `remaining_count` int          NOT NULL COMMENT '剩余数量',
    `created_at`      timestamp          DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      timestamp          DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_time` (`start_time`, `end_time`),
    KEY `idx_status` (`status`)
) COMMENT ='优惠券基本信息';
/*
switch coupon.Type {
    case 1: // 满减
        if orderAmount >= coupon.MinAmount {
            totalDiscount += coupon.Value
        }
    case 2: // 折扣（存储为90表示9折）
        discount := orderAmount * (100 - coupon.Value) / 100
        totalDiscount += discount
    case 3: // 立减
        totalDiscount += coupon.Value
}
 */
