# vim: filetype=mysql

CREATE TABLE IF NOT EXISTS todo_items (
  id           BIGINT       PRIMARY KEY AUTO_INCREMENT,
  description  VARCHAR(255) UNIQUE NOT NULL,
  completed    BOOL         NOT NULL DEFAULT false
);

INSERT IGNORE INTO todo_items (description) VALUES
  ('Hello World!'),
  ('Chaos Testing!'),
  ('Todo List!'),
  ('Johny Appleseed!');

