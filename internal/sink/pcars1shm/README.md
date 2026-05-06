# Project CARS 1 Shared Memory Sink

Writes telemetry to the Project CARS 1 shared memory region. Windows only.

## Configuration

```yaml
sinks:
  my_pcars1_shm:
    type: pcars1_shm
    inputs:
      - my_source
```

## Platform

Windows only. The shared memory file is named `$pcars$`. On non-Windows platforms, creating this sink returns an error.
