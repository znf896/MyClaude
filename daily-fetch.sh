#!/bin/bash

# Daily GitHub Top Projects Fetcher
# Automatically fetch top projects for different topics and save to dated markdown file

# Configuration
BINARY="./github-top"
OUTPUT_DIR="/Users/zhangzhanghaimin/Documents/48-AI项目/output"
GITHUB_TOKEN="${GITHUB_TOKEN:-}"

# Topics to fetch
TOPICS=(
    "ai artificial-intelligence machine-learning"
    "ai-agent autonomous-agents"
    "claude-code claude-ai"
    "ai-testing automation-testing"
)

echo "=== Daily GitHub Top Projects Fetcher ==="
echo "Date: $(date)"
echo "Output directory: $OUTPUT_DIR"
echo "Number of topics: ${#TOPICS[@]}"
echo ""

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: $BINARY not found. Please build it first with:"
    echo "  go build -o github-top ./cmd/github-top/"
    exit 1
fi

# Create output directory if not exists
mkdir -p "$OUTPUT_DIR"

# Fetch each topic
for i in "${!TOPICS[@]}"; do
    topic="${TOPICS[$i]}"
    echo "[$((i+1))/${#TOPICS[@]}] Fetching topic: $topic"

    if ./github-top -count 10 -query "$topic" -min-stars 1000; then
        echo "✓ Completed topic: $topic"
    else
        echo "✗ Failed to fetch topic: $topic"
    fi
    echo ""
done

echo "=== All done! ==="
echo "Results saved to: $OUTPUT_DIR/github-top-$(date +%Y-%m-%d).md"
