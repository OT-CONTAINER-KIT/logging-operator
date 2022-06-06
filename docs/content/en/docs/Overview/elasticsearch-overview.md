---
title: "Elasticsearch"
weight: 2
description: >
    An overview of Elasticsearch database/search-engine, types of setup and architecture design
---

Elasticsearch is a distributed open-source search and analytics engine built on Java and Apache Lucene. It allows us to store, search and analyze huge chunk of data with nearly real time and high performance. It is a REST API based system on which we can easily write and query the data, in easy words we can say that Elasticsearch is a server that can process JSON requests and returns JSON response.

There are different use cases for elasticsearch like:-

- NoSQL database
- Logs storage and searching
- Real time and time series analysis

<div align="center">
    <img src="https://miro.medium.com/max/558/1*AYP0Mg_MwJMm3Kbx8Xa8lQ.png">
</div>

## Features

- **Scalability:** It is scalable across multiple nodes. This means we can start with less number of nodes and in case our workload increases then we can scale across multiple nodes. It is easily scalable.
- **Fast:** It is really fast in terms of performance when compared to other search engines that are available.
- **Multilingual:** It supports various languages.
- **Document Oriented:** Instead of schemas and tables, the data is stored in documents. All the data is stored in JSON format. JSON is the widely accepted web format due to which we can easily integrate the generated output in other applications if required.
- **Auto-completion:** It returns documents that contain a specific prefix in a provided field.

## Elasticsearch Architecture

<div align="center">
    <img src="https://static.packt-cdn.com/products/9781789957754/graphics/assets/664ba9c8-5e54-42a4-ba90-635dc8d82276.png">
</div>
