CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
    "name" TEXT NOT NULL,
    "username" varchar(255) NOT NULL,
    "email" varchar(255),
    "password" varchar(255) NOT NULL
)