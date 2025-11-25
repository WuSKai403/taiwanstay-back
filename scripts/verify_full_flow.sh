#!/bin/bash

# Configuration
API_URL="http://localhost:8080/api/v1"
GUEST_EMAIL="guest_$(date +%s)@example.com"
HOST_EMAIL="host_$(date +%s)@example.com"
PASSWORD="password123"

# Helper function to parse JSON using python3
json_val() {
    python3 -c "import sys, json; print(json.load(sys.stdin)$1)"
}

echo "=== Starting Verification Flow ==="

# 1. Register Guest
echo "1. Registering Guest ($GUEST_EMAIL)..."
curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Guest User\",\"email\":\"$GUEST_EMAIL\",\"password\":\"$PASSWORD\"}" > /dev/null

# Login Guest
echo "   Logging in Guest..."
GUEST_TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"loginType\":\"password\",\"email\":\"$GUEST_EMAIL\",\"password\":\"$PASSWORD\"}" | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

if [ "$GUEST_TOKEN" == "" ]; then
    echo "Failed to login guest"
    exit 1
fi
echo "   Guest Token: ${GUEST_TOKEN:0:10}..."

# 2. Register Host
echo "2. Registering Host ($HOST_EMAIL)..."
curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Host User\",\"email\":\"$HOST_EMAIL\",\"password\":\"$PASSWORD\"}" > /dev/null

# Login Host
echo "   Logging in Host..."
HOST_TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"loginType\":\"password\",\"email\":\"$HOST_EMAIL\",\"password\":\"$PASSWORD\"}" | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

if [ "$HOST_TOKEN" == "" ]; then
    echo "Failed to login host"
    exit 1
fi
echo "   Host Token: ${HOST_TOKEN:0:10}..."

# Create Host Profile
echo "   Creating Host Profile..."
curl -s -X POST "$API_URL/hosts" \
  -H "Authorization: Bearer $HOST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Happy Farm","description":"Best farm in Taiwan"}' > /dev/null

# 3. Create Opportunity
echo "3. Creating Opportunity..."
OPP_ID=$(curl -s -X POST "$API_URL/opportunities" \
  -H "Authorization: Bearer $HOST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Rice Planting Helper",
    "description": "Help us plant rice",
    "type": "Farming",
    "location": {
        "city": "Hualien",
        "country": "Taiwan",
        "coordinates": {"type": "Point", "coordinates": [121.6, 23.9]}
    },
    "hasTimeSlots": true,
    "timeSlots": [
        {"startDate": "2023-07-01", "endDate": "2023-07-31", "status": "OPEN", "defaultCapacity": 2}
    ]
  }' | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")

if [ "$OPP_ID" == "" ]; then
    echo "Failed to create opportunity"
    exit 1
fi
echo "   Opportunity Created: $OPP_ID"

# 4. Search Opportunity
echo "4. Searching Opportunity..."
SEARCH_RESULT=$(curl -s -G "$API_URL/opportunities/search" \
    --data-urlencode "q=Rice" \
    --data-urlencode "city=Hualien" \
    --data-urlencode "startDate=2023-07-05" \
    --data-urlencode "endDate=2023-07-10")

FOUND_ID=$(echo $SEARCH_RESULT | python3 -c "import sys, json; print(json.load(sys.stdin)['data'][0]['id'])")
if [ "$FOUND_ID" != "$OPP_ID" ]; then
    echo "Search failed: Expected $OPP_ID, got $FOUND_ID"
    exit 1
fi
echo "   Found Opportunity: $FOUND_ID"

# 5. Apply for Opportunity
echo "5. Applying for Opportunity..."
APP_ID=$(curl -s -X POST "$API_URL/applications" \
  -H "Authorization: Bearer $GUEST_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"opportunityId\": \"$OPP_ID\",
    \"applicationDetails\": {
        \"message\": \"I love rice!\",
        \"startDate\": \"2023-07-05\",
        \"endDate\": \"2023-07-10\"
    }
  }" | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")

if [ "$APP_ID" == "" ]; then
    echo "Failed to apply"
    exit 1
fi
echo "   Application Created: $APP_ID"

# 6. Host Reviews Application
echo "6. Host Reviewing Application..."
# List Applications
# Need Host ID first
HOST_ID=$(curl -s -X GET "$API_URL/hosts/me" -H "Authorization: Bearer $HOST_TOKEN" | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")

LIST_RESULT=$(curl -s -G "$API_URL/applications" \
    -H "Authorization: Bearer $HOST_TOKEN" \
    --data-urlencode "hostId=$HOST_ID")

FOUND_APP_ID=$(echo $LIST_RESULT | python3 -c "import sys, json; print(json.load(sys.stdin)['data'][0]['id'])")
if [ "$FOUND_APP_ID" != "$APP_ID" ]; then
    echo "List failed: Expected $APP_ID, got $FOUND_APP_ID"
    exit 1
fi

# Accept Application
echo "   Accepting Application..."
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$API_URL/applications/$APP_ID" \
  -H "Authorization: Bearer $HOST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "ACCEPTED", "note": "Welcome!"}')

if [ "$STATUS_CODE" != "200" ]; then
    echo "Failed to accept application. Status: $STATUS_CODE"
    exit 1
fi

# 7. Verify Status
echo "7. Verifying Status..."
FINAL_STATUS=$(curl -s -X GET "$API_URL/applications/$APP_ID" \
  -H "Authorization: Bearer $GUEST_TOKEN" | python3 -c "import sys, json; print(json.load(sys.stdin)['status'])")

if [ "$FINAL_STATUS" == "ACCEPTED" ]; then
    echo "=== SUCCESS: Full Flow Verified! ==="
else
    echo "=== FAILED: Status is $FINAL_STATUS ==="
    exit 1
fi

# 8. Verify Notifications
echo "8. Verifying Notifications..."
# Guest should have a notification about status change
NOTIFS=$(curl -s -X GET "$API_URL/users/me/notifications" \
  -H "Authorization: Bearer $GUEST_TOKEN")

NOTIF_COUNT=$(echo $NOTIFS | python3 -c "import sys, json; print(json.load(sys.stdin)['total'])")
if [ "$NOTIF_COUNT" -gt 0 ]; then
    echo "   Guest has $NOTIF_COUNT notifications."
    TITLE=$(echo $NOTIFS | python3 -c "import sys, json; print(json.load(sys.stdin)['data'][0]['title'])")
    echo "   Latest Notification: $TITLE"
else
    echo "   WARNING: Guest has no notifications!"
fi
