PRAGMA journal_mode = WAL;

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
    list_id INTEGER DEFAULT 0, -- task should belong to only one list
    parent_id INTEGER DEFAULT 0, -- if value exist it means that this task is a sub-task of another
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    colour TEXT DEFAULT "", -- colour selected for list
    icon TEXT DEFAULT "", -- icon selected for list
    filter_by TEXT DEFAULT "", -- if filter exists the list becomes a smart list, where filter dictates which tasks to show
    group_id INTEGER DEFAULT 0, -- a list should belong to only one group
    pinned BOOLEAN DEFAULT false, -- show list in highlighted lists area
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS group_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Initialise default lists
INSERT INTO
    list (name, colour, icon, pinned)
VALUES
    (
        "Today",
        "--colour-fresh-blue-500",
        "today-icon",
        true
    );

INSERT INTO
    list (name, colour, icon, pinned)
VALUES
    (
        "Scheduled",
        "--colour-dust-red-500",
        "today-icon",
        true
    );

INSERT INTO
    list (name, colour, icon, pinned)
VALUES
    ("All", "--colour-indigo-500", "icon-tasks", true);

INSERT INTO
    list (name, colour, icon, pinned)
VALUES
    (
        "Important",
        "--colour-volcano-400",
        "icon-star",
        true
    );

INSERT INTO
    list (name, colour, icon, pinned)
VALUES
    (
        "Completed",
        "--colour-orange-300",
        "icon-square",
        true
    );

INSERT INTO
    list (name, colour, icon)
VALUES
    ("Inbox", "--colour-fresh-blue-500", "today-icon");
