CREATE TABLE IF NOT EXISTS task (
    task_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT "",
    completed BOOLEAN DEFAULT false,
    important BOOLEAN DEFAULT false,
    my_day BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO task (title) Values('first task')
