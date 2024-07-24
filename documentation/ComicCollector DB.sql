CREATE TABLE `users` (
  `id` integer PRIMARY KEY,
  `username` string,
  `password_hash` string,
  `permissions` integer
);

CREATE TABLE `owners` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `books` (
  `id` integer PRIMARY KEY,
  `title` integer,
  `number` integer,
  `release_date` integer,
  `added_date` integer,
  `cover_image` string,
  `description` string,
  `notes` string,
  `publisher` integer,
  `authors` integer,
  `book_type` integer UNIQUE,
  `book_edition` integer,
  `printing` integer,
  `isbn` string,
  `price` string,
  `count` integer,
  `location` string,
  `owner` integer
);

CREATE TABLE `permissions` (
  `id` integer PRIMARY KEY,
  `name` string,
  `description` string
);

CREATE TABLE `roles` (
  `id` integer PRIMARY KEY,
  `name` string,
  `description` string
);

CREATE TABLE `user_roles` (
  `user_id` integer,
  `role_id` integer
);

CREATE TABLE `user_permissions` (
  `user_id` integer,
  `permission_id` integer
);

CREATE TABLE `role_permissions` (
  `role_id` integer,
  `permission_id` integer
);

CREATE TABLE `publishers` (
  `id` integer PRIMARY KEY,
  `name` string,
  `website_url` string,
  `country` string
);

CREATE TABLE `authors` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `book_types` (
  `id` integer PRIMARY KEY,
  `name` string
);

CREATE TABLE `book_editions` (
  `id` integer PRIMARY KEY,
  `name` string
);

ALTER TABLE `books` ADD FOREIGN KEY (`publisher`) REFERENCES `publishers` (`id`);

ALTER TABLE `books` ADD FOREIGN KEY (`authors`) REFERENCES `authors` (`id`);

ALTER TABLE `books` ADD FOREIGN KEY (`book_type`) REFERENCES `book_types` (`id`);

ALTER TABLE `books` ADD FOREIGN KEY (`book_edition`) REFERENCES `book_editions` (`id`);

ALTER TABLE `books` ADD FOREIGN KEY (`owner`) REFERENCES `owners` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `user_permissions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `permissions` ADD FOREIGN KEY (`id`) REFERENCES `user_permissions` (`permission_id`);

ALTER TABLE `role_permissions` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `role_permissions` ADD FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`);
