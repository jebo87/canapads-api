CREATE OR REPLACE FUNCTION update_last_updated_field
()
	returns trigger as
$BODY$
BEGIN
    UPDATE ***REMOVED*** SET last_updated = now() where id = NEW.id;
    PERFORM pg_notify
    (
    'events',
    json_build_object
    (
      'operation', TG_OP,
      'record', row_to_json
    (NEW)
    )::text
  );
RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

CREATE TRIGGER ad_updated
AFTER
INSERT OR
UPDATE
ON ***REMOVED***
FOR EACH ROW
EXECUTE PROCEDURE update_last_updated_field
();