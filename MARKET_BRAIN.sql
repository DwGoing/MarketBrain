/*
 Navicat Premium Data Transfer
 
 Source Server         : develop
 Source Server Type    : MySQL
 Source Server Version : 80033
 Source Host           : 10.0.0.30:30743
 Source Schema         : MARKET_BRAIN
 
 Target Server Type    : MySQL
 Target Server Version : 80033
 File Encoding         : 65001
 
 Date: 29/09/2023 02:43:44
 */
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- ----------------------------
-- Table structure for CONFIG
-- ----------------------------
DROP TABLE IF EXISTS `CONFIG`;
CREATE TABLE `CONFIG` (
  `KEY` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `VALUE` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`KEY`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
-- ----------------------------
-- Records of CONFIG
-- ----------------------------
BEGIN;
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES (
    'CHAIN_CONFIGS',
    '{\"TRON\":{\"usdt\":\"TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj\",\"rpcNodes\":[\"grpc.nile.trongrid.io:50051\"],\"httpNodes\":[\"https://nile.trongrid.io\"],\"apiKeys\":[\"d9b77ec9-39e0-4765-98d8-2c59188344a0\"]}}'
  );
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES ('EXPIRE_TIME', '15');
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES ('MIN_GAS_THRESHOLD', '20');
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES (
    'MNEMONIC',
    '\"math absorb sweet shrimp time smoke net pulp carbon gorilla expand payment\"'
  );
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES ('TRANSFER_GAS_AMOUNT', '50');
INSERT INTO `CONFIG` (`KEY`, `VALUE`)
VALUES ('WALLET_COLLECT_THRESHOLD', '50');
COMMIT;
-- ----------------------------
-- Table structure for RECHARGE_ORDER
-- ----------------------------
DROP TABLE IF EXISTS `RECHARGE_ORDER`;
CREATE TABLE `RECHARGE_ORDER` (
  `ID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CREATED_AT` datetime NOT NULL,
  `UPDATED_AT` datetime NOT NULL,
  `EXTERNAL_IDENTITY` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `EXTERNAL_DATA` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `CALLBACK_URL` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CHAIN_TYPE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AMOUNT` decimal(64, 18) NOT NULL,
  `WALLET_INDEX` bigint NOT NULL,
  `WALLET_ADDRESS` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `STATUS` varchar(255) CHARACTER SET utf8mb3 NOT NULL,
  `EXPIRE_AT` datetime NOT NULL,
  `TX_HASH` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`ID`) USING BTREE,
  KEY `INDEX_TX_HASH` (`TX_HASH`) USING BTREE,
  KEY `INDEX_EXTERNAL_IDENTITY` (`EXTERNAL_IDENTITY`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;
-- ----------------------------
-- Records of RECHARGE_ORDER
-- ----------------------------
BEGIN;
COMMIT;
-- ----------------------------
-- Table structure for TRANSFER
-- ----------------------------
DROP TABLE IF EXISTS `TRANSFER`;
CREATE TABLE `TRANSFER` (
  `ID` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `CREATED_AT` datetime NOT NULL,
  `UPDATED_AT` datetime NOT NULL,
  `TOKEN` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CHAIN_TYPE` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `FROM_INDEX` bigint NOT NULL,
  `FROM_ADDRESS` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `TO` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  `AMOUNT` decimal(64, 18) NOT NULL,
  `STATUS` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `ERROR` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `REMARKS` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`ID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;
-- ----------------------------
-- Records of TRANSFER
-- ----------------------------
BEGIN;
COMMIT;
SET FOREIGN_KEY_CHECKS = 1;