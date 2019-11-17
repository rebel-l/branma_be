-- up
CREATE TABLE IF NOT EXISTS branches (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
    ticket_id VARCHAR(50) NOT NULL DEFAULT '',
    parent_ticket_id VARCHAR(50) NOT NULL DEFAULT '',
    ticket_summary VARCHAR(250) NOT NULL,
    ticket_status VARCHAR(100) NOT NULL,
    ticket_type VARCHAR(50) NOT NULL,
    branch_name VARCHAR(250) NOT NULL DEFAULT '',
    closed INTEGER(1) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS branches_after_update AFTER UPDATE ON branches BEGIN
    UPDATE branches SET modified_at = datetime('now') WHERE id = NEW.id;
end;

CREATE TABLE IF NOT EXISTS versions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    version VARCHAR(50) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS versions_after_update AFTER UPDATE ON versions BEGIN
    UPDATE versions SET modified_at = datetime('now') WHERE id = NEW.id;
end;

CREATE TABLE IF NOT EXISTS branch_versions (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    branch_id INTEGER NOT NULL,
    version_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (branch_id) REFERENCES branches(id),
    FOREIGN KEY (version_id) REFERENCES versions(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS branch_versions_idx ON branch_versions(branch_id,version_id);

CREATE TRIGGER IF NOT EXISTS branch_versions_after_update AFTER UPDATE ON branch_versions BEGIN
    UPDATE branch_versions SET modified_at = datetime('now') WHERE id = NEW.id;
end;

CREATE TABLE IF NOT EXISTS commits (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    commit_hash VARCHAR(50) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS commits_after_update AFTER UPDATE ON commits BEGIN
    UPDATE commits SET modified_at = datetime('now') WHERE id = NEW.id;
end;

CREATE TABLE IF NOT EXISTS branch_commits (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    branch_id INTEGER NOT NULL,
    commit_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (branch_id) REFERENCES branches(id),
    FOREIGN KEY (commit_id) REFERENCES commits(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS branch_commits_idx ON branch_commits(branch_id,commit_id);

CREATE TRIGGER IF NOT EXISTS branch_commits_after_update AFTER UPDATE ON branch_commits BEGIN
    UPDATE branch_commits SET modified_at = datetime('now') WHERE id = NEW.id;
end;

-- TODO: Upgrade & Downgrade