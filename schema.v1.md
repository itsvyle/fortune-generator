# V1 Schema

-   1 byte: version info
-   1 byte: number of files
-   Files
    -   1 byte: file id (incremental probably)
    -   1 byte: file weight (probably not going to be used in v1)
    -   1 byte: file name length
    -   n bytes: file name
-   10 bytes: padding
-   4 bytes: number of entries
-   n entries
    -   1 byte: file id
    -   4 bytes: entry offset
    -   4 bytes: entry length
    -   1 empty byte
