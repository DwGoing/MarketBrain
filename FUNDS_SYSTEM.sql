/*
 Navicat Premium Data Transfer

 Source Server         : develop
 Source Server Type    : MySQL
 Source Server Version : 80033
 Source Host           : 10.0.0.30:30743
 Source Schema         : FUNDS_SYSTEM

 Target Server Type    : MySQL
 Target Server Version : 80033
 File Encoding         : 65001

 Date: 01/08/2023 15:19:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for RECHARGE_RECORD
-- ----------------------------
DROP TABLE IF EXISTS `RECHARGE_RECORD`;
CREATE TABLE `RECHARGE_RECORD` (
  `ID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CREATED_AT` datetime NOT NULL,
  `UPDATED_AT` datetime NOT NULL,
  `EXTERNAL_IDENTITY` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `EXTERNAL_DATA` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `CALLBACK_URL` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AMOUNT` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `TOKEN` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `WALLET_INDEX` bigint NOT NULL,
  `WALLET_ADDRESS` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `BDFORE_BALANCE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AFTER_BALANCE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `STATUS` tinyint NOT NULL,
  `EXPIRE_AT` datetime NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for TRANSFER_RECORD
-- ----------------------------
DROP TABLE IF EXISTS `TRANSFER_RECORD`;
CREATE TABLE `TRANSFER_RECORD` (
  `ID` varchar(255) NOT NULL,
  `CREATED_AT` datetime NOT NULL,
  `UPDATED_AT` datetime NOT NULL,
  `FROM_INDEX` bigint NOT NULL,
  `FROM_ADDRESS` varchar(255) NOT NULL,
  `TO` varchar(255) NOT NULL,
  `TOKEN` varchar(255) NOT NULL,
  `AMOUNT` varchar(255) NOT NULL,
  `STATUS` tinyint NOT NULL,
  `ERROR` text,
  `REMARKS` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;

SET FOREIGN_KEY_CHECKS = 1;
