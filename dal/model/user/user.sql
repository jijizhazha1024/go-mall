CREATE TABLE users
(
    user_id       INT AUTO_INCREMENT COMMENT '主键，自增，用户 ID',
    username      VARCHAR(255) UNIQUE COMMENT '用户名，唯一',
    email         VARCHAR(255) UNIQUE COMMENT '邮箱，唯一',
    password_hash VARCHAR(255) COMMENT '密码哈希值',
    phone_number  VARCHAR(20) UNIQUE COMMENT '手机号，唯一',
    avatar_url    VARCHAR(255) COMMENT '头像图片 URL',
    created_at    DATETIME COMMENT '用户注册时间',
    updated_at    DATETIME COMMENT '用户信息最后更新时间',
    last_login_at DATETIME COMMENT '上次登录时间',
    status        TINYINT COMMENT '用户状态，0 - 正常，1 - 封禁',
    primary key (user_id)
);
