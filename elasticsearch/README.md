# Elasticsearch Hands-On Guide

This guide walks you through the **core concepts of Elasticsearch** using simple API calls (`curl`) or Kibana Dev Tools.
By running these step by step, you’ll see how clusters, nodes, indices, documents, shards, inverted indices, and analyzers actually work.

---

## 0. Test Elasticsearch

Check if Elasticsearch is running:

```bash
curl http://localhost:9200
```

-   Returns cluster details like:

    -   Cluster name
    -   UUID
    -   Elasticsearch version

---

## 1. Cluster

A **cluster** = group of nodes working together.

```bash
curl -X GET "localhost:9200/"
```

-   Confirms the cluster is alive
-   Shows:

    -   `cluster_name`
    -   `cluster_uuid`
    -   Elasticsearch `version`

---

## 2. Nodes

A **node** = a single Elasticsearch server inside the cluster.

```bash
curl -X GET "localhost:9200/_cat/nodes?v"
```

-   Returns all nodes in the cluster
-   Shows:

    -   Node name
    -   IP and port
    -   Roles (master, data, ingest, coordinating)
    -   Heap usage and uptime

-   In single-node mode → you’ll see just **1 node**

---

## 3. Index

An **index** = like a database table where documents are stored.

Create a new index:

```bash
curl -X PUT "localhost:9200/my-first-index"
```

-   Creates an empty index named `my-first-index`
-   Uses default settings (1 primary shard, 1 replica)

List all indices:

```bash
curl -X GET "localhost:9200/_cat/indices?v"
```

-   Shows:

    -   Index name
    -   Health (green/yellow/red)
    -   Document count
    -   Primary/replica shard count

---

## 4. Document

A **document** = like a row in a table, stored as JSON.

Add a new document:

```bash
curl -X POST "localhost:9200/my-first-index/_doc/1" -H 'Content-Type: application/json' -d '{
  "title": "Hello Elasticsearch",
  "content": "Elasticsearch makes search super fast!"
}'
```

-   Inserts a document with `_id=1`
-   Fields: `title` and `content`
-   Stored inside the `my-first-index` index

Retrieve the document:

```bash
curl -X GET "localhost:9200/my-first-index/_doc/1"
```

-   Fetches the JSON document by its ID
-   Returns stored data plus metadata (`_index`, `_id`, `_version`)

---

## 5. Shards

An **index is split into shards** for scalability.
Each shard is a self-contained Lucene instance.

Check shards:

```bash
curl -X GET "localhost:9200/_cat/shards?v"
```

-   Shows all shards in the cluster
-   Columns include:

    -   Index name
    -   Shard number
    -   Role (`p` = primary, `r` = replica)
    -   State (started, initializing, etc.)

-   In single-node mode → only **1 primary shard**, no replicas active

---

## 6. Inverted Index

The **inverted index** maps words → documents, making search lightning-fast.

Search for a word:

```bash
curl -X GET "localhost:9200/my-first-index/_search" -H 'Content-Type: application/json' -d '{
  "query": {
    "match": {
      "content": "search"
    }
  }
}'
```

-   Runs a full-text search on `content` field
-   Elasticsearch looks up the term `search` in its inverted index
-   Returns all documents containing that term

---

## 7. Analyzers

**Analyzers** process text using tokenizers and filters before indexing.

Test the default analyzer:

```bash
curl -X GET "localhost:9200/_analyze" -H 'Content-Type: application/json' -d '{
  "analyzer": "standard",
  "text": "The Quick Brown Foxes Jumped"
}'
```

-   Breaks text into tokens (words)
-   Lowercases everything
-   Result: `the`, `quick`, `brown`, `foxes`, `jumped`

---

## Quick Summary

-   **Cluster info** → `GET /`
-   **Nodes** → `GET /_cat/nodes`
-   **Index create/list** → `PUT /my-first-index`, `GET /_cat/indices`
-   **Documents add/get** → `POST /my-first-index/_doc`, `GET /_doc/{id}`
-   **Shards** → `GET /_cat/shards`
-   **Inverted Index search** → `/_search`
-   **Analyzers** → `/_analyze`

---

✅ Run these in order → you’ll see how Elasticsearch works under the hood.
