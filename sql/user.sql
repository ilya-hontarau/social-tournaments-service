CREATE DATABASE tournament_db;
USE tournament_db;
CREATE TABLE user (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL,
    balance INT(10) UNSIGNED NOT NULL,
    PRIMARY KEY (ID)
);
