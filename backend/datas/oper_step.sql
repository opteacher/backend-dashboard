# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.15)
# Database: backend
# Generation Time: 2019-08-14 09:20:59 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table oper_step
# ------------------------------------------------------------

LOCK TABLES `oper_step` WRITE;
/*!40000 ALTER TABLE `oper_step` DISABLE KEYS */;

INSERT INTO `oper_step` (`id`, `oper_key`, `requires`, `desc`, `inputs`, `outputs`, `code`, `api_name`, `symbol`)
VALUES
	(1,'json_marshal','encoding/json','将收到的请求参数编码成JSON字节数组','OBJECT:','bytes','bytes, err := json.Marshal(%OBJECT%)\nif err != nil {\n	return nil, fmt.Errorf(\"转JSON失败：%v\", err)\n}\n',NULL,NULL),
	(2,'json_unmarshal','%PACKAGE%/internal/utils','将JSON字节数组转成Map键值对','OBJ_TYPE:','omap','omap, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*%OBJ_TYPE%)(nil)).Elem())\nif err != nil {\n	return nil, fmt.Errorf(\"从JSON转回失败：%v\", err)\n}\n',NULL,NULL),
	(3,'database_beginTx',NULL,'开启数据库事务',NULL,'tx','tx, err := s.dao.BeginTx(ctx)\nif err != nil {\n	return nil, fmt.Errorf(\"开启事务失败：%v\", err)\n}\n',NULL,NULL),
	(4,'database_commitTx',NULL,'提交数据库事务',NULL,NULL,'err := s.dao.CommitTx(tx)\nif err != nil {\n	return nil, fmt.Errorf(\"提交事务失败：%v\", err)\n}\n',NULL,NULL),
	(5,'assignment',NULL,'将%SOURCE%赋值给%TARGET%','SOURCE:,TARGET:',NULL,'%TARGET% = %SOURCE%\n',NULL,NULL),
	(6,'assignment_append',NULL,'将%NEW_ADD%添加进%ARRAY%','ARRAY:,NEW_ADD:',NULL,'%ARRAY% = append(%ARRAY%, %NEW_ADD%)\n',NULL,NULL),
	(7,'assignment_create',NULL,'创建%TARGET%并用%SOURCE%初始化','SOURCE:,TARGET:',NULL,'%TARGET% := %SOURCE%\n',NULL,NULL),
	(8,'for_each',NULL,'循环遍历%SET%','KEY:,VALUE:,SET:',NULL,'for %KEY%, %VALUE% := range %SET%',NULL,1),
	(9,'return_succeed',NULL,'成功返回%RETURN%','RETURN:',NULL,'return %RETURN%, nil\n',NULL,4),
	(10,'database_insertTx',NULL,'做数据库插入操作','TABLE_NAME:,OBJ_MAP:','id','id, err := s.dao.InsertTx(tx, \"%TABLE_NAME%\", %OBJ_MAP%)\nif err != nil {\n	return nil, fmt.Errorf(\"插入数据表失败：%v\", err)\n}\n',NULL,NULL),
	(11,'database_queryTx',NULL,'做数据库查询操作（事务）','TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:','res','res, err := s.dao.QueryTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n	return nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n',NULL,NULL),
	(12,'database_query',NULL,'做数据库查询操作（会话）','TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:','res','res, err := s.dao.Query(ctx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n	return nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n',NULL,NULL),
	(13,'database_deleteTx',NULL,'做数据库删除操作','TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:',NULL,'_, err := s.dao.DeleteTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n	return nil, fmt.Errorf(\"删除数据表记录失败：%v\", err)\n}\n',NULL,NULL),
	(14,'database_updateTx',NULL,'做数据库更新操作','TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:,OBJ_MAP:','id','id, err := s.dao.SaveTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%, %OBJ_MAP%)\nif err != nil {\n	return nil, fmt.Errorf(\"更新数据表记录失败：%v\", err)\n}\n',NULL,NULL);

/*!40000 ALTER TABLE `oper_step` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
