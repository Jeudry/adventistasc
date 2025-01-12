DROP INDEX idx_comments_content;
DROP INDEX idx_comments_post_id;

DROP INDEX idx_posts_title;
DROP INDEX idx_posts_content;
DROP INDEX idx_posts_user_id;

DROP INDEX idx_users_name;

DROP EXTENSION pg_trgm;