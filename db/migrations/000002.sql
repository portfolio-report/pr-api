-- Create Tables
CREATE TABLE "tags" (
  "name" TEXT NOT NULL,

  PRIMARY KEY ("name")
);

CREATE TABLE "securities_tags" (
  "security_uuid" UUID NOT NULL,
  "tag_name" TEXT NOT NULL,

  PRIMARY KEY ("security_uuid", "tag_name")
);

-- Create Index
CREATE INDEX "securities_tags.tag_name" ON "securities_tags"("tag_name");

-- Add Foreign Keys
ALTER TABLE "securities_tags" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_tags" ADD FOREIGN KEY ("tag_name") REFERENCES "tags"("name") ON DELETE CASCADE ON UPDATE CASCADE;
