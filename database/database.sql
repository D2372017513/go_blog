CREATE DATABASE IF NOT EXISTS goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE goblog;

CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
)