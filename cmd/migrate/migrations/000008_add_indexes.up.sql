CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin(content gin_trgm_ops);
CREATE INDEX idx_comments_post_id ON comments (post_id);

CREATE INDEX idx_posts_title ON posts USING gin(title gin_trgm_ops);
CREATE INDEX idx_posts_content ON posts USING gin(content gin_trgm_ops);
CREATE INDEX idx_posts_user_id ON posts (user_id);

CREATE INDEX idx_users_name ON users (user_name);