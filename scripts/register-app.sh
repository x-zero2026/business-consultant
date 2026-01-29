#!/bin/bash

# Business Consultant - Register App Script
# This script registers the Business Consultant app in the DID Login system

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Business Consultant - App Registration${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
API_URL="${API_URL:-https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod}"
PROJECT_NAME="${PROJECT_NAME:-system}"
APP_NAME="商业顾问"
APP_URL="${APP_URL:-http://localhost:5173}"
APP_DESCRIPTION="AI驱动的一人公司商业咨询服务，帮助规划资源配置和预算"
APP_EMOJI="${APP_EMOJI:-👔}"
IS_GLOBAL="${IS_GLOBAL:-true}"

# Check if JWT token is provided
if [ -z "$JWT_TOKEN" ]; then
  echo -e "${YELLOW}JWT_TOKEN not set. Please provide your JWT token.${NC}"
  echo ""
  echo "You can get your token by:"
  echo "1. Login to DID Login UI (https://main.d2fozf421c6ftf.amplifyapp.com)"
  echo "2. Open browser DevTools (F12)"
  echo "3. Go to Application → Local Storage"
  echo "4. Copy the 'token' value"
  echo ""
  read -p "Enter your JWT token: " JWT_TOKEN
  echo ""
fi

# Validate token
if [ -z "$JWT_TOKEN" ]; then
  echo -e "${RED}Error: JWT token is required${NC}"
  exit 1
fi

echo -e "${BLUE}Configuration:${NC}"
echo "  API URL: $API_URL"
echo "  Project Name: $PROJECT_NAME"
echo "  App Name: $APP_NAME"
echo "  App URL: $APP_URL"
echo "  Description: $APP_DESCRIPTION"
echo "  Emoji: $APP_EMOJI"
echo "  Is Global: $IS_GLOBAL"
echo ""

# Get project ID by name
echo -e "${BLUE}Finding project '$PROJECT_NAME'...${NC}"

PROJECTS_RESPONSE=$(curl -s -X GET "$API_URL/api/projects" \
  -H "Authorization: Bearer $JWT_TOKEN")

PROJECT_ID=$(echo "$PROJECTS_RESPONSE" | jq -r ".data[] | select(.project_name == \"$PROJECT_NAME\") | .project_id")

if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" == "null" ]; then
  echo -e "${RED}Error: Project '$PROJECT_NAME' not found${NC}"
  echo ""
  echo -e "${YELLOW}Available projects:${NC}"
  echo "$PROJECTS_RESPONSE" | jq -r '.data[] | "  - \(.project_name) (\(.project_id))"'
  exit 1
fi

echo -e "${GREEN}✓ Found project: $PROJECT_ID${NC}"
echo ""

# Check if app already exists
echo -e "${BLUE}Checking if app already exists...${NC}"

EXISTING_APPS=$(curl -s -X GET "$API_URL/api/projects/$PROJECT_ID/apps" \
  -H "Authorization: Bearer $JWT_TOKEN")

APP_EXISTS=$(echo "$EXISTING_APPS" | jq -r ".data[] | select(.app_name == \"$APP_NAME\") | .app_id")

if [ -n "$APP_EXISTS" ] && [ "$APP_EXISTS" != "null" ]; then
  echo -e "${YELLOW}⚠ App '$APP_NAME' already exists (ID: $APP_EXISTS)${NC}"
  echo ""
  read -p "Do you want to update it? (y/n): " UPDATE_CHOICE
  if [ "$UPDATE_CHOICE" != "y" ]; then
    echo -e "${BLUE}Aborted.${NC}"
    exit 0
  fi
  
  # Update existing app
  echo -e "${BLUE}Updating app...${NC}"
  
  RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL/api/apps/$APP_EXISTS" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"app_name\": \"$APP_NAME\",
      \"url\": \"$APP_URL\",
      \"app_description\": \"$APP_DESCRIPTION\",
      \"emoji\": \"$APP_EMOJI\"
    }")
  
  HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
  BODY=$(echo "$RESPONSE" | sed '$d')
  
  if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ]; then
    echo -e "${GREEN}✓ App updated${NC}"
    
    # Set as global
    if [ "$IS_GLOBAL" == "true" ]; then
      echo -e "${BLUE}Setting app as global...${NC}"
      
      SET_GLOBAL_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/apps/$APP_EXISTS/set-global" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"is_global\": true}")
      
      SET_GLOBAL_HTTP_CODE=$(echo "$SET_GLOBAL_RESPONSE" | tail -n1)
      
      if [ "$SET_GLOBAL_HTTP_CODE" -eq 200 ]; then
        echo -e "${GREEN}✓ App set as global${NC}"
      else
        echo -e "${YELLOW}⚠ Could not set app as global${NC}"
      fi
    fi
    
    echo ""
    echo -e "${GREEN}✅ Success! App updated successfully.${NC}"
  else
    echo -e "${RED}❌ Error: Failed to update app (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
    exit 1
  fi
else
  # Create new app
  echo -e "${BLUE}Creating app...${NC}"
  
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/apps" \
    -H "Authorization: Bearer $JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"project_id\": \"$PROJECT_ID\",
      \"app_name\": \"$APP_NAME\",
      \"url\": \"$APP_URL\",
      \"app_description\": \"$APP_DESCRIPTION\",
      \"emoji\": \"$APP_EMOJI\"
    }")
  
  HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
  BODY=$(echo "$RESPONSE" | sed '$d')
  
  if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ]; then
    echo -e "${GREEN}✓ App created${NC}"
    
    # Extract app_id from response
    NEW_APP_ID=$(echo "$BODY" | jq -r '.data.app_id')
    
    # Set as global
    if [ "$IS_GLOBAL" == "true" ] && [ -n "$NEW_APP_ID" ] && [ "$NEW_APP_ID" != "null" ]; then
      echo -e "${BLUE}Setting app as global...${NC}"
      
      SET_GLOBAL_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/apps/$NEW_APP_ID/set-global" \
        -H "Authorization: Bearer $JWT_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"is_global\": true}")
      
      SET_GLOBAL_HTTP_CODE=$(echo "$SET_GLOBAL_RESPONSE" | tail -n1)
      
      if [ "$SET_GLOBAL_HTTP_CODE" -eq 200 ]; then
        echo -e "${GREEN}✓ App set as global${NC}"
      else
        echo -e "${YELLOW}⚠ Could not set app as global${NC}"
      fi
    fi
    
    echo ""
    echo -e "${GREEN}✅ Success! App registered successfully.${NC}"
  else
    echo -e "${RED}❌ Error: Failed to create app (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
    exit 1
  fi
fi

echo ""
echo -e "${BLUE}Response:${NC}"
echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "1. Go to DID Login Dashboard: https://main.d2fozf421c6ftf.amplifyapp.com"
echo "2. Select any project"
echo "3. You should see '$APP_NAME $APP_EMOJI' in the apps list (visible to all users)"
echo "4. Click on it to open: $APP_URL"
echo ""
echo -e "${BLUE}Note:${NC} This app is set as global (is_global=true), so it will be"
echo "visible to all users across all projects."
