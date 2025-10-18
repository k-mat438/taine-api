-- Rename table from user_organization_memberships to organization_members
ALTER TABLE user_organization_memberships RENAME TO organization_members;

-- Rename indexes to match new table name
ALTER INDEX idx_uom_user_id RENAME TO idx_org_members_user_id;
ALTER INDEX idx_uom_org_id RENAME TO idx_org_members_org_id;
ALTER INDEX idx_uom_role RENAME TO idx_org_members_role;