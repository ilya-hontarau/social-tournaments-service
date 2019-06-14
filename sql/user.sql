CREATE DATABASE social_tournament;
USE social_tournament;

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL,
    balance INT(10) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);
