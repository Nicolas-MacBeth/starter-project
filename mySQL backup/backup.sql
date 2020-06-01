-- MySQL dump 10.17  Distrib 10.3.22-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: starterproject
-- ------------------------------------------------------
-- Server version	10.3.22-MariaDB-1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ingredients`
--

DROP TABLE IF EXISTS `ingredients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ingredients` (
  `ingredient_name` varchar(50) NOT NULL,
  PRIMARY KEY (`ingredient_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ingredients`
--

LOCK TABLES `ingredients` WRITE;
/*!40000 ALTER TABLE `ingredients` DISABLE KEYS */;
INSERT INTO `ingredients` VALUES ('avocado'),('banana'),('carrot'),('chili'),('corn'),('pepper');
/*!40000 ALTER TABLE `ingredients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vendor`
--

DROP TABLE IF EXISTS `vendor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vendor` (
  `id` int(11) NOT NULL,
  `name` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vendor`
--

LOCK TABLES `vendor` WRITE;
/*!40000 ALTER TABLE `vendor` DISABLE KEYS */;
INSERT INTO `vendor` VALUES (1,'MacBeth delicious foods'),(2,'MacDonald official vendor'),(3,'Better foods'),(4,'Google micro-kitchen exclusive vendor'),(5,'Super ingredients, Inc.');
/*!40000 ALTER TABLE `vendor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vendor_ingredient`
--

DROP TABLE IF EXISTS `vendor_ingredient`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vendor_ingredient` (
  `vendor_id` int(11) NOT NULL,
  `ingredient_name` varchar(50) NOT NULL,
  `price` int(11) DEFAULT NULL,
  `inventory` int(11) DEFAULT NULL,
  `vendor_name` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`vendor_id`,`ingredient_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vendor_ingredient`
--

LOCK TABLES `vendor_ingredient` WRITE;
/*!40000 ALTER TABLE `vendor_ingredient` DISABLE KEYS */;
INSERT INTO `vendor_ingredient` VALUES (1,'avocado',60,37,'MacBeth delicious foods'),(1,'banana',14,40,'MacBeth delicious foods'),(1,'carrot',10,56,'MacBeth delicious foods'),(1,'chili',35,24,'MacBeth delicious foods'),(1,'corn',5,120,'MacBeth delicious foods'),(1,'pepper',16,75,'MacBeth delicious foods'),(2,'avocado',100,5,'MacDonald official vendor'),(2,'chili',29,99,'MacDonald official vendor'),(2,'pepper',20,200,'MacDonald official vendor'),(3,'avocado',55,34,'Better foods'),(3,'carrot',12,200,'Better foods'),(3,'pepper',24,140,'Better foods'),(4,'banana',14,156,'Google micro-kitchen exclusive vendor'),(4,'chili',22,500,'Google micro-kitchen exclusive vendor'),(4,'corn',7,450,'Google micro-kitchen exclusive vendor'),(5,'avocado',200,2,'Super ingredients, Inc.'),(5,'banana',13,25,'Super ingredients, Inc.'),(5,'corn',14,36,'Super ingredients, Inc.');
/*!40000 ALTER TABLE `vendor_ingredient` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-06-01 19:48:21
