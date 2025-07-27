CREATE TABLE users
(
   id INT AUTO_INCREMENT PRIMARY KEY,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null,
    UNIQUE (id)
);

CREATE TABLE todo_lists
(
  id INT AUTO_INCREMENT PRIMARY KEY,
    title       varchar(255) not null,
    description varchar(255),
    UNIQUE (id)
);

CREATE TABLE users_lists
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    list_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists (id) ON DELETE CASCADE,
    UNIQUE (id)
);

CREATE TABLE todo_items
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean      not null default false,
    UNIQUE (id)
);

CREATE TABLE lists_items
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT NOT NULL,
    list_id INT NOT NULL,
    FOREIGN KEY (item_id) REFERENCES todo_items (id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists (id) ON DELETE CASCADE,
    UNIQUE (id)
);