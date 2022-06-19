-- Create Index
CREATE UNIQUE INDEX "tags.name_lower_unique" ON tags((LOWER(name)));

-- Alter columns
UPDATE portfolios_securities SET attributes = '[]' WHERE attributes IS NULL;
UPDATE portfolios_securities SET events = '[]' WHERE events IS NULL;
UPDATE portfolios_securities SET properties = '[]' WHERE properties IS NULL;
ALTER TABLE portfolios_securities
  ALTER COLUMN attributes SET DEFAULT '[]',
  ALTER COLUMN attributes SET NOT NULL;
ALTER TABLE portfolios_securities
  ALTER COLUMN events SET DEFAULT '[]',
  ALTER COLUMN events SET NOT NULL;
ALTER TABLE portfolios_securities
  ALTER COLUMN properties SET DEFAULT '[]',
  ALTER COLUMN properties SET NOT NULL;

-- Add columns
ALTER TABLE securities ADD COLUMN extras JSONB NOT NULL DEFAULT '{}';
ALTER TABLE securities_markets ADD COLUMN extras JSONB NOT NULL DEFAULT '{}';
