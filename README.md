# qdb

qdb is an SQL database wrapper that uses query structs for interacting with the underlying database. 
Any errors that a query produces are wrapped with the specific SQL query and arguments, making inspection and
debugging easier.

### Testing
qdb allows you to mock a database that responds to named queries with specific values. See the mockdb pacakge.