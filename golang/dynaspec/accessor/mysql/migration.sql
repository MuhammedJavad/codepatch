CREATE TABLE trees (
  id UNSIGNED INT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  name VARCHAR(255) NOT NULL,
  start_time TIMESTAMP NULL,
  end_time TIMESTAMP NULL,
  result JSON NOT NULL,
  structure JSON NOT NULL,
  INDEX idx_trees_created_at ON trees(created_at),
  INDEX idx_trees_active ON trees(active)
);