-- CREATE DATABASE IF NOT EXISTS "db-rest-demo";

DROP TABLE IF EXISTS `books`;

USE db-rest-demo;

CREATE TABLE books(
    isbn TEXT(13) UNIQUE NOT NULL,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    price FLOAT NOT NULL,
    page INT NOT NULL,
    PRIMARY KEY(`isbn`)
);
CREATE TABLE authors(
    fname VARCHAR(50) NOT NULL,  
    lname VARCHAR(50) NOT NULL,
    dob DATETIME,
    qualification VARCHAR(100)
);
