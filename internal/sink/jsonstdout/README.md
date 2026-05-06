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

One JSON object per line (JSON Lines / NDJSON format). Each line contains the full `TelemetryFrame` struct serialized as JSON. Fields not provided by the source are omitted.

Every frame includes `source_name` (the config key, e.g. `forza_xbox`) and `source_type` (the source type, e.g. `forza`).

```json
{"timestamp":"...","source_name":"forza_xbox","source_type":"forza","is_race_on":true,"engine_rpm":1120.62,...}
```
