
-- +migrate Up

CREATE TABLE scan_history (
  repo_id VARCHAR(255) PRIMARY KEY NOT NULL,
  commit_hash VARCHAR(255) NULL,
  scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE IF EXISTS scan_history;
