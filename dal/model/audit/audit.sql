create table `audit`
(
    id           int UNSIGNED auto_increment comment '主键',
    user_id      int UNSIGNED not null comment '用户id',
    action_type  varchar(64)  not null comment '操作类型',
    action_desc  text comment '操作描述',
    old_data     json comment '旧数据',
    new_data     json comment '新数据',
    service_name varchar(64)  not null comment '服务名称',
    target_table varchar(64)  not null comment '目标表',
    target_id    int UNSIGNED not null comment '目标id',
    client_ip    varchar(45)  not null comment 'ip地址',
    trace_id     varchar(36)  not null comment 'traceid', -- 用于关联跟踪 （但是可能不到64字长）
    span_id      varchar(36)  not null comment 'spanid',  -- 用于关联跟踪
    created_at   TIMESTAMP default CURRENT_TIMESTAMP COMMENT '创建时间',

    primary key (id),
    UNIQUE idx_trace (trace_id),
    INDEX idx_user (user_id),
    INDEX idx_service (service_name),
    INDEX idx_action (action_type),
    INDEX idx_target (target_table, target_id),
    INDEX idx_time (created_at)
);
