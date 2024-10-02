-- +migrate Up
-- MYSQL 8.0 set the role value to `mysqluser' for all records
-- TODO find better solution instead of keeping the order
ALTER TABLE `users` ADD COLUMN `role` ENUM('mysqluser','admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`;