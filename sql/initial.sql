DROP DATABASE IF EXISTS chatZoo;

CREATE DATABASE chatZoo;

USE chatZoo;


SET NAMES utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;


-- ----------------------------

-- Table structure for User

-- ----------------------------

DROP TABLE IF EXISTS `User`;

CREATE TABLE `User`  (

    `ID` char(16) NOT NULL,
	`Pwd` char(16) NOT NULL,

    `Data` mediumblob NULL,
	
    PRIMARY KEY (`ID`) USING BTREE

) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;