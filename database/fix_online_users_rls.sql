-- Fix RLS for Online Users Counter
-- This allows the online users counter to work properly
-- Run this in Supabase SQL Editor

-- Allow anonymous users to read the count of online users
-- This is needed for the online users counter to work
CREATE POLICY "Allow anonymous to count online users"
    ON unlocked_users
    FOR SELECT
    TO anon
    USING (
        verified = true AND 
        last_active IS NOT NULL AND
        last_active >= (NOW() - INTERVAL '5 minutes')
    )
    WITH CHECK (false); -- anon can't insert/update

-- Also allow authenticated users to see online count
CREATE POLICY "Allow authenticated to count online users"
    ON unlocked_users
    FOR SELECT
    TO authenticated
    USING (
        verified = true AND 
        last_active IS NOT NULL AND
        last_active >= (NOW() - INTERVAL '5 minutes')
    );

-- Note: If you want to show ALL verified users (not just active), use this instead:
-- CREATE POLICY "Allow anonymous to count verified users"
--     ON unlocked_users
--     FOR SELECT
--     TO anon
--     USING (verified = true)
--     WITH CHECK (false);
