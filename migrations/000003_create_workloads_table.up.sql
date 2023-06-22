CREATE TYPE workloadtype AS ENUM ('deployment', 'deploymentconfig', 'replicaset', 'replicationcontroller', 'statefulset', 'daemonset');

CREATE TABLE IF NOT EXISTS workloads(
   id BIGSERIAL PRIMARY KEY,
   cluster_id BIGINT NOT NULL,
   experiment_name TEXT NOT NULL,
   namespace TEXT NOT NULL,
   workload_type workloadtype NOT NULL,
   workload_name TEXT NOT NULL,
   containers TEXT[] NOT NULL,
   metrics_upload_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_workloads_containers ON workloads USING gin(containers);

ALTER TABLE workloads
ADD CONSTRAINT fk_workloads_cluster FOREIGN KEY (cluster_id) REFERENCES clusters (id)
ON DELETE CASCADE;

ALTER TABLE workloads
ADD UNIQUE (cluster_id, experiment_name);
