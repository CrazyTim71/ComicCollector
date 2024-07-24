CREATE TABLE `user` (
  `id` integer PRIMARY KEY,
  `username` string,
  `password_hash` string
);

CREATE TABLE `owner` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `book` (
  `id` integer PRIMARY KEY,
  `title` integer,
  `number` integer,
  `release_date` integer,
  `added_date` integer,
  `cover_image` string,
  `description` string,
  `notes` string,
  `publishers` integer,
  `authors` integer,
  `book_type` integer UNIQUE,
  `book_edition` integer,
  `printing` integer,
  `isbn` string,
  `price` string,
  `count` integer,
  `location` integer,
  `owner` integer
);

CREATE TABLE `permission` (
  `id` integer PRIMARY KEY,
  `name` string,
  `description` string
);

CREATE TABLE `role` (
  `id` integer PRIMARY KEY,
  `name` string,
  `description` string
);

CREATE TABLE `user_role` (
  `id` integer PRIMARY KEY,
  `user_id` integer,
  `role_id` integer
);

CREATE TABLE `user_permission` (
  `id` integer PRIMARY KEY,
  `user_id` integer,
  `permission_ids` integer
);

CREATE TABLE `role_permission` (
  `id` integer PRIMARY KEY,
  `role_id` integer,
  `permission_id` integer
);

CREATE TABLE `publisher` (
  `id` integer PRIMARY KEY,
  `name` string,
  `website_url` string,
  `country` string
);

CREATE TABLE `author` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `book_type` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `location` (
  `id` integer PRIMARY KEY,
  `name` string,
  `description` string
);

CREATE TABLE `book_edition` (
  `id` integer PRIMARY KEY,
  `name` string
);

ALTER TABLE `book` ADD FOREIGN KEY (`publishers`) REFERENCES `publisher` (`id`);

ALTER TABLE `book` ADD FOREIGN KEY (`authors`) REFERENCES `author` (`id`);

ALTER TABLE `book` ADD FOREIGN KEY (`book_type`) REFERENCES `book_type` (`id`);

ALTER TABLE `book` ADD FOREIGN KEY (`book_edition`) REFERENCES `book_edition` (`id`);

ALTER TABLE `book` ADD FOREIGN KEY (`location`) REFERENCES `location` (`id`);

ALTER TABLE `book` ADD FOREIGN KEY (`owner`) REFERENCES `owner` (`id`);

ALTER TABLE `user_role` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `user_role` ADD FOREIGN KEY (`role_id`) REFERENCES `role` (`id`);

ALTER TABLE `user_permission` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `permission` ADD FOREIGN KEY (`id`) REFERENCES `user_permission` (`permission_ids`);

ALTER TABLE `role_permission` ADD FOREIGN KEY (`role_id`) REFERENCES `role` (`id`);

ALTER TABLE `role_permission` ADD FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id`);
