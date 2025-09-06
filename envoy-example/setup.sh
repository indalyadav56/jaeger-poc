#!/bin/bash

# Create directory structure
mkdir -p html1 html2

# Create sample HTML for web1
cat > html1/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Web Service 1</title>
    <style>
        body { font-family: Arial; text-align: center; padding: 50px; background: #e8f4fd; }
        h1 { color: #2196F3; }
    </style>
</head>
<body>
    <h1>Web Service 1</h1>
    <p>This is served by the first nginx container (web1)</p>
    <p>Load balanced by Envoy Proxy</p>
</body>
</html>
EOF

# Create sample HTML for web2
cat > html2/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>Web Service 2</title>
    <style>
        body { font-family: Arial; text-align: center; padding: 50px; background: #f3e5f5; }
        h1 { color: #9C27B0; }
    </style>
</head>
<body>
    <h1>Web Service 2</h1>
    <p>This is served by the second nginx container (web2)</p>
    <p>Load balanced by Envoy Proxy</p>
</body>
</html>
EOF

echo "Setup complete!"
echo ""
echo "To run the example:"
echo "1. chmod +x setup.sh && ./setup.sh"
echo "2. docker-compose up"
echo ""
echo "Then visit:"
echo "- http://localhost:8080 (Load balanced web services)"
echo "- http://localhost:8080/api (API service)"
echo "- http://localhost:9901 (Envoy admin interface)"