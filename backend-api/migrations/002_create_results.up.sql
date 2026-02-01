-- Create results table
CREATE TABLE IF NOT EXISTS results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    page_number INTEGER NOT NULL,
    original_path VARCHAR(512) NOT NULL,
    translated_path VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(request_id, page_number)
);

-- Create index for foreign key lookups
CREATE INDEX IF NOT EXISTS idx_results_request_id ON results(request_id);
