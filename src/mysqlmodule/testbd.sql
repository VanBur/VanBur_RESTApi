CREATE TABLE `content` (`id` int(11) NOT NULL, `protection_system_id` int(11) NOT NULL, `content_key` varchar(32) NOT NULL, `payload` varchar(64) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `devices` ( `id` int(11) NOT NULL, `name` varchar(50) NOT NULL, `protection_system_id` int(11) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `devices` (`id`, `name`, `protection_system_id`) VALUES (1, 'Android', 1), (2, 'Samsung', 2), (3, 'iOS', 1), (4, 'LG', 2);
CREATE TABLE `protection_systems` (`id` int(11) NOT NULL, `name` varchar(20) NOT NULL, `encryption_mode` varchar(30) NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `protection_systems` (`id`, `name`, `encryption_mode`) VALUES (1, 'AES 1', 'AES + ECB'),(2, 'AES 2', 'AES + CBC');
ALTER TABLE `content` ADD UNIQUE KEY `id` (`id`);
ALTER TABLE `devices` ADD UNIQUE KEY `id` (`id`), ADD KEY `protection_system` (`protection_system_id`);
ALTER TABLE `protection_systems` ADD PRIMARY KEY (`id`);
ALTER TABLE `content` MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;