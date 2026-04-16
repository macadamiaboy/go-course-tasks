CREATE OR REPLACE FUNCTION process_transfer()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE accounts 
    SET coins = coins - NEW.amount 
    WHERE id = NEW.source;

    UPDATE accounts 
    SET coins = coins + NEW.amount 
    WHERE id = NEW.target;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_after_transfer_insert ON transfers;

CREATE TRIGGER trg_after_transfer_insert
AFTER INSERT ON transfers
FOR EACH ROW
EXECUTE FUNCTION process_transfer();