DO $$
BEGIN
  IF current_setting('app.env', true) = 'prod' THEN
    ALTER TABLE sources
        ADD CONSTRAINT egress_url_security_check
        CHECK (
            -- Require HTTP/HTTPS and block obvious internal/loopback/metadata targets for SSRF mitigation
            egress_url ~ '^https?://' AND
            egress_url !~* '^https?://(localhost|127\.0\.0\.1|0\.0\.0\.0|\[?::1\]?)(/|:|$)' AND
            egress_url !~* '^https?://10\.' AND
            egress_url !~* '^https?://192\.168\.' AND
            egress_url !~* '^https?://172\.(1[6-9]|2[0-9]|3[0-1])\.' AND
            egress_url !~* '^https?://169\.254\.169\.254(/|:|$)' AND
            CHAR_LENGTH(egress_url) <= 2048
        );
  END IF;
END
$$;
