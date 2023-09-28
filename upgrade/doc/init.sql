CREATE TABLE `us_file_element` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `path` varchar(255) NOT NULL COMMENT '包内相对路径',
  `version` varchar(100) NOT NULL COMMENT '所属版本号',
  `storage_path` varchar(255) NOT NULL COMMENT '磁盘存储路径',
  `release_date` datetime NOT NULL COMMENT '发布时间',
  `checksum` varchar(32) NOT NULL COMMENT 'MD5校验值',
  `size` bigint(20) NOT NULL COMMENT '文件大小字节数',
  `release_id` bigint(20) NOT NULL COMMENT '发布ID',
  PRIMARY KEY (`id`),
  KEY `idx_file_rid` (`release_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `us_release_notes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `release_id` bigint(20) NOT NULL COMMENT '软件包发布标识',
  `lang` varchar(20) NOT NULL COMMENT '语言',
  `release_notes` text COMMENT '升级内容说明',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_notes_rid` (`release_id`, `lang`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `us_site_release` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `release_id` bigint(20) NOT NULL COMMENT '软件包发布标识',
  `release_version` varchar(50) NOT NULL COMMENT '发布版本',
  `site_id` varchar(50) NOT NULL COMMENT '产品标识',
  `application_id` bigint(20) NOT NULL COMMENT '应用标识',
  `client_type` varchar(50) NOT NULL COMMENT '终端类型',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '发布状态,0:未发布，1:已发布',
  `extend_attr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '扩展属性',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_release_id` (`release_id`) USING BTREE,
  KEY `idx_release_version` (`release_version`) USING BTREE,
  KEY `idx_release_appid` (`application_id`) USING BTREE,
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `us_diff_patch` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `storage_path` varchar(255) NOT NULL COMMENT '差异包存储路径',
  `checksum` varchar(128) NOT NULL COMMENT '差异包文件内容校验和',
  `size` bigint(20) NOT NULL COMMENT '差异包字节大小',
  `file_hash` varchar(128) NOT NULL COMMENT '差异文件的MD5值',
  PRIMARY KEY (`id`),
  KEY `idx_file_hash` (`file_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `us_ticket` (
  `key_name` varchar(100) NOT NULL,
  `max_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`key_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `test`.`us_ticket` (`key_name`, `max_id`) VALUES ('release', 1);