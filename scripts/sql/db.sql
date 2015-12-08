CREATE DATABASE beego_demo;
USE beego_demo;
CREATE TABLE roles (
  id		BIGINT PRIMARY KEY,
  name		VARCHAR(255),
  password	VARCHAR(255),
  reg_date	DATETIME
);
INSERT INTO roles VALUES(1, 'admin', 'admin', now());
