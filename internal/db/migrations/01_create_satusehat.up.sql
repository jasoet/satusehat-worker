CREATE TABLE satusehat
(
    visit_id             TEXT PRIMARY KEY,
    visit_date           DATETIME NOT NULL,
    satusehat_patient_id TEXT     NOT NULL,
    visit_detail         TEXT     NOT NULL,
    vital_sign           TEXT     NOT NULL,
    diagnosis            TEXT,
    lab                  TEXT,
    radiology            TEXT,
    medication_request   TEXT,
    medication_dispense  TEXT,
    medical_procedure    TEXT,
    publish_date         DATETIME,
    publish_request      TEXT,
    publish_response     TEXT,
    mapping_errors       TEXT,
    publish_status       TEXT     NOT NULL,
    mapping_status       TEXT     NOT NULL
);