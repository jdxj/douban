/*
 Navicat Premium Data Transfer

 Source Server         : mysql.aaronkir.xyz
 Source Server Type    : MySQL
 Source Server Version : 80017
 Source Host           : mysql.aaronkir.xyz:3306
 Source Schema         : douban

 Target Server Type    : MySQL
 Target Server Version : 80017
 File Encoding         : 65001

 Date: 23/11/2019 17:37:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for book
-- ----------------------------
DROP TABLE IF EXISTS `book`;
CREATE TABLE `book` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL COMMENT '书名',
  `author` varchar(100) DEFAULT NULL COMMENT '作者',
  `press` varchar(100) DEFAULT NULL COMMENT '出版社',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=127 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for book_url
-- ----------------------------
DROP TABLE IF EXISTS `book_url`;
CREATE TABLE `book_url` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=62341 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for flag
-- ----------------------------
DROP TABLE IF EXISTS `flag`;
CREATE TABLE `flag` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tag` int(11) NOT NULL COMMENT '标签',
  `type` int(11) NOT NULL COMMENT '书影音',
  `ref` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `flag_FK` (`tag`),
  KEY `flag_FK_1` (`type`),
  CONSTRAINT `flag_FK` FOREIGN KEY (`tag`) REFERENCES `tag` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `flag_FK_1` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='记录书影音都有哪些标签';

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
  `id` int(15) NOT NULL AUTO_INCREMENT,
  `log` int(15) NOT NULL DEFAULT '0' COMMENT '抓到哪行的行号',
  `type` int(15) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `type_id` (`type`),
  CONSTRAINT `type_id` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for opinion
-- ----------------------------
DROP TABLE IF EXISTS `opinion`;
CREATE TABLE `opinion` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `score` double NOT NULL DEFAULT '0' COMMENT '评分',
  `amount` int(11) NOT NULL DEFAULT '0' COMMENT '评分人数',
  `one` double NOT NULL DEFAULT '0' COMMENT '各星级百分比, 单位: %',
  `two` double NOT NULL DEFAULT '0',
  `three` double NOT NULL DEFAULT '0',
  `four` double NOT NULL DEFAULT '0',
  `five` double NOT NULL DEFAULT '0',
  `type` int(11) NOT NULL COMMENT '对什么评论, 书影音',
  `ref` int(11) NOT NULL COMMENT '书影音 ID',
  PRIMARY KEY (`id`),
  KEY `opinion_FK` (`type`),
  CONSTRAINT `opinion_FK` FOREIGN KEY (`type`) REFERENCES `type` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=121 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评价表';

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '标签名',
  `url` varchar(100) DEFAULT NULL COMMENT '点击该 url 能够获取该标签下的书',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='记录书影音等所有标签';

-- ----------------------------
-- Table structure for type
-- ----------------------------
DROP TABLE IF EXISTS `type`;
CREATE TABLE `type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '类型名',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='豆瓣里的种类, 书影音, 可能还有其他';

SET FOREIGN_KEY_CHECKS = 1;
