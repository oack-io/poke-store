#!/usr/bin/env bash
set -euo pipefail

# Generate pixel art Pokemon images via OpenAI gpt-image-1 API.
# Usage: OPENAI_API_KEY=sk-... ./scripts/gen-pixel-art.sh

OUTDIR="web/public/pokemon"
MODEL="gpt-image-1"
SIZE="1024x1024"

POKEMON=(
  articuno blastoise bulbasaur charizard charmander dragonite eevee gengar
  geodude gyarados jigglypuff lapras magikarp meowth mew mewtwo moltres
  pikachu psyduck raichu snorlax squirtle venusaur vulpix zapdos
)

if [[ -z "${OPENAI_API_KEY:-}" ]]; then
  echo "Error: OPENAI_API_KEY is not set" >&2
  exit 1
fi

mkdir -p "$OUTDIR"

for name in "${POKEMON[@]}"; do
  outfile="$OUTDIR/$name.png"
  echo "Generating $name..."

  prompt="Pixel art sprite of a cute fantasy creature called $name on a plain white background. 16-bit retro game style, clean pixels, no anti-aliasing, centered, full body visible."

  response=$(curl -s "https://api.openai.com/v1/images/generations" \
    -H "Authorization: Bearer $OPENAI_API_KEY" \
    -H "Content-Type: application/json" \
    -d "{
      \"model\": \"$MODEL\",
      \"prompt\": \"$prompt\",
      \"n\": 1,
      \"size\": \"$SIZE\",
      \"output_format\": \"png\"
    }")

  # Check for errors
  error=$(echo "$response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('error',{}).get('message',''))" 2>/dev/null || echo "")
  if [[ -n "$error" ]]; then
    echo "  ERROR: $error" >&2
    continue
  fi

  # Extract base64 and decode to file
  echo "$response" | python3 -c "
import sys, json, base64
data = json.load(sys.stdin)
b64 = data['data'][0]['b64_json']
sys.stdout.buffer.write(base64.b64decode(b64))
" > "$outfile"

  echo "  Saved $outfile ($(du -h "$outfile" | cut -f1))"
done

echo "Done! Generated ${#POKEMON[@]} images."
