CREATE TABLE "users" (
                         "user_id" bigserial PRIMARY KEY,
                         "name" varchar NOT NULL,
                         "username" varchar NOT NULL,
                         "email" varchar NOT NULL,
                         "password" varchar NOT NULL,
                         "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "directory" (
                             "dir_id" bigserial PRIMARY KEY,
                             "name" varchar NOT NULL,
                             "user_id" bigint NOT NULL,
                             "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "directory" ("user_id");