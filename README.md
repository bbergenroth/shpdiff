# shpdiff

Outputs features from one shapefile that are different (based on bounding box) or don't exist in another based on a key field

```
Usage: shpdiff file1.shp file2.shp diff.shp 
  -h    help
  -k string
        unique key to use as comparison (default: id)
  -v    verbose
Example: shpdiff new.shp old.shp diff.shp
```

TODO
- tests
- add deep compare
