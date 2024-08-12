CREATE TABLE "user" (
  "id" integer PRIMARY KEY,
  "username" string,
  "password" string,
  "roles" integer,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "owner" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book" (
  "id" integer PRIMARY KEY,
  "title" string,
  "number" integer,
  "release_date" integer,
  "cover_image" binary,
  "description" string,
  "notes" string,
  "authors" integer,
  "publisher" integer,
  "location" integer,
  "owners" integer,
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
  "permissions" integer,
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

CREATE TABLE "author" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_type" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_location" (
  "id" integer PRIMARY KEY,
  "name" string,
  "description" string,
  "created_at" integer,
  "updated_at" integer
);

CREATE TABLE "book_edition" (
  "id" integer PRIMARY KEY,
  "name" string,
  "created_at" integer,
  "updated_at" integer
);

ALTER TABLE "user" ADD FOREIGN KEY ("roles") REFERENCES "role" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("authors") REFERENCES "author" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("publisher") REFERENCES "publisher" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("location") REFERENCES "book_location" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("owners") REFERENCES "owner" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("book_type") REFERENCES "book_type" ("id");

ALTER TABLE "book" ADD FOREIGN KEY ("book_edition") REFERENCES "book_edition" ("id");

ALTER TABLE "role" ADD FOREIGN KEY ("permissions") REFERENCES "permission" ("id");
