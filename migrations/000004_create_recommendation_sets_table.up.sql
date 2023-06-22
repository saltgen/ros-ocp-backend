CREATE TABLE IF NOT EXISTS recommendation_sets(
   id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
   workload_id BIGINT NOT NULL,
   container_name TEXT NOT NULL,
   monitoring_start_time TIMESTAMP WITH TIME ZONE NOT NULL,
   monitoring_end_time TIMESTAMP WITH TIME ZONE NOT NULL,
   recommendations jsonb NOT NULL,
   updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_recommendation_sets_workload_id_container_name_end_time ON recommendation_sets (workload_id, container_name, monitoring_end_time);

ALTER TABLE recommendation_sets
ADD CONSTRAINT fk_recommendation_sets_workload FOREIGN KEY (workload_id) REFERENCES workloads (id)
ON DELETE CASCADE;
