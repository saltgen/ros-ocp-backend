CREATE TABLE IF NOT EXISTS rh_accounts(
   id BIGSERIAL PRIMARY KEY,
   account TEXT UNIQUE,
   org_id TEXT UNIQUE NOT NULL
);

CREATE INDEX idx_rh_accounts_org_id ON rh_accounts (org_id);