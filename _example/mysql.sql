DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` INTEGER unsigned NOT NULL PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
