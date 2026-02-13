-- Indexes for task table
CREATE INDEX IF NOT EXISTS idx_task_list_id ON task(list_id);
CREATE INDEX IF NOT EXISTS idx_task_completed ON task(completed);
CREATE INDEX IF NOT EXISTS idx_task_important ON task(important);
CREATE INDEX IF NOT EXISTS idx_task_parent_id ON task(parent_id);
CREATE INDEX IF NOT EXISTS idx_task_position ON task(position);

-- Indexes for list table
CREATE INDEX IF NOT EXISTS idx_list_group_id ON list(group_id);
CREATE INDEX IF NOT EXISTS idx_list_position ON list(position);

-- Indexes for group_list table
CREATE INDEX IF NOT EXISTS idx_group_list_position ON group_list(position);
