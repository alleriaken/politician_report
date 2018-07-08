CREATE SCHEMA `politicians` DEFAULT CHARACTER SET utf8mb4 ;

USE `politicians`;

CREATE TABLE `keywords` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `keyword` VARCHAR(255) NOT NULL DEFAULT "",
  `time_of_article` VARCHAR(64) NOT NULL DEFAULT "w",
  `last_crawled_at` INT(11) NULL DEFAULT 0,
  `crawl_delay_time` INT(11) NULL DEFAULT 0,
  PRIMARY KEY (`id`)
);

CREATE TABLE `crawled` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `url` TEXT NOT NULL,
  `host` TEXT NOT NULL,
  `title` TEXT NOT NULL,
  `google_title` TEXT NOT NULL,
  `preview_content` TEXT NOT NULL,
  `crawled_at` INT(11) NOT NULL DEFAULT 0,
  `keyword_id` INT(11) NOT NULL DEFAULT 0,
  `negative` INT(11) NOT NULL DEFAULT 0,
  `positive` INT(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
);

CREATE TABLE `email_delivering` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `keyword_id` INT(11) NULL DEFAULT NULL ,
  `email` TEXT NOT NULL,
  PRIMARY KEY (`id`)
);