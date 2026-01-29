#!/bin/bash

# Business Consultant - Check App Registration
# This script checks if the Business Consultant app is registered in the system

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Business Consultant - Check Registration${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
API_URL="${API_URL:-https://i149gvmuh8.execute-api.us-east-1.amazonaws.com/prod}"
APP_NAME="商业顾问"

# Check if JWT token is provided
if [ -z "$JWT_TOKEN" ]; then
  echo -e "${YELLOW}JWT_TOKEN not set. Please provide your JWT token.${NC}"
  read -p "Enter your JWT token: " JWT_TOKEN
  echo ""
fi

if [ -z "$JWT_TOKEN" ]; then
  echo -e "${RED}Error: JWT token is required${NC}"
  exit 1
fi

echo -e "${BLUE}Checking app registration...${NC}"
echo ""

# Get all apps
RESPONSE=$(curl -s -X GET "$API_URL/api/apps" \
  -H "Authorization: Bearer $JWT_TOKEN")

# Check if request was successful
if [ $? -ne 0 ]; then
  echo -e "${RED}❌ Error: Failed to connect to API${NC}"
  exit 1
fi

# Parse response
APP_INFO=$(echo "$RESPONSE" | jq -r ".data[] | select(.app_name == \"$APP_NAME\")")

if [ -z "$APP_INFO" ] || [ "$APP_INFO" == "null" ]; then
  echo -e "${RED}❌ App '$APP_NAME' is NOT registered${NC}"
  echo ""
  echo -e "${YELLOW}To register the app, run:${NC}"
  echo "  cd business-consultant/scripts"
  echo "  ./register-app.sh"
  exit 1
fi

# Extract app details
APP_ID=$(echo "$APP_INFO" | jq -r '.app_id')
APP_URL=$(echo "$APP_INFO" | jq -r '.url')
APP_EMOJI=$(echo "$APP_INFO" | jq -r '.emoji')
APP_DESC=$(echo "$APP_INFO" | jq -r '.app_description')
IS_GLOBAL=$(echo "$APP_INFO" | jq -r '.is_global')
CREATED_AT=$(echo "$APP_INFO" | jq -r '.created_at')

echo -e "${GREEN}✅ App is registered!${NC}"
echo ""
echo -e "${BLUE}App Details:${NC}"
echo "  Name: $APP_NAME $APP_EMOJI"
echo "  ID: $APP_ID"
echo "  URL: $APP_URL"
echo "  Description: $APP_DESC"
echo "  Is Global: $IS_GLOBAL"
echo "  Created: $CREATED_AT"
echo ""

if [ "$IS_GLOBAL" == "true" ]; then
  echo -e "${GREEN}✓ App is set as global (visible to all users)${NC}"
else
  echo -e "${YELLOW}⚠ App is NOT global (only visible to specific projects)${NC}"
  echo ""
  echo -e "${YELLOW}To set as global, run:${NC}"
  echo "  curl -X POST \"$API_URL/api/apps/$APP_ID/set-global\" \\"
  echo "    -H \"Authorization: Bearer \$JWT_TOKEN\" \\"
  echo "    -H \"Content-Type: application/json\" \\"
  echo "    -d '{\"is_global\": true}'"
fi

echo ""
echo -e "${BLUE}Access the app:${NC}"
echo "1. Go to: https://main.d2fozf421c6ftf.amplifyapp.com"
echo "2. Login and select any project"
echo "3. Click on '$APP_NAME $APP_EMOJI' in the apps list"
echo "4. You will be redirected to: $APP_URL"
