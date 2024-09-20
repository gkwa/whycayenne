# DynamoDB Tutorial: Single-Table Design and Efficient Querying

Welcome to this friendly, no-stress tutorial on Amazon DynamoDB! We'll explore single-table design and efficient querying using a real-world example of storing and searching for pepper products.

## Introduction to DynamoDB

Amazon DynamoDB is a fully managed NoSQL database service that provides fast and predictable performance with seamless scalability. It's designed to run high-performance applications at any scale.

## Single-Table Design

In DynamoDB, a single-table design means storing all your application's data in one table. This approach can lead to more efficient queries and better performance, especially for complex data models.

## Sample Data

Let's consider the sample data you provided:

```json
{
  "name": "Vigo Pepperoncini Greek Peppers",
  "price": 33.78,
  "price_per_lb": 33.78,
  "price_per_oz": 2.11,
  "store": "Kroger",
  "volume": "16 oz",
  "weight": 1.0,
  "on_sale": true,
  "datetime": "2024-09-19T12:00:00Z"
}

{
  "name": "Jeffs Garden Pepper Whole Greek Peperoncini",
  "price": 46.99,
  "price_per_lb": 10.44,
  "price_per_oz": 0.65,
  "store": "Kroger",
  "volume": "72 oz",
  "weight": 4.5,
  "on_sale": true,
  "datetime": "2024-09-19T12:00:00Z"
}
```

## Table Structure

To design an efficient single-table structure for this data, we need to consider our access patterns. Let's assume we want to:

1. Find products by name
2. Search for peppers
3. Get products by store

Here's a proposed table structure:

- Partition Key: `PK` (String)
- Sort Key: `SK` (String)

We'll use a composite key strategy:

- For products: `PK = "PRODUCT#<product_name>"`, `SK = "METADATA#<store>"
- For categories: `PK = "CATEGORY#<category_name>"`, `SK = "PRODUCT#<product_name>"

## Querying for Peppers

To efficiently query for peppers without performing a full table scan, we can use a Global Secondary Index (GSI) with an inverted index pattern:

1. Create a GSI with:

   - Partition Key: `GSI1PK` (String)
   - Sort Key: `GSI1SK` (String)

2. When inserting items, add these attributes:
   - `GSI1PK = "CATEGORY#pepper"`
   - `GSI1SK = "PRODUCT#<product_name>"`

Now, to search for peppers, you can query the GSI using:

```
GSI1PK = "CATEGORY#pepper"
```

This query will return all products categorized as peppers without performing a full table scan.

## Avoiding Table Scans

As you mentioned, table scans are expensive and should be avoided when possible. Here are some strategies to avoid scans:

1. Use appropriate partition and sort keys
2. Utilize Global Secondary Indexes (GSIs) for alternative access patterns
3. Use the `Query` operation instead of `Scan` whenever possible
4. If you must use `Scan`, consider using parallel scans to improve performance

## Best Practices

1. Denormalize your data to fit the single-table model
2. Use consistent naming conventions for keys
3. Use sparse indexes to reduce storage and improve query performance
4. Consider using DynamoDB Streams for real-time data processing
5. Regularly review and optimize your access patterns

By following these practices and designing your table structure thoughtfully, you can make the most of DynamoDB's performance capabilities while avoiding expensive table scans.

Remember, the key to efficient DynamoDB usage is understanding your data access patterns and designing your table structure accordingly.
