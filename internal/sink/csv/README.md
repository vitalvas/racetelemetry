# CSV Sink

Writes telemetry frames as CSV rows to a file. Compatible with MoTeC i2, Race Studio, Excel, and other analysis tools.

## Configuration

```yaml
sinks:
  telemetry_log:
    type: csv
    inputs:
      - my_source
    file_path: "telemetry.csv"
    csv_format: default
```

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `file_path` | yes | | Path to the output CSV file (created or overwritten) |
| `csv_format` | no | `default` | Column format: `default` or `forza` |

## Formats

### `default`

Snake_case column names with all unified model fields. Includes source metadata columns (`timestamp`, `source_name`, `source_type`) and fields from all source types.

### `forza`

PascalCase column names matching the official Forza packet specification. Compatible with:

- [richstokes/Forza-data-tools](https://github.com/richstokes/Forza-data-tools)
- [austinbaccus/forza-telemetry](https://github.com/austinbaccus/forza-telemetry)
- Other Forza community analysis tools

Columns follow the Forza packet field order: `IsRaceOn`, `TimestampMS`, `EngineMaxRpm`, `EngineIdleRpm`, `CurrentEngineRpm`, etc.

## Output Details

- First row is the header with column names
- One data row per telemetry frame
- Empty string for fields not provided by the source
- Per-wheel data is expanded to separate columns
- File is flushed after each row for real-time access
