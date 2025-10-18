-- Rename indexes back to original names
ALTER INDEX idx_org_members_user_id RENAME TO idx_uom_user_id;
ALTER INDEX idx_org_members_org_id RENAME TO idx_uom_org_id;
ALTER INDEX idx_org_members_role RENAME TO idx_uom_role;

-- Rename table back to user_organization_memberships
ALTER TABLE organization_members RENAME TO user_organization_memberships;