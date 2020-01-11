-- up
CREATE TABLE IF NOT EXISTS repositories (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    url VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS repositories_after_update AFTER UPDATE ON repositories BEGIN
    UPDATE repositories SET modified_at = datetime('now') WHERE id = NEW.id;
end;

CREATE UNIQUE INDEX IF NOT EXISTS repositories_idx ON repositories(url, name);

CREATE TABLE IF NOT EXISTS branches (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    ticket_id VARCHAR(50) NOT NULL DEFAULT '',
    parent_ticket_id VARCHAR(50) NOT NULL DEFAULT '',
    repository_id INTEGER NOT NULL,
    ticket_summary VARCHAR(250) NOT NULL,
    ticket_status VARCHAR(100) NOT NULL,
    ticket_type VARCHAR(50) NOT NULL,
    branch_name VARCHAR(250) NOT NULL,
    closed INTEGER(1) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repository_id) REFERENCES repositories(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS branches_idx ON branches(branch_name, repository_id);

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


-- down
DROP TRIGGER IF EXISTS branch_commits_after_update;
DROP INDEX IF EXISTS branch_commits_idx;
DROP TABLE IF EXISTS branch_commits;

DROP TRIGGER IF EXISTS commits_after_update;
DROP TABLE IF EXISTS commits;

DROP TRIGGER IF EXISTS branch_versions_after_update;
DROP INDEX IF EXISTS branch_versions_idx;
DROP TABLE IF EXISTS branch_versions;

DROP TRIGGER IF EXISTS versions_after_update;
DROP TABLE IF EXISTS versions;

DROP TRIGGER IF EXISTS branches_after_update;
DROP TABLE IF EXISTS branches;

DROP TRIGGER IF EXISTS repositories_after_update;
DROP INDEX IF EXISTS repositories_idx;
DROP TABLE IF EXISTS repositories;
