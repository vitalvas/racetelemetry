# Project CARS 2/3 Shared Memory Sink

Writes telemetry to the Project CARS 2 shared memory region. Windows only. Also compatible with Project CARS 3.

## Configuration

```yaml
sinks:
  my_pcars2_shm:
    type: pcars2_shm
    inputs:
      - my_source
```

## Platform

Windows only. The shared memory file is named `$pcars2$`. On non-Windows platforms, creating this sink returns an error.

CrewChief V4 reads from this shared memory region. The sequence number protocol is used: set to odd before writing, even after writing, so readers can detect torn frames.
