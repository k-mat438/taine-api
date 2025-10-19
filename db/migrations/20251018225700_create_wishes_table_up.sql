CREATE TABLE IF NOT EXISTS wishes (
  id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id uuid        NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  title           text        NOT NULL,
  note            text        NOT NULL DEFAULT '',
  order_no        int         NOT NULL DEFAULT 0,       -- 並び順/優先度
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  deleted_at      timestamptz
);

CREATE INDEX idx_wishes_org ON wishes(organization_id);
CREATE INDEX idx_wishes_priority ON wishes(organization_id, order_no DESC);