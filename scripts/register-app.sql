-- Business Consultant - Register App via SQL
-- Run this in your Supabase SQL Editor if the API method doesn't work

-- Step 1: Insert the app
-- Note: Replace 'your-did-here' with your actual DID
INSERT INTO apps (
  app_id,
  app_name,
  app_description,
  emoji,
  url,
  is_global,
  created_by_did,
  created_at,
  updated_at
) VALUES (
  gen_random_uuid(),
  '商业顾问',
  'AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算',
  '👔',  -- Tie emoji for business consultant
  'http://localhost:5173',  -- Change to production URL when deploying
  true,  -- Set to true to make it visible to all users across all projects
  'your-did-here',  -- IMPORTANT: Replace with your actual DID
  NOW(),
  NOW()
)
ON CONFLICT (app_name) DO UPDATE SET
  url = EXCLUDED.url,
  app_description = EXCLUDED.app_description,
  emoji = EXCLUDED.emoji,
  is_global = EXCLUDED.is_global,
  updated_at = NOW();

-- Step 2: Verify the app was created
SELECT 
  app_id,
  app_name,
  app_description,
  emoji,
  url,
  is_global,
  created_by_did,
  created_at
FROM apps
WHERE app_name = '商业顾问';

-- Step 3: Check if app is visible globally
-- Since is_global = true, this app should be visible to all users
-- No need to add to app_projects table

-- Step 4: View all global apps
SELECT 
  app_id,
  app_name,
  emoji,
  url,
  is_global,
  created_by_did
FROM apps
WHERE is_global = true
ORDER BY created_at DESC;

-- Step 5: Get your DID (if you don't know it)
-- Run this to find your DID:
/*
SELECT did, username, email 
FROM users 
WHERE username = 'your-username';
*/

-- Step 6 (Optional): Update production URL after deployment
-- Uncomment and run after deploying to Amplify
/*
UPDATE apps
SET 
  url = 'https://your-amplify-url.amplifyapp.com',
  updated_at = NOW()
WHERE app_name = '商业顾问';
*/

-- Step 7: Verify all apps in the system
SELECT 
  app_name,
  emoji,
  is_global,
  url,
  created_at
FROM apps
ORDER BY 
  is_global DESC,
  created_at DESC;
