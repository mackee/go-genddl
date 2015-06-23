DROP TABLE IF EXISTS project;

CREATE TABLE project (
    `user` INTEGER unsigned NOT NULL,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `user_id` INTEGER unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`user`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
DROP TABLE IF EXISTS user;

CREATE TABLE user (
    `user` INTEGER unsigned NOT NULL,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (`user`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
