CREATE TABLE "user" (
  "id" integer PRIMARY KEY,
  "username" string,
  "password" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "owner" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_owner" (
  "id" integer PRIMARY KEY,
  "book_id" integer,
  "owner_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book" (
  "id" integer PRIMARY KEY,
  "title" string,
  "number" integer,
  "release_date" integer,
  "cover_image" string,
  "description" string,
  "notes" string,
  "book_type" integer UNIQUE,
  "book_edition" integer UNIQUE,
  "printing" integer,
  "isbn" string,
  "price" string,
  "count" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "permission" (
  "id" integer PRIMARY KEY,
  "name" string,
  "description" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "role" (
  "id" integer PRIMARY KEY,
  "name" string,
  "description" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "user_role" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "role_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "user_permission" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "permission_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "role_permission" (
  "id" integer PRIMARY KEY,
  "role_id" integer,
  "permission_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "publisher" (
  "id" integer PRIMARY KEY,
  "name" string,
  "website_url" string,
  "country" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_publisher" (
  "id" integer PRIMARY KEY,
  "book_id" integer,
  "publisher_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "author" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_author" (
  "id" integer PRIMARY KEY,
  "book_id" integer,
  "author_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_type" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "location" (
  "id" integer PRIMARY KEY,
  "name" string,
  "description" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_location" (
  "id" integer PRIMARY KEY,
  "book_id" integer,
  "location_id" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_edition" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

ALTER TABLE "book_owner" ADD FOREIGN KEY ("book_id") REFERENCES "book" ("id");

ALTER TABLE "book_owner" ADD FOREIGN KEY ("owner_id") REFERENCES "owner" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("book_type") REFERENCES "book_type" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("book_edition") REFERENCES "book_edition" ("id");

ALTER TABLE "user_role" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_role" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE "user_permission" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "permission" ADD FOREIGN KEY ("id") REFERENCES "user_permission" ("permission_id");

ALTER TABLE "role_permission" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE "role_permission" ADD FOREIGN KEY ("permission_id") REFERENCES "permission" ("id");

ALTER TABLE "book_publisher" ADD FOREIGN KEY ("book_id") REFERENCES "book" ("id");

ALTER TABLE "book_publisher" ADD FOREIGN KEY ("publisher_id") REFERENCES "publisher" ("id");

ALTER TABLE "book_author" ADD FOREIGN KEY ("book_id") REFERENCES "book" ("id");

ALTER TABLE "book_author" ADD FOREIGN KEY ("author_id") REFERENCES "author" ("id");

ALTER TABLE "book_location" ADD FOREIGN KEY ("book_id") REFERENCES "book" ("id");

ALTER TABLE "book_location" ADD FOREIGN KEY ("location_id") REFERENCES "location" ("id");
