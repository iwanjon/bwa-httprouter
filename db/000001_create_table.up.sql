# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.11)
# Database: bwastartup
# Generation Time: 2020-09-23 13:46:56 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table campaign_images
# ------------------------------------------------------------

DROP TABLE IF EXISTS campaign_images;

CREATE TABLE campaign_images (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  campaign_id int(11) DEFAULT NULL,
  file_name varchar(255) DEFAULT NULL,
  is_primary tinyint(4) DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table campaigns
# ------------------------------------------------------------

DROP TABLE IF EXISTS campaigns;

CREATE TABLE campaigns (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  user_id int(11) DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  short_description varchar(255) DEFAULT NULL,
  description text,
  perks text,
  backer_count int(11) DEFAULT NULL,
  goal_amount int(11) DEFAULT NULL,
  current_amount int(11) DEFAULT NULL,
  slug varchar(255) DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table transactions
# ------------------------------------------------------------

DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  campaign_id int(11) DEFAULT NULL,
  user_id int(11) DEFAULT NULL,
  amount int(11) DEFAULT NULL,
  status varchar(255) DEFAULT NULL,
  code varchar(255) DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) DEFAULT NULL,
  occupation varchar(255) DEFAULT NULL,
  email varchar(255) DEFAULT NULL,
  password_hash varchar(255) DEFAULT NULL,
  avatar_file_name varchar(255) DEFAULT NULL,
  role varchar(255) DEFAULT NULL,
  token varchar(255) DEFAULT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


INSERT INTO bwastartup.users
(name, occupation, email, password_hash, avatar_file_name, `role`, token, created_at, updated_at)
VALUES('joko1', 'jokojob1', 'joko1@gmail.com', 'joko1hash', 'joko1.jpg', 'user', 'joko1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('joko2', 'jokojob2', 'joko2@gmail.com', 'joko2hash', 'joko2.jpg', 'user', 'joko2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('joko3', 'jokojob3', 'joko3@gmail.com', 'joko3hash', 'joko3.jpg', 'user', 'joko3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


# db, err := sql.Open("mysql", "root:a@tcp(localhost:5555)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local")
# "mysql://root:a@tcp(localhost:5555)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
# mysqlvolume
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
