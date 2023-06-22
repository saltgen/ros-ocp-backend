CREATE TABLE IF NOT EXISTS clusters(
   id BIGSERIAL PRIMARY KEY,
   tenant_id BIGINT NOT NULL,
   source_id TEXT NOT NULL,
   cluster_uuid TEXT NOT NULL,
   cluster_alias TEXT NOT NULL,
   last_reported_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_clusters_tenant_id ON clusters (tenant_id);
CREATE INDEX idx_clusters_last_reported ON clusters (last_reported_at);

ALTER TABLE clusters
ADD CONSTRAINT fk_clusters_rh_account FOREIGN KEY (tenant_id) REFERENCES rh_accounts (id)
ON DELETE CASCADE;

ALTER TABLE clusters
ADD UNIQUE (tenant_id, source_id, cluster_uuid, cluster_alias);
