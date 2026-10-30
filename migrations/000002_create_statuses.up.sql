CREATE TABLE statuses (
    code text PRIMARY KEY,
    title text NOT NULL
);

INSERT INTO statuses (code, title) VALUES
    ('active', 'В работе'),
    ('pending', 'Ожидает'),
    ('closed', 'Закрыта');
