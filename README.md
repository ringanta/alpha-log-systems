# Alpha Log Systems

Alpha log system consists of two components. They are AlphaClient and AlphaServer.
AlphaClient monitors ssh log files and send the event to AlphaServer.
AlphaServer receives events from AlphaClient and displays number of ssh attempt based on event from AlphaClient.

## System Architecture

Here is system architecture for the solution ![System architecture](./docs/images/alpha-system.png)

Some of notable implementations decisions:
- AlphaClient sends raw SSH attempt event to AlphaServer. AlphaClient using HTTP POST with basic token authentication for simplicity of implementation.
- AlphaServer responsible for deduplication of event since the limited functionality provided by library that's available to monitor SSH attempt log file.
- AlphaServer stores SSH attempt event to Postgres database for simplicity on the implementor side.
- AlphaServer queries databases to calculate SSH attempt for each hosts and display attempt metrics on a simple static web page for simplicity of implementation.

