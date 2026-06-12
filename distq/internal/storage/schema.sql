CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    paylaod TEXT NOT NULL,
    
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    
    lease_id TEXT,
    leased_by TEXT,
    lease_until DATETIME,
    
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    
    last_error TEXT
);

CREATE INDEX IF NOT EXISTS idx_jobs_pick
ON jobs(status, created_at);

CREATE INDEX IF NOT EXISTS idx_jobs_lease_until
ON jobs(status, lease_until);
