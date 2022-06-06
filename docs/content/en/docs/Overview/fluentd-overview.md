---
title: "Fluentd"
weight: 3
description: >
    A detailed guide for designing the setup of Fluentd architecture
---

Fluentd is an open source data collector, which lets you unify the data collection and consumption for a better use and understanding of data.
Fluentd tries to structure data as JSON as much as possible: this allows Fluentd to unify all facets of processing log data: collecting, filtering, buffering, and outputting logs across multiple sources and destinations.

<div align="center">
    <img src="https://i.pinimg.com/originals/c2/c9/6b/c2c96be47b8b2abb758833628088808a.png" width="500" height="200">
</div>

## Features

- **JSON Logging:** Fluentd tries to structure data as JSON as much as possible: this allows Fluentd to unify all facets of processing log data: collecting, filtering, buffering, and outputting logs across multiple sources and destinations.
- **Pluggable Architecture:** Fluentd has a flexible plugin system that allows the community to extend its functionality. Our 500+ community-contributed plugins connect dozens of data sources and data outputs.
- **Minimum Resources Required:** Fluentd is written in a combination of C language and Ruby, and requires very little system resource. The vanilla instance runs on 30-40MB of memory and can process 13,000 events/second/core.
- **Built-in Reliability:** Fluentd supports memory- and file-based buffering to prevent inter-node data loss. Fluentd also supports robust failover and can be set up for high availability.

## Architecture

<div align="center">
    <img src="https://fluentbit.io/images/blog/blog-forwarder-aggregator.png">
</div>
