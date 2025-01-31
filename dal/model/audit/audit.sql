create table `audit`
(
    id                 int auto_increment comment '主键',
    user_id            int          not null comment '用户id',
    username           varchar(255) not null comment '用户名',
    action_type        varchar(50)  not null comment '操作类型',
    action_description text comment '操作描述',
    target_table       varchar(255) not null comment '目标表',
    old_data           json comment '旧数据',
    new_data           json comment '新数据',
    target_id          int          not null comment '目标id',
    ip_address         varchar(45)  not null comment 'ip地址',
    trace_id           varchar(255) comment 'traceid', -- 用于关联跟踪 （但是可能不到255字长）
    span_id            varchar(255) comment 'spanid',  -- 用于关联跟踪
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    primary key (id)
);
