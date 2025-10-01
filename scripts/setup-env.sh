#!/bin/bash

# Clone the repository if not already cloned
if [ ! -d ".git" ]; then
    git clone https://github.com/vi88i/eeye.git .
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

# Setup environment variables
echo "Setting up environment variables..."
if [ ! -f ".env" ]; then
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo "Created .env file from template"
        echo "Please edit .env and set your GROWW_ACCESS_TOKEN"
    else
        echo "Error: .env.example not found"
        exit 1
    fi
else
    echo ".env file already exists"
fi

echo "Setup complete!"
echo "Don't forget to set your GROWW_ACCESS_TOKEN in the .env file"