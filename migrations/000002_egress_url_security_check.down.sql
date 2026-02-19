DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'egress_url_security_check'
  ) THEN
    ALTER TABLE sources
    DROP CONSTRAINT egress_url_security_check;
  END IF;
END
$$;
