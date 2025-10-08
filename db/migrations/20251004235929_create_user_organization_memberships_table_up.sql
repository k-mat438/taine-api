CREATE TABLE IF NOT EXISTS user_organization_memberships (
  id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         uuid        NOT NULL REFERENCES users(id),
  organization_id uuid        NOT NULL REFERENCES organizations(id),
  role            text        NOT NULL,  -- 'owner'|'admin'|'member' etc.
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (user_id, organization_id)
);
CREATE INDEX IF NOT EXISTS idx_uom_user_id ON user_organization_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_uom_org_id  ON user_organization_memberships(organization_id);
CREATE INDEX IF NOT EXISTS idx_uom_role    ON user_organization_memberships(role);
