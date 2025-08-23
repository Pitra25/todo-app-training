CREATE TABLE users
(
   id SERIAL PRIMARY KEY,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null,
    email         varchar(255) not null,
    status        varchar(50) not null default 'not_confirmed'
);

CREATE TABLE todo_lists
(
  id SERIAL PRIMARY KEY,
    title       varchar(255) not null,
    description varchar(255)
);

CREATE TABLE users_lists
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    list_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists (id) ON DELETE CASCADE
);

CREATE TABLE todo_items
(
    id SERIAL PRIMARY KEY,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean      not null default false
);

CREATE TABLE lists_items
(
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    list_id INT NOT NULL,
    FOREIGN KEY (item_id) REFERENCES todo_items (id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists (id) ON DELETE CASCADE
);

CREATE TABLE users_code
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    code varchar(255) not null,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (user_id, code)
);