-- Создание таблицы комментариев
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    news_id INT,
    text TEXT,
    parent_comment_id INT
);
