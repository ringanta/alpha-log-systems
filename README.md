# Alpha Log Systems

Alpha log system consists of two components. They are AlphaClient and AlphaServer.
AlphaClient monitors ssh log files and send the event to AlphaServer.
AlphaServer receives events from AlphaClient and displays number of ssh attempt based on event from AlphaClient.

## System Architecture

Frameworks and libraries used:
- [Fiber](https://github.com/gofiber/fiber)
