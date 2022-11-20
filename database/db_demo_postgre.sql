-- 调整后的sql语句
create table tb_auth_casbin_rule
(
    id         bigserial
        primary key,
    ptype      varchar(100),
    v0         varchar(100),
    v1         varchar(100),
    v2         varchar(100),
    v3         varchar(100),
    v4         varchar(100),
    v5         varchar(100),
    v6         varchar(25),
    v7         varchar(25),
    created_at text,
    updated_at text,
    deleted_at timestamp with time zone
);

alter table tb_auth_casbin_rule
    owner to postgres;

create unique index idx_tb_auth_casbin_rule
    on tb_auth_casbin_rule (ptype, v0, v1, v2, v3, v4, v5, v6, v7);

create table tb_auth_users
(
    id            bigserial
        primary key,
    created_at    text,
    updated_at    text,
    deleted_at    timestamp with time zone,
    user_name     text,
    pass          text,
    phone         text,
    real_name     text,
    status        bigint,
    token         text,
    last_login_ip text,
    avatar        text
);

alter table tb_auth_users
    owner to postgres;

create table tb_auth_menu
(
    id                    bigserial
        primary key,
    created_at            text,
    updated_at            text,
    deleted_at            timestamp with time zone,
    name                  text,
    parent_name           text,
    parent_id             bigint,
    order_no              integer,
    path                  text,
    icon                  text,
    hide_menu             boolean,
    component             text,
    is_iframe             boolean,
    frame_src             text,
    is_cache              boolean,
    menu_type             text,
    title                 text,
    redirect              text,
    hide_children_in_menu boolean,
    current_active_menu   text,
    hide_breadcrumb       boolean
);

alter table tb_auth_menu
    owner to postgres;

create table tb_auth_role
(
    id         bigserial
        primary key,
    created_at text,
    updated_at text,
    deleted_at timestamp with time zone,
    role_key   text,
    role_name  text
);

alter table tb_auth_role
    owner to postgres;

create table tb_auth_role_menus
(
    role_model_id bigint not null
        constraint fk_tb_auth_role_menus_role_model
            references tb_auth_role,
    menu_model_id bigint not null
        constraint fk_tb_auth_role_menus_menu_model
            references tb_auth_menu,
    primary key (role_model_id, menu_model_id)
);

alter table tb_auth_role_menus
    owner to postgres;