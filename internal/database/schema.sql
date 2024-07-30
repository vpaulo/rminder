CREATE TABLE IF NOT EXISTS task (
    task_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT "",
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO task (title) Values('first task')
