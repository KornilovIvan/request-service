ALTER TABLE requests
    ADD CONSTRAINT requests_customer_id_fkey
    FOREIGN KEY (customer_id) REFERENCES customers (id);

ALTER TABLE requests
    ADD CONSTRAINT requests_status_fkey
    FOREIGN KEY (status) REFERENCES statuses (code);

ALTER TABLE documents
    ADD CONSTRAINT documents_request_id_fkey
    FOREIGN KEY (request_id) REFERENCES requests (id);
