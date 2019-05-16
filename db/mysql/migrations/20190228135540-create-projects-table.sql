
-- +migrate Up

CREATE TABLE scan_history (
  id CHAR(36) PRIMARY KEY NOT NULL,
  repo_id VARCHAR(255) NOT NULL,
  commit_hash VARCHAR(255) NULL,
  scanned_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS scan_history;
