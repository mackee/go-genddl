-- generated by github.com/mackee/go-genddl. DO NOT EDIT!!!

CREATE SEQUENCE IF NOT EXISTS seq_location_id START WITH 1 INCREMENT BY 1;
DROP TABLE IF EXISTS "location";

CREATE TABLE "location" (
    "id" UBIGINT NOT NULL PRIMARY KEY DEFAULT nextval('seq_location_id'),
    "name" VARCHAR NOT NULL,
    "description" VARCHAR NOT NULL,
    "place" BLOB NOT NULL
) ;


CREATE SEQUENCE IF NOT EXISTS seq_product_id START WITH 1 INCREMENT BY 1;
DROP TABLE IF EXISTS "product";

CREATE TABLE "product" (
    "id" UINTEGER NOT NULL PRIMARY KEY DEFAULT nextval('seq_product_id'),
    "name" VARCHAR NOT NULL,
    "type" UINTEGER NOT NULL,
    "user_id" UINTEGER NOT NULL,
    "received_user_id" UINTEGER NULL,
    "description" VARCHAR NOT NULL,
    "full_description" VARCHAR NOT NULL,
    "size" SMALLINT NULL,
    "status" UTINYINT NOT NULL,
    "category" TINYINT NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NULL,
    "removed_at" DATETIME NULL,
    UNIQUE ("user_id", "type")
) ;
CREATE INDEX product_user_id_created_at ON "product" ("user_id", "created_at");


CREATE SEQUENCE IF NOT EXISTS seq_user_id START WITH 1 INCREMENT BY 1;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    "id" UINTEGER NOT NULL PRIMARY KEY DEFAULT nextval('seq_user_id'),
    "name" VARCHAR NOT NULL UNIQUE,
    "age" BIGINT NULL,
    "message" VARCHAR NULL,
    "icon_image" BLOB NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NULL
) ;

CREATE VIEW "user_product"
  ("u_id", "u_name", "ru_name", "p_id", "p_type") AS
  SELECT u.id, u.name, ru.name, p.id, p.type FROM product AS p
    INNER JOIN user AS u ON p.user_id = u.id
    LEFT JOIN user AS ru ON p.received_user_id = ru.id;

CREATE VIEW "user_product_structured"
  ("p_id", "p_name", "p_type", "p_user_id", "p_received_user_id", "p_description", "p_full_description", "p_size", "p_status", "p_category", "p_created_at", "p_updated_at", "p_removed_at", "u_id", "u_name", "u_age", "u_message", "u_icon_image", "u_created_at", "u_updated_at", "ru_id", "ru_name", "ru_age", "ru_message", "ru_icon_image", "ru_created_at", "ru_updated_at") AS
  SELECT p.*, u.*, ru.* FROM product AS p
    INNER JOIN user AS u ON p.user_id = u.id
    LEFT JOIN user AS ru ON p.received_user_id = ru.id;

