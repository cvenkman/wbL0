CREATE DATABASE cvenkman; /* or change database name in configs/config.toml */

DROP TABLE IF EXISTS delivery;

CREATE TABLE delivery (
    id          varchar(128) PRIMARY KEY,
    content     text NOT NULL
);