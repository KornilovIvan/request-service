ALTER TABLE documents DROP CONSTRAINT documents_request_id_fkey;

ALTER TABLE requests DROP CONSTRAINT requests_status_fkey;

ALTER TABLE requests DROP CONSTRAINT requests_customer_id_fkey;
