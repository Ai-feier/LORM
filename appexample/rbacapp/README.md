# Lorm 示例项目 kube-rbacauth
定位: 基于 rbac 授权模式的 k8s 纯后端认证授权平台

## 目的
- 测试 lorm 框架性能
- 基于第三方管理 k8s 认证授权信息

## 系统架构
![](../../assets/img/system%20structure.png)

## 数据库设计
![](../../assets/img/model%20info.png)

### mysql DDL
```mysql
CREATE TABLE IF NOT EXISTS cluster_role_bindings
(
    id      BIGINT AUTO_INCREMENT PRIMARY KEY,
    name    LONGTEXT NOT NULL,
    users   LONGTEXT NOT NULL,
    role_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_role_sub_refs
(
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    cluster_role_id BIGINT NOT NULL,
    verbs           LONGTEXT NOT NULL,
    resources       LONGTEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_roles
(
    id   BIGINT AUTO_INCREMENT PRIMARY KEY,
    name LONGTEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_bindings
(
    id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    name      LONGTEXT NOT NULL,
    namespace LONGTEXT NOT NULL,
    users     LONGTEXT NOT NULL,
    role_id   BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_sub_refs
(
    id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    role_id   BIGINT NOT NULL,
    verbs     LONGTEXT NOT NULL,
    resources LONGTEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS roles
(
    id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    name      LONGTEXT NOT NULL,
    namespace LONGTEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name   LONGTEXT NOT NULL,
    password    LONGTEXT NOT NULL,
    create_time BIGINT NOT NULL,
    update_time BIGINT NOT NULL
);

# create table if not exists cluster_role_bindings
# (
#     id      bigint auto_increment
#         primary key,
#     name    longtext not null,
#     users   longtext not null,
#     role_id bigint   not null
# );
# 
# create table if not exists cluster_role_sub_refs
# (
#     id              bigint auto_increment
#         primary key,
#     cluster_role_id bigint   not null,
#     verbs           longtext not null,
#     resources       longtext not null
# );
# 
# create table if not exists cluster_roles
# (
#     id   bigint auto_increment
#         primary key,
#     name longtext not null
# );
# 
# create table if not exists role_bindings
# (
#     id        bigint auto_increment
#         primary key,
#     name      longtext not null,
#     namespace longtext not null,
#     users     longtext not null,
#     role_id   bigint   not null
# );
# 
# create table if not exists role_sub_refs
# (
#     id        bigint auto_increment
#         primary key,
#     role_id   bigint   not null,
#     verbs     longtext not null,
#     resources longtext not null
# );
# 
# create table if not exists roles
# (
#     id        bigint auto_increment
#         primary key,
#     name      longtext not null,
#     namespace longtext not null
# );
# 
# create table if not exists users
# (
#     id          bigint auto_increment
#         primary key,
#     user_name   longtext not null,
#     password    longtext not null,
#     create_time bigint   not null,
#     update_time bigint   not null
# );
# 

```

### postgre
1. PostgreSQL不支持auto_increment，而是使用SERIAL关键字来实现自增列。
2. PostgreSQL使用VARCHAR代替MySQL中的LONGTEXT。
3. PostgreSQL使用SERIAL关键字来创建自增列。
4. 将bigint替换为BIGINT。
5. 将longtext替换为VARCHAR。
```postgresql
CREATE TABLE IF NOT EXISTS cluster_role_bindings
(
    id      BIGINT SERIAL PRIMARY KEY,
    name    VARCHAR NOT NULL,
    users   VARCHAR NOT NULL,
    role_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_role_sub_refs
(
    id              BIGINT SERIAL PRIMARY KEY,
    cluster_role_id BIGINT NOT NULL,
    verbs           VARCHAR NOT NULL,
    resources       VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_roles
(
    id   BIGINT SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS role_bindings
(
    id        BIGINT SERIAL PRIMARY KEY,
    name      VARCHAR NOT NULL,
    namespace VARCHAR NOT NULL,
    users     VARCHAR NOT NULL,
    role_id   BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_sub_refs
(
    id        BIGINT SERIAL PRIMARY KEY,
    role_id   BIGINT NOT NULL,
    verbs     VARCHAR NOT NULL,
    resources VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS roles
(
    id        BIGINT SERIAL PRIMARY KEY,
    name      VARCHAR NOT NULL,
    namespace VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id          BIGINT SERIAL PRIMARY KEY,
    user_name   VARCHAR NOT NULL,
    password    VARCHAR NOT NULL,
    create_time BIGINT NOT NULL,
    update_time BIGINT NOT NULL
);

```

### sqlite3
1. SQLite支持AUTOINCREMENT来实现自增列。
2. SQLite使用TEXT代替MySQL中的LONGTEXT。
3. SQLite中不需要指定列的大小，因此可以省略VARCHAR后的大小限制。
```sqlite
CREATE TABLE IF NOT EXISTS cluster_role_bindings
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT NOT NULL,
    users   TEXT NOT NULL,
    role_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_role_sub_refs
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_role_id INTEGER NOT NULL,
    verbs           TEXT NOT NULL,
    resources       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster_roles
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_bindings
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,
    namespace TEXT NOT NULL,
    users     TEXT NOT NULL,
    role_id   INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS role_sub_refs
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id   INTEGER NOT NULL,
    verbs     TEXT NOT NULL,
    resources TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS roles
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,
    namespace TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name   TEXT NOT NULL,
    password    TEXT NOT NULL,
    create_time INTEGER NOT NULL,
    update_time INTEGER NOT NULL
);
```
