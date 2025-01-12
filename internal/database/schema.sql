PRAGMA journal_mode = WAL;

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS group_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    colour TEXT DEFAULT "", -- colour selected for list
    icon TEXT DEFAULT "", -- icon selected for list
    filter_by TEXT DEFAULT "", -- if filter exists the list becomes a smart list, where filter dictates which tasks to show
    group_id INTEGER DEFAULT 0, -- REFERENCES group_list (id) ON UPDATE CASCADE ON DELETE SET DEFAULT,
    pinned BOOLEAN DEFAULT false, -- show list in highlighted lists area
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT "",
    completed BOOLEAN DEFAULT false,
    important BOOLEAN DEFAULT false,
    priority INTEGER DEFAULT 0, -- None=0, Low=1, Medium=2, High=3
    position INTEGER DEFAULT 0,
    start_at TIMESTAMP DEFAULT "",
    end_at TIMESTAMP DEFAULT "",
    list_id INTEGER DEFAULT 0 REFERENCES list (id) ON UPDATE CASCADE ON DELETE SET DEFAULT,
    parent_id INTEGER DEFAULT 0, -- if value exist it means that this task is a sub-task of another
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Initialise default lists
INSERT INTO
    list (name, colour, icon, pinned, position)
VALUES
    (
        "Today",
        "--colour-fresh-blue-500",
        "today-icon",
        true,
        1
    ),
    (
        "Scheduled",
        "--colour-cyan-700",
        "days-icon",
        true,
        2
    ),
    ("All", "--base-colour", "icon-tasks", true, 3),
    (
        "Important",
        "--colour-volcano-400",
        "icon-star",
        true,
        4
    ),
    (
        "Completed",
        "--colour-lime-700",
        "icon-check-square",
        true,
        5
    ),
    ("Inbox", "--base-colour", "today-icon", false, 6);
