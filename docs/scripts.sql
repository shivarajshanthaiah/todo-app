
-- Create table users
CREATE TABLE users (
    id VARCHAR(63) NOT NULL,
    username VARCHAR(63) NOT NULL,                     
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    email VARCHAR(63) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX emailusername ON users (lower(email));


CREATE TABLE tasks (
  id SERIAL, -- can be uuid genereated by DB itself, here im keeping it simple
  user_id VARCHAR(63) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  title VARCHAR(119) NOT NULL,
  description TEXT,
  priority INT NOT NULL DEFAULT 1,
  status INT NOT NULL DEFAULT 1,
  due_at TIMESTAMP
);

CREATE INDEX idx_tasks_user_id ON tasks (user_id); -- to make the query excecute faster