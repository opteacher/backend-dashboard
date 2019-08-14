# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.15)
# Database: backend
# Generation Time: 2019-08-14 09:22:34 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table dao_categories
# ------------------------------------------------------------

LOCK TABLES `dao_categories` WRITE;
/*!40000 ALTER TABLE `dao_categories` DISABLE KEYS */;

INSERT INTO `dao_categories` (`id`, `name`, `desc`, `lang`)
VALUES
	(1,'database_notx','数据库（无事务）','golang'),
	(2,'database_tx','数据库（带事务）','golang'),
    (3,'cache','缓存','golang');

/*!40000 ALTER TABLE `dao_categories` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `dao_interfaces` WRITE;
/*!40000 ALTER TABLE `dao_interfaces` DISABLE KEYS */;

INSERT INTO `dao_interfaces` (`id`, `name`, `category`, `desc`)
VALUES
	(1,'Close','database_tx','关闭数据库'),
	(2,'Ping','database_tx','测试数据库通讯状态'),
    (3,'BeginTx','database_tx','开启事务');

/*!40000 ALTER TABLE `dao_interfaces` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `dao_interface_params_mapper` WRITE;
/*!40000 ALTER TABLE `dao_interface_params_mapper` DISABLE KEYS */;

INSERT INTO `dao_interface_params_mapper` (`id`, `name`, `type`, `dao_interface_id`)
VALUES
	(1,'ctx','context.Context',2),
	(2,'ctx','context.Context',3);

/*!40000 ALTER TABLE `dao_interface_params_mapper` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `dao_interface_requires_mapper` WRITE;
/*!40000 ALTER TABLE `dao_interface_requires_mapper` DISABLE KEYS */;

INSERT INTO `dao_interface_requires_mapper` (`id`, `require`, `dao_interface_id`)
VALUES
	(1,'context',2),
	(1,'context',2);

/*!40000 ALTER TABLE `dao_interface_requires_mapper` ENABLE KEYS */;
UNLOCK TABLES;

LOCK TABLES `dao_interface_returns_mapper` WRITE;
/*!40000 ALTER TABLE `dao_interface_returns_mapper` DISABLE KEYS */;

INSERT INTO `dao_interface_returns_mapper` (`id`, `return`, `dao_interface_id`)
VALUES
	(1,'error',2);

/*!40000 ALTER TABLE `dao_interface_returns_mapper` ENABLE KEYS */;
UNLOCK TABLES;


/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
