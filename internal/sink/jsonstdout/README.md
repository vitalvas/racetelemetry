# JSON Stdout Sink

Writes each telemetry frame as a JSON line to stdout. Useful for debugging, piping to other tools, or logging.

## Configuration

```yaml
sinks:
  debug:
    type: json_stdout
    inputs:
      - my_source
```

No additional fields required.

## Output Format

One JSON object per line (JSON Lines / NDJSON format). Each line contains the full `TelemetryFrame` struct serialized as JSON.
