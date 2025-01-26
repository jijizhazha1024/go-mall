
CREATE TABLE users
(
    user_id       INT AUTO_INCREMENT COMMENT '主键，自增，用户 ID',
    username      VARCHAR(255) DEFAULT NULL COMMENT '用户名，可空',
    email         VARCHAR(255) UNIQUE COMMENT '邮箱，唯一',
    password_hash VARCHAR(512) COMMENT '密码哈希值',
    avatar_url    VARCHAR(255) DEFAULT NULL COMMENT '头像图片 URL',
    created_at    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    user_deleted  BOOLEAN DEFAULT FALSE COMMENT '用户是否已删除',
    updated_at    TIMESTAMP     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (user_id)
);