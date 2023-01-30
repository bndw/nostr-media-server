HOST=localhost:9000

fn="${1:-example.gif}"

curl \
  -F "filename=$fn" \
  -F "file=@./$fn" \
  "${HOST}/upload"
