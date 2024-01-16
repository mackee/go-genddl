-- generated by github.com/mackee/go-genddl. DO NOT EDIT!!!

DROP TABLE IF EXISTS "location";

CREATE TABLE "location" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "place" GEOMETRY NOT NULL,
    SPATIAL KEY place ("place"),
    FULLTEXT KEY description ("description")
) ;


DROP TABLE IF EXISTS "product";

CREATE TABLE "product" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL,
    "type" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "received_user_id" INTEGER NULL,
    "description" TEXT NOT NULL,
    "full_description" TEXT NOT NULL,
    "size" INTEGER NULL,
    "status" INTEGER NOT NULL,
    "category" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NULL,
    UNIQUE ("user_id", "type"),
    FOREIGN KEY ("user_id") REFERENCES user("id") ON DELETE CASCADE ON UPDATE CASCADE
) ;
CREATE INDEX product_user_id_created_at ON product ("user_id", "created_at");


DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL UNIQUE,
    "age" INTEGER NULL,
    "message" TEXT NULL,
    "icon_image" BLOB NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NULL
) ;

CREATE VIEW "user_product"
  ("p_id", "u_name", "ru_name", "p_id", "p_type") AS 
  SELECT p.id, u.name, ru.name, p.id, p.type FROM product AS p
    INNER JOIN user AS u ON p.user_id = u.id
    LEFT JOIN user AS ru ON p.received_user_id = ru.id;

CREATE VIEW "user_product_structured"
  ("p_id", "p_name", "p_type", "p_user_id", "p_received_user_id", "p_description", "p_full_description", "p_size", "p_status", "p_category", "p_created_at", "p_updated_at", "u_id", "u_name", "u_age", "u_message", "u_icon_image", "u_created_at", "u_updated_at", "ru_id", "ru_name", "ru_age", "ru_message", "ru_icon_image", "ru_created_at", "ru_updated_at") AS
  SELECT p.*, u.*, ru.* FROM product AS p
    INNER JOIN user AS u ON p.user_id = u.id
    LEFT JOIN user AS ru ON p.received_user_id = ru.id;

