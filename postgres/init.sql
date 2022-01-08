CREATE TABLE "public"."users" (
    "user_id" serial4,
    "name" varchar,
    "username" varchar,
    "password" varchar,
    "email" varchar,
    "created_at" timestamp DEFAULT now(),
    PRIMARY KEY ("user_id")
);