
CREATE TABLE posts (
                       id SERIAL PRIMARY KEY UNIQUE,
                       title VARCHAR(255) NOT NULL,
                       content TEXT NOT NULL,
                       allow_comments BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TABLE comments (
                          id SERIAL PRIMARY KEY UNIQUE,
                          post_id INT NOT NULL,
                          parent_id INT,
                          author VARCHAR(255) NOT NULL,
                          content TEXT NOT NULL CHECK (LENGTH(content) <= 2000),
                          CONSTRAINT fk_post
                              FOREIGN KEY(post_id)
                                  REFERENCES posts(id)
                                  ON DELETE CASCADE,
                          CONSTRAINT fk_parent
                              FOREIGN KEY(parent_id)
                                  REFERENCES comments(id)
                                  ON DELETE CASCADE
);

CREATE INDEX idx_post_id ON comments(post_id);
CREATE INDEX idx_parent_id ON comments(parent_id);