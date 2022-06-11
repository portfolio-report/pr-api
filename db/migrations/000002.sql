-- Create Tables
CREATE TABLE "tags" (
  "uuid" UUID NOT NULL,
  "name" TEXT NOT NULL,

  PRIMARY KEY ("uuid")
);

CREATE TABLE "securities_tags" (
  "security_uuid" UUID NOT NULL,
  "tag_uuid" UUID NOT NULL,

  PRIMARY KEY ("security_uuid", "tag_uuid")
);

-- Create Indexes
CREATE INDEX "tags.name" ON "tags"("name");
CREATE INDEX "securities_tags.tag_uuid" ON "securities_tags"("tag_uuid");

-- Add Foreign Keys
ALTER TABLE "securities_tags" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE CASCADE;
ALTER TABLE "securities_tags" ADD FOREIGN KEY ("tag_uuid") REFERENCES "tags"("uuid") ON DELETE CASCADE;
