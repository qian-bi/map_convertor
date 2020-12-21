# Map convertor

As customer request, we need to convert STIF map to the format as below.

1. INKLESS list file

    - The list file 1-8 is the header. The header information is the same as the CPC Map. Please see the note at the back of the //. TOTAL PASS DIE is the number of Lot Pass IC.
    - Line 10 is the field name. Copy it, separated by at least one space, all uppercase.
    - Field description:  
        - MAPPING_FILE_NAME: Inkless file name for each wafer of this Lot  
        - WAFER_ID: Consistent with CPC Map, cannot contain Lot ID  
        - GOOD: Number of Good Dies that can be packaged  
        - YIELD: corresponding yield of wafer

    ![list file](./doc/list%20file.png)

2. INKLESS MAP file \([sample](./doc/NPY651-01.map)\)

    - Line 1-15 is the header. The header information is the same as CPC Map. Please see the comment at the back.
    - Line 18 is the MAP start identifier
    - The beginning of line 19 is the content of the Inkless map. The Inkless Map states:
        - "1" means PASS DIE, to be encapsulated
        - "A, B, C, D . . . ", Partial Good Bin
        - "X"表示FAIL DIE, Not encapsulated
        - "." Indicates blank

    ![Map file](./doc/Map%20file.png)
