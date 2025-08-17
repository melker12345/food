-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create extension for array operations
CREATE EXTENSION IF NOT EXISTS "intarray";

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE food_app TO postgres;

-- Note: Tables will be created by GORM auto-migration
-- This file can be extended with custom indexes, triggers, or initial data
