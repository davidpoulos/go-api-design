CREATE TABLE "user"(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    password VARCHAR (70) NOT NULL,
    email VARCHAR (50) UNIQUE NOT NULL,
    date_created TIMESTAMP,
    "role" VARCHAR (50)
);

